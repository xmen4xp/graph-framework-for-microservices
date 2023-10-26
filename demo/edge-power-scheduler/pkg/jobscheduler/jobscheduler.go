package jobscheduler

import (
	"context"
	nexus_client "powerschedulermodel/build/nexus-client"
	"time"

	edcv1 "powerschedulermodel/build/apis/edgedc.intel.com/v1"
	jiv1 "powerschedulermodel/build/apis/jobmgmt.intel.com/v1"
	jsv1 "powerschedulermodel/build/apis/jobscheduler.intel.com/v1"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// look at all the jobs that are pending execution
// sort by creation time
// schedule the jobs on different nodes, may be load all the edges.
// periodically collect stats and update the jobs

type JobScheduler struct {
	cfg       *nexus_client.ConfigConfig
	inv       *nexus_client.InventoryInventory
	dcfg      *nexus_client.DesiredconfigDesiredEdgeConfig
	jobSplits uint32
	log       *logrus.Entry
}

func New(cfg *nexus_client.ConfigConfig, inv *nexus_client.InventoryInventory,
	dcfg *nexus_client.DesiredconfigDesiredEdgeConfig) *JobScheduler {
	log := logrus.WithFields(logrus.Fields{
		"module": "jobscheduler",
	})
	return &JobScheduler{cfg, inv, dcfg, 4, log}
}

// check if an edge is free
func (js *JobScheduler) isEdgeFree(ctx context.Context, edge *nexus_client.EdgeEdge) bool {
	dcEdge, e := js.dcfg.GetEdgesDC(ctx, edge.DisplayName())
	if e != nil {
		if nexus_client.IsChildNotFound(e) {
			return true
		}
		js.log.Error(e)
		return false
	}
	jobs, e := dcEdge.GetAllJobsInfo(ctx)
	for _, job := range jobs {
		if job.Status.State.Progress != 100 {
			return false
		}
	}
	return true
}

// find the next job in line for execution sorted based on arrival time.
func (js *JobScheduler) nextInLineJob(jlist []*nexus_client.JobschedulerJob) (bool, int) {
	first := true
	foundIdx := 0
	foundTimeStamp := int64(0)
	for idx, job := range jlist {
		preq := uint32(0)
		for _, pex := range job.Status.State.Execution {
			preq += pex.PowerRequested
		}
		if job.Spec.PowerNeeded > preq {
			if job.Spec.CreationTime < foundTimeStamp || first {
				foundTimeStamp = job.Spec.CreationTime
				foundIdx = idx
				first = false
			}
		}
	}
	return !first, foundIdx
}

// allocate the job on the free edges.
func (js *JobScheduler) allocateJob(ctx context.Context, job *nexus_client.JobschedulerJob, edgeList []*nexus_client.EdgeEdge, freeEdge []int) {
	// job.Spec.PowerNeeded
	preq := uint32(0)
	for _, ex := range job.Status.State.Execution {
		preq += ex.PowerRequested

	}
	// split poewrneeded into x parts
	psplit := job.Spec.PowerNeeded / uint32(js.jobSplits)
	// edge case for rounding error
	psplitLast := job.Spec.PowerNeeded - psplit*(js.jobSplits-1)

	for _, fedge := range freeEdge {
		pneeded := psplit
		if preq+psplitLast == job.Spec.PowerNeeded {
			pneeded = psplitLast
		}
		preq += pneeded // allocated power

		edge := edgeList[fedge]
		edgeName := edge.DisplayName()
		js.log.Infof("Allocationg a job name = %s on edge %s %s", job.DisplayName(), edgeName, edge.Name)
		curTime := time.Now().Unix()
		// create the desired state config
		edc, e := js.dcfg.GetEdgesDC(ctx, edgeName)
		if e != nil {
			newEdge := &edcv1.EdgeDC{}
			newEdge.SetName(edgeName)
			edc, e = js.dcfg.AddEdgesDC(ctx, newEdge)
			if e != nil {
				js.log.Error("Creating Desired edge config", e)
			}
		}
		jobName := job.DisplayName() + "-" + edgeName
		jobInfo := &jiv1.JobInfo{}
		jobInfo.SetName(jobName)
		jobInfo.Spec.RequestorJob = job.DisplayName()
		jobInfo.Spec.PowerNeeded = pneeded
		_, e = edc.AddJobsInfo(ctx, jobInfo)
		if e != nil {
			js.log.Error("When added job info in edgedc", e)
		}
		// update the requester
		if job.Status.State.Execution == nil {
			job.Status.State.Execution = make(map[string]jsv1.NodeExecutionStatus)
		}
		if len(job.Status.State.Execution) == 0 {
			job.Status.State.StartTime = curTime
		}
		job.Status.State.Execution[edgeName] = jsv1.NodeExecutionStatus{
			PowerRequested: pneeded,
			StartTime:      0,
			EndTime:        0,
			Progress:       0,
		}
		job.Status.State.EndTime = curTime
		// commit the data back
		if e := edc.Update(ctx); e != nil {
			js.log.Error("When updating job spec", e)
		}
		if e := job.SetState(ctx, &job.Status.State); e != nil {
			js.log.Error("When updating jop status", e)
		}
		if preq >= job.Spec.PowerNeeded {
			return // done with this job allocation
		}
	}
}

// reconcile the pending job to see if it can be allocated to edge
func (js *JobScheduler) reconcileRequest(ctx context.Context) {
	// look through all the job requests that are not 100% scheduled
	// each job should be scheduled to "x" runners
	// find the next candidate and schedule the runner tasks
	// update the job with the scheduling info.
	// get list of edges
	edges, e := js.inv.GetAllEdges(ctx)
	if e != nil {
		js.log.Error(e)
		return
	}
	var freeEdges []int
	for idx, edge := range edges {
		if js.isEdgeFree(ctx, edge) {
			freeEdges = append(freeEdges, idx)
		}
	}
	// find out how many are idling
	// if all busy then nonthing to do return
	if len(freeEdges) == 0 {
		return // nothing to do
	}

	// find out the next job
	jlist, e := js.cfg.GetAllJobs(ctx)
	if e != nil {
		js.log.Error(e)
		return
	}
	jobFound, jobIdx := js.nextInLineJob(jlist)
	if !jobFound {
		return // nothing to do
	}
	// schedule the job on to the free edges
	js.allocateJob(ctx, jlist[jobIdx], edges, freeEdges)

	// do another round of reconcile so as to capture
	// partially filled requests.
	js.reconcileRequest(ctx)
}

// update comming form execution to be fed back to the origin job status
func (js *JobScheduler) executorUpdate(ctx context.Context, obj *nexus_client.JobmgmtJobInfo) {
	// get the requestor job
	// put the requestor job id
	// update the progress back
	// and recompute total progress
	jobName := obj.Spec.RequestorJob
	edge, e := obj.GetParent(ctx)
	if e != nil {
		js.log.Error("when getting edge", e)
	}
	edgeName := edge.DisplayName()
	job, e := js.cfg.GetJobs(ctx, jobName)
	if e != nil {
		js.log.Error("Unable to get job", e)
	}
	einfo, ok := job.Status.State.Execution[edgeName]
	if !ok {
		js.log.Errorf("Invalid  Edge or EdgeNot found %s", edgeName)
	}
	einfo.Progress = obj.Status.State.Progress
	einfo.StartTime = obj.Status.State.StartTime
	einfo.EndTime = obj.Status.State.EndTime
	totalPowerExecuted := uint32(0)
	for _, pex := range job.Status.State.Execution {
		totalPowerExecuted += pex.PowerRequested * pex.Progress / 100
	}
	job.Status.State.PercentCompleted = totalPowerExecuted * 100 / job.Spec.PowerNeeded
	if e := job.SetState(ctx, &job.Status.State); e != nil {
		js.log.Error("When updating job status", e)
	}
}

func (js *JobScheduler) Start(nexusClient *nexus_client.Clientset, g *errgroup.Group, gctx context.Context) {
	nexusClient.RootPowerScheduler().Config().Jobs("*").RegisterAddCallback(
		func(obj *nexus_client.JobschedulerJob) {
			js.reconcileRequest(gctx)
		})
	// nexusClient.RootPowerScheduler().DesiredEdgeConfig().EdgesDC("*").JobsInfo("*").RegisterAddCallback(
	// 	func(obj *nexus_client.JobmgmtJobInfo) {
	// 	})
	nexusClient.RootPowerScheduler().DesiredEdgeConfig().EdgesDC("*").JobsInfo("*").RegisterUpdateCallback(
		func(oldObj *nexus_client.JobmgmtJobInfo, newObj *nexus_client.JobmgmtJobInfo) {
			js.executorUpdate(gctx, newObj)
		})
}

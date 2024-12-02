package jobscheduler

import (
	"context"
	"fmt"
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
	nexusClient *nexus_client.Clientset
	jobSplits   uint32
	log         *logrus.Entry
}

func New(ncli *nexus_client.Clientset) *JobScheduler {
	log := logrus.WithFields(logrus.Fields{
		"module": "jobscheduler",
	})
	return &JobScheduler{ncli, 4, log}
}

// check if an edge is free
func (js *JobScheduler) isEdgeFree(ctx context.Context, edge *nexus_client.EdgeEdge) bool {
	dcfg := js.nexusClient.RootPowerScheduler().DesiredEdgeConfig()
	dcEdge, e := dcfg.GetEdgesDC(ctx, edge.DisplayName())
	if e != nil {
		if nexus_client.IsNotFound(e) {
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

	dcfg := js.nexusClient.RootPowerScheduler().DesiredEdgeConfig()

	for _, fedge := range freeEdge {
		pneeded := psplit
		if preq+psplitLast == job.Spec.PowerNeeded {
			pneeded = psplitLast
		}
		preq += pneeded // allocated power

		edge := edgeList[fedge]
		edgeName := edge.DisplayName()
		jobName := fmt.Sprintf("%s-%s-%d", job.DisplayName(), edgeName, preq)

		js.log.Infof("Allocationg part of job name = %s on edge %s JOBName %s PowerAllocated %d, TotalAllocations %d / %d ",
			job.DisplayName(), edgeName, jobName, pneeded, preq, job.Spec.PowerNeeded)
		curTime := time.Now().Unix()
		// create the desired state config
		edc, e := dcfg.GetEdgesDC(ctx, edgeName)
		if nexus_client.IsNotFound(e) {
			newEdge := &edcv1.EdgeDC{}
			newEdge.SetName(edgeName)
			// js.log.Infof("Creating the EdgeDC %s", edgeName)
			edc, e = dcfg.AddEdgesDC(ctx, newEdge)
			if e != nil {
				js.log.Error("Creating Desired edge config", e)
			}
		} else if e != nil {
			js.log.Error(e)
		}

		jobInfo := &jiv1.JobInfo{}
		jobInfo.SetName(jobName)
		jobInfo.Spec.RequestorJob = job.DisplayName()
		jobInfo.Spec.PowerNeeded = pneeded
		_, e = edc.AddJobsInfo(ctx, jobInfo)
		if e != nil {
			js.log.Errorf("When adding job info %s in edgedc %v", jobName, e)
		}
		// update the requester
		if job.Status.State.Execution == nil {
			job.Status.State.Execution = make(map[string]jsv1.NodeExecutionStatus)
		}
		if len(job.Status.State.Execution) == 0 {
			job.Status.State.StartTime = curTime
		}
		job.Status.State.Execution[jobName] = jsv1.NodeExecutionStatus{
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
	inv, e := js.nexusClient.RootPowerScheduler().GetInventory(ctx)
	if e != nil {
		js.log.Error(e)
		return
	}
	cfg, e := js.nexusClient.RootPowerScheduler().GetConfig(ctx)
	if e != nil {
		js.log.Error(e)
		return
	}
	edges, e := inv.GetAllEdges(ctx)
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
	jlist, e := cfg.GetAllJobs(ctx)
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
func (js *JobScheduler) executorUpdate(ctx context.Context, oldObj *nexus_client.JobmgmtJobInfo, obj *nexus_client.JobmgmtJobInfo) {
	// get the requestor job
	// put the requestor job id
	// update the progress back
	// and recompute total progress
	reqJobName := obj.Spec.RequestorJob
	jobName := obj.DisplayName()

	js.log.Infof("Got executor Update for job %s Requestro %s %+v", jobName, reqJobName, obj.Status.State)
	cfg := js.nexusClient.RootPowerScheduler().Config()
	reqJob, e := cfg.GetJobs(ctx, reqJobName)
	if e != nil {
		js.log.Error("Unable to get job", e)
	}
	einfo, ok := reqJob.Status.State.Execution[jobName]
	if !ok {
		js.log.Errorf("Invalid  Job or JobNot found %s", jobName)
	}
	einfo.Progress = obj.Status.State.Progress
	einfo.StartTime = obj.Status.State.StartTime
	einfo.EndTime = obj.Status.State.EndTime
	reqJob.Status.State.Execution[jobName] = einfo
	totalPowerExecuted := uint32(0)
	for _, pex := range reqJob.Status.State.Execution {
		totalPowerExecuted += pex.PowerRequested * pex.Progress / 100
	}
	reqJob.Status.State.PercentCompleted = totalPowerExecuted * 100 / reqJob.Spec.PowerNeeded
	if e := reqJob.SetState(ctx, &reqJob.Status.State); e != nil {
		js.log.Error("When updating job status", e)
	} else {
		// js.log.Infof("Updating Job %s to State %+v", reqJobName, reqJob.Status.State)
	}
	if oldObj.Status.State.Progress != 100 && obj.Status.State.Progress == 100 {
		js.reconcileRequest(ctx)
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
			js.executorUpdate(gctx, oldObj, newObj)
		})
}

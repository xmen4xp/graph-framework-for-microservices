package jobscheduler

import (
	"context"
	"fmt"
	"log"
	edcv1 "powerschedulermodel/build/apis/edgedc.intel.com/v1"
	jiv1 "powerschedulermodel/build/apis/jobmgmt.intel.com/v1"
	jsv1 "powerschedulermodel/build/apis/jobscheduler.intel.com/v1"
	nexus_client "powerschedulermodel/build/nexus-client"
	"sync"
	"time"

	"k8s.io/utils/keymutex"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

// This module is used for scheduling pending jobs on available edges
// the jobs are collected form the config/jobscheduler tree
// the available edges are collectected form inventory/edge tree
// the output of this modules is allocation of jobse in desiredconfig/edgedc
// as all the edges work on the workload they update the status in the edgedc node
// this module will merge all the updated status back into the config/jobscheduler node.
// jobs are executed in arrival order
//
// This module show a purely event based implementation using the data model.
// the event based architecture should allow this module to handle large number of edges / jobs.
// and to horizontally scale with the addtion of load sharing capability.

type JobScheduler struct {
	nexusClient    *nexus_client.Clientset
	jobSplits      uint32
	statusMutex    sync.Mutex
	reconcileMutex sync.Mutex
	log            *logrus.Entry
	edgeLock       keymutex.KeyMutex
	allocatedJobs  sync.Map //  map[string]bool
	allocatedEdges sync.Map
}
type Decision struct {
	jobIdx      int
	edgeName    string
	jobName     string
	powerneeded uint32
}

func New(ncli *nexus_client.Clientset) *JobScheduler {
	log := logrus.WithFields(logrus.Fields{
		"module": "jobscheduler",
	})
	js := &JobScheduler{}
	js.nexusClient = ncli
	js.jobSplits = 4
	js.log = log
	js.edgeLock = keymutex.NewHashed(0)
	js.statusMutex = sync.Mutex{}
	// js.allocatedJobs = make(map[string]bool)
	return js
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
	if e != nil {
		js.log.Error(e)
	}
	for _, job := range jobs {
		if job.Status.State.Progress != 100 {
			return false
		}
	}
	if _, ok := js.allocatedEdges.Load(edge.DisplayName()); ok {
		return false
	}
	return true
}

// find the next job in line for execution sorted based on arrival time.
func (js *JobScheduler) nextInLineJob(jlist []*nexus_client.JobschedulerJob) (bool, int, int) {
	first := true
	foundIdx := 0
	foundTimeStamp := int64(0)
	availableJobs := 0
	for idx, job := range jlist {
		preq := uint32(0)
		for _, pex := range job.Status.State.Execution {
			preq += pex.PowerRequested
		}
		if job.Spec.PowerNeeded > preq {
			availableJobs += 1
			if job.Spec.CreationTime < foundTimeStamp || first {
				foundTimeStamp = job.Spec.CreationTime
				foundIdx = idx
				first = false
			}
		}
	}
	return !first, foundIdx, availableJobs
}

func (js *JobScheduler) createJobInfo(ctx context.Context, job *nexus_client.JobschedulerJob, edgeName string, dec *Decision) {

	dcfg := js.nexusClient.RootPowerScheduler().DesiredEdgeConfig()

	edc, e := dcfg.GetEdgesDC(ctx, edgeName)
	if nexus_client.IsNotFound(e) {
		newEdge := &edcv1.EdgeDC{}
		newEdge.SetName(edgeName)
		js.log.Infof("Creating the EdgeDC %s", edgeName)
		edc, e = dcfg.AddEdgesDC(ctx, newEdge)
		if e != nil {
			js.log.Error("Creating Desired edge config", e)
		}
	} else if e != nil {
		js.log.Error(e)
	}

	jobInfo := &jiv1.JobInfo{}
	jobInfo.SetName(dec.jobName)
	jobInfo.Spec.RequestorJob = job.DisplayName()
	jobInfo.Spec.PowerNeeded = dec.powerneeded
	curTime := time.Now().Unix()
	_, e = edc.AddJobsInfo(ctx, jobInfo)
	if e != nil {
		js.log.Errorf("When adding job info %s in edgedc %v", dec.jobName, e)
	}
	// update the requester
	if job.Status.State.Execution == nil {
		job.Status.State.Execution = make(map[string]jsv1.NodeExecutionStatus)
	}
	if len(job.Status.State.Execution) == 0 {
		job.Status.State.StartTime = curTime
	}
	job.Status.State.Execution[dec.jobName] = jsv1.NodeExecutionStatus{
		PowerRequested: dec.powerneeded,
		StartTime:      0,
		EndTime:        0,
		Progress:       0,
	}
	job.Status.State.EndTime = curTime
	//	js.allocatedJobs[job.DisplayName()] = true
	// commit the data back
	js.log.Infof("Creating DC JobInfo on edge %s with job %s => %s", edgeName, jobInfo.DisplayName(), jobInfo.Spec.RequestorJob)
	if e := edc.Update(ctx); e != nil {
		js.log.Fatal("When updating job spec", e)
	}
	if e := job.SetState(ctx, &job.Status.State); e != nil {
		js.log.Fatal("When updating jop status", e)
	}
}

// allocate the job on the free edges.
func (js *JobScheduler) allocateJob(ctx context.Context, jobList []*nexus_client.JobschedulerJob,
	edgeList []*nexus_client.EdgeEdge, freeEdge []int) {
	decisions := make(map[string]Decision) // job decisions for each edge
	jobAllocationMap := make(map[int]bool)
	for _, fedge := range freeEdge {
		// find the next job
		first := true
		jobIdx := 0
		foundTimeStamp := int64(0)
		for idx, job := range jobList {
			preq := uint32(0)
			for _, pex := range job.Status.State.Execution {
				preq += pex.PowerRequested
			}
			// skip jobs that are allocated
			if _, ok := jobAllocationMap[idx]; ok {
				continue
			}
			if _, ok := js.allocatedJobs.Load(job.DisplayName()); ok {
				continue
			}
			if job.Spec.PowerNeeded > preq {
				if job.Spec.CreationTime < foundTimeStamp || first {
					foundTimeStamp = job.Spec.CreationTime
					jobIdx = idx
					first = false
				}
			}
		}
		if first {
			// no jobs found
			break
		}

		edge := edgeList[fedge]
		edgeName := edge.DisplayName()
		job := jobList[jobIdx]
		jobAllocationMap[jobIdx] = true
		jobName := fmt.Sprintf("%s-%s", job.DisplayName(), edgeName)
		pneeded := job.Spec.PowerNeeded
		js.log.Infof("Allocating job name = %s on edge %s JOBName %s PowerAllocated %d, TotalAllocations %d ",
			job.DisplayName(), edgeName, jobName, pneeded, job.Spec.PowerNeeded)

		decisions[edgeName] = Decision{
			jobIdx:      jobIdx,
			jobName:     jobName,
			powerneeded: pneeded,
		}
	}

	var wg sync.WaitGroup
	sem := semaphore.NewWeighted(int64(5)) // run 10 at a time.
	for edgeName, dec := range decisions {
		wg.Add(1)
		js.allocatedJobs.Store(dec.jobName, true)
		js.allocatedEdges.Store(edgeName, jobList[dec.jobIdx].DisplayName())

		go func(curDec Decision, eName string) {
			defer wg.Done()
			if e := sem.Acquire(ctx, 1); e != nil {
				js.log.Errorf("Failed to acquire semaphore: %v", e)
				return
			}
			defer sem.Release(1)
			job := *jobList[curDec.jobIdx]
			js.createJobInfo(ctx, &job, eName, &curDec)
		}(dec, edgeName)
	}
	wg.Wait()
}

// reconcile the pending job to see if it can be allocated to edge
func (js *JobScheduler) reconcileRequest(ctx context.Context) {
	// look through all the job requests that are not 100% scheduled
	// each job should be scheduled to "x" runners
	// find the next candidate and schedule the runner tasks
	// update the job with the scheduling info.
	// get list of edges
	js.reconcileMutex.Lock()
	defer js.reconcileMutex.Unlock()

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
		js.log.Infof("Reconciler free Edges:0")
		return // nothing to do
	}

	// find out the next job
	jlist, e := cfg.GetAllJobs(ctx)
	if e != nil {
		js.log.Error(e)
		return
	}
	jobFound, jobIdx, availableJobs := js.nextInLineJob(jlist)
	js.log.Infof("Reconciler job found:%t.%d. TotalsjobsPending:%d FreeEdges:%d %v", jobFound, jobIdx, availableJobs, len(freeEdges), freeEdges)
	if !jobFound {
		return // nothing to do
	}
	// schedule the job on to the free edges
	js.allocateJob(ctx, jlist, edges, freeEdges)

	// do another round of reconcile so as to capture
	// partially filled requests.
	//  go func() { js.reconcileRequest(ctx) }()
}

// update comming form execution to be fed back to the origin job status
func (js *JobScheduler) executorUpdate(ctx context.Context, oldObj *nexus_client.JobmgmtJobInfo, obj *nexus_client.JobmgmtJobInfo) {
	// get the requestor job
	// put the requestor job id
	// update the progress back
	// and recompute total progress
	reqJobName := obj.Spec.RequestorJob
	jobName := obj.DisplayName()
	// get edge name
	edgedc, e := obj.GetParent(ctx)
	if e != nil {
		log.Fatal(e)
	}
	edgeName := edgedc.DisplayName()
	js.edgeLock.LockKey(edgeName)
	defer js.edgeLock.UnlockKey(edgeName)
	js.log.Infof("Got executor Update for job %s Requestor %s Power %d Status %+v",
		jobName, reqJobName, obj.Spec.PowerNeeded, obj.Status.State)
	if obj.Status.State.Progress == 100 {
		if jn, ok := js.allocatedEdges.Load(edgeName); ok {
			if jn == reqJobName {
				js.allocatedEdges.Delete(edgeName)
			}
		}
	}
	cfg := js.nexusClient.RootPowerScheduler().Config()
	reqJob, e := cfg.GetJobs(ctx, reqJobName)
	if e != nil {
		js.log.Error("Unable to get job", e)
	}
	einfo, ok := reqJob.Status.State.Execution[jobName]
	if !ok {
		js.log.Errorf("Invalid  Job or JobNot found %s", jobName)
	}
	einfo.PowerRequested = obj.Spec.PowerNeeded
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

	// Update Stats when job completes
	if oldObj.Status.State.Progress != 100 && obj.Status.State.Progress == 100 {
		go func() { js.reconcileRequest(ctx) }()
		js.statusMutex.Lock()
		defer js.statusMutex.Unlock()
		jc := js.getJSStatus(ctx)
		jc.Status.State.TotalJobsExecuted.TotalJobsProcessed++
		jc.Status.State.TotalJobsExecuted.TotalPowerProcessed += obj.Spec.PowerNeeded
		if jc.Status.State.TotalJobsExecutedPerEdge == nil {
			jc.Status.State.TotalJobsExecutedPerEdge = make(map[string]jsv1.ExectuedJobStats)
		}
		ej, ok := jc.Status.State.TotalJobsExecutedPerEdge[edgeName]
		if !ok {
			ej = jsv1.ExectuedJobStats{TotalJobsProcessed: 0, TotalPowerProcessed: 0}
			jc.Status.State.TotalJobsExecutedPerEdge[edgeName] = ej
		}
		ej.TotalJobsProcessed++
		ej.TotalPowerProcessed += obj.Spec.PowerNeeded
		jc.Status.State.TotalJobsExecutedPerEdge[edgeName] = ej
		// delete(jc.Status.State.CurrentJobsExecuting, edgeName)
		js.setJSStatus(ctx, jc)
		// delete the obj
		e = edgedc.DeleteJobsInfo(ctx, obj.DisplayName())
		if e != nil {
			js.log.Error("Error when cleaning up desired config jobInfo ", e)
		}
		js.allocatedJobs.Delete(reqJobName)
	}
	// update stats when job starts or updates
	if obj.Status.State.Progress < 100 {
		js.statusMutex.Lock()
		defer js.statusMutex.Unlock()
		jc := js.getJSStatus(ctx)
		if jc.Status.State.CurrentJobsExecuting == nil {
			jc.Status.State.CurrentJobsExecuting = make(map[string]jsv1.ExecutingJobStatus)
		}
		ej, ok := jc.Status.State.CurrentJobsExecuting[edgeName]
		if !ok {
			ej = jsv1.ExecutingJobStatus{}
			jc.Status.State.CurrentJobsExecuting[edgeName] = ej
		}
		ej.PowerRequested = obj.Spec.PowerNeeded
		ej.StartTime = obj.Status.State.StartTime
		ej.EndTime = obj.Status.State.EndTime
		ej.Progress = obj.Status.State.Progress
		jc.Status.State.CurrentJobsExecuting[edgeName] = ej
		js.setJSStatus(ctx, jc)
	}
}

func (js *JobScheduler) getJSStatus(ctx context.Context) *nexus_client.JobschedulerSchedulerConfig {
	jscfg, e := js.nexusClient.RootPowerScheduler().Config().GetScheduler(ctx)
	if nexus_client.IsNotFound(e) {
		js.log.Info("Creating the Scheduler ...", e)
		itm := jsv1.SchedulerConfig{}
		itm.Status.State.CurrentJobsExecuting = make(map[string]jsv1.ExecutingJobStatus)
		itm.Status.State.TotalJobsExecutedPerEdge = make(map[string]jsv1.ExectuedJobStats)
		jscfg, e = js.nexusClient.RootPowerScheduler().Config().AddScheduler(ctx, &itm)
	}
	if e != nil {
		js.log.Fatal("when getting Scheduler status  ", e)
	}
	if jscfg.Status.State.CurrentJobsExecuting == nil {
		jscfg.Status.State.CurrentJobsExecuting = make(map[string]jsv1.ExecutingJobStatus)
	}
	if jscfg.Status.State.TotalJobsExecutedPerEdge == nil {
		jscfg.Status.State.TotalJobsExecutedPerEdge = make(map[string]jsv1.ExectuedJobStats)
	}
	return jscfg
}
func (js *JobScheduler) setJSStatus(ctx context.Context, jscfg *nexus_client.JobschedulerSchedulerConfig) {
	if e := jscfg.SetState(ctx, &jscfg.Status.State); e != nil {
		js.log.Fatal("when setting Scheduler status ", e)
	}
}

func (js *JobScheduler) Start(nexusClient *nexus_client.Clientset, g *errgroup.Group, gctx context.Context) {
	// subscribe so the read events are processed in the local cache.
	// and local cache is kept up to date from change events.
	nexusClient.RootPowerScheduler().Inventory().Edges("*").Subscribe()
	nexusClient.RootPowerScheduler().Config().Jobs("*").Subscribe()
	nexusClient.RootPowerScheduler().DesiredEdgeConfig().EdgesDC("*").Subscribe()
	nexusClient.RootPowerScheduler().DesiredEdgeConfig().EdgesDC("*").JobsInfo("*").Subscribe()

	// call back registration for nodes
	nexusClient.RootPowerScheduler().Config().Jobs("*").RegisterAddCallback(
		func(obj *nexus_client.JobschedulerJob) {
			js.reconcileRequest(gctx)
		})
	nexusClient.RootPowerScheduler().DesiredEdgeConfig().EdgesDC("*").JobsInfo("*").RegisterUpdateCallback(
		func(oldObj *nexus_client.JobmgmtJobInfo, newObj *nexus_client.JobmgmtJobInfo) {
			js.executorUpdate(gctx, oldObj, newObj)
		})

	// timer task to check on stuck jobs and reschedule them.
}

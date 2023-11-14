package jobscheduler

import (
	"context"
	"fmt"
	"log"
	edcv1 "powerschedulermodel/build/apis/edgedc.intel.com/v1"
	jiv1 "powerschedulermodel/build/apis/jobmgmt.intel.com/v1"
	jsv1 "powerschedulermodel/build/apis/jobscheduler.intel.com/v1"
	nexus_client "powerschedulermodel/build/nexus-client"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"k8s.io/utils/keymutex"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

// This module is used for scheduling pending jobs on available edges
// the jobs are collected form the config/job scheduler tree
// the available edges are collected form inventory/edge tree
// the output of this modules is allocation of jobs in desiredconfig/edgedc
// as all the edges work on the workload they update the status in the edgedc node
// this module will merge all the updated status back into the config/jobscheduler node.
// jobs are executed in arrival order
//
// This module show a purely event based implementation using the data model.
// the event based architecture should allow this module to handle large number of edges / jobs.
// and to horizontally scale with the addtion of load sharing capability.

type JobScheduler struct {
	nexusClient       *nexus_client.Clientset
	jobSplits         uint32
	statusMutex       sync.Mutex
	reconcileMutex    sync.Mutex
	log               *logrus.Entry
	edgeLock          keymutex.KeyMutex
	allocatedEdges    sync.Map
	jobsDone60Sec     uint32
	jobsDone60SecTime int64
	reconcilerWaitCnt int32
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
	js.jobsDone60Sec = 0
	js.jobsDone60SecTime = time.Now().Unix()
	js.reconcilerWaitCnt = 0
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

func (js *JobScheduler) createJobInfo(
	ctx context.Context,
	job *nexus_client.JobschedulerJob,
	edge *nexus_client.EdgeEdge) {

	dcfg := js.nexusClient.RootPowerScheduler().DesiredEdgeConfig()
	edgeName := edge.DisplayName()
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
	jobName := fmt.Sprintf("%s-%s", job.DisplayName(), edgeName)

	jobInfo := &jiv1.JobInfo{}
	jobInfo.SetName(jobName)
	jobInfo.Spec.RequestorJob = job.DisplayName()
	jobInfo.Spec.PowerNeeded = job.Spec.PowerNeeded
	curTime := time.Now().Unix()
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
		PowerRequested: job.Spec.PowerNeeded,
		StartTime:      0,
		EndTime:        0,
		Progress:       0,
	}
	job.Status.State.EndTime = curTime
	// commit the data back
	js.log.Infof("Creating DC JobInfo on edge %s with job %s => %s", edgeName, jobInfo.DisplayName(), jobInfo.Spec.RequestorJob)
	if e := edc.Update(ctx); e != nil {
		js.log.Fatal("When updating job spec", e)
	}
	if e := job.SetState(ctx, &job.Status.State); e != nil {
		js.log.Fatal("When updating jop status", e)
	}
}

// reconcile the pending job to see if it can be allocated to edge
func (js *JobScheduler) reconcileRequest(ctx context.Context) {
	if js.reconcilerWaitCnt > 1 {
		// no need to pile on reconcilers back to back ...
		// if one is waiting to start execution drop subsequent ones.
		return
	}
	atomic.AddInt32(&js.reconcilerWaitCnt, 1)
	js.reconcileMutex.Lock()
	defer func() {
		js.reconcileMutex.Unlock()
		atomic.AddInt32(&js.reconcilerWaitCnt, -1)
	}()

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
		// js.log.Infof("Reconciler free Edges:0 cnt=%d", js.reconcilerWaitCnt)
		return // nothing to do
	}

	// find out the next job
	alljoblist, e := cfg.GetAllJobs(ctx)
	if e != nil {
		js.log.Error(e)
		return
	}
	// identify list of unallocated jobs
	var jobList []*nexus_client.JobschedulerJob
	for _, j := range alljoblist {
		if j.Status.State.PercentCompleted >= 100 {
			continue
		}
		preq := uint32(0)
		for _, pex := range j.Status.State.Execution {
			preq += pex.PowerRequested
		}
		if j.Spec.PowerNeeded > preq {
			jobList = append(jobList, j)
		}
	}
	sort.Slice(jobList, func(i, j int) bool {
		return jobList[i].Spec.CreationTime < jobList[j].Spec.CreationTime
	})
	if len(jobList) == 0 {
		return // nothing to do
	}
	// schedule the jobs on to the free edges
	jobMap, e := js.simpleScheduler(ctx, jobList, edges, freeEdges)
	if e != nil {
		js.log.Error(e)
	}
	// parallel write using "x" number of workers.
	var wg sync.WaitGroup
	sem := semaphore.NewWeighted(int64(5)) // run 10 at a time.
	for jobIdx, edgeIdx := range jobMap {
		edge := edges[edgeIdx]
		job := jobList[jobIdx]

		wg.Add(1)
		go func(jobItm *nexus_client.JobschedulerJob, edgeItem *nexus_client.EdgeEdge) {
			defer wg.Done()
			if e := sem.Acquire(ctx, 1); e != nil {
				js.log.Errorf("Failed to acquire semaphore: %v", e)
				return
			}
			defer sem.Release(1)
			js.createJobInfo(ctx, jobItm, edgeItem)
		}(job, edge)
	}
	wg.Wait()
}

// update coming form execution to be fed back to the origin job status
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
		e = edgedc.DeleteJobsInfo(ctx, obj.DisplayName())
		if e != nil {
			js.log.Error("Error when cleaning up desired config jobInfo ", e)
		}
	}
}

func (js *JobScheduler) Start(nexusClient *nexus_client.Clientset, g *errgroup.Group, gctx context.Context) {
	// subscribe so the read events are processed in the local cache.
	// and local cache is kept uptodate with database.
	nexusClient.RootPowerScheduler().Inventory().Edges("*").Subscribe()
	nexusClient.RootPowerScheduler().Config().Jobs("*").Subscribe()
	nexusClient.RootPowerScheduler().DesiredEdgeConfig().EdgesDC("*").Subscribe()
	nexusClient.RootPowerScheduler().DesiredEdgeConfig().EdgesDC("*").JobsInfo("*").Subscribe()

	// call back registration for nodes change events
	nexusClient.RootPowerScheduler().Config().Jobs("*").RegisterAddCallback(
		func(obj *nexus_client.JobschedulerJob) {
			js.reconcileRequest(gctx)
		})
	nexusClient.RootPowerScheduler().DesiredEdgeConfig().EdgesDC("*").JobsInfo("*").RegisterUpdateCallback(
		func(oldObj *nexus_client.JobmgmtJobInfo, newObj *nexus_client.JobmgmtJobInfo) {
			js.executorUpdate(gctx, oldObj, newObj)
		})

	// timer task to update stats every 10 sec
	g.Go(func() error {
		tickerJob := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-tickerJob.C:
				e := js.updateStats(gctx)
				if e != nil {
					js.log.Error("check and create jobs ", e)
				}
			case <-gctx.Done():
				return gctx.Err()
			}
		}
	})
}

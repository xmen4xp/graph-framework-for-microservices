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
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"k8s.io/utils/keymutex"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
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
	nexusClient             *nexus_client.Clientset
	jobSplits               uint32
	statusMutex             sync.Mutex
	reconcileMutex          sync.Mutex
	jobsDCReconcileMutex    sync.Mutex
	log                     *logrus.Entry
	edgeLock                keymutex.KeyMutex
	allocatedEdges          sync.Map
	jobsDone60Sec           uint32
	jobsDoneTotal           uint32
	jobsDone60SecTime       int64
	reconcilerWaitCnt       int32
	jobsDCReconcilerWaitCnt int32
	lockWaitOnEdge          int32
	edgeFreeMap             sync.Map // string -> *nexus_client.EdgeEdge
	jobFreeMap              sync.Map // string -> *nexus_client.JobschedulerJob // config-job
	numSchedulerInstances   int
	schedulerInstanceId     int
}

func New(ncli *nexus_client.Clientset, numSchedulerInstances, schedulerInstanceId int) *JobScheduler {
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
	js.jobsDoneTotal = 0
	js.jobsDone60SecTime = time.Now().Unix()
	js.reconcilerWaitCnt = 0
	js.lockWaitOnEdge = 0
	js.jobsDCReconcilerWaitCnt = 0
	js.numSchedulerInstances = numSchedulerInstances
	js.schedulerInstanceId = schedulerInstanceId
	log.Infof("Initialiazing scheduler: InstanceId %d of %d", js.schedulerInstanceId, js.numSchedulerInstances)
	return js
}

// check if an edge is free
func (js *JobScheduler) isEdgeFree(ctx context.Context, siteName string, edge *nexus_client.EdgeEdge) (bool, int) {
	numRunningJobs := 0
	edgeFree := true
	site, e := js.nexusClient.RootPowerScheduler().DesiredSiteConfig().GetSitesDC(ctx, siteName)
	if e != nil {
		js.log.Fatalf("Error getting site %v", e)
	}

	dcEdge, e := site.GetEdgesDC(ctx, edge.DisplayName())
	if e != nil {
		if nexus_client.IsNotFound(e) {
			return false, numRunningJobs
		}
		js.log.Error(e)
		return false, numRunningJobs
	}
	jobs := dcEdge.GetAllJobsInfoIter(ctx)

	// if len(jobs.JobsInfo) > 0 {
	//	for job, _ := jobs.Next(ctx); job != nil; job, _ = jobs.Next(ctx) {
	//		js.log.Debugf("--> isEdgeFree %s : %s: %d", dcEdge.DisplayName(), job.DisplayName(), job.Status.State.Progress)
	//	}
	//	return false
	//}

	for job, _ := jobs.Next(ctx); job != nil; job, _ = jobs.Next(ctx) {
		//js.log.Infof("--> isEdgeFree %s : %s: %d", dcEdge.DisplayName(), job.DisplayName(), job.Status.State.Progress)
		if job.Status.State.Progress < 100 {
			edgeFree = false
			numRunningJobs++
		}
	}

	// if _, ok := js.allocatedEdges.Load(edge.DisplayName()); ok {
	// 	return false
	// }
	return edgeFree, numRunningJobs
}

func (js *JobScheduler) isSiteOwned(siteName string) bool {

	if (js.numSchedulerInstances <= 1) || (js.schedulerInstanceId >= js.numSchedulerInstances) {
		//js.log.Infof("edge %s treated as owned", edgeName)
		return true
	}

	s := strings.Split(siteName, "-")
	if len(s) > 0 {
		if siteId, err := strconv.Atoi(s[len(s)-1]); err == nil {
			mappedSiteId := siteId % js.numSchedulerInstances
			if mappedSiteId == js.schedulerInstanceId {
				// js.log.Infof("site %s determined as owned", jobName)
				return true
			} else {
				return false
			}
		}
	}
	log.Fatalf("site name %s is not in expected format", siteName)
	return false
}

func (js *JobScheduler) assignJobGroupToSite(ctx context.Context, invSite *nexus_client.SiteSite) bool {

	site, e := js.nexusClient.RootPowerScheduler().DesiredSiteConfig().GetSitesDC(ctx, invSite.DisplayName())
	if e != nil {
		js.log.Errorf("Update of find site %s in desired config. Error %v Will retry... ", invSite.DisplayName(), e)
		return false
	}

	// find out the next job
	cfg, e := js.nexusClient.RootPowerScheduler().GetConfig(ctx)
	if e != nil {
		js.log.Fatal(e)
	}

	allJobGroupList := cfg.GetAllJobGroupsIter(ctx)

	var jg *nexus_client.JobgroupJobgroup
	ownedJobGroup := ""
	for jg, _ = allJobGroupList.Next(context.Background()); jg != nil; jg, _ = allJobGroupList.Next(context.Background()) {
		if js.isJobGroupOwned(jg.DisplayName()) == true {
			ownedJobGroup = jg.DisplayName()
			break
		}
	}

	if ownedJobGroup == "" {
		return false
	}

	currJobGroup, ok := site.Labels["jobgroups.jobgroup.intel.com"]
	if !ok || currJobGroup != ownedJobGroup {
		site.Labels["jobgroups.jobgroup.intel.com"] = ownedJobGroup
		updateErr := site.Update(ctx)
		if updateErr != nil {
			js.log.Fatalf("Update of jobgroup label %s on site %s failed with error %v", ownedJobGroup, site.DisplayName(), updateErr)
		}
	}
	return true
}

/*
func (js *JobScheduler) isJobOwned(jobName string) bool {
	if (js.numSchedulerInstances <= 1) || (js.schedulerInstanceId >= js.numSchedulerInstances) {
		js.log.Infof("job %s treated as owned", jobName)
		return true
	}

	h := sha256.New()
	h.Write([]byte(jobName))
	hash := h.Sum(nil)

	if (int(hash[0]) % js.numSchedulerInstances) == js.schedulerInstanceId {
		// js.log.Infof("job %s determined as owned", jobName)
		return true
	}

	// js.log.Infof("job %s determined as NOT owned", jobName)
	return false
} */

func (js *JobScheduler) isJobGroupOwned(jobGroupName string) bool {
	if js.numSchedulerInstances <= 1 {
		js.log.Infof("JOB %s treated as owned", jobGroupName)
		return true
	}

	s := strings.Split(jobGroupName, "-")
	if len(s) > 0 {
		if jobGroupId, err := strconv.Atoi(s[len(s)-1]); err == nil {
			if jobGroupId == js.schedulerInstanceId {
				// js.log.Infof("job group %s determined as owned", jobName)
				return true
			} else {
				return false
			}
		}
	}

	log.Fatalf("job group name %s is not in expected format", jobGroupName)
	return false

	/*
		h := sha256.New()
		h.Write([]byte(jobGroupName))
		hash := h.Sum(nil)

		if (int(hash[0]) % js.numSchedulerInstances) == js.schedulerInstanceId {
			// js.log.Infof("job group %s determined as owned", jobName)
			return true
		}
	*/
	// js.log.Infof("job group %s determined as NOT owned", jobName)
	// return false
}

func (js *JobScheduler) createJobInfo(
	ctx context.Context,
	job *nexus_client.JobschedulerJob,
	siteName string,
	edge *nexus_client.EdgeEdge) {

	siteDC, e := js.nexusClient.RootPowerScheduler().DesiredSiteConfig().GetSitesDC(ctx, siteName)
	if e != nil {
		js.log.Fatalf("Error getting sitedc %v", e)
	}
	edgeName := edge.DisplayName()
	edc, e := siteDC.GetEdgesDC(ctx, edgeName)
	if nexus_client.IsNotFound(e) {
		newEdge := &edcv1.EdgeDC{}
		newEdge.SetName(edgeName)
		js.log.Infof("Creating the EdgeDC %s", edgeName)
		edc, e = siteDC.AddEdgesDC(ctx, newEdge)
		if e != nil {
			js.log.Error("Creating Desired edge config", e)
		}
	} else if e != nil {
		js.log.Error("createJobInfo:edgeDC", e)
	}
	jobName := fmt.Sprintf("%s-%s", job.DisplayName(), edgeName)
	jobGroup, e := job.GetParent(ctx)
	if e != nil {
		js.log.Fatalf("Error when getting parent jobgroup %v", e)
	}

	jobInfo := &jiv1.JobInfo{}
	jobInfo.SetName(jobName)
	jobInfo.Spec.RequestorJob = job.DisplayName()
	jobInfo.Spec.RequestorJobGroup = jobGroup.DisplayName()
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
	// js.log.Infof("Creating DC JobInfo on edge %s with job %s => %s", edgeName, jobInfo.DisplayName(), jobInfo.Spec.RequestorJob)
	// if e := edc.Update(ctx); e != nil {
	//	js.log.Fatal("When updating job spec", e)
	// }
	if e := job.SetState(ctx, &job.Status.State); e != nil {
		js.log.Fatal("When updating jop status", e)
	}
	js.log.Infof("Creating DC JobInfo on edge %s with job %s DONE2", edgeName, jobInfo.DisplayName())
}

// reconcile the pending job to see if it can be allocated to edge
func (js *JobScheduler) reconcileRequest(ctx context.Context) {
	if js.reconcilerWaitCnt > 1 {
		// no need to pile on reconcilers back to back ...
		// if one is waiting to start execution drop subsequent ones.
		js.log.Infof("Reconciler ignoring event. %d", js.reconcilerWaitCnt)
		return
	}
	atomic.AddInt32(&js.reconcilerWaitCnt, 1)
	js.reconcileMutex.Lock()
	startTime := time.Now()
	defer func() {
		js.reconcileMutex.Unlock()
		atomic.AddInt32(&js.reconcilerWaitCnt, -1)
	}()

	inv, e := js.nexusClient.RootPowerScheduler().GetInventory(ctx)
	if e != nil {
		js.log.Error(e)
		return
	}

	st := ""
	var freeEdges []int
	var edges []*nexus_client.EdgeEdge
	var edgesSiteName []string
	idx := 0
	var totalRunningJobs int
	{
		var e error
		var site *nexus_client.SiteSite
		allSiteList := inv.GetAllSiteIter(ctx)
		for site, e = allSiteList.Next(ctx); site != nil && e == nil; site, e = allSiteList.Next(ctx) {

			if !js.isSiteOwned(site.DisplayName()) {
				continue
			} else {
				if _, ok := site.Labels["jobgroups.jobgroup.intel.com"]; !ok {
					if js.assignJobGroupToSite(ctx, site) == false {
						js.log.Errorf("Unable to assign site %s to any job group", site.DisplayName())
						continue
					}
				}
			}

			alledgelist := site.GetAllEdgesIter(ctx)
			var edge *nexus_client.EdgeEdge
			for edge, e = alledgelist.Next(ctx); edge != nil && e == nil; edge, e = alledgelist.Next(ctx) {

				//if js.isEdgeOwned(edge.DisplayName()) == false {
				// Edge is not owned by this scheduler instance. Ignore.
				//	continue
				//}

				edgeFree, runningJobs := js.isEdgeFree(ctx, site.DisplayName(), edge)
				totalRunningJobs += runningJobs
				if edgeFree == true {
					freeEdges = append(freeEdges, idx)
					st = st + ", " + edge.DisplayName()
				}
				edges = append(edges, edge)
				edgesSiteName = append(edgesSiteName, site.DisplayName())
				idx++
			}
			if e != nil {
				js.log.Fatalf("Error iterating on Edges %v", e)
			}
		}
		if e != nil {
			js.log.Fatalf("Error iterating on Sites %v", e)
		}
	}
	// js.log.Infof("DEBUG Total Edges: %d free %d freeEdges [%s]", len(edges), len(freeEdges), st)
	js.log.Infof("DEBUG Total Edges: %d Num freeEdges %d totalRunningJobs %d", len(edges), len(freeEdges), totalRunningJobs)

	// find out how many are idling
	// if all busy then nonthing to do return
	if len(freeEdges) == 0 {
		js.log.Infof("DEBUG Reconciler free Edges:0 cnt=%d", js.reconcilerWaitCnt)
		return // nothing to do
	}

	// find out the next job
	cfg, e := js.nexusClient.RootPowerScheduler().GetConfig(ctx)
	if e != nil {
		js.log.Error(e)
		return
	}
	allJobGroupList := cfg.GetAllJobGroupsIter(ctx)

	// identify list of unallocated jobs
	var jobList []*nexus_client.JobschedulerJob
	var j *nexus_client.JobschedulerJob
	var jg *nexus_client.JobgroupJobgroup
	var getJobsError, jgError error
	var numJobs int
	var jobids []string
	// var jobidsNotOwned []string

	for jg, jgError = allJobGroupList.Next(context.Background()); jg != nil && jgError == nil; jg, jgError = allJobGroupList.Next(context.Background()) {

		if js.isJobGroupOwned(jg.DisplayName()) == false {
			// Job group is not owned by this instance. Skip processing this job group.
			continue
		}

		alljoblist := jg.GetAllJobsIter(ctx)
		for j, getJobsError = alljoblist.Next(context.Background()); j != nil && getJobsError == nil; j, getJobsError = alljoblist.Next(context.Background()) {
			numJobs += 1
			if j.Status.State.PercentCompleted >= 100 {
				continue
			}
			preq := uint32(0)
			for _, pex := range j.Status.State.Execution {
				preq += pex.PowerRequested
			}

			///if js.isJobOwned(j.DisplayName()) == false {
			// Job is marked as not owned. Ignore.
			// notOwned += 1
			//jobidsNotOwned = append(jobidsNotOwned, j.DisplayName())
			//continue
			//}

			// Power requested will be zero if the job is not scheduled yet.
			if j.Spec.PowerNeeded > preq {
				jobList = append(jobList, j)
				jobids = append(jobids, j.DisplayName())
			}
		}

		if getJobsError != nil {
			js.log.Error(getJobsError)
			return
		}
	}
	// js.log.Infof("DEBUG Total Num Jobs %d NotOwned %d Completed %d NumRunning %d Free jobs %d Jobs %v", numJobs, notOwned, completed, numRunning, len(jobList), jobids)
	js.log.Infof("DEBUG Total Num Jobs In Config %d Unallocated jobs %d", numJobs, len(jobList))

	if getJobsError != nil {
		js.log.Error(getJobsError)
		return
	}

	sort.Slice(jobList, func(i, j int) bool {
		return jobList[i].Spec.CreationTime < jobList[j].Spec.CreationTime
	})
	jst := ""
	for _, j := range jobList {
		jst += "," + j.DisplayName()
	}
	js.log.Infof("JOB List [%d] %s", len(jobList), jst)
	if len(jobList) == 0 {
		return // nothing to do
	}
	// schedule the jobs on to the free edges
	jobMap, e := js.simpleScheduler(ctx, jobList, edges, freeEdges)
	if e != nil {
		js.log.Error(e)
	}

	midTime := time.Now()
	// parallel write using "x" number of workers.
	// var wg sync.WaitGroup
	// sem := semaphore.NewWeighted(int64(10)) // run 10 at a time.
	for jobIdx, edgeIdx := range jobMap {
		edge := edges[edgeIdx]
		siteName := edgesSiteName[edgeIdx]
		job := jobList[jobIdx]
		// wg.Add(1)
		func(jobItm *nexus_client.JobschedulerJob, sn string, edgeItem *nexus_client.EdgeEdge) {
			// defer wg.Done()
			// if e := sem.Acquire(ctx, 1); e != nil {
			// 	js.log.Errorf("Failed to acquire semaphore: %v", e)
			// 	return
			//}
			//defer sem.Release(1)
			js.createJobInfo(ctx, jobItm, sn, edgeItem)
		}(job, siteName, edge)
	}
	//wg.Wait()
	endTime := time.Now()
	js.log.Infof("DEBUG Schedular Reconcile time mid : %d end : %d", midTime.UnixMilli()-startTime.UnixMilli(), endTime.UnixMilli()-startTime.UnixMilli())
}

func (js *JobScheduler) executorUpdateReconcile(ctx context.Context, siteName, edgeName, jobGroupName, jobName, reqJobName string) {
	js.reconcilerWaitCnt++
	js.log.Infof("Got executor Update for job %s (%s) on edge %s waiting count = %d",
		jobName, reqJobName, edgeName, js.reconcilerWaitCnt)
	// get the current status of the job
	js.edgeLock.LockKey(edgeName)
	defer js.edgeLock.UnlockKey(edgeName)
	js.reconcilerWaitCnt--

	obj, e := js.nexusClient.RootPowerScheduler().DesiredSiteConfig().SitesDC(siteName).EdgesDC(edgeName).GetJobsInfo(ctx, jobName)
	if e != nil {
		if nexus_client.IsNotFound(e) {
			return // nop
		}
		js.log.Fatalf("executor update, getjobinfo ", e)
	}

	// if reqJob.Status.State.PercentCompleted >= 100 {
	if obj.Status.State.Progress >= 100 {
		// delete and trigger reconciler
		// go func() { js.reconcileRequest(ctx) }()

		reqJobGroup, e := js.nexusClient.RootPowerScheduler().Config().GetJobGroups(ctx, jobGroupName)
		if e != nil {
			js.log.Fatalf("Get job group %v failed with error ", jobGroupName, e)
		}
		reqJob, e := reqJobGroup.GetJobs(ctx, reqJobName)
		if e != nil || reqJob == nil {
			js.log.Fatalf("Get job  %v failed with error ", reqJobName, e)
		}

		if reqJob.Status.State.Execution == nil {
			reqJob.Status.State.Execution = make(map[string]jsv1.NodeExecutionStatus)
		}
		einfo, ok := reqJob.Status.State.Execution[jobName]
		if !ok {
			js.log.Errorf("Error: Invalid  Job or JobNot found %s", jobName)
		}
		einfo.PowerRequested = obj.Spec.PowerNeeded
		einfo.Progress = obj.Status.State.Progress
		einfo.StartTime = obj.Status.State.StartTime
		einfo.EndTime = obj.Status.State.EndTime
		if reqJob.Status.State.Execution == nil {
			reqJob.Status.State.Execution = make(map[string]jsv1.NodeExecutionStatus)
		}
		reqJob.Status.State.Execution[jobName] = einfo
		// totalPowerExecuted := uint32(0)
		// for _, pex := range reqJob.Status.State.Execution {
		// 	totalPowerExecuted += pex.PowerRequested * pex.Progress / 100
		// }
		reqJob.Status.State.PercentCompleted = 100 // totalPowerExecuted * 100 / reqJob.Spec.PowerNeeded
		if e := reqJob.SetState(ctx, &reqJob.Status.State); e != nil {
			js.log.Fatalf("updating status for job %v failed with error %v", reqJob.GetName(), e)
		} else {
			js.log.Infof("Updating Job %s to State %+v", reqJobName, reqJob.Status.State)
		}

		// Delete the job from desired config.
		edc, e := js.nexusClient.RootPowerScheduler().DesiredSiteConfig().SitesDC(siteName).GetEdgesDC(ctx, edgeName)
		if e != nil {
			js.log.Error("Error when gettin gEdgeDC", e)
		}
		e = edc.DeleteJobsInfo(ctx, obj.DisplayName())
		if e != nil {
			js.log.Error("Error when cleaning up desired config jobInfo ", e)
		}
		js.log.Infof("Compleated executor Update for job %s (%s) on edge %s waiting count = %d",
			jobName, reqJobName, edgeName, js.reconcilerWaitCnt)
		// Trigger reconciler.
		js.reconcileRequest(ctx)
	}

}

// update coming form execution to be fed back to the origin job status
func (js *JobScheduler) executorUpdate(ctx context.Context, oldObj *nexus_client.JobmgmtJobInfo, obj *nexus_client.JobmgmtJobInfo) {
	// get the requestor job
	// put the requestor job id
	// update the progress back
	// and recompute total progress
	reqJobName := obj.Spec.RequestorJob
	jobName := obj.DisplayName()
	reqJobGroupName := obj.Spec.RequestorJobGroup
	if oldObj.Status.State.Progress != 100 && obj.Status.State.Progress == 100 {
		js.jobsDoneTotal++
	}
	if js.isJobGroupOwned(reqJobGroupName) == false {
		// JobGroup is not owned by this scheduler instance. Ignore event.
		return
	}

	// get edge name
	edgedc, e := obj.GetParent(ctx)
	if e != nil {
		log.Fatalf("Getting edgedc %v", e)
	}
	edgeName := edgedc.DisplayName()

	sitedc, e := edgedc.GetParent(ctx)
	if e != nil {
		log.Fatalf("Getting sitedc %v", e)
	}
	siteName := sitedc.DisplayName()

	//if js.isEdgeOwned(edgeName) == false {
	// Edge is not owned by this scheduler instance. Ignore event.
	// js.log.Panicf("Job %s is owned, but edge %s on which the job is running is not owned. This should never happen", jobName, edgeName)
	//js.log.Errorf("Job %s is owned, but edge %s on which is not owned. Ignoring event.", jobName, edgeName)
	//return
	//}

	js.executorUpdateReconcile(ctx, siteName, edgeName, reqJobGroupName, jobName, reqJobName)

	return
	/*
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
		if reqJob.Status.State.Execution == nil {
			reqJob.Status.State.Execution = make(map[string]jsv1.NodeExecutionStatus)
		}
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
		}*/
}

// periodic call to updated the status
func (js *JobScheduler) garbageCollectJobsDC(ctx context.Context) {

	if js.jobsDCReconcilerWaitCnt > 1 {
		// no need to pile on reconcilers back to back ...
		// if one is waiting to start execution drop subsequent ones.
		js.log.Infof("[garbageCollectJobsDC] JobsDC Reconciler ignoring event. %d", js.jobsDCReconcilerWaitCnt)
		return
	}
	atomic.AddInt32(&js.jobsDCReconcilerWaitCnt, 1)
	js.jobsDCReconcileMutex.Lock()

	defer func() {
		js.jobsDCReconcileMutex.Unlock()
		atomic.AddInt32(&js.jobsDCReconcilerWaitCnt, -1)
	}()

	js.log.Infof("[garbageCollectJobsDC] JobsDC garbage collector running...")

	dcfg, e := js.nexusClient.RootPowerScheduler().GetDesiredSiteConfig(ctx)
	if e != nil {
		js.log.Errorf("[garbageCollectJobsDC] JobsDC garbage collector running.. Unable to read desired edge config. error %v", e)
		return
	}
	sites := dcfg.GetAllSitesDCIter(ctx)
	for site, _ := sites.Next(ctx); site != nil; site, _ = sites.Next(ctx) {
		edges := site.GetAllEdgesDCIter(ctx)
		for edge, _ := edges.Next(ctx); edge != nil; edge, _ = edges.Next(ctx) {
			js.log.Debugf("[garbageCollectJobsDC] evaluating edge %v", edge.DisplayName())
			ji := edge.GetAllJobsInfoIter(ctx)
			for job, _ := ji.Next(ctx); job != nil; job, _ = ji.Next(ctx) {
				if job.Status.State.Progress < 100 {
					// Job is in progress. Do not garbage collect.
					continue
				}

				elapsedTime := time.Now().Unix() - job.Status.State.EndTime
				if elapsedTime < 5*60 {
					js.log.Debugf("[garbageCollectJobsDC] completed job %v is not ready to be garbage collected yet. Elapsed time %v ", job.DisplayName(), elapsedTime)
					continue
				} else {
					js.log.Debugf("[garbageCollectJobsDC] completed job %+v in group %s looks stale. Needs garbage collection.", job.GetName(), job.Spec.RequestorJobGroup)
				}

				if js.isJobGroupOwned(job.Spec.RequestorJobGroup) {
					// Job group is owned by this scheduler.
					reqJobGroup, e := js.nexusClient.RootPowerScheduler().Config().GetJobGroups(ctx, job.Spec.RequestorJobGroup)
					if e != nil {
						js.log.Errorf("[garbageCollectJobsDC] unable to localte jobgroup %s for job %s", job.Spec.RequestorJobGroup, job.GetName())

						// Unable to locate jobgroup. Nothing to do.
						// TBD: should we force read to see if jog group actually exists ?
						// Its possible the job is from the jobgroup that has been deleted ?
						continue
					}

					reqJob, e := reqJobGroup.GetJobs(ctx, job.Spec.RequestorJob)
					if e != nil || reqJob == nil {
						js.log.Infof("[garbageCollectJobsDC] Deleting Jobs %v from edge %v, as config", job.GetName(), edge.GetName())
						// Job not found in config job group. Must be stale.
						e = edge.DeleteJobsInfo(ctx, job.DisplayName())
						if e != nil {
							js.log.Errorf("Delete of stale job %v from edge %v failed with error %v", job.GetName(), edge.GetName(), e)
						}
					} else {
						if reqJob != nil {
							js.log.Errorf("[garbageCollectJobsDC] Job %v exists", reqJob.GetName())
							if status, ok := reqJob.Status.State.Execution[job.DisplayName()]; ok && status.Progress < 100 {
								js.log.Errorf("[garbageCollectJobsDC] Invoke update for job %s as progress %d is stale", job.GetName(), status.Progress)
								js.executorUpdateReconcile(ctx, site.DisplayName(), edge.DisplayName(), job.Spec.RequestorJobGroup, job.DisplayName(), job.Spec.RequestorJob)
							}
						}
					}
				}
			}
		}
	}
}

func (js *JobScheduler) Start(nexusClient *nexus_client.Clientset, g *errgroup.Group, gctx context.Context) {
	// subscribe so the read events are processed in the local cache.
	// and local cache is kept update with database.
	nexusClient.RootPowerScheduler().Config().Subscribe()
	nexusClient.RootPowerScheduler().Inventory().Subscribe()
	nexusClient.RootPowerScheduler().DesiredSiteConfig().Subscribe()

	nexusClient.RootPowerScheduler().Inventory().Site("*").Subscribe()
	nexusClient.RootPowerScheduler().Inventory().Site("*").Edges("*").Subscribe()
	nexusClient.RootPowerScheduler().Config().JobGroups("*").Subscribe()
	nexusClient.RootPowerScheduler().Config().JobGroups("*").Jobs("*").Subscribe()
	nexusClient.RootPowerScheduler().Config().Scheduler().Subscribe()
	nexusClient.RootPowerScheduler().DesiredSiteConfig().SitesDC("*").Subscribe()
	nexusClient.RootPowerScheduler().DesiredSiteConfig().SitesDC("*").EdgesDC("*").Subscribe()
	nexusClient.RootPowerScheduler().DesiredSiteConfig().SitesDC("*").EdgesDC("*").JobsInfo("*").Subscribe()

	nexusClient.RootPowerScheduler().Inventory().Site("*").Edges("*").RegisterAddCallback(
		func(obj *nexus_client.EdgeEdge) {
			js.reconcileRequest(gctx)
		})
	// call back registration for nodes change events
	nexusClient.RootPowerScheduler().Config().JobGroups("*").Jobs("*").RegisterAddCallback(
		func(obj *nexus_client.JobschedulerJob) {
			js.reconcileRequest(gctx)
		})
	nexusClient.RootPowerScheduler().DesiredSiteConfig().SitesDC("*").EdgesDC("*").JobsInfo("*").RegisterUpdateCallback(
		func(oldObj *nexus_client.JobmgmtJobInfo, newObj *nexus_client.JobmgmtJobInfo) {
			js.executorUpdate(gctx, oldObj, newObj)
		})

	// timer task to run reconciler loop every 10 sec
	g.Go(func() error {
		tickerJob := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-tickerJob.C:
				js.reconcileRequest(gctx)
			case <-gctx.Done():
				return gctx.Err()
			}
		}
	})

	// timer task to update stats every 10 sec
	g.Go(func() error {
		tickerJob := time.NewTicker(4 * time.Second)
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
	g.Go(func() error {
		tickerJob := time.NewTicker(300 * time.Second)
		for {
			select {
			case <-tickerJob.C:
				js.garbageCollectJobsDC(gctx)
			case <-gctx.Done():
				return gctx.Err()
			}
		}
	})
}

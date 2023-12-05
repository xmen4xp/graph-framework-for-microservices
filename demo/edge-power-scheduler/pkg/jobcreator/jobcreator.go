package jobcreator

// periodically
// based on a policy create jobs and queue them into the pipeline
//
// periodically poll progress of completion of jobs and report
//
// periodically remove jobs that are older than 1hr from the system

// work on the config -> Jobs

import (
	"context"
	"fmt"
	"math/rand"
	jgv1 "powerschedulermodel/build/apis/jobgroup.intel.com/v1"
	jsv1 "powerschedulermodel/build/apis/jobscheduler.intel.com/v1"
	nexus_client "powerschedulermodel/build/nexus-client"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/sirupsen/logrus"

	"time"

	"golang.org/x/sync/errgroup"
)

// this module is used for generation of demo job requests
// each job will have a random power requirements
// The module will load up to MaxPendingJobs in the queue and
// will keep adding to the queue as jobs completes
// This module will also do garbage collection and will delete jobs that are
// more than and hour old.
// A variable LimitJobCnt can be used for debug and limit how many jobs are generated
// value of 0 indicates infinity.
//
// This module shows a timer based implementation of data model.
type JobCreator struct {
	Name               string
	MaxPendingJobs     uint32
	MaxActiveJobs      uint32
	MaxJobPower        uint32
	MinJobPower        uint32
	lastJobId          uint64
	LimitJobCnt        uint32
	PendingJobs        sync.Map
	PendingJobCnt      uint32
	reconcileMutex     sync.Mutex
	jobCreatationMutex sync.Mutex
	nexusClient        *nexus_client.Clientset
	log                *logrus.Entry
	reconcilerWaitCnt  int32
}

func New(ncli *nexus_client.Clientset, name string, maxPendingJobs uint32, maxActiveJobs uint32,
	maxJobPowerRange uint32, minJobPowerRange uint32, limitJobs uint32, startId uint64) *JobCreator {
	log := logrus.WithFields(logrus.Fields{
		"module": "jobcreator",
	})
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "15:04:05"
	logrus.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
	jc := &JobCreator{}
	jc.Name = name
	jc.MaxPendingJobs = maxPendingJobs
	jc.MaxActiveJobs = maxActiveJobs
	jc.MaxJobPower = maxJobPowerRange
	jc.MinJobPower = minJobPowerRange
	jc.lastJobId = startId
	jc.LimitJobCnt = limitJobs
	jc.nexusClient = ncli
	jc.log = log
	jc.reconcilerWaitCnt = 0
	return jc
}

func (jc *JobCreator) createAndAddJob(ctx context.Context, jGroup *nexus_client.JobgroupJobgroup, pendingJobList int, pendingJobNotStarted int) error {
	jc.jobCreatationMutex.Lock()
	j := &jsv1.Job{}
	j.Name = fmt.Sprintf("job-%d", jc.lastJobId)
	j.Spec.JobId = jc.lastJobId
	j.Spec.PowerNeeded = jc.MinJobPower + uint32(rand.Intn(int(jc.MaxJobPower-jc.MinJobPower)))
	j.Spec.CreationTime = time.Now().Unix()
	jc.lastJobId++
	// jc.log.Infof("Creating Job %s for requested power %d [Pending Jobs %d/%d]", j.Name, j.Spec.PowerNeeded, pendingJobList, pendingJobNotStarted)
	jc.log.Infof("Creating Job %s for requested power %d", j.Name, j.Spec.PowerNeeded)

	jc.PendingJobs.Store(j.Name, true)
	jc.PendingJobCnt++
	jc.jobCreatationMutex.Unlock()
	_, e := jGroup.AddJobs(ctx, j)

	return e
}

func (jc *JobCreator) cleanupCompletedJobs(ctx context.Context) error {
	cfg, e := jc.nexusClient.RootPowerScheduler().GetConfig(ctx)
	if e != nil {
		return e
	}

	jGroup, e := cfg.GetJobGroups(ctx, jc.Name)
	if e != nil {
		return e
	}

	allJobs := jGroup.GetAllJobsIter(ctx)
	var job *nexus_client.JobschedulerJob
	for job, _ = allJobs.Next(context.Background()); job != nil; job, _ = allJobs.Next(context.Background()) {
		if job.Status.State.PercentCompleted >= 100 {
			elapsedTime := time.Now().Unix() - job.Status.State.EndTime
			// delete jobs that compleated an hour back.
			if elapsedTime > 1 {
				if e = jGroup.DeleteJobs(ctx, job.DisplayName()); e != nil {
					jc.log.Errorf("CREATOR: Delete failed for completed job %s with error %v", job.DisplayName(), e)
				}
			}
			// completedJobCount += 1
		}
	}
	return nil
}

func (jc *JobCreator) checkAndCreateJobs(ctx context.Context) error {

	cfg, e := jc.nexusClient.RootPowerScheduler().GetConfig(ctx)
	if e != nil {
		return e
	}

	jGroup, e := cfg.GetJobGroups(ctx, jc.Name)
	if nexus_client.IsChildNotFound(e) {
		// add job group
		var j jgv1.Jobgroup
		j.SetName(jc.Name)
		jGroup, e = cfg.AddJobGroups(ctx, &j)
		if e != nil {
			jc.log.Fatalf("Error creating jobgroup %v", e)
		}
	} else if e != nil {
		jc.log.Fatalf("Error getting jobgroup %v", e)
	}
	allJobs := jGroup.GetAllJobsIter(ctx)

	// inPrgoressJobCnt := 0
	pendingJobNotStarted := 0
	// completedJobCount := 0
	pendingJobList := []string{}

	var job *nexus_client.JobschedulerJob
	var getJobsError error
	totalJobs := 0
	jobsActive := 0
	for job, getJobsError = allJobs.Next(context.Background()); job != nil && getJobsError == nil; job, getJobsError = allJobs.Next(context.Background()) {
		totalJobs++
		if job.Status.State.PercentCompleted >= 100 {
			elapsedTime := time.Now().Unix() - job.Status.State.EndTime
			// delete jobs that compleated an hour back.
			if elapsedTime > 1 {
				if e = jGroup.DeleteJobs(ctx, job.DisplayName()); e != nil {
					jc.log.Errorf("CREATOR: Aborting job iteration as delete failed for completed job %s with error %v", job.DisplayName(), e)
					return e
				}
			}
			// completedJobCount += 1
		} else {
			if job.Status.State.Execution == nil || len(job.Status.State.Execution) == 0 {
				pendingJobNotStarted += 1
			} else {
				// inPrgoressJobCnt += 1
				jobsActive++
			}
			pendingJobList = append(pendingJobList, job.DisplayName())
		}
	}
	if getJobsError != nil {
		jc.log.Errorf("getJobsError: %v", getJobsError)
		return getJobsError
	}

	// if int(jc.PendingJobCnt) > pendingJobCnt {
	// 	pendingJobNotStarted += int(jc.PendingJobCnt) - pendingJobCnt

	// }
	// jc.log.Infof("CREATOR: Pending JOB COUNTER: CompletedJobCnt %d InProgressJobCnt %d PendingJobNotStarted %d", completedJobCount, inPrgoressJobCnt, pendingJobNotStarted)
	jc.log.Infof("CREATOR: Pending JOB COUNTER: PendingJobNotStarted %d", pendingJobNotStarted)
	if jobsActive < int(jc.MaxActiveJobs) && pendingJobNotStarted < int(jc.MaxPendingJobs) && (jc.LimitJobCnt == 0 || uint64(jc.LimitJobCnt) > jc.lastJobId) {
		jobsCnt := int(jc.MaxPendingJobs) - pendingJobNotStarted
		jc.log.Infof("CREATOR Creating Job count: %d", jobsCnt)

		// var wg sync.WaitGroup
		for i := 0; i < jobsCnt; i++ {
			//wg.Add(1)
			func(idx int) {
				//	defer wg.Done()
				if e = jc.createAndAddJob(ctx, jGroup, len(pendingJobList)+idx, pendingJobNotStarted+idx); e != nil {
					jc.log.Error(e)
				}
			}(i)
		}
		// wg.Wait()
		jc.log.Infof("CREATOR Creating Job count: %d DONE", jobsCnt)

		// go func() { jc.checkAndCreateJobs(ctx) }()
	}
	return nil
}

func (jc *JobCreator) checkDCAndCreateJobs(ctx context.Context) error {

	if jc.reconcilerWaitCnt > 1 {
		// no need to pile on reconcilers back to back ...
		// if one is waiting to start execution drop subsequent ones.
		jc.log.Debugf("Reconciler ignoring event. %d", jc.reconcilerWaitCnt)
		return nil
	}

	atomic.AddInt32(&jc.reconcilerWaitCnt, 1)
	jc.reconcileMutex.Lock()
	defer func() {
		jc.reconcileMutex.Unlock()
		atomic.AddInt32(&jc.reconcilerWaitCnt, -1)
	}()

	dcfg, e := jc.nexusClient.RootPowerScheduler().GetDesiredSiteConfig(ctx)
	if e != nil {
		return e
	}

	ownedSites := []*nexus_client.SitedcSiteDC{}
	sideDCIter := dcfg.GetAllSitesDCIter(ctx)
	for site, _ := sideDCIter.Next(ctx); site != nil; site, _ = sideDCIter.Next(ctx) {
		if jg, ok := site.Labels["jobgroups.jobgroup.intel.com"]; !ok {
			jc.log.Infof("Site %s does not have a assigned job group. Ignoring the site", site.DisplayName())
			continue
		} else {
			if jg == jc.Name {
				ownedSites = append(ownedSites, site)
			}
		}
	}

	if len(ownedSites) == 0 {
		jc.log.Infof("No sites are assigned to job group %s. Create jobs by config nodes.", jc.Name)
		return jc.checkAndCreateJobs(ctx)
	}

	maxJobId := -1
	for _, site := range ownedSites {
		edgeIter := site.GetAllEdgesDCIter(ctx)
		for edge, _ := edgeIter.Next(ctx); edge != nil; edge, _ = edgeIter.Next(ctx) {
			ji := edge.GetAllJobsInfoIter(ctx)
			for job, _ := ji.Next(ctx); job != nil; job, _ = ji.Next(ctx) {
				jobName := job.Spec.RequestorJob
				if jobName == "" {
					jc.log.Fatal("Job name cannot be empty")
				}

				jobNameList := strings.Split(jobName, "-")
				if len(jobNameList) != 2 {
					jc.log.Fatalf("Job name %s not in expected format of job-xxx", jobName)
				}

				if jobId, err := strconv.Atoi(jobNameList[len(jobNameList)-1]); err == nil {
					if maxJobId < jobId {
						maxJobId = jobId
					}
				} else {
					jc.log.Fatalf("Error getting jobid from name %s. Error %v", jobName, err)
				}
			}
		}
	}

	if maxJobId == -1 {
		return jc.checkAndCreateJobs(ctx)
	}

	if maxJobId > int(jc.lastJobId) {
		jc.log.Fatalf("maxJobId in DC %d is greater than the last created jobId %d", maxJobId, jc.lastJobId)
	}
	pendingJobNotStarted := jc.lastJobId - uint64(maxJobId)
	jc.log.Infof("CREATOR: Pending JOB COUNTER: PendingJobNotStarted %d", pendingJobNotStarted)
	if pendingJobNotStarted < uint64(jc.MaxPendingJobs) {

		cfg, e := jc.nexusClient.RootPowerScheduler().GetConfig(ctx)
		if e != nil {
			return e
		}
		jGroup, e := cfg.GetJobGroups(ctx, jc.Name)
		if nexus_client.IsChildNotFound(e) {
			// add job group
			var j jgv1.Jobgroup
			j.SetName(jc.Name)
			jGroup, e = cfg.AddJobGroups(ctx, &j)
			if e != nil {
				jc.log.Fatalf("Error creating jobgroup %v", e)
			}
		} else if e != nil {
			jc.log.Fatalf("Error getting jobgroup %v", e)
		}

		jobsCnt := uint64(jc.MaxPendingJobs) - pendingJobNotStarted
		jc.log.Infof("CREATOR Creating Job count: %d", jobsCnt)
		for i := uint64(0); i < jobsCnt; i++ {
			if e = jc.createAndAddJob(ctx, jGroup, 0, 0); e != nil {
				jc.log.Error(e)
			}
		}
		jc.log.Infof("CREATOR Creating Job count: %d DONE", jobsCnt)
	}

	// cleanup completed jobs from config
	jc.cleanupCompletedJobs(ctx)
	return nil
}

func (jc *JobCreator) jobsCB(cbType string, obj *nexus_client.JobschedulerJob) {
	var execInfo string = ""
	e := obj.Status.State.Execution
	for n, v := range e {
		execInfo += fmt.Sprintf("[%s: %d, %d%%]", n, v.PowerRequested, v.Progress)
	}
	if obj.Status.State.PercentCompleted >= 100 {
		jc.log.Infof("Config/Job CBFN: Job %s [%s] PowerReq:%d Completion:%d Distribution:%s", obj.DisplayName(),
			cbType, obj.Spec.PowerNeeded, obj.Status.State.PercentCompleted, execInfo)
		if _, ok := jc.PendingJobs.Load(obj.DisplayName()); ok {
			jc.PendingJobs.Delete(obj.DisplayName())
			jc.PendingJobCnt--
		}
	}

	go func() {
		ctx := context.Background()
		e := jc.checkAndCreateJobs(ctx)
		if e != nil {
			jc.log.Error("check and create jobs on call back", e)
		}
	}()
}

func (jc *JobCreator) Start(nexusClient *nexus_client.Clientset, g *errgroup.Group, gctx context.Context) {
	nexusClient.RootPowerScheduler().Config().Subscribe()
	nexusClient.RootPowerScheduler().Config().JobGroups(jc.Name).Jobs("*").Subscribe()
	nexusClient.RootPowerScheduler().Config().JobGroups("*").Subscribe()
	nexusClient.RootPowerScheduler().DesiredSiteConfig().Subscribe()
	nexusClient.RootPowerScheduler().DesiredSiteConfig().SitesDC("*").Subscribe()
	nexusClient.RootPowerScheduler().DesiredSiteConfig().SitesDC("*").EdgesDC("*").Subscribe()
	nexusClient.RootPowerScheduler().DesiredSiteConfig().SitesDC("*").EdgesDC("*").JobsInfo("*").Subscribe()
	// nexusClient.RootPowerScheduler().Config().JobGroups(jc.Name).Jobs("*").RegisterAddCallback(
	//	func(obj *nexus_client.JobschedulerJob) {
	//		jc.jobsCB("Add", obj)
	//	})
	//nexusClient.RootPowerScheduler().Config().JobGroups(jc.Name).Jobs("*").RegisterUpdateCallback(
	//	func(oldObj *nexus_client.JobschedulerJob, newObj *nexus_client.JobschedulerJob) {
	//		jc.jobsCB("Upd", newObj)
	//	})
	g.Go(func() error {
		tickerJob := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-tickerJob.C:
				e := jc.checkDCAndCreateJobs(gctx)
				if e != nil {
					jc.log.Error("check and create jobs ", e)
				}
			case <-gctx.Done():
				return gctx.Err()
			}
		}
	})
}

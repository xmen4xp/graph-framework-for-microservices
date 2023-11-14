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
	jsv1 "powerschedulermodel/build/apis/jobscheduler.intel.com/v1"
	nexus_client "powerschedulermodel/build/nexus-client"
	"sync"

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
	MaxPendingJobs     uint32
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
}

func New(ncli *nexus_client.Clientset, maxPendingJobs uint32,
	maxJobPowerRange uint32, minJobPowerRange uint32, limitJobs uint32, startId uint64) *JobCreator {
	log := logrus.WithFields(logrus.Fields{
		"module": "jobcreator",
	})
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "15:04:05"
	logrus.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
	jc := &JobCreator{}
	jc.MaxPendingJobs = maxPendingJobs
	jc.MaxJobPower = maxJobPowerRange
	jc.MinJobPower = minJobPowerRange
	jc.lastJobId = startId
	jc.LimitJobCnt = limitJobs
	jc.nexusClient = ncli
	jc.log = log
	return jc
}

func (jc *JobCreator) createAndAddJob(ctx context.Context, cfg *nexus_client.ConfigConfig, pendingJobList int, pendingJobNotStarted int) error {
	jc.jobCreatationMutex.Lock()
	j := &jsv1.Job{}
	j.Name = fmt.Sprintf("job-%d", jc.lastJobId)
	j.Spec.JobId = jc.lastJobId
	j.Spec.PowerNeeded = jc.MinJobPower + uint32(rand.Intn(int(jc.MaxJobPower-jc.MinJobPower)))
	j.Spec.CreationTime = time.Now().Unix()
	jc.lastJobId++
	jc.log.Infof("Creating Job %s for requested power %d [Pending Jobs %d/%d]", j.Name, j.Spec.PowerNeeded, pendingJobList, pendingJobNotStarted)
	jc.PendingJobs.Store(j.Name, true)
	jc.PendingJobCnt++
	jc.jobCreatationMutex.Unlock()
	_, e := cfg.AddJobs(ctx, j)

	return e
}
func (jc *JobCreator) checkAndCreateJobs(ctx context.Context) error {
	jc.reconcileMutex.Lock()
	defer jc.reconcileMutex.Unlock()
	cfg, e := jc.nexusClient.RootPowerScheduler().GetConfig(ctx)
	if e != nil {
		return e
	}
	allJobs, e := cfg.GetAllJobs(ctx)
	if e != nil {
		return e
	}
	pendingJobCnt := 0
	pendingJobNotStarted := 0
	pendingJobList := []string{}
	for _, job := range allJobs {
		if job.Status.State.PercentCompleted >= 100 {
			elapsedTime := time.Now().Unix() - job.Status.State.EndTime
			// delete jobs that compleated an hour back.
			if elapsedTime > 60*60 {
				if e = cfg.DeleteJobs(ctx, job.DisplayName()); e != nil {
					return e
				}
			}
		} else {
			pendingJobCnt += 1
			if job.Status.State.Execution == nil || len(job.Status.State.Execution) == 0 {
				pendingJobNotStarted += 1
			}
			pendingJobList = append(pendingJobList, job.DisplayName())
		}
	}
	// if int(jc.PendingJobCnt) > pendingJobCnt {
	// 	pendingJobNotStarted += int(jc.PendingJobCnt) - pendingJobCnt

	// }
	//	jc.log.Infof("Pending JOB COUNTER:%d %d", pendingJobCnt, pendingJobNotStarted)
	if pendingJobNotStarted < int(jc.MaxPendingJobs) && (jc.LimitJobCnt == 0 || uint64(jc.LimitJobCnt) > jc.lastJobId) {
		jobsCnt := int(jc.MaxPendingJobs) - pendingJobNotStarted
		var wg sync.WaitGroup
		for i := 0; i < jobsCnt; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				if e = jc.createAndAddJob(ctx, cfg, len(pendingJobList)+idx, pendingJobNotStarted+idx); e != nil {
					jc.log.Error(e)
				}
			}(i)
		}
		wg.Wait()
		// go func() { jc.checkAndCreateJobs(ctx) }()
	}
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
	nexusClient.RootPowerScheduler().Config().Jobs("*").Subscribe()
	nexusClient.RootPowerScheduler().Config().Jobs("*").RegisterAddCallback(
		func(obj *nexus_client.JobschedulerJob) {
			jc.jobsCB("Add", obj)
		})
	nexusClient.RootPowerScheduler().Config().Jobs("*").RegisterUpdateCallback(
		func(oldObj *nexus_client.JobschedulerJob, newObj *nexus_client.JobschedulerJob) {
			jc.jobsCB("Upd", newObj)
		})
	g.Go(func() error {
		tickerJob := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-tickerJob.C:
				e := jc.checkAndCreateJobs(gctx)
				if e != nil {
					jc.log.Error("check and create jobs ", e)
				}
			case <-gctx.Done():
				return gctx.Err()
			}
		}
	})
}

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

	"github.com/sirupsen/logrus"

	"time"

	"golang.org/x/sync/errgroup"
)

// this module is used for generaton of demo job requests
// each job will have a random power requirements
// The module will load up to MaxPendingJobs in the queue and
// will keep adding to the queue as jobs completes
// The moudle will also do garbage collection and will delete jobs that are
// more than and hour old.
// A variable LimitJobCnt can be used for debug and limit how many jobs are generated
// value of 0 indicates infinity.
//
// This module shows a timer based implementation of data model.
type JobCreator struct {
	MaxPendingJobs uint32
	MaxJobPower    uint32
	MinJobPower    uint32
	lastJobId      uint64
	LimitJobCnt    uint32
	nexusClient    *nexus_client.Clientset
	log            *logrus.Entry
}

func New(ncli *nexus_client.Clientset, maxPendingJobs uint32,
	maxJobPowerRange uint32, minJobPowerRange uint32, limitJobs uint32) *JobCreator {
	log := logrus.WithFields(logrus.Fields{
		"module": "jobcreator",
	})
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "15:04:05"
	logrus.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true

	jc := &JobCreator{maxPendingJobs, maxJobPowerRange, minJobPowerRange, 1, limitJobs, ncli, log}
	return jc
}

func (jc *JobCreator) createAndAddJob(ctx context.Context, cfg *nexus_client.ConfigConfig, pendingJobs []string) error {
	j := &jsv1.Job{}
	j.Name = fmt.Sprintf("job-%d", jc.lastJobId)
	j.Spec.JobId = jc.lastJobId
	j.Spec.PowerNeeded = jc.MinJobPower + uint32(rand.Intn(int(jc.MaxJobPower-jc.MinJobPower)))
	j.Spec.CreationTime = time.Now().Unix()
	jc.lastJobId++
	jc.log.Infof("Creating Job %s for requested power %d [Pending Jobs %v]", j.Name, j.Spec.PowerNeeded, pendingJobs)
	_, e := cfg.AddJobs(ctx, j)
	return e
}
func (jc *JobCreator) checkAndCreateJobs(ctx context.Context) error {
	cfg, e := jc.nexusClient.RootPowerScheduler().GetConfig(ctx)
	if e != nil {
		return e
	}
	allJobs, e := cfg.GetAllJobs(ctx)
	if e != nil {
		return e
	}
	pendingJobCnt := 0
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
			pendingJobList = append(pendingJobList, job.DisplayName())
		}
	}
	if pendingJobCnt < int(jc.MaxPendingJobs) && (jc.LimitJobCnt == 0 || uint64(jc.LimitJobCnt) > jc.lastJobId) {
		if e = jc.createAndAddJob(ctx, cfg, pendingJobList); e != nil {
			return e
		}
	}
	return nil
}
func (jc *JobCreator) jobsCB(cbType string, obj *nexus_client.JobschedulerJob) {
	var execInfo string = ""
	e := obj.Status.State.Execution
	for n, v := range e {
		execInfo += fmt.Sprintf("[%s: %d, %d%%]", n, v.PowerRequested, v.Progress)
	}
	jc.log.Infof("Config/Job CBFN: Job %s [%s] PowerReq:%d Completion:%d Distribution:%s", obj.DisplayName(),
		cbType, obj.Spec.PowerNeeded, obj.Status.State.PercentCompleted, execInfo)
}

func (jc *JobCreator) Start(nexusClient *nexus_client.Clientset, g *errgroup.Group, gctx context.Context) {
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

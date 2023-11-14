package jobscheduler

import (
	jsv1 "powerschedulermodel/build/apis/jobscheduler.intel.com/v1"
	nexus_client "powerschedulermodel/build/nexus-client"

	"context"
	"time"
)

func (js *JobScheduler) getJSStatus(ctx context.Context) *nexus_client.JobschedulerSchedulerConfig {
	jscfg, e := js.nexusClient.RootPowerScheduler().Config().GetScheduler(ctx)
	if nexus_client.IsNotFound(e) {
		js.log.Info("Creating the Scheduler ...", e)
		itm := jsv1.SchedulerConfig{}
		jscfg, e = js.nexusClient.RootPowerScheduler().Config().AddScheduler(ctx, &itm)
	}
	if e != nil {
		js.log.Fatal("when getting Scheduler status  ", e)
	}
	return jscfg
}
func (js *JobScheduler) setJSStatus(ctx context.Context, jscfg *nexus_client.JobschedulerSchedulerConfig) {
	if e := jscfg.SetState(ctx, &jscfg.Status.State); e != nil {
		js.log.Fatal("when setting Scheduler status ", e)
	}
}

// periodic call to updated the status
func (js *JobScheduler) updateStats(ctx context.Context) error {
	js.statusMutex.Lock()
	defer js.statusMutex.Unlock()
	sch := js.getJSStatus(ctx)
	dcfg, e := js.nexusClient.RootPowerScheduler().GetDesiredEdgeConfig(ctx)
	if e != nil {
		return e
	}
	edges, e := dcfg.GetAllEdgesDC(ctx)
	if e != nil {
		return e
	}
	cfg, e := js.nexusClient.RootPowerScheduler().GetConfig(ctx)
	if e != nil {
		return e
	}
	jobs, e := cfg.GetAllJobs(ctx)
	if e != nil {
		return e
	}
	sch.Status.State.BusyEdges = 0
	sch.Status.State.FreeEdges = 0
	sch.Status.State.JobsRunning = 0
	sch.Status.State.JobsQueued = 0
	sch.Status.State.TotalJobsExecuted.TotalPowerProcessed = 0
	jobsDone := uint32(0)
	jobsMap := make(map[string]bool)
	for _, edge := range edges {
		ji, e := edge.GetAllJobsInfo(ctx)
		if e != nil {
			return e
		}
		found := false
		for _, job := range ji {
			if job.Status.State.Progress >= 100 {

			} else {
				found = true
				sch.Status.State.BusyEdges++
				jobsMap[job.Spec.RequestorJob] = true
				break
			}
		}
		if !found {
			sch.Status.State.FreeEdges++
		}
	}
	for _, job := range jobs {
		if job.Status.State.PercentCompleted >= 100 {
			jobsDone++
			sch.Status.State.TotalJobsExecuted.TotalPowerProcessed += job.Spec.PowerNeeded
		} else if _, ok := jobsMap[job.DisplayName()]; ok {
			sch.Status.State.JobsRunning++
		} else {
			sch.Status.State.JobsQueued++
		}
	}
	sch.Status.State.JobRate10Sec = jobsDone - sch.Status.State.TotalJobsExecuted.TotalJobsProcessed
	curTime := time.Now().Unix()
	if curTime-js.jobsDone60SecTime > 60 {
		sch.Status.State.JobRate1Min = jobsDone - js.jobsDone60Sec
		js.jobsDone60SecTime = curTime
		js.jobsDone60Sec = jobsDone
	}
	sch.Status.State.TotalJobsExecuted.TotalJobsProcessed = jobsDone
	js.setJSStatus(ctx, sch)
	return nil
}

package jobscheduler

import (
	jsv1 "powerschedulermodel/build/apis/jobgroup.intel.com/v1"
	nexus_client "powerschedulermodel/build/nexus-client"
	"time"

	"context"
)

func (js *JobScheduler) getJSStatus(ctx context.Context) *nexus_client.JobgroupSchedulerConfig {
	jscfg, e := js.nexusClient.RootPowerScheduler().Config().GetScheduler(ctx)
	if nexus_client.IsNotFound(e) {
		js.log.Info("Creating the Scheduler ...", e)
		itm := jsv1.SchedulerConfig{}
		jscfg, e = js.nexusClient.RootPowerScheduler().Config().AddScheduler(ctx, &itm)
	}
	if e != nil {
		js.log.Fatal("when getting Scheduler status ...  ", e)
	}
	return jscfg
}

func (js *JobScheduler) setJSStatus(ctx context.Context, jscfg *nexus_client.JobgroupSchedulerConfig) {
	if e := jscfg.SetState(ctx, &jscfg.Status.State); e != nil {
		js.log.Error("when setting Scheduler status ...... ", e)
	}
}

// periodic call to updated the status
func (js *JobScheduler) updateStats(ctx context.Context) error {

	if js.schedulerInstanceId != 0 {
		// Only let the first scheduler instance to update the scheduler stats.
		// Ideally this could be a seperate process. For now, let scheduler-0 do the job.
		return nil
	}
	js.log.Infof("Start Status Update.")

	js.statusMutex.Lock()
	defer js.statusMutex.Unlock()

	sch := js.getJSStatus(ctx)
	dcfg, e := js.nexusClient.RootPowerScheduler().GetDesiredSiteConfig(ctx)
	if e != nil {
		return e
	}

	cfg, e := js.nexusClient.RootPowerScheduler().GetConfig(ctx)
	if e != nil {
		return e
	}
	jobGroups := cfg.GetAllJobGroupsIter(ctx)
	numTotalEdges := 0
	{
		sites := dcfg.GetAllSitesDCIter(ctx)
		var e error
		var site *nexus_client.SitedcSiteDC
		sch.Status.State.BusyEdges = 0
		sch.Status.State.FreeEdges = 0
		sch.Status.State.JobsRunning = 0
		sch.Status.State.JobsQueued = 0
		sch.Status.State.TotalJobsExecuted.TotalPowerProcessed = 0
		sch.Status.State.TotalJobsExecuted.TotalJobsProcessed = js.jobsDoneTotal
		tm := time.Now().Unix()
		if (tm - js.jobsDone60SecTime) > 60 {
			sch.Status.State.JobRate1Min = uint32(float32(js.jobsDoneTotal-js.jobsDone60Sec) * 1000 / float32(tm-js.jobsDone60SecTime))
			js.jobsDone60Sec = js.jobsDoneTotal
			js.jobsDone60SecTime = tm
		}
		sch.Status.State.BusyEdgesPerSite = make(map[string]uint32)
		for site, e = sites.Next(ctx); site != nil && e == nil; site, e = sites.Next(ctx) {
			edges := site.GetAllEdgesDCIter(ctx)
			// jobsMap := make(map[string]bool)
			sn := site.DisplayName()
			sch.Status.State.BusyEdgesPerSite[sn] = 0
			for edge, _ := edges.Next(ctx); edge != nil; edge, _ = edges.Next(ctx) {
				numTotalEdges++
				ji := edge.GetAllJobsInfoIter(ctx)
				busyEdge := false
				for job, _ := ji.Next(ctx); job != nil; job, _ = ji.Next(ctx) {
					if job.Status.State.Progress < 100 {
						sch.Status.State.JobsRunning++
						if busyEdge == false {
							sch.Status.State.BusyEdges++
							busyEdge = true
							sch.Status.State.BusyEdgesPerSite[sn]++
						}
						// jobsMap[job.Spec.RequestorJob] = true
					}
				}
				if !busyEdge {
					sch.Status.State.FreeEdges++
				}
			}
		}
	}
	for jg, getJobGrpsError := jobGroups.Next(ctx); jg != nil && getJobGrpsError == nil; jg, getJobGrpsError = jobGroups.Next(ctx) {
		jobs := jg.GetAllJobsIter(ctx)
		var job *nexus_client.JobschedulerJob
		var getJobsError error
		for job, getJobsError = jobs.Next(context.Background()); job != nil && getJobsError == nil; job, getJobsError = jobs.Next(context.Background()) {

			if job.Status.State.PercentCompleted >= 100 {
				continue
				// sch.Status.State.TotalJobsExecuted.TotalPowerProcessed += job.Spec.PowerNeeded
			}

			preq := uint32(0)
			for _, pex := range job.Status.State.Execution {
				preq += pex.PowerRequested
			}

			if preq == 0 {
				sch.Status.State.JobsQueued++
			}
		}

		if getJobsError != nil {
			js.log.Infof("Error: %+v", getJobsError)

			return getJobsError
		}
	}
	js.log.Infof("DEBUG Total updateStats: numTotalEdges: %d BusyEdges %d JobsRunning %d JobsQueued %d", numTotalEdges, sch.Status.State.BusyEdges, sch.Status.State.JobsRunning, sch.Status.State.JobsQueued)

	// sch.Status.State.JobRate10Sec = jobsDone - sch.Status.State.TotalJobsExecuted.TotalJobsProcessed
	// curTime := time.Now().Unix()
	// if curTime-js.jobsDone60SecTime > 60 {
	// 	sch.Status.State.JobRate1Min = jobsDone - js.jobsDone60Sec
	// 	js.jobsDone60SecTime = curTime
	//	js.jobsDone60Sec = jobsDone
	//}
	//sch.Status.State.TotalJobsExecuted.TotalJobsProcessed = jobsDone
	js.setJSStatus(ctx, sch)
	js.log.Infof("Done Status Update.")

	return nil
}

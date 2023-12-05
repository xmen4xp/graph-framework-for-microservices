package edgeagent

import (
	"context"
	"log"
	"math/rand"
	nexus_client "powerschedulermodel/build/nexus-client"
	"sync"

	"time"

	"github.com/sirupsen/logrus"
)

// Edge agent that is going to run on each edge and will be receiving the
// job request and will mark it complete after the power requested was delivered.
// This module operates on the desiredconfig node and works on a single job
// there is a random selection of power available at any time and will be used to make progress on the job
//
// The module show a periodic implementation and interaction with the data model with time based reconciler
type EdgeAgent struct {
	nexusClient           *nexus_client.Clientset
	SiteName              string
	EdgeName              string
	AvailablePowerOptions [4]int32
	AvailableCurrentPower int32
	powerSwitchCnt        int32 // slot when available power will be switched
	statusMutex           *sync.Mutex
	log                   *logrus.Entry
	inProgressJobs        sync.Map
}

type jobProgressInfo struct {
	progress uint32
}

func New(siteName string, edgeName string, nclient *nexus_client.Clientset) *EdgeAgent {
	pm := int32(1)
	e := &EdgeAgent{nclient,
		siteName,
		edgeName,
		[4]int32{5 * pm, 10 * pm, 22 * pm, 37 * pm},
		5,
		0,
		&sync.Mutex{},
		logrus.WithFields(logrus.Fields{
			"module": edgeName,
		}),
		sync.Map{}}
	return e
}
func (ea *EdgeAgent) getEdgeStatus(ctx context.Context) *nexus_client.EdgeEdge {
	site, e := ea.nexusClient.RootPowerScheduler().Inventory().GetSite(ctx, ea.SiteName)
	if e != nil {
		ea.log.Fatal(e)
	}

	edge, e := site.GetEdges(ctx, ea.EdgeName)
	if e != nil {
		ea.log.Fatal(e)
	}
	return edge
}
func (ea *EdgeAgent) setEdgeStatus(ctx context.Context, edge *nexus_client.EdgeEdge) {
	// if e := edge.SetState(ctx, &edge.Status.State); e != nil {
	// 	ea.log.Fatal("when setting edge status ", e)
	// }
}

func (ea *EdgeAgent) updateJobStatus(ctx context.Context, job *nexus_client.JobmgmtJobInfo) error {
	updateJobStatus := false
	if job.Status.State.StartTime == 0 {
		job.Status.State.StartTime = time.Now().Unix()
		job.Status.State.Progress = 0
		job.Status.State.EndTime = job.Status.State.StartTime
		ea.log.Infof("Init Status on Job with %s powerRequested %d", job.DisplayName(), job.Spec.PowerNeeded)

		// Job is being initialized. Update the job status to capture start of job.
		updateJobStatus = true

		/*
			 		ea.statusMutex.Lock()
					defer ea.statusMutex.Unlock()
					edge := ea.getEdgeStatus(ctx)
					edge.Status.State.CurrentJob = job.DisplayName()
					edge.Status.State.CurrentJobPercentCompleted = 0
					edge.Status.State.CurrentJobPowerRequested = job.Spec.PowerNeeded
					ea.setEdgeStatus(ctx, edge)
		*/
	} else {
		var remainingPower int32 = int32(1.0 * job.Spec.PowerNeeded * (100.0 - job.Status.State.Progress) / 100.0)
		curTime := time.Now().Unix()
		var elapsed int32 = int32(curTime - job.Status.State.EndTime)
		remainingPower = int32(remainingPower) - elapsed*ea.AvailableCurrentPower
		if remainingPower <= 0 {
			job.Status.State.Progress = 100
			job.Status.State.EndTime = curTime
			ea.log.Infof("Completing the job id %s with requested power %d", job.DisplayName(), job.Spec.PowerNeeded)

			// Job has completed. Update the job status to capture end of job.
			updateJobStatus = true

			/* 			ea.statusMutex.Lock()
			   			defer ea.statusMutex.Unlock()
			   			edge := ea.getEdgeStatus(ctx)
			   			edge.Status.State.TotalJobsProcessed += 1
			   			edge.Status.State.TotalPowerServed += job.Spec.PowerNeeded
			   			edge.Status.State.CurrentJob = ""
			   			edge.Status.State.CurrentJobPercentCompleted = 0
			   			edge.Status.State.CurrentJobPowerRequested = 0
			   			ea.setEdgeStatus(ctx, edge) */
		} else {
			p := 1.0 * (job.Spec.PowerNeeded - uint32(remainingPower)) * 100.0 / job.Spec.PowerNeeded
			job.Status.State.Progress = p
			job.Status.State.EndTime = curTime
			ea.log.Infof("Updating the job id %s Progress %d pn %d rp %d available %d elaps %d", job.DisplayName(), p,
				job.Spec.PowerNeeded, remainingPower, ea.AvailableCurrentPower, elapsed)
			/* 			ea.statusMutex.Lock()
			   			defer ea.statusMutex.Unlock()
			   			edge := ea.getEdgeStatus(ctx)
			   			edge.Status.State.CurrentJob = job.DisplayName()
			   			edge.Status.State.CurrentJobPercentCompleted = p
			   			ea.setEdgeStatus(ctx, edge) */
		}
	}

	if updateJobStatus == true {
		if err := job.SetState(ctx, &job.Status.State); err != nil {
			log.Fatalf("Updating status of job %s failed with error %v", job.DisplayName(), err)
		}
	}

	if job.Status.State.Progress == 100 {
		// Job has finished. There is no need to track it in-memory.
		ea.inProgressJobs.Delete(job.DisplayName())
	} else {
		// Job is in progress. Track progress in-memory.
		ea.inProgressJobs.Store(job.DisplayName(), jobProgressInfo{progress: job.Status.State.Progress})
	}
	return nil
}

func (ea *EdgeAgent) JobPeriodicReconciler(ctx context.Context) {
	if ea.powerSwitchCnt == 0 {
		ea.powerSwitchCnt = rand.Int31n(20) // next switch slot
		ea.AvailableCurrentPower = ea.AvailablePowerOptions[rand.Intn(len(ea.AvailablePowerOptions))]
	} else {
		ea.powerSwitchCnt--
	}
	dcfg := ea.nexusClient.RootPowerScheduler().DesiredSiteConfig()

	sdc, e := dcfg.GetSitesDC(ctx, ea.SiteName)
	if nexus_client.IsNotFound(e) {
		// no requests yet so nothing to do
		return
	} else if e != nil {
		ea.log.Error("Error on getting SiteDC", e)
		return
	}
	edc, e := sdc.GetEdgesDC(ctx, ea.EdgeName)
	if nexus_client.IsNotFound(e) {
		// no requests yet so nothing to do
		return
	} else if e != nil {
		ea.log.Error("Error on getting EdgeDC", e)
		return
	}
	jobInfo := edc.GetAllJobsInfoIter(ctx)
	foundOneJob := false
	for job, _ := jobInfo.Next(ctx); job != nil; job, _ = jobInfo.Next(ctx) {
		if job.Status.State.Progress == 100 {
			continue
		}
		foundOneJob = true
		// starting a goroutine for a job would be there can be multiple jobs running at one time.
		// but the edge agent status is capable of being updated with onely one job info.
		go func(j *nexus_client.JobmgmtJobInfo) {
			if e := ea.updateJobStatus(ctx, j); e != nil {
				ea.log.Error("Updating job", e)
			}
		}(job)
	}
	if !foundOneJob {
		ea.log.Info("Did not find any job pending in this time period.")
	}

}
func (ea *EdgeAgent) Start(gctx context.Context) error {
	// ea.nexusClient.RootPowerScheduler().Inventory().Subscribe()
	// ea.nexusClient.RootPowerScheduler().DesiredEdgeConfig().Subscribe()
	// ea.nexusClient.RootPowerScheduler().DesiredEdgeConfig().EdgesDC("*").Subscribe()
	// ea.nexusClient.RootPowerScheduler().DesiredEdgeConfig().EdgesDC("*").JobsInfo("*").Subscribe()
	// ea.nexusClient.RootPowerScheduler().Inventory().Edges("*").Subscribe()
	ea.log.Infof("Starting Site:%s Edge:%s", ea.SiteName, ea.EdgeName)
	tickerJob := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-tickerJob.C:
			ea.JobPeriodicReconciler(gctx)
		case <-gctx.Done():
			ea.log.Debug("Closing ticker")
			return gctx.Err()
		}
	}
}

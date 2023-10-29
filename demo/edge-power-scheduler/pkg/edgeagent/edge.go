package edgeagent

import (
	"context"
	"math/rand"
	nexus_client "powerschedulermodel/build/nexus-client"
	"sync"

	"time"

	"github.com/sirupsen/logrus"
)

type EdgeAgent struct {
	nexusClient           *nexus_client.Clientset
	EdgeName              string
	AvailablePowerOptions [4]int32
	AvailableCurrentPower int32
	powerSwitchCnt        int32 // slot when available power will be switched
	statusMutex           *sync.Mutex
	log                   *logrus.Entry
}

func New(name string, nclient *nexus_client.Clientset) *EdgeAgent {
	e := &EdgeAgent{nclient,
		name,
		[4]int32{5, 10, 22, 37},
		5,
		0,
		&sync.Mutex{},
		logrus.WithFields(logrus.Fields{
			"module": "edgeAgent",
		})}
	return e
}
func (ea *EdgeAgent) getEdgeStatus(ctx context.Context) *nexus_client.EdgeEdge {
	edge, e := ea.nexusClient.RootPowerScheduler().Inventory().GetEdges(ctx, ea.EdgeName)
	if e != nil {
		ea.log.Fatal(e)
	}
	return edge
}
func (ea *EdgeAgent) setEdgeStatus(ctx context.Context, edge *nexus_client.EdgeEdge) {
	if e := edge.SetState(ctx, &edge.Status.State); e != nil {
		ea.log.Fatal("when setting edge status ", e)
	}
}

func (ea *EdgeAgent) updateJobStatus(ctx context.Context, job *nexus_client.JobmgmtJobInfo) error {
	if job.Status.State.StartTime == 0 {
		job.Status.State.StartTime = time.Now().Unix()
		job.Status.State.Progress = 0
		job.Status.State.EndTime = job.Status.State.StartTime
		ea.log.Infof("Init Status on Job with %s powerRequested %d", job.DisplayName(), job.Spec.PowerNeeded)
		ea.statusMutex.Lock()
		defer ea.statusMutex.Unlock()
		edge := ea.getEdgeStatus(ctx)
		edge.Status.State.CurrentJob = job.DisplayName()
		edge.Status.State.CurrentJobPercentCompleted = 0
		edge.Status.State.CurrentJobPowerRequested = job.Spec.PowerNeeded
		ea.setEdgeStatus(ctx, edge)
	} else {
		var remainingPower int32 = int32(1.0 * job.Spec.PowerNeeded * (100.0 - job.Status.State.Progress) / 100.0)
		curTime := time.Now().Unix()
		var elapsed int32 = int32(curTime - job.Status.State.EndTime)
		remainingPower = int32(remainingPower) - elapsed*ea.AvailableCurrentPower
		if remainingPower <= 0 {
			job.Status.State.Progress = 100
			job.Status.State.EndTime = curTime
			ea.log.Infof("Completing the job id %s with requested power %d", job.DisplayName(), job.Spec.PowerNeeded)
			ea.statusMutex.Lock()
			defer ea.statusMutex.Unlock()
			edge := ea.getEdgeStatus(ctx)
			edge.Status.State.TotalJobsProcessed += 1
			edge.Status.State.TotalPowerServed += job.Spec.PowerNeeded
			edge.Status.State.CurrentJob = ""
			edge.Status.State.CurrentJobPercentCompleted = 0
			edge.Status.State.CurrentJobPowerRequested = 0
			ea.setEdgeStatus(ctx, edge)
		} else {
			p := 1.0 * (job.Spec.PowerNeeded - uint32(remainingPower)) * 100.0 / job.Spec.PowerNeeded
			job.Status.State.Progress = p
			job.Status.State.EndTime = curTime
			ea.log.Infof("Updating the job id %s Progress %d pn %d rp %d available %d elaps %d", job.DisplayName(), p,
				job.Spec.PowerNeeded, remainingPower, ea.AvailableCurrentPower, elapsed)
			ea.statusMutex.Lock()
			defer ea.statusMutex.Unlock()
			edge := ea.getEdgeStatus(ctx)
			edge.Status.State.CurrentJob = job.DisplayName()
			edge.Status.State.CurrentJobPercentCompleted = p
			ea.setEdgeStatus(ctx, edge)
		}
	}
	return job.SetState(ctx, &job.Status.State)
}

func (ea *EdgeAgent) JobPeriodicReconciler(ctx context.Context) {
	if ea.powerSwitchCnt == 0 {
		ea.powerSwitchCnt = rand.Int31n(20) // next switch slot
		ea.AvailableCurrentPower = ea.AvailablePowerOptions[rand.Intn(len(ea.AvailablePowerOptions))]
	} else {
		ea.powerSwitchCnt--
	}
	dcfg := ea.nexusClient.RootPowerScheduler().DesiredEdgeConfig()

	edc, e := dcfg.GetEdgesDC(ctx, ea.EdgeName)
	if nexus_client.IsNotFound(e) {
		// no requests yet so nothing to do
		return
	} else if e != nil {
		ea.log.Error("Error on getting EdgeDC", e)
		return
	}
	jobInfo, e := edc.GetAllJobsInfo(ctx)
	if e != nil {
		ea.log.Error("Error getting JobInfo", e)
		return
	}
	foundOneJob := false
	for _, job := range jobInfo {
		if job.Status.State.Progress == 100 {
			continue
		}
		foundOneJob = true
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

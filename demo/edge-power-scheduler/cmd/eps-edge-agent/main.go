package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"powerscheduler/pkg/dminit"
	ev1 "powerschedulermodel/build/apis/edge.intel.com/v1"
	dev1 "powerschedulermodel/build/apis/edgedc.intel.com/v1"
	nexus_client "powerschedulermodel/build/nexus-client"

	"golang.org/x/sync/errgroup"

	"time"

	"github.com/sirupsen/logrus"

	"k8s.io/client-go/rest"
)

func generateTrackingId() int {
	return rand.Intn(10000-1) + 1
}

var nexusClient *nexus_client.Clientset
var AvailablePowerOptions [4]int32 = [4]int32{5, 10, 22, 37}
var AvailableCurrentPower int32 = AvailablePowerOptions[0]

func updateJobStatus(ctx context.Context, job *nexus_client.JobmgmtJobInfo) error {
	log := logrus.WithFields(logrus.Fields{
		"module": "updateJobStatus",
	})
	if job.Status.State.Progress == 100 {
		return nil
	}
	if job.Status.State.StartTime == 0 {
		job.Status.State.StartTime = time.Now().Unix()
		job.Status.State.Progress = 0
		job.Status.State.EndTime = job.Status.State.StartTime
		log.Infof("Init Status on Job with %s", job.Name)
	} else {
		var remainingPower int32 = int32(1.0 * job.Spec.PowerNeeded * (100.0 - job.Status.State.Progress) / 100.0)
		curTime := time.Now().Unix()
		var elapsed int32 = int32(curTime - job.Status.State.EndTime)
		remainingPower = int32(remainingPower) - elapsed*AvailableCurrentPower
		if remainingPower <= 0 {
			job.Status.State.Progress = 100
			job.Status.State.EndTime = curTime
			log.Infof("Completing the job id %s", job.DisplayName())
		} else {
			p := 1.0 * (job.Spec.PowerNeeded - uint32(remainingPower)) * 100.0 / job.Spec.PowerNeeded
			job.Status.State.Progress = p
			job.Status.State.EndTime = curTime
			log.Infof("Updating the job id %s Progress %d pn %d rp %d available %d elaps %d", job.DisplayName(), p,
				job.Spec.PowerNeeded, remainingPower, AvailableCurrentPower, elapsed)

		}
	}
	return job.SetState(ctx, &job.Status.State)
	//return job.Update(ctx)
}

func jobPeriodicReconciler(
	ctx context.Context,
	log *logrus.Logger,
	nexusClient *nexus_client.Clientset,
	dcfg *nexus_client.DesiredconfigDesiredEdgeConfig,
	edgeName string) {

	edc, e := dcfg.GetEdgesDC(ctx, edgeName)
	if e != nil {
		//log.Warn("Error on getting EdgeDC", e)
		return
	}
	jobInfo, e := edc.GetAllJobsInfo(ctx)
	if e != nil {
		log.Error("Error getting JobInfo", e)
		return
	}
	for _, job := range jobInfo {
		go func(j *nexus_client.JobmgmtJobInfo) {
			if e := updateJobStatus(ctx, j); e != nil {
				log.Error("Updating job", e)
			}
		}(job)
	}
}

func main() {
	edgeName := os.Getenv("EDGE_NAME")
	dmAPIGWPort := os.Getenv("DM_APIGW_PORT")
	if edgeName == "" {
		edgeName = "edge-unk"
	}
	if dmAPIGWPort == "" {
		dmAPIGWPort = "8100"
	}
	rand.Seed(time.Now().UnixNano())
	var log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "15:04:05.000"
	logrus.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true

	host := "localhost:" + dmAPIGWPort
	log.Infof("Edge %s Creating client to host at : %s", edgeName, host)
	var e error

	nexusClient, e = nexus_client.NewForConfig(&rest.Config{Host: host})
	if e != nil {
		log.Fatal("nexusClient", e)
	}
	ctx := context.Background()
	e = dminit.Init(ctx, nexusClient)
	if e != nil {
		log.Fatal(e)
	}
	// add inventory item for self
	inv := nexusClient.RootPowerScheduler().Inventory()
	newEdge, e := inv.GetEdges(ctx, edgeName)
	edge := ev1.Edge{}
	edge.Name = edgeName
	if nexus_client.IsNotFound(e) {
		newEdge, e = inv.AddEdges(ctx, &edge)
		if e != nil {
			log.Error("Creating Inv Edge ", e)
		}
	} else if e != nil {
		log.Error("when getting Edge ", e)
	}
	// add desired config node
	dc := nexusClient.RootPowerScheduler().DesiredEdgeConfig()
	_, e = dc.GetEdgesDC(ctx, edgeName)
	if nexus_client.IsNotFound(e) {
		dcedge := dev1.EdgeDC{}
		dcedge.Name = edgeName
		_, e = dc.AddEdgesDC(ctx, &dcedge)
		if e != nil {
			log.Error("Creating DCEdge ", e)
		}
	} else if e != nil {
		log.Error("when getting DCEdge ", e)
	}

	newEdge.Status.State.PowerInfo.TotalPowerAvailable = 100
	newEdge.Status.State.PowerInfo.FreePowerAvailable = 0
	if e = newEdge.SetState(ctx, &newEdge.Status.State); e != nil {
		log.Error("When setting state of the Edge", e)
	}
	powerChangeAt := 6
	powerChangeCnt := 0
	cnt := 0

	ctx, done := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)
	eagent := edgeagent.New(edgeName, nexusClient)
	// look for signal
	g.Go(func() error {
		return dminit.SignalInit(gctx, done, log)
	})

	// timer task
	g.Go(func() error { return eagent.Start(gctx) })

	err := g.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Debug("Context was cancelled")
		} else {
			log.Debugf("received error: %v", err)
		}
	}

	fmt.Println("exiting")
}

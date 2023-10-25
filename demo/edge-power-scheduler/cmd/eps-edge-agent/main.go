package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"powerscheduler/pkg/dminit"
	ev1 "powerschedulermodel/build/apis/edge.intel.com/v1"
	pmv1 "powerschedulermodel/build/apis/powermgmt.intel.com/v1"
	nexus_client "powerschedulermodel/build/nexus-client"

	"golang.org/x/sync/errgroup"

	"syscall"
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

	host := "localhost:" + dmAPIGWPort
	log.Infof("Edge %s Creating client to host at : %s", edgeName, host)
	var e error

	nexusClient, e = nexus_client.NewForConfig(&rest.Config{Host: host})
	if e != nil {
		log.Fatal("nexusClient", e)
	}
	ctx := context.Background()
	_, _, dcfg, inv, e := dminit.Init(ctx, nexusClient)
	if e != nil {
		log.Fatal(e)
	}
	// add inventory item for self

	newEdge, e := inv.GetEdges(ctx, edgeName)
	edge := ev1.Edge{}
	edge.Name = edgeName
	if e != nil {
		newEdge, e = inv.AddEdges(ctx, &edge)
		if e != nil {
			log.Warn("Creating Inv Edge", e)
		}
	}

	pi, e := newEdge.GetPowerInfo(ctx)
	if e != nil {
		pi := pmv1.PowerInfo{}
		pi.Spec.TotalPowerAvailable = 100
		pi.Spec.FreePowerAvailable = 0
		_, e = newEdge.AddPowerInfo(ctx, &pi)
		if e != nil {
			log.Warn("Creating Edge PowerInfo", e)
		}
	} else {
		pi.Spec.TotalPowerAvailable = 100
		pi.Spec.FreePowerAvailable = 0
		pi.Update(ctx)
	}
	powerChangeAt := 6
	powerChangeCnt := 0
	cnt := 0

	ctx, done := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)

	// look for signal
	g.Go(func() error {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-sigs:
			log.Infof("Received signal: %s", sig)
			done()
		case <-gctx.Done():
			log.Debug("Closing signal goroutine")
			return gctx.Err()
		}
		return nil
	})

	// timer task
	g.Go(func() error {
		tickerJob := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-tickerJob.C:
				powerChangeCnt += 1
				if powerChangeCnt == powerChangeAt {
					powerChangeCnt = 0
					AvailableCurrentPower = AvailablePowerOptions[rand.Intn(len(AvailablePowerOptions))]
				}
				// log.Infof("Iteration %d", cnt)
				cnt += 1
				jobPeriodicReconciler(ctx, log, nexusClient, dcfg, edgeName)
			case <-gctx.Done():
				log.Debug("Closing ticker")
				return gctx.Err()
			}
		}
	})

	err := g.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Debug("Context was cancled")
		} else {
			log.Debugf("received error: %v", err)
		}
	}

	fmt.Println("exiting")
}

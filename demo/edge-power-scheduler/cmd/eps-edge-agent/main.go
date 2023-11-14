package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"powerscheduler/pkg/dminit"
	"powerscheduler/pkg/edgeagent"
	ev1 "powerschedulermodel/build/apis/edge.intel.com/v1"
	dev1 "powerschedulermodel/build/apis/edgedc.intel.com/v1"
	nexus_client "powerschedulermodel/build/nexus-client"

	"golang.org/x/sync/errgroup"

	"time"

	"github.com/sirupsen/logrus"

	"k8s.io/client-go/rest"
)

func main() {
	edgeName := os.Getenv("EDGE_NAME")
	dmAPIGWPort := os.Getenv("DM_APIGW_PORT")
	if edgeName == "" {
		edgeName = "edge-unk"
	}
	if dmAPIGWPort == "" {
		dmAPIGWPort = "8000"
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

	nexusClient, e := nexus_client.NewForConfig(&rest.Config{Host: host})
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

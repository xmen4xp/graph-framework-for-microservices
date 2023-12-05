package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"powerscheduler/pkg/dminit"
	"powerscheduler/pkg/edgeagent"
	ev1 "powerschedulermodel/build/apis/edge.intel.com/v1"
	dev1 "powerschedulermodel/build/apis/edgedc.intel.com/v1"
	sv1 "powerschedulermodel/build/apis/site.intel.com/v1"
	sdcv1 "powerschedulermodel/build/apis/sitedc.intel.com/v1"
	nexus_client "powerschedulermodel/build/nexus-client"
	"strconv"

	"golang.org/x/sync/errgroup"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	"time"

	"github.com/sirupsen/logrus"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func createEdgeAgent(
	ctx context.Context, log *logrus.Logger,
	nexusClient *nexus_client.Clientset,
	siteInv *nexus_client.SiteSite,
	siteDC *nexus_client.SitedcSiteDC,
	edgeName string) *edgeagent.EdgeAgent {

	edge := ev1.Edge{}
	edge.Name = edgeName
	newEdge, e := siteInv.AddEdges(ctx, &edge)
	if e != nil && k8serrors.IsAlreadyExists(e) == false {
		log.Fatalf("Creating Inv Edge %s failed with error %v", edgeName, e)
	}

	// add desired config node
	_, e = siteDC.GetEdgesDC(ctx, edgeName)
	if nexus_client.IsNotFound(e) || nexus_client.IsChildNotFound(e) {
		dcedge := dev1.EdgeDC{}
		dcedge.Name = edgeName
		_, e = siteDC.AddEdgesDC(ctx, &dcedge)
		if e != nil {
			log.Fatalf("Creating DCEdge %s failed with error %v", edgeName, e)
		}
	} else if e != nil {
		log.Fatalf("when getting DCEdge: err %+v", e)
	}

	newEdge.Status.State.PowerInfo.TotalPowerAvailable = 100
	newEdge.Status.State.PowerInfo.FreePowerAvailable = 0
	if e = newEdge.SetState(ctx, &newEdge.Status.State); e != nil {
		log.Error("When setting state of the Edge", e)
	}

	eagent := edgeagent.New(siteInv.DisplayName(), edgeName, nexusClient)
	return eagent
}

func main() {
	var kubeconfig string
	flag.StringVar(&kubeconfig, "k", "", "Absolute path to the kubeconfig file. Defaults to ~/.kube/config.")
	flag.Parse()

	siteName := os.Getenv("EDGE_NAME")
	dmAPIGWPort := os.Getenv("DM_APIGW_PORT")
	if siteName == "" {
		siteName = "edge-unk"
	}
	if dmAPIGWPort == "" {
		dmAPIGWPort = "8000"
	}
	edgeCount, e := strconv.Atoi(os.Getenv("EDGE_COUNT"))
	if e != nil {
		edgeCount = 1
	}
	rand.Seed(time.Now().UnixNano())
	var log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "15:04:05.000"
	logrus.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true

	host := "localhost:" + dmAPIGWPort
	log.Infof("Edge:: %s Creating client to host at : %s", siteName, host)

	var config *rest.Config
	if len(kubeconfig) != 0 {
		var err error
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err)
		}
	} else {
		config = &rest.Config{Host: host}
	}
	config.Burst = 1000
	config.QPS = 500
	nexusClient, e := nexus_client.NewForConfig(config)
	if e != nil {
		log.Fatal("nexusClient", e)
	}
	ctx := context.Background()
	e = dminit.Init(ctx, nexusClient)
	if e != nil {
		log.Fatal(e)
	}

	ctx, done := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)
	// look for signal
	g.Go(func() error {
		return dminit.SignalInit(gctx, done, log)
	})

	nexusClient.RootPowerScheduler().Inventory().Subscribe()
	nexusClient.RootPowerScheduler().DesiredSiteConfig().Subscribe()
	nexusClient.RootPowerScheduler().DesiredSiteConfig().SitesDC(siteName).Subscribe()
	nexusClient.RootPowerScheduler().DesiredSiteConfig().SitesDC(siteName).EdgesDC("*").Subscribe()
	nexusClient.RootPowerScheduler().DesiredSiteConfig().SitesDC(siteName).EdgesDC("*").JobsInfo("*").Subscribe()
	nexusClient.RootPowerScheduler().Inventory().Site(siteName).Edges("*").Subscribe()
	log.Infof("Creating site %s", siteName)
	site, e := nexusClient.RootPowerScheduler().Inventory().GetSite(ctx, siteName)
	if nexus_client.IsNotFound(e) {
		s := &sv1.Site{}
		s.SetName(siteName)
		site, e = nexusClient.RootPowerScheduler().Inventory().AddSite(ctx, s)
		if e != nil {
			log.Fatalf("Error creating site %v", e)
		}
	} else if e != nil {
		log.Fatalf("Error getting site %v", e)
	}

	siteDC, e := nexusClient.RootPowerScheduler().DesiredSiteConfig().GetSitesDC(ctx, siteName)
	if nexus_client.IsNotFound(e) {
		s := &sdcv1.SiteDC{}
		s.SetName(siteName)
		siteDC, e = nexusClient.RootPowerScheduler().DesiredSiteConfig().AddSitesDC(ctx, s)
		if e != nil {
			log.Fatalf("Error creating sitedc %v", e)
		}
	} else if e != nil {
		log.Fatalf("Error getting sitedc %v", e)
	}

	if edgeCount == 1 {
		eagent := createEdgeAgent(ctx, log, nexusClient, site, siteDC, siteName)
		// timer task
		g.Go(func() error { return eagent.Start(gctx) })
	} else {
		for eCount := 0; eCount < edgeCount; eCount++ {
			eName := fmt.Sprintf("%s-%d", siteName, eCount)
			log.Infof("Creating edge %s", eName)
			eagent := createEdgeAgent(ctx, log, nexusClient, site, siteDC, eName)
			go eagent.Start(gctx)
		}
	}
	log.Infof("Done Creating all edges.")
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

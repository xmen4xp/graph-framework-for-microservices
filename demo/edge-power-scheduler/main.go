package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	cfgv1 "powerschedulermodel/build/apis/config.intel.com/v1"
	dcfgv1 "powerschedulermodel/build/apis/desiredconfig.intel.com/v1"
	invv1 "powerschedulermodel/build/apis/inventory.intel.com/v1"
	rootv1 "powerschedulermodel/build/apis/root.intel.com/v1"
	nexus_client "powerschedulermodel/build/nexus-client"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"k8s.io/client-go/rest"
)

func generateTrackingId() int {
	return rand.Intn(10000-1) + 1
}

func ProcessJobs(job *nexus_client.JobschedulerJob) {
	fmt.Printf("Got a call back  for JobschedulerJob name: %s displayName: %s spec: %v",
		job.Name, job.DisplayName(), job.Spec)
}

var nexusClient *nexus_client.Clientset

func main() {

	rand.Seed(time.Now().UnixNano())
	var log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	host := "localhost:8101"
	log.Info("Creating client to host at : ", host)
	var e error
	nexusClient, e = nexus_client.NewForConfig(&rest.Config{Host: host})
	if e != nil {
		log.Fatal("nexusClient", e)
	}
	ctx := context.Background()
	// ps := v1.PowerScheduler{}
	// nexusClient.RootSockShop().PO("*").RegisterAddCallback(ProcessOrder)
	// create some nodes that are needed for the system to function.
	rps, e := nexusClient.GetRootPowerScheduler(ctx)
	if e != nil {
		log.Info("Get on RootPowerScheduler node resulted in error:", e)
		log.Info("Will create a new RootPowerScheduler node")
		rps, e = nexusClient.AddRootPowerScheduler(ctx, &rootv1.PowerScheduler{})
		if e != nil {
			log.Fatal("RootPowerScheduler", e)
		}
	}
	_, e = rps.GetConfig(ctx)
	if e != nil {
		log.Info("Get on Config node resulted in error:", e)
		log.Info("Will create a new Config node")
		_, e = rps.AddConfig(ctx, &cfgv1.Config{})
		if e != nil {
			log.Fatal("Config", e)
		}
	}
	_, e = rps.GetDesiredEdgeConfig(ctx)
	if e != nil {
		log.Info("Get on DesiredConfig node resulted in error:", e)
		log.Info("Will create a new DesiredConfig node")
		_, e = rps.AddDesiredEdgeConfig(ctx, &dcfgv1.DesiredEdgeConfig{})
		if e != nil {
			log.Fatal("DesiredEdgeConfig", e)
		}
	}
	_, e = rps.GetInventory(ctx)
	if e != nil {
		log.Info("Get on Inventory node resulted in error:", e)
		log.Info("Will create a new Inventory node")
		_, e = rps.AddInventory(ctx, &invv1.Inventory{})
		if e != nil {
			log.Fatal("Inventory", e)
		}
	}
	// create subscriptions
	nexusClient.RootPowerScheduler().Config().Jobs("*").RegisterAddCallback(ProcessJobs)
	//	nexusClient.RootPowerScheduler().Inventory().Edges("*").RegisterAddCallback(ProcessEdges)
	//	nexusClient.RootPowerScheduler().Inventory().Edges("*").PowerInfo("*").RegisterAddCallback(ProcessEdgePowerInfo)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		<-sigs
		done <- true
	}()
	<-done
	fmt.Println("exiting")
}

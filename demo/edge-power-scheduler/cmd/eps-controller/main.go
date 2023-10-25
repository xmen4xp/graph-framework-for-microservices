package main

import (
	"context"
	"fmt"
	"math/rand"
	"powerscheduler/pkg/dminit"
	nexus_client "powerschedulermodel/build/nexus-client"

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

	host := "localhost:8100"
	log.Info("Creating client to host at : ", host)
	var e error

	nexusClient, e = nexus_client.NewForConfig(&rest.Config{Host: host})
	if e != nil {
		log.Fatal("nexusClient", e)
	}
	ctx := context.Background()
	_, _, _, _, e = dminit.Init(ctx, nexusClient)
	if e != nil {
		log.Fatal(e)
	}
	// create subscriptions
	// nexusClient.RootPowerScheduler().Config().Jobs("*").RegisterAddCallback(ProcessJobs)
	// nexusClient.RootPowerScheduler().DesiredEdgeConfig().EdgesDC("*").JobsInfo("*").RegisterEventHandler()
	// nexusClient.RootPowerScheduler().Inventory().Edges("*").RegisterAddCallback(ProcessEdges)
	// nexusClient.RootPowerScheduler().Inventory().Edges("*").PowerInfo("*").RegisterAddCallback(ProcessEdgePowerInfo)

	fmt.Println("exiting")
}

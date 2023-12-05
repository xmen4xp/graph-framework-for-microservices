package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"powerscheduler/pkg/dminit"
	"powerscheduler/pkg/jobscheduler"
	nexus_client "powerschedulermodel/build/nexus-client"
	"strconv"
	"strings"

	"golang.org/x/sync/errgroup"

	"time"

	"github.com/sirupsen/logrus"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	//var e error

	var kubeconfig string
	flag.StringVar(&kubeconfig, "k", "", "Absolute path to the kubeconfig file. Defaults to ~/.kube/config.")
	flag.Parse()

	dmAPIGWPort := os.Getenv("DM_APIGW_PORT")
	if dmAPIGWPort == "" {
		dmAPIGWPort = "8000"
	}

	// There is at least 1 scheduler running.
	numSchedulerInstances := 1
	schedulerInstances := os.Getenv("NUM_SCHEDULER_REPLICAS")
	if schedulerInstances != "" {
		if num, err := strconv.Atoi(schedulerInstances); err == nil {
			numSchedulerInstances = num
		}
	}

	instanceId := 0
	instanceName := os.Getenv("HOSTNAME")
	if instanceName != "" {
		s := strings.Split(instanceName, "-")
		if len(s) > 0 {
			if id, err := strconv.Atoi(s[len(s)-1]); err == nil && instanceId < numSchedulerInstances {
				instanceId = id
			}
		}
	}

	if instanceId >= numSchedulerInstances {
		log.Fatalf("Scheduler instance id %d is greater than configure number of scheduler instances %d."+
			" Set NUM_SCHEDULER_REPLICAS env to reflect the accurate number of running scheduler instances.", instanceId, numSchedulerInstances)
	}

	rand.Seed(time.Now().UnixNano())
	var log = logrus.New()
	log.SetLevel(logrus.InfoLevel)

	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "15:04:05.000"
	logrus.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true

	host := "localhost:" + dmAPIGWPort
	log.Infof("Power Scheduler Creating Client to host at : %s", host)

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

	config.Burst = 2500
	config.QPS = 2000
	nexusClient, e := nexus_client.NewForConfig(config)
	if e != nil {
		log.Fatal("nexusClient", e)
	}
	ctx := context.Background()
	e = dminit.Init(ctx, nexusClient)
	if e != nil {
		log.Fatal(e)
	}
	// add inventory item for self

	ctx, done := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)

	jobScheduler := jobscheduler.New(nexusClient, numSchedulerInstances, instanceId)

	// look for signal
	g.Go(func() error {
		return dminit.SignalInit(gctx, done, log)
	})

	jobScheduler.Start(nexusClient, g, gctx)
	err := g.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Debug("Context was cancelled")
		} else {
			log.Debugf("received error: %v", err)
		}
	}

	fmt.Println("exiting Scheduler.")
}

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"powerscheduler/pkg/dminit"
	"powerscheduler/pkg/jobcreator"
	nexus_client "powerschedulermodel/build/nexus-client"
	"strconv"

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

	groupName := "DefaultGroup"
	instanceName := os.Getenv("HOSTNAME")
	if instanceName != "" {
		groupName = instanceName
	}

	dmAPIGWPort := os.Getenv("DM_APIGW_PORT")
	if dmAPIGWPort == "" {
		dmAPIGWPort = "8000"
	}
	maxJobs, e := strconv.Atoi(os.Getenv("MAXJOBS"))
	if e != nil {
		maxJobs = 0
	}
	maxPower, e := strconv.Atoi(os.Getenv("JOB_MAXPOWER"))
	if e != nil {
		maxPower = 500
	}
	minPower, e := strconv.Atoi(os.Getenv("JOB_MINPOWER"))
	if e != nil {
		minPower = 300
	}
	pendingJobs, e := strconv.Atoi(os.Getenv("JOB_PENDING"))
	if e != nil {
		pendingJobs = 10
	}
	maxActiveJobs, e := strconv.Atoi(os.Getenv("MAX_ACTIVE_JOBS"))
	if e != nil {
		maxActiveJobs = 100
	}
	rand.Seed(time.Now().UnixNano())
	var log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "15:04:05.000"
	logrus.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true

	host := "localhost:" + dmAPIGWPort

	var config *rest.Config
	if len(kubeconfig) != 0 {
		var err error
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err)
		}
	} else {
		log.Infof("Job Generator Creating Client to host at : %s", host)
		config = &rest.Config{Host: host}
	}

	config.QPS = 2000
	config.Burst = 2500

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
	startJobID, e := strconv.Atoi(os.Getenv("START_JOBID"))
	if e != nil {
		startJobID = 0
		cfg, e := nexusClient.RootPowerScheduler().GetConfig(ctx)
		if e != nil {
			log.Error(e)
		}
		jg, e := cfg.GetJobGroups(ctx, groupName)
		if e == nil {
			jobs := jg.GetAllJobsIter(ctx)
			var jb *nexus_client.JobschedulerJob
			var getJobsError error
			for jb, getJobsError = jobs.Next(context.Background()); jb != nil && getJobsError == nil; jb, getJobsError = jobs.Next(context.Background()) {
				if startJobID < int(jb.Job.Spec.JobId) {
					startJobID = int(jb.Job.Spec.JobId)
				}
			}
			if getJobsError != nil {
				log.Error(getJobsError)
			}
		}
		startJobID++
	}

	log.Infof("Starting job with id: %d", startJobID)
	ctx, done := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)
	jobCreator := jobcreator.New(nexusClient, groupName,
		uint32(pendingJobs), uint32(maxActiveJobs), uint32(maxPower), uint32(minPower),
		uint32(maxJobs), uint64(startJobID))

	// look for signal
	g.Go(func() error {
		return dminit.SignalInit(gctx, done, log)
	})
	jobCreator.Start(nexusClient, g, gctx)
	err := g.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Debug("Context was cancelled")
		} else {
			log.Debugf("received error: %v", err)
		}
	}

	fmt.Println("exiting job generator.")
}

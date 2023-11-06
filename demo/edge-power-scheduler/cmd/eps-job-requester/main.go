package main

import (
	"context"
	"errors"
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
)

func main() {
	//var e error

	dmAPIGWPort := os.Getenv("DM_APIGW_PORT")
	if dmAPIGWPort == "" {
		dmAPIGWPort = "8100"
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
	rand.Seed(time.Now().UnixNano())
	var log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "15:04:05.000"
	logrus.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true

	host := "localhost:" + dmAPIGWPort
	log.Infof("Job Generator Creating Client to host at : %s", host)

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
	startJobID, e := strconv.Atoi(os.Getenv("START_JOBID"))
	if e != nil {
		startJobID = 0
		cfg, e := nexusClient.RootPowerScheduler().GetConfig(ctx)
		if e != nil {
			log.Error(e)
		}
		jobs, e := cfg.GetAllJobs(ctx)
		if e != nil {
			log.Error(e)
		}
		for _, jb := range jobs {
			if startJobID < int(jb.Job.Spec.JobId) {
				startJobID = int(jb.Job.Spec.JobId)
			}
		}
		startJobID++
	}
	log.Infof("Starting job with id: %d", startJobID)
	ctx, done := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)
	jobCreator := jobcreator.New(nexusClient, 10, uint32(maxPower), uint32(minPower), uint32(maxJobs), uint64(startJobID))

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

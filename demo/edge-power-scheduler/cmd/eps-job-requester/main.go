package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"powerscheduler/pkg/dminit"
	"powerscheduler/pkg/jobcreator"
	nexus_client "powerschedulermodel/build/nexus-client"
	"strconv"

	"golang.org/x/sync/errgroup"

	"syscall"
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
	rand.Seed(time.Now().UnixNano())
	var log = logrus.New()
	log.SetLevel(logrus.InfoLevel)

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

	ctx, done := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)

	jobCreator := jobcreator.New(nexusClient, 10, 1000, 100, uint32(maxJobs))

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
	jobCreator.Start(nexusClient, g, gctx)
	err := g.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Debug("Context was cancled")
		} else {
			log.Debugf("received error: %v", err)
		}
	}

	fmt.Println("exiting job generator.")
}

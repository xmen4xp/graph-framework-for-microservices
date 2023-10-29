package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"powerscheduler/pkg/dminit"
	"powerscheduler/pkg/jobscheduler"
	nexus_client "powerschedulermodel/build/nexus-client"

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
	rand.Seed(time.Now().UnixNano())
	var log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "15:04:05.000"
	logrus.SetFormatter(customFormatter)
	//logrus.Info("Hello Walrus before FullTimestamp=true")
	customFormatter.FullTimestamp = true
	//logrus.Info("Hello Walrus after FullTimestamp=true")

	host := "localhost:" + dmAPIGWPort
	log.Infof("Power Scheduler Creating Client to host at : %s", host)

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

	jobScheduler := jobscheduler.New(nexusClient)

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
	jobScheduler.Start(nexusClient, g, gctx)
	err := g.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Debug("Context was cancled")
		} else {
			log.Debugf("received error: %v", err)
		}
	}

	fmt.Println("exiting Scheduler.")
}

package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	nexus_client "trafficlight/build/nexus-client"

	"golang.org/x/sync/errgroup"

	"time"

	"trafficlightdemo/pkg/dminit"
	"trafficlightdemo/pkg/light"

	"github.com/sirupsen/logrus"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig string
	flag.StringVar(&kubeconfig, "k", "", "Absolute path to the kubeconfig file. Defaults to ~/.kube/config.")
	flag.Parse()

	dmAPIGWPort := os.Getenv("DM_APIGW_PORT")
	if dmAPIGWPort == "" {
		dmAPIGWPort = "8000"
	}

	instanceId := 0
	instanceName := os.Getenv("HOSTNAME")
	if instanceName != "" {
		s := strings.Split(instanceName, "-")
		if len(s) > 0 {
			if id, err := strconv.Atoi(s[len(s)-1]); err == nil {
				instanceId = id
			}
		}
	}

	numLights := 1
	numLightsStr := os.Getenv("NUM_LIGHTS")
	if numLightsStr != "" {
		if id, err := strconv.Atoi(numLightsStr); err == nil {
			numLights = id
		}
	}

	rand.Seed(time.Now().UnixNano())

	lvl, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		lvl = "info"
	}
	// parse string, this is built-in feature of logrus
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.DebugLevel
	}
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "15:04:05.000"
	customFormatter.FullTimestamp = true
	var log = logrus.New()
	log.SetLevel(ll)
	log.SetFormatter(customFormatter)

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

	ctx, cancelFn := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)

	dmInitError := dminit.Init(gctx, nexusClient)
	if dmInitError != nil {
		log.Fatal(dmInitError)
	}

	light.LightGroupInit(nexusClient, instanceId, numLights, log, gctx, cancelFn)

	// look for signal
	g.Go(func() error {
		return signalInit(gctx, cancelFn, log)
	})

	<-gctx.Done()

	fmt.Println("exiting light ...")
}

func signalInit(gctx context.Context, cancelFn context.CancelFunc, log *logrus.Logger) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-sigs:
		log.Infof("Received signal: %s", sig)
		cancelFn()
	}
	return nil

}

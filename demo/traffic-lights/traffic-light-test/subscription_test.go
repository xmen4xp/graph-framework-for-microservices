package subscription

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
	nexus_client "trafficlight/build/nexus-client"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var nexusClient *nexus_client.Clientset
var nexusRestConfig *rest.Config

func getNexusClient() *nexus_client.Clientset {
	var err error
	var nclient *nexus_client.Clientset

	nclient, err = nexus_client.NewForConfig(nexusRestConfig)
	if err != nil {
		log.Fatal("nexusClient", err)
	}
	return nclient
}

func retry[arg1 *nexus_client.RootRoot | *nexus_client.ConfigConfig | *nexus_client.InventoryInventory](fn func(context.Context) (arg1, error)) (arg1, error) {
	retryCount := 10
	count := 0
	for {
		obj, err := fn(context.Background())
		if err == nil {
			return obj, nil
		}

		if count < retryCount {
			count += 1
			time.Sleep(time.Second)
		} else {
			return nil, fmt.Errorf("failed to get object with error: %v", err)
		}
	}
}

func TestShouldSubscribeToNodeTypeSuccessfully(t *testing.T) {
	nexusClient.RootRoot().Subscribe()
	assert.True(t, nexusClient.RootRoot().IsSubscribed())
	assert.False(t, nexusClient.RootRoot().Config().IsSubscribed())
	assert.False(t, nexusClient.RootRoot().DesiredConfig().IsSubscribed())
	assert.False(t, nexusClient.RootRoot().Inventory().IsSubscribed())
}

func TestShouldNotReturnConfigAsItIsNotSubscribed(t *testing.T) {
	root, err := retry(nexusClient.GetRootRoot)
	assert.NoError(t, err)

	config, err := retry(root.GetConfig)
	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestShouldReturnConfigWhenSubscribed(t *testing.T) {
	nexusClient.RootRoot().Config().Subscribe()
	root, err := retry(nexusClient.GetRootRoot)
	assert.NoError(t, err)

	config, err := retry(root.GetConfig)
	assert.NoError(t, err)
	assert.NotNil(t, config)
}

func TestShouldSubscribeAllSuccessfuly(t *testing.T) {
	nc := getNexusClient()
	nc.SubscribeAll()
	assert.True(t, nc.RootRoot().IsSubscribed())
	assert.True(t, nc.RootRoot().Config().IsSubscribed())
	assert.True(t, nc.RootRoot().Config().PerfTest().IsSubscribed())
	assert.True(t, nc.RootRoot().Config().Group("*").IsSubscribed())
	assert.True(t, nc.RootRoot().Config().Group("*").IsSubscribed())

	assert.True(t, nc.RootRoot().DesiredConfig().IsSubscribed())
	assert.True(t, nc.RootRoot().DesiredConfig().Lights("*").IsSubscribed())
	assert.True(t, nc.RootRoot().DesiredConfig().Lights("*").Status("*").IsSubscribed())

	assert.True(t, nc.RootRoot().Inventory().IsSubscribed())
	assert.True(t, nc.RootRoot().DesiredConfig().Lights("*").IsSubscribed())

	root, err := retry(nexusClient.GetRootRoot)
	assert.NoError(t, err)

	inventory, err := retry(root.GetInventory)
	assert.NoError(t, err)
	assert.NotNil(t, inventory)
}

func TestMain(m *testing.M) {

	var kubeconfig string
	var apiGWHostPort string
	var err error
	defaultApiGWHostPort := "localhost:8082"

	flag.StringVar(&kubeconfig, "k", "", "Absolute path to the kubeconfig file. Defaults to ~/.kube/config.")
	flag.StringVar(&apiGWHostPort, "a", defaultApiGWHostPort, fmt.Sprintf("hostname:port where api gw is reachable. Defaults to %s", defaultApiGWHostPort))
	flag.Parse()

	if len(kubeconfig) != 0 {
		nexusRestConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err)
		}
	} else {
		nexusRestConfig = &rest.Config{Host: apiGWHostPort}
	}

	nexusRestConfig.Burst = 2500
	nexusRestConfig.QPS = 2000
	nexusClient = getNexusClient()

	os.Exit(m.Run())
}

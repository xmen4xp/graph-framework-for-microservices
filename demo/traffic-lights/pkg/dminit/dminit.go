package dminit

import (
	"context"
	cfgv1 "trafficlight/build/apis/config.trafficlight.com/v1"
	dcfgv1 "trafficlight/build/apis/desiredconfig.trafficlight.com/v1"
	invv1 "trafficlight/build/apis/inventory.trafficlight.com/v1"
	rootv1 "trafficlight/build/apis/root.trafficlight.com/v1"
	nexus_client "trafficlight/build/nexus-client"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Initialized some of the base node that are needed to anchor the graph at the top level
// These are singleton and will be mostly static.
func Init(ctx context.Context, nexusClient *nexus_client.Clientset) error {
	// A common pattern is to re-use fields between logging statements by re-using
	// the logrus.Entry returned from WithFields()
	log := logrus.WithFields(logrus.Fields{
		"module": "dminit",
	})

	// create some nodes that are needed for the system to function.
	root, e := nexusClient.GetRootRoot(ctx)
	if nexus_client.IsNotFound(e) {
		log.Info("Node not found. Will create a new Root node")
		if root, e = nexusClient.AddRootRoot(ctx, &rootv1.Root{}); e != nil {
			return errors.WithMessage(e, "When creating Root Node")
		}
	} else if e != nil {
		log.Error("Get on Root node resulted in error:", e)
	}
	if root == nil {
		log.Fatalf("Unable to get Root Node.")
	}

	_, e = root.GetConfig(ctx)
	if nexus_client.IsChildNotFound(e) {
		log.Info("Node not found. Will create a new Config node")
		if _, e = root.AddConfig(ctx, &cfgv1.Config{}); e != nil {
			return errors.WithMessage(e, "When creating Config Node")
		}
	} else if e != nil {
		log.Panicf("Get on Config node resulted in error:", e)
		// root.AddConfig(ctx, &cfgv1.Config{})
	}
	_, e = root.GetDesiredConfig(ctx)
	if nexus_client.IsChildNotFound(e) {
		log.Info("Node not found. Will create a new DesiredConfig node")
		if _, e = root.AddDesiredConfig(ctx, &dcfgv1.DesiredConfig{}); e != nil {
			return errors.WithMessage(e, "When creating DesiredConfig Node")
		}
	} else if e != nil {
		log.Panicf("Get on DesiredConfig node resulted in error:", e)
		// root.AddDesiredConfig(ctx, &dcfgv1.DesiredConfig{})

	}
	_, e = root.GetInventory(ctx)
	if nexus_client.IsChildNotFound(e) {
		log.Info("Node not found. Will create a new Inventory node")
		if _, e = root.AddInventory(ctx, &invv1.Inventory{}); e != nil {
			return errors.WithMessage(e, "When creating Inventory Node")
		}
	} else if e != nil {
		log.Panicf("Get on Inventory node resulted in error:", e)
		//root.AddInventory(ctx, &invv1.Inventory{})
	}

	return nil
}

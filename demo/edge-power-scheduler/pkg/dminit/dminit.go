package dminit

import (
	"context"
	cfgv1 "powerschedulermodel/build/apis/config.intel.com/v1"
	dcfgv1 "powerschedulermodel/build/apis/desiredconfig.intel.com/v1"
	invv1 "powerschedulermodel/build/apis/inventory.intel.com/v1"
	rootv1 "powerschedulermodel/build/apis/root.intel.com/v1"
	nexus_client "powerschedulermodel/build/nexus-client"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func Init(ctx context.Context, nexusClient *nexus_client.Clientset) error {
	// A common pattern is to re-use fields between logging statements by re-using
	// the logrus.Entry returned from WithFields()
	log := logrus.WithFields(logrus.Fields{
		"module": "dminit",
	})

	// create some nodes that are needed for the system to function.
	rps, e := nexusClient.GetRootPowerScheduler(ctx)
	if nexus_client.IsNotFound(e) {
		log.Info("Node not found. Will create a new RootPowerScheduler node")
		if rps, e = nexusClient.AddRootPowerScheduler(ctx, &rootv1.PowerScheduler{}); e != nil {
			return errors.WithMessage(e, "When creating RootPowerScheduler Node")
		}
	} else if e != nil {
		log.Error("Get on RootPowerScheduler node resulted in error:", e)
	}
	_, e = rps.GetConfig(ctx)
	if nexus_client.IsChildNotFound(e) {
		log.Info("Node not found. Will create a new Config node")
		if _, e = rps.AddConfig(ctx, &cfgv1.Config{}); e != nil {
			return errors.WithMessage(e, "When creating Config Node")
		}
	} else if e != nil {
		log.Error("Get on Config node resulted in error:", e)
	}
	_, e = rps.GetDesiredEdgeConfig(ctx)
	if nexus_client.IsChildNotFound(e) {
		log.Info("Node not found. Will create a new DesiredConfig node")
		if _, e = rps.AddDesiredEdgeConfig(ctx, &dcfgv1.DesiredEdgeConfig{}); e != nil {
			return errors.WithMessage(e, "When creating DesiredEdgeConfig Node")
		}
	} else if e != nil {
		log.Error("Get on DesiredConfig node resulted in error:", e)
	}
	_, e = rps.GetInventory(ctx)
	if nexus_client.IsChildNotFound(e) {
		log.Info("Node not found. Will create a new Inventory node")
		if _, e = rps.AddInventory(ctx, &invv1.Inventory{}); e != nil {
			return errors.WithMessage(e, "When creating Inventory Node")
		}
	} else if e != nil {
		log.Error("Get on Inventory node resulted in error:", e)
	}

	return nil
}

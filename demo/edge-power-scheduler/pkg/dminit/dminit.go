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

func Init(ctx context.Context, nexusClient *nexus_client.Clientset) (
	*nexus_client.RootPowerScheduler,
	*nexus_client.ConfigConfig,
	*nexus_client.DesiredconfigDesiredEdgeConfig,
	*nexus_client.InventoryInventory, error) {
	// A common pattern is to re-use fields between logging statements by re-using
	// the logrus.Entry returned from WithFields()
	log := logrus.WithFields(logrus.Fields{
		"module": "dminit",
	})

	// create some nodes that are needed for the system to function.
	rps, e := nexusClient.GetRootPowerScheduler(ctx)
	if e != nil {
		log.Info("Get on RootPowerScheduler node resulted in error:", e)
		log.Info("Will create a new RootPowerScheduler node")
		if rps, e = nexusClient.AddRootPowerScheduler(ctx, &rootv1.PowerScheduler{}); e != nil {
			return nil, nil, nil, nil, errors.WithMessage(e, "When creating RootPowerScheduler Node")
		}
	}
	cfg, e := rps.GetConfig(ctx)
	if e != nil {
		log.Info("Get on Config node resulted in error:", e)
		log.Info("Will create a new Config node")
		if cfg, e = rps.AddConfig(ctx, &cfgv1.Config{}); e != nil {
			return nil, nil, nil, nil, errors.WithMessage(e, "When creating Config Node")
		}
	}
	dcfg, e := rps.GetDesiredEdgeConfig(ctx)
	if e != nil {
		log.Info("Get on DesiredConfig node resulted in error:", e)
		log.Info("Will create a new DesiredConfig node")
		if dcfg, e = rps.AddDesiredEdgeConfig(ctx, &dcfgv1.DesiredEdgeConfig{}); e != nil {
			return nil, nil, nil, nil, errors.WithMessage(e, "When creating DesiredEdgeConfig Node")
		}
	}
	inv, e := rps.GetInventory(ctx)
	if e != nil {
		log.Info("Get on Inventory node resulted in error:", e)
		log.Info("Will create a new Inventory node")
		if inv, e = rps.AddInventory(ctx, &invv1.Inventory{}); e != nil {
			return nil, nil, nil, nil, errors.WithMessage(e, "When creating Inventory Node")
		}
	}
	return rps, cfg, dcfg, inv, nil
}

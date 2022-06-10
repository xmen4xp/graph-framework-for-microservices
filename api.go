package apis

import (
	"golang-appnet.eng.vmware.com/nexus-sdk/api/config"
	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"
)

// Api is the root node for Nexus infra/runtime datamodel.
//
// This hosts the graph that will consist of user configuration,
// runtime state, inventory and other state essential to the
// functioning of Nexus SDK and runtime.
type Api struct {
	nexus.Node

	// Configuration.
	Config config.Config `nexus:"child"`
}

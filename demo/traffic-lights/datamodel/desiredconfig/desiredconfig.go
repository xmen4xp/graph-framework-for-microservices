package desiredconfig

import (
	"trafficlight/desiredconfig/desiredconfiglight"

	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

type DesiredConfig struct {
	nexus.SingletonNode

	Lights desiredconfiglight.LightConfig `nexus:"children"`
}

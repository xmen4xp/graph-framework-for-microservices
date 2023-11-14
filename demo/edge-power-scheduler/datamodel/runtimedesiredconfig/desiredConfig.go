package runtimedesiredconfig

import (
	"powerschedulermodel/runtimedesiredconfig/edgedc"

	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus"
)

type DesiredEdgeConfig struct {
	nexus.SingletonNode
	EdgesDC edgedc.EdgeDC `nexus:"children"`
}

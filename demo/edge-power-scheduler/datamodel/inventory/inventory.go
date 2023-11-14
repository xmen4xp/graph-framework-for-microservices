package inventory

import (
	"powerschedulermodel/inventory/edge"

	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus"
)

type Inventory struct {
	nexus.SingletonNode
	Edges edge.Edge `nexus:"children"`
}

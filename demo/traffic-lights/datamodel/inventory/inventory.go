package inventory

import (
	"trafficlight/inventory/light"

	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

type Inventory struct {
	nexus.SingletonNode

	Lights light.Light `nexus:"children"`
}

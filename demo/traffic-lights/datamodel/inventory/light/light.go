package light

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

type Light struct {
	nexus.Node

	GroupNumber int
	LightNumber int
}

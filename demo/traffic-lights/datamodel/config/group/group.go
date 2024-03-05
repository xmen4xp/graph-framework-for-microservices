package group

import "github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"

type Group struct {
	nexus.Node
	Status GroupStatus `nexus:"status"`
}

type GroupStatus struct {
	TotalLights uint32
}

package edge

import (
	"powerschedulermodel/inventory/edge/powermgmt"

	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus"
)

var EdgeRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/v1alpha1/edges/{edge.Edge}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
		{
			Uri:     "/v1alpha1/edges",
			Methods: nexus.HTTPListResponse,
		},
	},
}

// nexus-rest-api-gen:EdgeRestAPISpec
// nexus-description: This is for adding a edge to inventory
// in the control plane.
type Edge struct {
	nexus.Node
	PowerInfo powermgmt.PowerInfo `nexus:"child"`
}

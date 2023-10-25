package powermgmt

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus"
)

var PIRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/v1alpha1/edges/{edge.Edge}/power-info",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
	},
}

// nexus-rest-api-gen:PIRestAPISpec
// nexus-description: This is for adding a edge power related info
type PowerInfo struct {
	nexus.SingletonNode
	TotalPowerAvailable uint32
	FreePowerAvailable  uint32
}

package edgedc

import (
	"powerschedulermodel/desiredconfig/edgedc/jobmgmt"

	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus"
)

var EdgeDCRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/v1alpha1/edge-desired-config/{edgedc.EdgeDC}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
		{
			Uri:     "/v1alpha1/edge-desired-config",
			Methods: nexus.HTTPListResponse,
		},
	},
}

// nexus-rest-api-gen:EdgeDCRestAPISpec
// nexus-description: This is for adding a desired edge config
type EdgeDC struct {
	nexus.Node
	JobsInfo jobmgmt.JobInfo `nexus:"children"`
}

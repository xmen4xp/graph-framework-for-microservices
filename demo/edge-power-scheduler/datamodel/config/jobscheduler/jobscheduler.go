package jobscheduler

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus"
)

var JOBRequesterAPI = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/v1alpha1/job-requests/{jobscheduler.Job}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
		{
			Uri:     "/v1alpha1/job-requests",
			Methods: nexus.HTTPListResponse,
		},
	},
}

// nexus-rest-api-gen:JOBRequesterAPI
// nexus-description: Job requets api's
type Job struct {
	nexus.Node
	JobId        uint64
	PowerNeeded  uint32 // watt-second
	CreationTime int64
	State        JobStatus `nexus:"status"`
}

type JobStatus struct {
	PercentCompleted uint32
	StartTime        int64
	EndTime          int64
	Execution        map[string]NodeExecutionStatus
}

type NodeExecutionStatus struct {
	PowerRequested uint32
	StartTime      int64
	EndTime        int64
	Progress       uint32
}

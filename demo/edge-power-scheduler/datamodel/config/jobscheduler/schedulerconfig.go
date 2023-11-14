package jobscheduler

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus"
)

var JOBRequesterAPI = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/v1alpha1/schedular",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
	},
}

// nex us-rest-api-gen:JOBRequesterAPI
// nexus-description: Job request api's
type SchedulerConfig struct {
	nexus.SingletonNode
	State JobSchedulerStatus `nexus:"status"`
}

type JobSchedulerStatus struct {
	TotalJobsExecuted ExectuedJobStats
	JobsQueued        uint32
	JobsRunning       uint32
	BusyEdges         uint32
	FreeEdges         uint32
	JobRate1Min       uint32
	JobRate10Sec      uint32
}

type ExectuedJobStats struct {
	TotalJobsProcessed  uint32
	TotalPowerProcessed uint32
}

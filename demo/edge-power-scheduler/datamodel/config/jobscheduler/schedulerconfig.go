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

// nexus-rest-api-gen:JOBRequesterAPI
// nexus-description: Job requets api's
type SchedulerConfig struct {
	nexus.SingletonNode
	State JobSchedulerStatus `nexus:"status"`
}

type JobSchedulerStatus struct {
	TotalJobsExecuted        ExectuedJobStats
	CurrentJobsExecuting     map[string]ExecutingJobStatus
	TotalJobsExecutedPerEdge map[string]ExectuedJobStats
}

type ExecutingJobStatus struct {
	PowerRequested uint32
	StartTime      int64
	EndTime        int64
	Progress       uint32
}

type ExectuedJobStats struct {
	TotalJobsProcessed  uint32
	TotalPowerProcessed uint32
}

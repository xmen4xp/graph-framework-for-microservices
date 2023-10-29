package config

import (
	"powerschedulermodel/config/jobscheduler"

	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus"
)

type Config struct {
	nexus.SingletonNode
	Jobs      jobscheduler.Job             `nexus:"children"`
	Scheduler jobscheduler.SchedulerConfig `nexus:"child"`
}

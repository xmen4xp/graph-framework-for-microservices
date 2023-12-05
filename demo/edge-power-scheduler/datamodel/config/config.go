package config

import (
	"powerschedulermodel/config/jobgroup"

	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus"
)

type Config struct {
	nexus.SingletonNode
	JobGroups jobgroup.Jobgroup        `nexus:"children"`
	Scheduler jobgroup.SchedulerConfig `nexus:"child"`
}

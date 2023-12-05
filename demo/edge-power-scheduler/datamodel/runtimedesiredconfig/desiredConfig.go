package runtimedesiredconfig

import (
	"powerschedulermodel/runtimedesiredconfig/sitedc"

	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus"
)

type DesiredSiteConfig struct {
	nexus.SingletonNode
	SitesDC sitedc.SiteDC `nexus:"children"`
}

package config

import (
	"nexus/admin/api/apigateway"
	tenantconfig "nexus/admin/api/config/tenant"
	"nexus/admin/api/config/user"
	"nexus/admin/api/connect"
	"nexus/admin/api/route"

	"nexus/base/nexus"
)

// Config holds the Nexus configuration.
// Configuration in Nexus is intent-driven.
type Config struct {
	nexus.Node

	// Gateway configuration.
	ApiGateway apigateway.ApiGateway `nexus:"child"`

	// API extensions configuration.
	Routes route.Route `nexus:"children"`

	// Nexus Connect configuration.
	Connect      connect.Connect     `nexus:"child"`
	Tenant       tenantconfig.Tenant `nexus:"children" json:"tenant,omitempty"`
	TenantPolicy tenantconfig.Policy `nexus:"children" json:"tenant_policy,omitempty"`

	User user.User `nexus:"children" json:"user,omitempty"`
}

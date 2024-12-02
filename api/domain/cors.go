package domain

import (
	"nexus/base/nexus"
)

// CORSConfig contains the properties of an CORS Domain configuration
// Adding it as node for user identification
type CORSConfig struct {
	nexus.Node
	// adding DomainNames as array here , in echo Domain it allows user to configure multiple origin/domains in single cors
	Origins []string `json:"origins"`
	// making customHeaders an array and making it optional
	Headers []string `json:"headers,omitempty"`
}

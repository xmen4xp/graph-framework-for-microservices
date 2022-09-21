package gns

import (
	"net/http"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/nexus"
)

var FooCustomMethodsResponses = nexus.HTTPMethodsResponses{
	http.MethodDelete: nexus.HTTPCodesResponse{
		http.StatusOK:              nexus.HTTPResponse{Description: "ok"},
		http.StatusNotFound:        nexus.HTTPResponse{Description: http.StatusText(http.StatusNotFound)},
		nexus.DefaultHTTPErrorCode: nexus.DefaultHTTPError,
	},
}

var BarCustomCodesResponses = nexus.HTTPCodesResponse{
	http.StatusBadRequest: nexus.HTTPResponse{Description: "Bad Request"},
}

var BarCustomMethodsResponses = nexus.HTTPMethodsResponses{
	http.MethodPatch: BarCustomCodesResponses,
}

var GNSRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri: "/v1alpha2/global-namespace/{gns.Gns}",
			QueryParams: []string{
				"config.Config",
			},
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
		{
			Uri: "/v1alpha2/global-namespaces",
			QueryParams: []string{
				"config.Config",
			},
			Methods: nexus.HTTPListResponse,
		},
		{
			Uri: "/test-foo",
			QueryParams: []string{
				"config.Config",
			},
			Methods: FooCustomMethodsResponses,
		},
		{
			Uri: "/test-bar",
			QueryParams: []string{
				"config.Config",
			},
			Methods: BarCustomMethodsResponses,
		},
	},
}

type Port uint16

// Host the IP of the endpoint
type Host string

type HostPort struct {
	Host Host
	Port Port
}

type Instance float32
type AliasArr []int

// Gns struct.
// nexus-description: this is my awesome node
// specification of GNS.
type Gns struct {
	nexus.SingletonNode
	//nexus-validation: MaxLength=8, MinLength=2
	//nexus-validation: Pattern=abc
	Domain           string
	UseSharedGateway bool
	Mydesc           Description
	FooLink          Bar `nexus:"link"`
	FooLinks         Bar `nexus:"links"`
	FooChild         Bar `nexus:"child"`
	FooChildren      Bar `nexus:"children"`
	HostPort         HostPort
	TestArray        EmptyData
	Instance         Instance
	Array1           float32
	Array2           []Description
	Array3           []Bar
	Array4           []Instance
	Array5           AliasArr
}

// This is Description struct.
type Description struct {
	Color     string
	Version   string
	ProjectID []string
	TestAns   []Answer
	Instance  Instance
	HostPort  HostPort
}

type Bar struct {
	nexus.Node
	Name uint64
}

type EmptyData struct {
	nexus.Node
}

type Answer struct {
	Name string
}

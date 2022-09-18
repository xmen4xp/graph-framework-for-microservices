// Code generated by nexus. DO NOT EDIT.

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/crd_generated/common"
)

// +k8s:openapi-gen=true
type Child struct {
	Group string `json:"group" yaml:"group"`
	Kind  string `json:"kind" yaml:"kind"`
	Name  string `json:"name" yaml:"name"`
}

// +k8s:openapi-gen=true
type Link struct {
	Group string `json:"group" yaml:"group"`
	Kind  string `json:"kind" yaml:"kind"`
	Name  string `json:"name" yaml:"name"`
}

// +k8s:openapi-gen=true
type NexusStatus struct {
	SourceGeneration int64 `json:"sourceGeneration" yaml:"sourceGeneration"`
	RemoteGeneration int64 `json:"remoteGeneration" yaml:"remoteGeneration"`
}

/* ------------------- CRDs definitions ------------------- */

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type Config struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              ConfigSpec        `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status            ConfigNexusStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// +k8s:openapi-gen=true
type ConfigNexusStatus struct {
	Nexus NexusStatus `json:"nexus,omitempty" yaml:"nexus,omitempty"`
}

func (c *Config) CRDName() string {
	return "configs.config.tsm.tanzu.vmware.com"
}

func (c *Config) DisplayName() string {
	if c.GetLabels() != nil {
		return c.GetLabels()[common.DISPLAY_NAME_LABEL]
	}
	return ""
}

// +k8s:openapi-gen=true
type ConfigSpec struct {
	ConfigName string            `json:"configName" yaml:"configName"`
	Cluster    Cluster           `json:"cluster" yaml:"cluster"`
	FooA       AMap              `json:"fooA" yaml:"fooA"`
	FooMap     map[string]string `json:"fooMap" yaml:"fooMap"`
	FooB       BArray            `json:"fooB" yaml:"fooB"`
	FooC       CInt              `nexus-graphql:"ignore:true"`
	FooD       DFloat            `nexus-graphql:"type:string"`
	GNSGvk     *Child            `json:"gNSGvk,omitempty" yaml:"gNSGvk,omitempty" nexus:"child"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ConfigList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []Config `json:"items" yaml:"items"`
}

// +k8s:openapi-gen=true
type Cluster struct {
	Name string
	MyID int
}

type AMap map[string]string
type BArray []string
type CInt uint8
type DFloat float32

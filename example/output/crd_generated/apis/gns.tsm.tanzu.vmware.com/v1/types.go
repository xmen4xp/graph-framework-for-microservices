// Code generated by nexus. DO NOT EDIT.

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/crd_generated/helper"
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

/* ------------------- CRDs definitions ------------------- */

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type Gns struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              GnsSpec  `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status            GnsState `json:"status,omitempty" yaml:"status,omitempty"`
}

func (c *Gns) CRDName() string {
	return "gnses.gns.tsm.tanzu.vmware.com"
}

func (c *Gns) DisplayName() string {
	if c.GetLabels() != nil {
		return c.GetLabels()[helper.DISPLAY_NAME_LABEL]
	}
	return ""
}

// +k8s:openapi-gen=true
type GnsSpec struct {
	Domain                    string           `json:"domain" yaml:"domain"`
	UseSharedGateway          bool             `json:"useSharedGateway" yaml:"useSharedGateway"`
	Description               Description      `json:"description" yaml:"description"`
	GnsServiceGroupsGvk       map[string]Child `json:"gnsServiceGroupsGvk,omitempty" yaml:"gnsServiceGroupsGvk,omitempty" nexus:"child"`
	GnsAccessControlPolicyGvk *Child           `json:"gnsAccessControlPolicyGvk,omitempty" yaml:"gnsAccessControlPolicyGvk,omitempty" nexus:"child"`
	DnsGvk                    *Link            `json:"dnsGvk,omitempty" yaml:"dnsGvk,omitempty" nexus:"link"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type GnsList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []Gns `json:"items" yaml:"items"`
}

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type Dns struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
}

func (c *Dns) CRDName() string {
	return "dnses.gns.tsm.tanzu.vmware.com"
}

func (c *Dns) DisplayName() string {
	if c.GetLabels() != nil {
		return c.GetLabels()[helper.DISPLAY_NAME_LABEL]
	}
	return ""
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DnsList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []Dns `json:"items" yaml:"items"`
}

// +k8s:openapi-gen=true
type Description struct {
	Color     string
	Version   string
	ProjectId string
}

// +k8s:openapi-gen=true
type GnsState struct {
	Working     bool
	Temperature int
}

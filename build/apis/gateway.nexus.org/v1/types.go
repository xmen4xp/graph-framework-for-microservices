// Code generated by nexus. DO NOT EDIT.

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/api.git/build/common"
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
type Gateway struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              GatewaySpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

func (c *Gateway) CRDName() string {
	return "gateways.gateway.nexus.org"
}

func (c *Gateway) DisplayName() string {
	if c.GetLabels() != nil {
		return c.GetLabels()[common.DISPLAY_NAME_LABEL]
	}
	return ""
}

// +k8s:openapi-gen=true
type GatewaySpec struct {
	Config   GatewayConfig `json:"config" yaml:"config"`
	AuthnGvk *Link         `json:"authnGvk,omitempty" yaml:"authnGvk,omitempty" nexus:"link"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type GatewayList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []Gateway `json:"items" yaml:"items"`
}

// +k8s:openapi-gen=true
type GatewayConfig struct {
}

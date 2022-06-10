// Code generated by nexus. DO NOT EDIT.

package v1

import (
	confignexusorgv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/config.nexus.org/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
type Api struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              ApiSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

func (c *Api) CRDName() string {
	return "apis.apis.nexus.org"
}

// +k8s:openapi-gen=true
type ApiSpec struct {
	Config    confignexusorgv1.Config `json:"-" yaml:"-"`
	ConfigGvk Child                   `json:"configGvk,omitempty" yaml:"configGvk,omitempty" nexus:"child"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ApiList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []Api `json:"items" yaml:"items"`
}

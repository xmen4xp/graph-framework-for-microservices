// Code generated by nexus. DO NOT EDIT.

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"/build/common"
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
type SyncerStatus struct {
	EtcdVersion    int64 `json:"etcdVersion, omitempty" yaml:"etcdVersion, omitempty"`
	CRGenerationId int64 `json:"cRGenerationId, omitempty" yaml:"cRGenerationId, omitempty"`
}

// +k8s:openapi-gen=true
type NexusStatus struct {
	SourceGeneration int64        `json:"sourceGeneration, omitempty" yaml:"sourceGeneration, omitempty"`
	RemoteGeneration int64        `json:"remoteGeneration, omitempty" yaml:"remoteGeneration, omitempty"`
	SyncerStatus     SyncerStatus `json:"syncerStatus, omitempty" yaml:"syncerStatus, omitempty"`
}

/* ------------------- CRDs definitions ------------------- */

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type Root struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              RootSpec        `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status            RootNexusStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// +k8s:openapi-gen=true
type RootNexusStatus struct {
	Nexus NexusStatus `json:"nexus,omitempty" yaml:"nexus,omitempty"`
}

func (c *Root) CRDName() string {
	return "roots.root.tsm.tanzu.vmware.com"
}

func (c *Root) DisplayName() string {
	if c.GetLabels() != nil {
		return c.GetLabels()[common.DisplayNameLabel]
	}
	return ""
}

// +k8s:openapi-gen=true
type RootSpec struct {
	ConfigGvk *Child `json:"configGvk,omitempty" yaml:"configGvk,omitempty" nexus:"child"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type RootList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []Root `json:"items" yaml:"items"`
}

// +k8s:openapi-gen=true
type NonNexusType struct {
	Test int `json:"test" yaml:"test"`
}

// +k8s:openapi-gen=true
type queryFilters struct {
	StartTime           string `json:"startTime" yaml:"startTime"`
	EndTime             string `json:"endTime" yaml:"endTime"`
	Interval            string `json:"interval" yaml:"interval"`
	IsServiceDeployment bool   `json:"isServiceDeployment" yaml:"isServiceDeployment"`
	StartVal            int    `json:"startVal" yaml:"startVal"`
}

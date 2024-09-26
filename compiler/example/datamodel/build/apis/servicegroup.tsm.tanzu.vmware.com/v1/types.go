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
type SvcGroup struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              SvcGroupSpec        `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status            SvcGroupNexusStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// +k8s:openapi-gen=true
type SvcGroupNexusStatus struct {
	Nexus NexusStatus `json:"nexus,omitempty" yaml:"nexus,omitempty"`
}

func (c *SvcGroup) CRDName() string {
	return "svcgroups.servicegroup.tsm.tanzu.vmware.com"
}

func (c *SvcGroup) DisplayName() string {
	if c.GetLabels() != nil {
		return c.GetLabels()[common.DISPLAY_NAME_LABEL]
	}
	return ""
}

// +k8s:openapi-gen=true
type SvcGroupSpec struct {
	DisplayName string `json:"displayName" yaml:"displayName"`
	Description string `json:"description" yaml:"description"`
	Color       string `json:"color" yaml:"color"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SvcGroupList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []SvcGroup `json:"items" yaml:"items"`
}

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type SvcGroupLinkInfo struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              SvcGroupLinkInfoSpec        `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status            SvcGroupLinkInfoNexusStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// +k8s:openapi-gen=true
type SvcGroupLinkInfoNexusStatus struct {
	Nexus NexusStatus `json:"nexus,omitempty" yaml:"nexus,omitempty"`
}

func (c *SvcGroupLinkInfo) CRDName() string {
	return "svcgrouplinkinfos.servicegroup.tsm.tanzu.vmware.com"
}

func (c *SvcGroupLinkInfo) DisplayName() string {
	if c.GetLabels() != nil {
		return c.GetLabels()[common.DISPLAY_NAME_LABEL]
	}
	return ""
}

// +k8s:openapi-gen=true
type SvcGroupLinkInfoSpec struct {
	ClusterName string `json:"clusterName" yaml:"clusterName"`
	DomainName  string `json:"domainName" yaml:"domainName"`
	ServiceName string `json:"serviceName" yaml:"serviceName"`
	ServiceType string `json:"serviceType" yaml:"serviceType"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SvcGroupLinkInfoList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []SvcGroupLinkInfo `json:"items" yaml:"items"`
}

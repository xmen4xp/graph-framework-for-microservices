// Code generated by nexus. DO NOT EDIT.

package v1

import (
	servicegrouptsmtanzuvmwarecomv1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/_generated/apis/servicegroup.tsm.tanzu.vmware.com/v1"
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
type AccessControlPolicy struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              AccessControlPolicySpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

func (c *AccessControlPolicy) CRDName() string {
	return "accesscontrolpolicies.policy.tsm.tanzu.vmware.com"
}

// +k8s:openapi-gen=true
type AccessControlPolicySpec struct {
	PolicyConfigs    map[string]ACPConfig `json:"-" yaml:"-"`
	PolicyConfigsGvk map[string]Child     `json:"policyConfigsGvk,omitempty" yaml:"policyConfigsGvk,omitempty" nexus:"child"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AccessControlPolicyList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []AccessControlPolicy `json:"items" yaml:"items"`
}

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type ACPConfig struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              ACPConfigSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

func (c *ACPConfig) CRDName() string {
	return "acpconfigs.policy.tsm.tanzu.vmware.com"
}

// +k8s:openapi-gen=true
type ACPConfigSpec struct {
	DisplayName        string                                              `json:"displayName" yaml:"displayName"`
	Gns                string                                              `json:"gns" yaml:"gns"`
	Description        string                                              `json:"description" yaml:"description"`
	Tags               []string                                            `json:"tags" yaml:"tags"`
	ProjectId          string                                              `json:"projectId" yaml:"projectId"`
	Conditions         []string                                            `json:"conditions" yaml:"conditions"`
	DestSvcGroups      map[string]servicegrouptsmtanzuvmwarecomv1.SvcGroup `json:"-" yaml:"-"`
	DestSvcGroupsGvk   map[string]Link                                     `json:"destSvcGroupsGvk,omitempty" yaml:"destSvcGroupsGvk,omitempty" nexus:"link"`
	SourceSvcGroups    map[string]servicegrouptsmtanzuvmwarecomv1.SvcGroup `json:"-" yaml:"-"`
	SourceSvcGroupsGvk map[string]Link                                     `json:"sourceSvcGroupsGvk,omitempty" yaml:"sourceSvcGroupsGvk,omitempty" nexus:"link"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ACPConfigList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []ACPConfig `json:"items" yaml:"items"`
}

// +k8s:openapi-gen=true
type ResourceGroupRef struct {
	Name string
	Type string
}

// +k8s:openapi-gen=true
type ACPSvcGroupLinkInfo struct {
	ServiceName string
	ServiceType string
}

// +k8s:openapi-gen=true
type PolicyCfgAction struct {
	Action PolicyActionType `json:"action" mapstructure:"action"`
}

// +k8s:openapi-gen=true
type ResourceGroupID struct {
	Name string `json:"name" mapstruction:"name"`
	Type string `json:"type" mapstruction:"type"`
}

type PolicyActionType string
type PolicyCfgActions []PolicyCfgAction
type ResourceGroupIDs []ResourceGroupID

// Code generated by nexus. DO NOT EDIT.

package v1

import (
	cartv1 "github.com/vmware-tanzu/cartographer/pkg/apis/v1alpha1"
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/generated/common"
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
type Foo struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              FooSpec        `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status            FooNexusStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// +k8s:openapi-gen=true
type FooNexusStatus struct {
	Nexus NexusStatus `json:"nexus,omitempty" yaml:"nexus,omitempty"`
}

func (c *Foo) CRDName() string {
	return "foos.gns.tsm.tanzu.vmware.com"
}

func (c *Foo) DisplayName() string {
	if c.GetLabels() != nil {
		return c.GetLabels()[common.DISPLAY_NAME_LABEL]
	}
	return ""
}

// +k8s:openapi-gen=true
type FooSpec struct {
	Password string `json:"password" yaml:"password"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type FooList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []Foo `json:"items" yaml:"items"`
}

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type Gns struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              GnsSpec        `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status            GnsNexusStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// +k8s:openapi-gen=true
type GnsNexusStatus struct {
	State GnsState    `json:"state,omitempty" yaml:"state,omitempty"`
	Nexus NexusStatus `json:"nexus,omitempty" yaml:"nexus,omitempty"`
}

func (c *Gns) CRDName() string {
	return "gnses.gns.tsm.tanzu.vmware.com"
}

func (c *Gns) DisplayName() string {
	if c.GetLabels() != nil {
		return c.GetLabels()[common.DISPLAY_NAME_LABEL]
	}
	return ""
}

// +k8s:openapi-gen=true
type GnsSpec struct {
	//nexus-validation: MaxLength=8, MinLength=2
	//nexus-validation: Pattern=abc
	Domain                    string                       `json:"domain" yaml:"domain"`
	UseSharedGateway          bool                         `json:"useSharedGateway" yaml:"useSharedGateway"`
	Annotations               nexus.NexusGenericObject     `nexus-graphql-jsonencoded:""`
	TargetPort                intstr.IntOrString           `json:"targetPort,omitempty" mapstructure:"targetPort,omitempty"`
	Description               Description                  `json:"description" yaml:"description"`
	Meta                      string                       `json:"meta" yaml:"meta"`
	IntOrString               []intstr.IntOrString         `nexus-graphql-type-name:"IntOrString" json:"intOrString,omitempty" mapstructure:"intOrString,omitempty"`
	Port                      *int                         `json:"port" yaml:"port"`
	OtherDescription          *Description                 `json:"otherDescription" yaml:"otherDescription"`
	MapPointer                *map[string]string           `json:"mapPointer" yaml:"mapPointer"`
	SlicePointer              *[]string                    `json:"slicePointer" yaml:"slicePointer"`
	WorkloadSpec              cartv1.WorkloadSpec          `json:"workloadSpec" yaml:"workloadSpec"`
	DifferentSpec             *cartv1.WorkloadSpec         `json:"differentSpec" yaml:"differentSpec"`
	ServiceSegmentRef         ServiceSegmentRef            `json:"serviceSegmentRef,omitempty"`
	ServiceSegmentRefPointer  *ServiceSegmentRef           `json:"serviceSegmentRefPointer,omitempty"`
	ServiceSegmentRefs        []ServiceSegmentRef          `json:"serviceSegmentRefs,omitempty"`
	ServiceSegmentRefMap      map[string]ServiceSegmentRef `json:"serviceSegmentRefMap,omitempty"`
	GnsServiceGroupsGvk       map[string]Child             `json:"gnsServiceGroupsGvk,omitempty" yaml:"gnsServiceGroupsGvk,omitempty" nexus:"children"`
	GnsAccessControlPolicyGvk *Child                       `json:"gnsAccessControlPolicyGvk,omitempty" yaml:"gnsAccessControlPolicyGvk,omitempty" nexus:"child"`
	FooChildGvk               *Child                       `json:"fooChildGvk,omitempty" yaml:"fooChildGvk,omitempty" nexus:"child"`
	IgnoreChildGvk            *Child                       `json:"ignoreChildGvk,omitempty" yaml:"ignoreChildGvk,omitempty" nexus:"child"`
	FooGvk                    *Child                       `json:"fooGvk,omitempty" yaml:"fooGvk,omitempty" nexus:"child"`
	DnsGvk                    *Link                        `json:"dnsGvk,omitempty" yaml:"dnsGvk,omitempty" nexus:"link"`
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
type BarChild struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              BarChildSpec        `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status            BarChildNexusStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// +k8s:openapi-gen=true
type BarChildNexusStatus struct {
	Nexus NexusStatus `json:"nexus,omitempty" yaml:"nexus,omitempty"`
}

func (c *BarChild) CRDName() string {
	return "barchilds.gns.tsm.tanzu.vmware.com"
}

func (c *BarChild) DisplayName() string {
	if c.GetLabels() != nil {
		return c.GetLabels()[common.DISPLAY_NAME_LABEL]
	}
	return ""
}

// +k8s:openapi-gen=true
type BarChildSpec struct {
	Name string `json:"name" yaml:"name"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type BarChildList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []BarChild `json:"items" yaml:"items"`
}

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type IgnoreChild struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              IgnoreChildSpec        `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status            IgnoreChildNexusStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// +k8s:openapi-gen=true
type IgnoreChildNexusStatus struct {
	Nexus NexusStatus `json:"nexus,omitempty" yaml:"nexus,omitempty"`
}

func (c *IgnoreChild) CRDName() string {
	return "ignorechilds.gns.tsm.tanzu.vmware.com"
}

func (c *IgnoreChild) DisplayName() string {
	if c.GetLabels() != nil {
		return c.GetLabels()[common.DISPLAY_NAME_LABEL]
	}
	return ""
}

// +k8s:openapi-gen=true
type IgnoreChildSpec struct {
	Name string `json:"name" yaml:"name"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type IgnoreChildList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []IgnoreChild `json:"items" yaml:"items"`
}

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type Dns struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`

	Status DnsNexusStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// +k8s:openapi-gen=true
type DnsNexusStatus struct {
	Nexus NexusStatus `json:"nexus,omitempty" yaml:"nexus,omitempty"`
}

func (c *Dns) CRDName() string {
	return "dnses.gns.tsm.tanzu.vmware.com"
}

func (c *Dns) DisplayName() string {
	if c.GetLabels() != nil {
		return c.GetLabels()[common.DISPLAY_NAME_LABEL]
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
type RandomDescription struct {
	DiscriptionA string `json:"discriptionA" yaml:"discriptionA"`
	DiscriptionB string `json:"discriptionB" yaml:"discriptionB"`
	DiscriptionC string `json:"discriptionC" yaml:"discriptionC"`
	DiscriptionD string `json:"discriptionD" yaml:"discriptionD"`
}

// +k8s:openapi-gen=true
type RandomStatus struct {
	StatusX int `json:"statusX" yaml:"statusX"`
	StatusY int `json:"statusY" yaml:"statusY"`
}

// +k8s:openapi-gen=true
type HostPort struct {
	Host Host `json:"host" yaml:"host"`
	Port Port `json:"port" yaml:"port"`
}

// +k8s:openapi-gen=true
type ReplicationSource struct {
	Kind SourceKind `json:"kind" yaml:"kind"`
}

// +k8s:openapi-gen=true
type gnsQueryFilters struct {
	StartTime           string `json:"startTime" yaml:"startTime"`
	EndTime             string `json:"endTime" yaml:"endTime"`
	Interval            string `json:"interval" yaml:"interval"`
	IsServiceDeployment bool   `json:"isServiceDeployment" yaml:"isServiceDeployment"`
	StartVal            int    `json:"startVal" yaml:"startVal"`
}

// +k8s:openapi-gen=true
type metricsFilers struct {
	StartTime    string `json:"startTime" yaml:"startTime"`
	EndTime      string `json:"endTime" yaml:"endTime"`
	TimeInterval string `json:"timeInterval" yaml:"timeInterval"`
	SomeUserArg1 string `json:"someUserArg1" yaml:"someUserArg1"`
	SomeUserArg2 int    `json:"someUserArg2" yaml:"someUserArg2"`
	SomeUserArg3 bool   `json:"someUserArg3" yaml:"someUserArg3"`
}

// +k8s:openapi-gen=true
type ServiceSegmentRef struct {
	Field1 string `json:"field1" yaml:"field1"`
	Field2 string `json:"field2" yaml:"field2"`
}

// +k8s:openapi-gen=true
type Description struct {
	Color     string   `json:"color" yaml:"color"`
	Version   string   `json:"version" yaml:"version"`
	ProjectId string   `json:"projectId" yaml:"projectId"`
	TestAns   []Answer `json:"testAns" yaml:"testAns"`
	Instance  Instance `json:"instance" yaml:"instance"`
	HostPort  HostPort `json:"hostPort" yaml:"hostPort"`
}

// +k8s:openapi-gen=true
type Answer struct {
	Name string `json:"name" yaml:"name"`
}

// +k8s:openapi-gen=true
type GnsState struct {
	Working     bool `json:"working" yaml:"working"`
	Temperature int  `json:"temperature" yaml:"temperature"`
}

// +k8s:openapi-gen=true
type AdditionalDescription struct {
	DiscriptionA string `json:"discriptionA" yaml:"discriptionA"`
	DiscriptionB string `json:"discriptionB" yaml:"discriptionB"`
	DiscriptionC string `json:"discriptionC" yaml:"discriptionC"`
	DiscriptionD string `json:"discriptionD" yaml:"discriptionD"`
}

// +k8s:openapi-gen=true
type AdditionalStatus struct {
	StatusX int `json:"statusX" yaml:"statusX"`
	StatusY int `json:"statusY" yaml:"statusY"`
}

type RandomConst1 string
type RandomConst2 string
type RandomConst3 string
type MyConst string
type SourceKind string
type Port uint16
type Host string
type Instance float32
type AliasArr []int
type MyStr string
type TempConst1 string
type TempConst2 string
type TempConst3 string

const (
	MyConst3 RandomConst3 = "Const3"
	MyConst2 RandomConst2 = "Const2"
	MyConst1 RandomConst1 = "Const1"
	Object   SourceKind   = "Object"
	Type     SourceKind   = "Type"
	XYZ      MyConst      = "xyz"
	Const3   TempConst3   = "Const3"
	Const2   TempConst2   = "Const2"
	Const1   TempConst1   = "Const1"
)

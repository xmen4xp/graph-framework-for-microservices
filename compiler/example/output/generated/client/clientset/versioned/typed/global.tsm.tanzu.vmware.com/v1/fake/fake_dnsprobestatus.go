/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"
	globaltsmtanzuvmwarecomv1 "nexustempmodule/apis/global.tsm.tanzu.vmware.com/v1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeDNSProbeStatuses implements DNSProbeStatusInterface
type FakeDNSProbeStatuses struct {
	Fake *FakeGlobalTsmV1
}

var dnsprobestatusesResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "dnsprobestatuses"}

var dnsprobestatusesKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "DNSProbeStatus"}

// Get takes name of the dNSProbeStatus, and returns the corresponding dNSProbeStatus object, and an error if there is any.
func (c *FakeDNSProbeStatuses) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.DNSProbeStatus, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(dnsprobestatusesResource, name), &globaltsmtanzuvmwarecomv1.DNSProbeStatus{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.DNSProbeStatus), err
}

// List takes label and field selectors, and returns the list of DNSProbeStatuses that match those selectors.
func (c *FakeDNSProbeStatuses) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.DNSProbeStatusList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(dnsprobestatusesResource, dnsprobestatusesKind, opts), &globaltsmtanzuvmwarecomv1.DNSProbeStatusList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.DNSProbeStatusList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.DNSProbeStatusList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.DNSProbeStatusList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested dNSProbeStatuses.
func (c *FakeDNSProbeStatuses) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(dnsprobestatusesResource, opts))
}

// Create takes the representation of a dNSProbeStatus and creates it.  Returns the server's representation of the dNSProbeStatus, and an error, if there is any.
func (c *FakeDNSProbeStatuses) Create(ctx context.Context, dNSProbeStatus *globaltsmtanzuvmwarecomv1.DNSProbeStatus, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.DNSProbeStatus, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(dnsprobestatusesResource, dNSProbeStatus), &globaltsmtanzuvmwarecomv1.DNSProbeStatus{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.DNSProbeStatus), err
}

// Update takes the representation of a dNSProbeStatus and updates it. Returns the server's representation of the dNSProbeStatus, and an error, if there is any.
func (c *FakeDNSProbeStatuses) Update(ctx context.Context, dNSProbeStatus *globaltsmtanzuvmwarecomv1.DNSProbeStatus, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.DNSProbeStatus, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(dnsprobestatusesResource, dNSProbeStatus), &globaltsmtanzuvmwarecomv1.DNSProbeStatus{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.DNSProbeStatus), err
}

// Delete takes name of the dNSProbeStatus and deletes it. Returns an error if one occurs.
func (c *FakeDNSProbeStatuses) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(dnsprobestatusesResource, name), &globaltsmtanzuvmwarecomv1.DNSProbeStatus{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeDNSProbeStatuses) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(dnsprobestatusesResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.DNSProbeStatusList{})
	return err
}

// Patch applies the patch and returns the patched dNSProbeStatus.
func (c *FakeDNSProbeStatuses) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.DNSProbeStatus, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(dnsprobestatusesResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.DNSProbeStatus{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.DNSProbeStatus), err
}

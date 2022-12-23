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

// FakeExternalDNSInventoryHealthChecks implements ExternalDNSInventoryHealthCheckInterface
type FakeExternalDNSInventoryHealthChecks struct {
	Fake *FakeGlobalTsmV1
}

var externaldnsinventoryhealthchecksResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "externaldnsinventoryhealthchecks"}

var externaldnsinventoryhealthchecksKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "ExternalDNSInventoryHealthCheck"}

// Get takes name of the externalDNSInventoryHealthCheck, and returns the corresponding externalDNSInventoryHealthCheck object, and an error if there is any.
func (c *FakeExternalDNSInventoryHealthChecks) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.ExternalDNSInventoryHealthCheck, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(externaldnsinventoryhealthchecksResource, name), &globaltsmtanzuvmwarecomv1.ExternalDNSInventoryHealthCheck{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ExternalDNSInventoryHealthCheck), err
}

// List takes label and field selectors, and returns the list of ExternalDNSInventoryHealthChecks that match those selectors.
func (c *FakeExternalDNSInventoryHealthChecks) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.ExternalDNSInventoryHealthCheckList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(externaldnsinventoryhealthchecksResource, externaldnsinventoryhealthchecksKind, opts), &globaltsmtanzuvmwarecomv1.ExternalDNSInventoryHealthCheckList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.ExternalDNSInventoryHealthCheckList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.ExternalDNSInventoryHealthCheckList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.ExternalDNSInventoryHealthCheckList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested externalDNSInventoryHealthChecks.
func (c *FakeExternalDNSInventoryHealthChecks) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(externaldnsinventoryhealthchecksResource, opts))
}

// Create takes the representation of a externalDNSInventoryHealthCheck and creates it.  Returns the server's representation of the externalDNSInventoryHealthCheck, and an error, if there is any.
func (c *FakeExternalDNSInventoryHealthChecks) Create(ctx context.Context, externalDNSInventoryHealthCheck *globaltsmtanzuvmwarecomv1.ExternalDNSInventoryHealthCheck, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.ExternalDNSInventoryHealthCheck, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(externaldnsinventoryhealthchecksResource, externalDNSInventoryHealthCheck), &globaltsmtanzuvmwarecomv1.ExternalDNSInventoryHealthCheck{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ExternalDNSInventoryHealthCheck), err
}

// Update takes the representation of a externalDNSInventoryHealthCheck and updates it. Returns the server's representation of the externalDNSInventoryHealthCheck, and an error, if there is any.
func (c *FakeExternalDNSInventoryHealthChecks) Update(ctx context.Context, externalDNSInventoryHealthCheck *globaltsmtanzuvmwarecomv1.ExternalDNSInventoryHealthCheck, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.ExternalDNSInventoryHealthCheck, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(externaldnsinventoryhealthchecksResource, externalDNSInventoryHealthCheck), &globaltsmtanzuvmwarecomv1.ExternalDNSInventoryHealthCheck{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ExternalDNSInventoryHealthCheck), err
}

// Delete takes name of the externalDNSInventoryHealthCheck and deletes it. Returns an error if one occurs.
func (c *FakeExternalDNSInventoryHealthChecks) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(externaldnsinventoryhealthchecksResource, name), &globaltsmtanzuvmwarecomv1.ExternalDNSInventoryHealthCheck{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeExternalDNSInventoryHealthChecks) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(externaldnsinventoryhealthchecksResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.ExternalDNSInventoryHealthCheckList{})
	return err
}

// Patch applies the patch and returns the patched externalDNSInventoryHealthCheck.
func (c *FakeExternalDNSInventoryHealthChecks) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.ExternalDNSInventoryHealthCheck, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(externaldnsinventoryhealthchecksResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.ExternalDNSInventoryHealthCheck{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ExternalDNSInventoryHealthCheck), err
}

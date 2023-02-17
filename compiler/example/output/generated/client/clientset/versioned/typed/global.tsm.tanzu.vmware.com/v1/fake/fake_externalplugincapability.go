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

// FakeExternalPluginCapabilities implements ExternalPluginCapabilityInterface
type FakeExternalPluginCapabilities struct {
	Fake *FakeGlobalTsmV1
}

var externalplugincapabilitiesResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "externalplugincapabilities"}

var externalplugincapabilitiesKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "ExternalPluginCapability"}

// Get takes name of the externalPluginCapability, and returns the corresponding externalPluginCapability object, and an error if there is any.
func (c *FakeExternalPluginCapabilities) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.ExternalPluginCapability, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(externalplugincapabilitiesResource, name), &globaltsmtanzuvmwarecomv1.ExternalPluginCapability{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ExternalPluginCapability), err
}

// List takes label and field selectors, and returns the list of ExternalPluginCapabilities that match those selectors.
func (c *FakeExternalPluginCapabilities) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.ExternalPluginCapabilityList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(externalplugincapabilitiesResource, externalplugincapabilitiesKind, opts), &globaltsmtanzuvmwarecomv1.ExternalPluginCapabilityList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.ExternalPluginCapabilityList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.ExternalPluginCapabilityList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.ExternalPluginCapabilityList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested externalPluginCapabilities.
func (c *FakeExternalPluginCapabilities) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(externalplugincapabilitiesResource, opts))
}

// Create takes the representation of a externalPluginCapability and creates it.  Returns the server's representation of the externalPluginCapability, and an error, if there is any.
func (c *FakeExternalPluginCapabilities) Create(ctx context.Context, externalPluginCapability *globaltsmtanzuvmwarecomv1.ExternalPluginCapability, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.ExternalPluginCapability, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(externalplugincapabilitiesResource, externalPluginCapability), &globaltsmtanzuvmwarecomv1.ExternalPluginCapability{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ExternalPluginCapability), err
}

// Update takes the representation of a externalPluginCapability and updates it. Returns the server's representation of the externalPluginCapability, and an error, if there is any.
func (c *FakeExternalPluginCapabilities) Update(ctx context.Context, externalPluginCapability *globaltsmtanzuvmwarecomv1.ExternalPluginCapability, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.ExternalPluginCapability, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(externalplugincapabilitiesResource, externalPluginCapability), &globaltsmtanzuvmwarecomv1.ExternalPluginCapability{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ExternalPluginCapability), err
}

// Delete takes name of the externalPluginCapability and deletes it. Returns an error if one occurs.
func (c *FakeExternalPluginCapabilities) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(externalplugincapabilitiesResource, name), &globaltsmtanzuvmwarecomv1.ExternalPluginCapability{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeExternalPluginCapabilities) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(externalplugincapabilitiesResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.ExternalPluginCapabilityList{})
	return err
}

// Patch applies the patch and returns the patched externalPluginCapability.
func (c *FakeExternalPluginCapabilities) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.ExternalPluginCapability, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(externalplugincapabilitiesResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.ExternalPluginCapability{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ExternalPluginCapability), err
}
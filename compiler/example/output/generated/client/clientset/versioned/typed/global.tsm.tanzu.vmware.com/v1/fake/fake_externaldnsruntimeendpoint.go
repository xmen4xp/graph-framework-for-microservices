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

// FakeExternalDNSRuntimeEndpoints implements ExternalDNSRuntimeEndpointInterface
type FakeExternalDNSRuntimeEndpoints struct {
	Fake *FakeGlobalTsmV1
}

var externaldnsruntimeendpointsResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "externaldnsruntimeendpoints"}

var externaldnsruntimeendpointsKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "ExternalDNSRuntimeEndpoint"}

// Get takes name of the externalDNSRuntimeEndpoint, and returns the corresponding externalDNSRuntimeEndpoint object, and an error if there is any.
func (c *FakeExternalDNSRuntimeEndpoints) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.ExternalDNSRuntimeEndpoint, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(externaldnsruntimeendpointsResource, name), &globaltsmtanzuvmwarecomv1.ExternalDNSRuntimeEndpoint{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ExternalDNSRuntimeEndpoint), err
}

// List takes label and field selectors, and returns the list of ExternalDNSRuntimeEndpoints that match those selectors.
func (c *FakeExternalDNSRuntimeEndpoints) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.ExternalDNSRuntimeEndpointList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(externaldnsruntimeendpointsResource, externaldnsruntimeendpointsKind, opts), &globaltsmtanzuvmwarecomv1.ExternalDNSRuntimeEndpointList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.ExternalDNSRuntimeEndpointList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.ExternalDNSRuntimeEndpointList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.ExternalDNSRuntimeEndpointList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested externalDNSRuntimeEndpoints.
func (c *FakeExternalDNSRuntimeEndpoints) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(externaldnsruntimeendpointsResource, opts))
}

// Create takes the representation of a externalDNSRuntimeEndpoint and creates it.  Returns the server's representation of the externalDNSRuntimeEndpoint, and an error, if there is any.
func (c *FakeExternalDNSRuntimeEndpoints) Create(ctx context.Context, externalDNSRuntimeEndpoint *globaltsmtanzuvmwarecomv1.ExternalDNSRuntimeEndpoint, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.ExternalDNSRuntimeEndpoint, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(externaldnsruntimeendpointsResource, externalDNSRuntimeEndpoint), &globaltsmtanzuvmwarecomv1.ExternalDNSRuntimeEndpoint{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ExternalDNSRuntimeEndpoint), err
}

// Update takes the representation of a externalDNSRuntimeEndpoint and updates it. Returns the server's representation of the externalDNSRuntimeEndpoint, and an error, if there is any.
func (c *FakeExternalDNSRuntimeEndpoints) Update(ctx context.Context, externalDNSRuntimeEndpoint *globaltsmtanzuvmwarecomv1.ExternalDNSRuntimeEndpoint, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.ExternalDNSRuntimeEndpoint, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(externaldnsruntimeendpointsResource, externalDNSRuntimeEndpoint), &globaltsmtanzuvmwarecomv1.ExternalDNSRuntimeEndpoint{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ExternalDNSRuntimeEndpoint), err
}

// Delete takes name of the externalDNSRuntimeEndpoint and deletes it. Returns an error if one occurs.
func (c *FakeExternalDNSRuntimeEndpoints) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(externaldnsruntimeendpointsResource, name), &globaltsmtanzuvmwarecomv1.ExternalDNSRuntimeEndpoint{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeExternalDNSRuntimeEndpoints) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(externaldnsruntimeendpointsResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.ExternalDNSRuntimeEndpointList{})
	return err
}

// Patch applies the patch and returns the patched externalDNSRuntimeEndpoint.
func (c *FakeExternalDNSRuntimeEndpoints) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.ExternalDNSRuntimeEndpoint, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(externaldnsruntimeendpointsResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.ExternalDNSRuntimeEndpoint{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ExternalDNSRuntimeEndpoint), err
}

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

// FakeOutboundAuthenticationModes implements OutboundAuthenticationModeInterface
type FakeOutboundAuthenticationModes struct {
	Fake *FakeGlobalTsmV1
}

var outboundauthenticationmodesResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "outboundauthenticationmodes"}

var outboundauthenticationmodesKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "OutboundAuthenticationMode"}

// Get takes name of the outboundAuthenticationMode, and returns the corresponding outboundAuthenticationMode object, and an error if there is any.
func (c *FakeOutboundAuthenticationModes) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.OutboundAuthenticationMode, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(outboundauthenticationmodesResource, name), &globaltsmtanzuvmwarecomv1.OutboundAuthenticationMode{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.OutboundAuthenticationMode), err
}

// List takes label and field selectors, and returns the list of OutboundAuthenticationModes that match those selectors.
func (c *FakeOutboundAuthenticationModes) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.OutboundAuthenticationModeList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(outboundauthenticationmodesResource, outboundauthenticationmodesKind, opts), &globaltsmtanzuvmwarecomv1.OutboundAuthenticationModeList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.OutboundAuthenticationModeList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.OutboundAuthenticationModeList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.OutboundAuthenticationModeList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested outboundAuthenticationModes.
func (c *FakeOutboundAuthenticationModes) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(outboundauthenticationmodesResource, opts))
}

// Create takes the representation of a outboundAuthenticationMode and creates it.  Returns the server's representation of the outboundAuthenticationMode, and an error, if there is any.
func (c *FakeOutboundAuthenticationModes) Create(ctx context.Context, outboundAuthenticationMode *globaltsmtanzuvmwarecomv1.OutboundAuthenticationMode, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.OutboundAuthenticationMode, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(outboundauthenticationmodesResource, outboundAuthenticationMode), &globaltsmtanzuvmwarecomv1.OutboundAuthenticationMode{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.OutboundAuthenticationMode), err
}

// Update takes the representation of a outboundAuthenticationMode and updates it. Returns the server's representation of the outboundAuthenticationMode, and an error, if there is any.
func (c *FakeOutboundAuthenticationModes) Update(ctx context.Context, outboundAuthenticationMode *globaltsmtanzuvmwarecomv1.OutboundAuthenticationMode, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.OutboundAuthenticationMode, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(outboundauthenticationmodesResource, outboundAuthenticationMode), &globaltsmtanzuvmwarecomv1.OutboundAuthenticationMode{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.OutboundAuthenticationMode), err
}

// Delete takes name of the outboundAuthenticationMode and deletes it. Returns an error if one occurs.
func (c *FakeOutboundAuthenticationModes) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(outboundauthenticationmodesResource, name), &globaltsmtanzuvmwarecomv1.OutboundAuthenticationMode{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeOutboundAuthenticationModes) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(outboundauthenticationmodesResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.OutboundAuthenticationModeList{})
	return err
}

// Patch applies the patch and returns the patched outboundAuthenticationMode.
func (c *FakeOutboundAuthenticationModes) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.OutboundAuthenticationMode, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(outboundauthenticationmodesResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.OutboundAuthenticationMode{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.OutboundAuthenticationMode), err
}
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

// FakeGlobalNses implements GlobalNsInterface
type FakeGlobalNses struct {
	Fake *FakeGlobalTsmV1
}

var globalnsesResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "globalnses"}

var globalnsesKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "GlobalNs"}

// Get takes name of the globalNs, and returns the corresponding globalNs object, and an error if there is any.
func (c *FakeGlobalNses) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.GlobalNs, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(globalnsesResource, name), &globaltsmtanzuvmwarecomv1.GlobalNs{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.GlobalNs), err
}

// List takes label and field selectors, and returns the list of GlobalNses that match those selectors.
func (c *FakeGlobalNses) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.GlobalNsList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(globalnsesResource, globalnsesKind, opts), &globaltsmtanzuvmwarecomv1.GlobalNsList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.GlobalNsList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.GlobalNsList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.GlobalNsList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested globalNses.
func (c *FakeGlobalNses) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(globalnsesResource, opts))
}

// Create takes the representation of a globalNs and creates it.  Returns the server's representation of the globalNs, and an error, if there is any.
func (c *FakeGlobalNses) Create(ctx context.Context, globalNs *globaltsmtanzuvmwarecomv1.GlobalNs, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.GlobalNs, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(globalnsesResource, globalNs), &globaltsmtanzuvmwarecomv1.GlobalNs{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.GlobalNs), err
}

// Update takes the representation of a globalNs and updates it. Returns the server's representation of the globalNs, and an error, if there is any.
func (c *FakeGlobalNses) Update(ctx context.Context, globalNs *globaltsmtanzuvmwarecomv1.GlobalNs, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.GlobalNs, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(globalnsesResource, globalNs), &globaltsmtanzuvmwarecomv1.GlobalNs{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.GlobalNs), err
}

// Delete takes name of the globalNs and deletes it. Returns an error if one occurs.
func (c *FakeGlobalNses) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(globalnsesResource, name), &globaltsmtanzuvmwarecomv1.GlobalNs{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeGlobalNses) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(globalnsesResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.GlobalNsList{})
	return err
}

// Patch applies the patch and returns the patched globalNs.
func (c *FakeGlobalNses) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.GlobalNs, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(globalnsesResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.GlobalNs{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.GlobalNs), err
}

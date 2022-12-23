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

// FakeKnativeIngresses implements KnativeIngressInterface
type FakeKnativeIngresses struct {
	Fake *FakeGlobalTsmV1
}

var knativeingressesResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "knativeingresses"}

var knativeingressesKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "KnativeIngress"}

// Get takes name of the knativeIngress, and returns the corresponding knativeIngress object, and an error if there is any.
func (c *FakeKnativeIngresses) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.KnativeIngress, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(knativeingressesResource, name), &globaltsmtanzuvmwarecomv1.KnativeIngress{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.KnativeIngress), err
}

// List takes label and field selectors, and returns the list of KnativeIngresses that match those selectors.
func (c *FakeKnativeIngresses) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.KnativeIngressList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(knativeingressesResource, knativeingressesKind, opts), &globaltsmtanzuvmwarecomv1.KnativeIngressList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.KnativeIngressList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.KnativeIngressList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.KnativeIngressList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested knativeIngresses.
func (c *FakeKnativeIngresses) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(knativeingressesResource, opts))
}

// Create takes the representation of a knativeIngress and creates it.  Returns the server's representation of the knativeIngress, and an error, if there is any.
func (c *FakeKnativeIngresses) Create(ctx context.Context, knativeIngress *globaltsmtanzuvmwarecomv1.KnativeIngress, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.KnativeIngress, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(knativeingressesResource, knativeIngress), &globaltsmtanzuvmwarecomv1.KnativeIngress{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.KnativeIngress), err
}

// Update takes the representation of a knativeIngress and updates it. Returns the server's representation of the knativeIngress, and an error, if there is any.
func (c *FakeKnativeIngresses) Update(ctx context.Context, knativeIngress *globaltsmtanzuvmwarecomv1.KnativeIngress, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.KnativeIngress, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(knativeingressesResource, knativeIngress), &globaltsmtanzuvmwarecomv1.KnativeIngress{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.KnativeIngress), err
}

// Delete takes name of the knativeIngress and deletes it. Returns an error if one occurs.
func (c *FakeKnativeIngresses) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(knativeingressesResource, name), &globaltsmtanzuvmwarecomv1.KnativeIngress{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeKnativeIngresses) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(knativeingressesResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.KnativeIngressList{})
	return err
}

// Patch applies the patch and returns the patched knativeIngress.
func (c *FakeKnativeIngresses) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.KnativeIngress, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(knativeingressesResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.KnativeIngress{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.KnativeIngress), err
}

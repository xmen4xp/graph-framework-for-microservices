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

	configtsmtanzuvmwarecomv1 "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/generated/apis/config.tsm.tanzu.vmware.com/v1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeFooTypeABCs implements FooTypeABCInterface
type FakeFooTypeABCs struct {
	Fake *FakeConfigTsmV1
}

var footypeabcsResource = schema.GroupVersionResource{Group: "config.tsm.tanzu.vmware.com", Version: "v1", Resource: "footypeabcs"}

var footypeabcsKind = schema.GroupVersionKind{Group: "config.tsm.tanzu.vmware.com", Version: "v1", Kind: "FooTypeABC"}

// Get takes name of the fooTypeABC, and returns the corresponding fooTypeABC object, and an error if there is any.
func (c *FakeFooTypeABCs) Get(ctx context.Context, name string, options v1.GetOptions) (result *configtsmtanzuvmwarecomv1.FooTypeABC, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(footypeabcsResource, name), &configtsmtanzuvmwarecomv1.FooTypeABC{})
	if obj == nil {
		return nil, err
	}
	return obj.(*configtsmtanzuvmwarecomv1.FooTypeABC), err
}

// List takes label and field selectors, and returns the list of FooTypeABCs that match those selectors.
func (c *FakeFooTypeABCs) List(ctx context.Context, opts v1.ListOptions) (result *configtsmtanzuvmwarecomv1.FooTypeABCList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(footypeabcsResource, footypeabcsKind, opts), &configtsmtanzuvmwarecomv1.FooTypeABCList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &configtsmtanzuvmwarecomv1.FooTypeABCList{ListMeta: obj.(*configtsmtanzuvmwarecomv1.FooTypeABCList).ListMeta}
	for _, item := range obj.(*configtsmtanzuvmwarecomv1.FooTypeABCList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested fooTypeABCs.
func (c *FakeFooTypeABCs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(footypeabcsResource, opts))
}

// Create takes the representation of a fooTypeABC and creates it.  Returns the server's representation of the fooTypeABC, and an error, if there is any.
func (c *FakeFooTypeABCs) Create(ctx context.Context, fooTypeABC *configtsmtanzuvmwarecomv1.FooTypeABC, opts v1.CreateOptions) (result *configtsmtanzuvmwarecomv1.FooTypeABC, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(footypeabcsResource, fooTypeABC), &configtsmtanzuvmwarecomv1.FooTypeABC{})
	if obj == nil {
		return nil, err
	}
	return obj.(*configtsmtanzuvmwarecomv1.FooTypeABC), err
}

// Update takes the representation of a fooTypeABC and updates it. Returns the server's representation of the fooTypeABC, and an error, if there is any.
func (c *FakeFooTypeABCs) Update(ctx context.Context, fooTypeABC *configtsmtanzuvmwarecomv1.FooTypeABC, opts v1.UpdateOptions) (result *configtsmtanzuvmwarecomv1.FooTypeABC, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(footypeabcsResource, fooTypeABC), &configtsmtanzuvmwarecomv1.FooTypeABC{})
	if obj == nil {
		return nil, err
	}
	return obj.(*configtsmtanzuvmwarecomv1.FooTypeABC), err
}

// Delete takes name of the fooTypeABC and deletes it. Returns an error if one occurs.
func (c *FakeFooTypeABCs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(footypeabcsResource, name), &configtsmtanzuvmwarecomv1.FooTypeABC{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeFooTypeABCs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(footypeabcsResource, listOpts)

	_, err := c.Fake.Invokes(action, &configtsmtanzuvmwarecomv1.FooTypeABCList{})
	return err
}

// Patch applies the patch and returns the patched fooTypeABC.
func (c *FakeFooTypeABCs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *configtsmtanzuvmwarecomv1.FooTypeABC, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(footypeabcsResource, name, pt, data, subresources...), &configtsmtanzuvmwarecomv1.FooTypeABC{})
	if obj == nil {
		return nil, err
	}
	return obj.(*configtsmtanzuvmwarecomv1.FooTypeABC), err
}
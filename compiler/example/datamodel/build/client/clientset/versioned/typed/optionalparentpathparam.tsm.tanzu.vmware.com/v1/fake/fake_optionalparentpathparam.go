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
	optionalparentpathparamtsmtanzuvmwarecomv1 "/build/apis/optionalparentpathparam.tsm.tanzu.vmware.com/v1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeOptionalParentPathParams implements OptionalParentPathParamInterface
type FakeOptionalParentPathParams struct {
	Fake *FakeOptionalparentpathparamTsmV1
}

var optionalparentpathparamsResource = schema.GroupVersionResource{Group: "optionalparentpathparam.tsm.tanzu.vmware.com", Version: "v1", Resource: "optionalparentpathparams"}

var optionalparentpathparamsKind = schema.GroupVersionKind{Group: "optionalparentpathparam.tsm.tanzu.vmware.com", Version: "v1", Kind: "OptionalParentPathParam"}

// Get takes name of the optionalParentPathParam, and returns the corresponding optionalParentPathParam object, and an error if there is any.
func (c *FakeOptionalParentPathParams) Get(ctx context.Context, name string, options v1.GetOptions) (result *optionalparentpathparamtsmtanzuvmwarecomv1.OptionalParentPathParam, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(optionalparentpathparamsResource, name), &optionalparentpathparamtsmtanzuvmwarecomv1.OptionalParentPathParam{})
	if obj == nil {
		return nil, err
	}
	return obj.(*optionalparentpathparamtsmtanzuvmwarecomv1.OptionalParentPathParam), err
}

// List takes label and field selectors, and returns the list of OptionalParentPathParams that match those selectors.
func (c *FakeOptionalParentPathParams) List(ctx context.Context, opts v1.ListOptions) (result *optionalparentpathparamtsmtanzuvmwarecomv1.OptionalParentPathParamList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(optionalparentpathparamsResource, optionalparentpathparamsKind, opts), &optionalparentpathparamtsmtanzuvmwarecomv1.OptionalParentPathParamList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &optionalparentpathparamtsmtanzuvmwarecomv1.OptionalParentPathParamList{ListMeta: obj.(*optionalparentpathparamtsmtanzuvmwarecomv1.OptionalParentPathParamList).ListMeta}
	for _, item := range obj.(*optionalparentpathparamtsmtanzuvmwarecomv1.OptionalParentPathParamList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested optionalParentPathParams.
func (c *FakeOptionalParentPathParams) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(optionalparentpathparamsResource, opts))
}

// Create takes the representation of a optionalParentPathParam and creates it.  Returns the server's representation of the optionalParentPathParam, and an error, if there is any.
func (c *FakeOptionalParentPathParams) Create(ctx context.Context, optionalParentPathParam *optionalparentpathparamtsmtanzuvmwarecomv1.OptionalParentPathParam, opts v1.CreateOptions) (result *optionalparentpathparamtsmtanzuvmwarecomv1.OptionalParentPathParam, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(optionalparentpathparamsResource, optionalParentPathParam), &optionalparentpathparamtsmtanzuvmwarecomv1.OptionalParentPathParam{})
	if obj == nil {
		return nil, err
	}
	return obj.(*optionalparentpathparamtsmtanzuvmwarecomv1.OptionalParentPathParam), err
}

// Update takes the representation of a optionalParentPathParam and updates it. Returns the server's representation of the optionalParentPathParam, and an error, if there is any.
func (c *FakeOptionalParentPathParams) Update(ctx context.Context, optionalParentPathParam *optionalparentpathparamtsmtanzuvmwarecomv1.OptionalParentPathParam, opts v1.UpdateOptions) (result *optionalparentpathparamtsmtanzuvmwarecomv1.OptionalParentPathParam, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(optionalparentpathparamsResource, optionalParentPathParam), &optionalparentpathparamtsmtanzuvmwarecomv1.OptionalParentPathParam{})
	if obj == nil {
		return nil, err
	}
	return obj.(*optionalparentpathparamtsmtanzuvmwarecomv1.OptionalParentPathParam), err
}

// Delete takes name of the optionalParentPathParam and deletes it. Returns an error if one occurs.
func (c *FakeOptionalParentPathParams) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(optionalparentpathparamsResource, name, opts), &optionalparentpathparamtsmtanzuvmwarecomv1.OptionalParentPathParam{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeOptionalParentPathParams) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(optionalparentpathparamsResource, listOpts)

	_, err := c.Fake.Invokes(action, &optionalparentpathparamtsmtanzuvmwarecomv1.OptionalParentPathParamList{})
	return err
}

// Patch applies the patch and returns the patched optionalParentPathParam.
func (c *FakeOptionalParentPathParams) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *optionalparentpathparamtsmtanzuvmwarecomv1.OptionalParentPathParam, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(optionalparentpathparamsResource, name, pt, data, subresources...), &optionalparentpathparamtsmtanzuvmwarecomv1.OptionalParentPathParam{})
	if obj == nil {
		return nil, err
	}
	return obj.(*optionalparentpathparamtsmtanzuvmwarecomv1.OptionalParentPathParam), err
}

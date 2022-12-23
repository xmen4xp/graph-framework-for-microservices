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

// FakeServiceConfigs implements ServiceConfigInterface
type FakeServiceConfigs struct {
	Fake *FakeGlobalTsmV1
}

var serviceconfigsResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "serviceconfigs"}

var serviceconfigsKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "ServiceConfig"}

// Get takes name of the serviceConfig, and returns the corresponding serviceConfig object, and an error if there is any.
func (c *FakeServiceConfigs) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.ServiceConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(serviceconfigsResource, name), &globaltsmtanzuvmwarecomv1.ServiceConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceConfig), err
}

// List takes label and field selectors, and returns the list of ServiceConfigs that match those selectors.
func (c *FakeServiceConfigs) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.ServiceConfigList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(serviceconfigsResource, serviceconfigsKind, opts), &globaltsmtanzuvmwarecomv1.ServiceConfigList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.ServiceConfigList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.ServiceConfigList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.ServiceConfigList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested serviceConfigs.
func (c *FakeServiceConfigs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(serviceconfigsResource, opts))
}

// Create takes the representation of a serviceConfig and creates it.  Returns the server's representation of the serviceConfig, and an error, if there is any.
func (c *FakeServiceConfigs) Create(ctx context.Context, serviceConfig *globaltsmtanzuvmwarecomv1.ServiceConfig, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.ServiceConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(serviceconfigsResource, serviceConfig), &globaltsmtanzuvmwarecomv1.ServiceConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceConfig), err
}

// Update takes the representation of a serviceConfig and updates it. Returns the server's representation of the serviceConfig, and an error, if there is any.
func (c *FakeServiceConfigs) Update(ctx context.Context, serviceConfig *globaltsmtanzuvmwarecomv1.ServiceConfig, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.ServiceConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(serviceconfigsResource, serviceConfig), &globaltsmtanzuvmwarecomv1.ServiceConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceConfig), err
}

// Delete takes name of the serviceConfig and deletes it. Returns an error if one occurs.
func (c *FakeServiceConfigs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(serviceconfigsResource, name), &globaltsmtanzuvmwarecomv1.ServiceConfig{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeServiceConfigs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(serviceconfigsResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.ServiceConfigList{})
	return err
}

// Patch applies the patch and returns the patched serviceConfig.
func (c *FakeServiceConfigs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.ServiceConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(serviceconfigsResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.ServiceConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceConfig), err
}

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

// FakeServiceEntryConfigs implements ServiceEntryConfigInterface
type FakeServiceEntryConfigs struct {
	Fake *FakeGlobalTsmV1
}

var serviceentryconfigsResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "serviceentryconfigs"}

var serviceentryconfigsKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "ServiceEntryConfig"}

// Get takes name of the serviceEntryConfig, and returns the corresponding serviceEntryConfig object, and an error if there is any.
func (c *FakeServiceEntryConfigs) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.ServiceEntryConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(serviceentryconfigsResource, name), &globaltsmtanzuvmwarecomv1.ServiceEntryConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceEntryConfig), err
}

// List takes label and field selectors, and returns the list of ServiceEntryConfigs that match those selectors.
func (c *FakeServiceEntryConfigs) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.ServiceEntryConfigList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(serviceentryconfigsResource, serviceentryconfigsKind, opts), &globaltsmtanzuvmwarecomv1.ServiceEntryConfigList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.ServiceEntryConfigList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.ServiceEntryConfigList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.ServiceEntryConfigList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested serviceEntryConfigs.
func (c *FakeServiceEntryConfigs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(serviceentryconfigsResource, opts))
}

// Create takes the representation of a serviceEntryConfig and creates it.  Returns the server's representation of the serviceEntryConfig, and an error, if there is any.
func (c *FakeServiceEntryConfigs) Create(ctx context.Context, serviceEntryConfig *globaltsmtanzuvmwarecomv1.ServiceEntryConfig, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.ServiceEntryConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(serviceentryconfigsResource, serviceEntryConfig), &globaltsmtanzuvmwarecomv1.ServiceEntryConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceEntryConfig), err
}

// Update takes the representation of a serviceEntryConfig and updates it. Returns the server's representation of the serviceEntryConfig, and an error, if there is any.
func (c *FakeServiceEntryConfigs) Update(ctx context.Context, serviceEntryConfig *globaltsmtanzuvmwarecomv1.ServiceEntryConfig, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.ServiceEntryConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(serviceentryconfigsResource, serviceEntryConfig), &globaltsmtanzuvmwarecomv1.ServiceEntryConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceEntryConfig), err
}

// Delete takes name of the serviceEntryConfig and deletes it. Returns an error if one occurs.
func (c *FakeServiceEntryConfigs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(serviceentryconfigsResource, name), &globaltsmtanzuvmwarecomv1.ServiceEntryConfig{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeServiceEntryConfigs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(serviceentryconfigsResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.ServiceEntryConfigList{})
	return err
}

// Patch applies the patch and returns the patched serviceEntryConfig.
func (c *FakeServiceEntryConfigs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.ServiceEntryConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(serviceentryconfigsResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.ServiceEntryConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceEntryConfig), err
}

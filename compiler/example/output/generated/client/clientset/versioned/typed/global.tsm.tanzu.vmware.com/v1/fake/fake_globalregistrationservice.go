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

// FakeGlobalRegistrationServices implements GlobalRegistrationServiceInterface
type FakeGlobalRegistrationServices struct {
	Fake *FakeGlobalTsmV1
}

var globalregistrationservicesResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "globalregistrationservices"}

var globalregistrationservicesKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "GlobalRegistrationService"}

// Get takes name of the globalRegistrationService, and returns the corresponding globalRegistrationService object, and an error if there is any.
func (c *FakeGlobalRegistrationServices) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.GlobalRegistrationService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(globalregistrationservicesResource, name), &globaltsmtanzuvmwarecomv1.GlobalRegistrationService{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.GlobalRegistrationService), err
}

// List takes label and field selectors, and returns the list of GlobalRegistrationServices that match those selectors.
func (c *FakeGlobalRegistrationServices) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.GlobalRegistrationServiceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(globalregistrationservicesResource, globalregistrationservicesKind, opts), &globaltsmtanzuvmwarecomv1.GlobalRegistrationServiceList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.GlobalRegistrationServiceList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.GlobalRegistrationServiceList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.GlobalRegistrationServiceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested globalRegistrationServices.
func (c *FakeGlobalRegistrationServices) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(globalregistrationservicesResource, opts))
}

// Create takes the representation of a globalRegistrationService and creates it.  Returns the server's representation of the globalRegistrationService, and an error, if there is any.
func (c *FakeGlobalRegistrationServices) Create(ctx context.Context, globalRegistrationService *globaltsmtanzuvmwarecomv1.GlobalRegistrationService, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.GlobalRegistrationService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(globalregistrationservicesResource, globalRegistrationService), &globaltsmtanzuvmwarecomv1.GlobalRegistrationService{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.GlobalRegistrationService), err
}

// Update takes the representation of a globalRegistrationService and updates it. Returns the server's representation of the globalRegistrationService, and an error, if there is any.
func (c *FakeGlobalRegistrationServices) Update(ctx context.Context, globalRegistrationService *globaltsmtanzuvmwarecomv1.GlobalRegistrationService, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.GlobalRegistrationService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(globalregistrationservicesResource, globalRegistrationService), &globaltsmtanzuvmwarecomv1.GlobalRegistrationService{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.GlobalRegistrationService), err
}

// Delete takes name of the globalRegistrationService and deletes it. Returns an error if one occurs.
func (c *FakeGlobalRegistrationServices) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(globalregistrationservicesResource, name), &globaltsmtanzuvmwarecomv1.GlobalRegistrationService{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeGlobalRegistrationServices) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(globalregistrationservicesResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.GlobalRegistrationServiceList{})
	return err
}

// Patch applies the patch and returns the patched globalRegistrationService.
func (c *FakeGlobalRegistrationServices) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.GlobalRegistrationService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(globalregistrationservicesResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.GlobalRegistrationService{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.GlobalRegistrationService), err
}

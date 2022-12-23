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

// FakeServiceInstanceContainers implements ServiceInstanceContainerInterface
type FakeServiceInstanceContainers struct {
	Fake *FakeGlobalTsmV1
}

var serviceinstancecontainersResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "serviceinstancecontainers"}

var serviceinstancecontainersKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "ServiceInstanceContainer"}

// Get takes name of the serviceInstanceContainer, and returns the corresponding serviceInstanceContainer object, and an error if there is any.
func (c *FakeServiceInstanceContainers) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.ServiceInstanceContainer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(serviceinstancecontainersResource, name), &globaltsmtanzuvmwarecomv1.ServiceInstanceContainer{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceInstanceContainer), err
}

// List takes label and field selectors, and returns the list of ServiceInstanceContainers that match those selectors.
func (c *FakeServiceInstanceContainers) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.ServiceInstanceContainerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(serviceinstancecontainersResource, serviceinstancecontainersKind, opts), &globaltsmtanzuvmwarecomv1.ServiceInstanceContainerList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.ServiceInstanceContainerList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.ServiceInstanceContainerList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.ServiceInstanceContainerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested serviceInstanceContainers.
func (c *FakeServiceInstanceContainers) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(serviceinstancecontainersResource, opts))
}

// Create takes the representation of a serviceInstanceContainer and creates it.  Returns the server's representation of the serviceInstanceContainer, and an error, if there is any.
func (c *FakeServiceInstanceContainers) Create(ctx context.Context, serviceInstanceContainer *globaltsmtanzuvmwarecomv1.ServiceInstanceContainer, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.ServiceInstanceContainer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(serviceinstancecontainersResource, serviceInstanceContainer), &globaltsmtanzuvmwarecomv1.ServiceInstanceContainer{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceInstanceContainer), err
}

// Update takes the representation of a serviceInstanceContainer and updates it. Returns the server's representation of the serviceInstanceContainer, and an error, if there is any.
func (c *FakeServiceInstanceContainers) Update(ctx context.Context, serviceInstanceContainer *globaltsmtanzuvmwarecomv1.ServiceInstanceContainer, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.ServiceInstanceContainer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(serviceinstancecontainersResource, serviceInstanceContainer), &globaltsmtanzuvmwarecomv1.ServiceInstanceContainer{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceInstanceContainer), err
}

// Delete takes name of the serviceInstanceContainer and deletes it. Returns an error if one occurs.
func (c *FakeServiceInstanceContainers) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(serviceinstancecontainersResource, name), &globaltsmtanzuvmwarecomv1.ServiceInstanceContainer{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeServiceInstanceContainers) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(serviceinstancecontainersResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.ServiceInstanceContainerList{})
	return err
}

// Patch applies the patch and returns the patched serviceInstanceContainer.
func (c *FakeServiceInstanceContainers) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.ServiceInstanceContainer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(serviceinstancecontainersResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.ServiceInstanceContainer{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceInstanceContainer), err
}

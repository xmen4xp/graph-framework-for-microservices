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

// FakeServiceDeploymentContainers implements ServiceDeploymentContainerInterface
type FakeServiceDeploymentContainers struct {
	Fake *FakeGlobalTsmV1
}

var servicedeploymentcontainersResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "servicedeploymentcontainers"}

var servicedeploymentcontainersKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "ServiceDeploymentContainer"}

// Get takes name of the serviceDeploymentContainer, and returns the corresponding serviceDeploymentContainer object, and an error if there is any.
func (c *FakeServiceDeploymentContainers) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.ServiceDeploymentContainer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(servicedeploymentcontainersResource, name), &globaltsmtanzuvmwarecomv1.ServiceDeploymentContainer{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceDeploymentContainer), err
}

// List takes label and field selectors, and returns the list of ServiceDeploymentContainers that match those selectors.
func (c *FakeServiceDeploymentContainers) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.ServiceDeploymentContainerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(servicedeploymentcontainersResource, servicedeploymentcontainersKind, opts), &globaltsmtanzuvmwarecomv1.ServiceDeploymentContainerList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.ServiceDeploymentContainerList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.ServiceDeploymentContainerList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.ServiceDeploymentContainerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested serviceDeploymentContainers.
func (c *FakeServiceDeploymentContainers) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(servicedeploymentcontainersResource, opts))
}

// Create takes the representation of a serviceDeploymentContainer and creates it.  Returns the server's representation of the serviceDeploymentContainer, and an error, if there is any.
func (c *FakeServiceDeploymentContainers) Create(ctx context.Context, serviceDeploymentContainer *globaltsmtanzuvmwarecomv1.ServiceDeploymentContainer, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.ServiceDeploymentContainer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(servicedeploymentcontainersResource, serviceDeploymentContainer), &globaltsmtanzuvmwarecomv1.ServiceDeploymentContainer{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceDeploymentContainer), err
}

// Update takes the representation of a serviceDeploymentContainer and updates it. Returns the server's representation of the serviceDeploymentContainer, and an error, if there is any.
func (c *FakeServiceDeploymentContainers) Update(ctx context.Context, serviceDeploymentContainer *globaltsmtanzuvmwarecomv1.ServiceDeploymentContainer, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.ServiceDeploymentContainer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(servicedeploymentcontainersResource, serviceDeploymentContainer), &globaltsmtanzuvmwarecomv1.ServiceDeploymentContainer{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceDeploymentContainer), err
}

// Delete takes name of the serviceDeploymentContainer and deletes it. Returns an error if one occurs.
func (c *FakeServiceDeploymentContainers) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(servicedeploymentcontainersResource, name), &globaltsmtanzuvmwarecomv1.ServiceDeploymentContainer{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeServiceDeploymentContainers) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(servicedeploymentcontainersResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.ServiceDeploymentContainerList{})
	return err
}

// Patch applies the patch and returns the patched serviceDeploymentContainer.
func (c *FakeServiceDeploymentContainers) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.ServiceDeploymentContainer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(servicedeploymentcontainersResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.ServiceDeploymentContainer{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceDeploymentContainer), err
}

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

// FakeNodeStatuses implements NodeStatusInterface
type FakeNodeStatuses struct {
	Fake *FakeGlobalTsmV1
}

var nodestatusesResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "nodestatuses"}

var nodestatusesKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "NodeStatus"}

// Get takes name of the nodeStatus, and returns the corresponding nodeStatus object, and an error if there is any.
func (c *FakeNodeStatuses) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.NodeStatus, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(nodestatusesResource, name), &globaltsmtanzuvmwarecomv1.NodeStatus{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.NodeStatus), err
}

// List takes label and field selectors, and returns the list of NodeStatuses that match those selectors.
func (c *FakeNodeStatuses) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.NodeStatusList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(nodestatusesResource, nodestatusesKind, opts), &globaltsmtanzuvmwarecomv1.NodeStatusList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.NodeStatusList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.NodeStatusList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.NodeStatusList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested nodeStatuses.
func (c *FakeNodeStatuses) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(nodestatusesResource, opts))
}

// Create takes the representation of a nodeStatus and creates it.  Returns the server's representation of the nodeStatus, and an error, if there is any.
func (c *FakeNodeStatuses) Create(ctx context.Context, nodeStatus *globaltsmtanzuvmwarecomv1.NodeStatus, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.NodeStatus, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(nodestatusesResource, nodeStatus), &globaltsmtanzuvmwarecomv1.NodeStatus{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.NodeStatus), err
}

// Update takes the representation of a nodeStatus and updates it. Returns the server's representation of the nodeStatus, and an error, if there is any.
func (c *FakeNodeStatuses) Update(ctx context.Context, nodeStatus *globaltsmtanzuvmwarecomv1.NodeStatus, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.NodeStatus, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(nodestatusesResource, nodeStatus), &globaltsmtanzuvmwarecomv1.NodeStatus{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.NodeStatus), err
}

// Delete takes name of the nodeStatus and deletes it. Returns an error if one occurs.
func (c *FakeNodeStatuses) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(nodestatusesResource, name), &globaltsmtanzuvmwarecomv1.NodeStatus{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeNodeStatuses) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(nodestatusesResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.NodeStatusList{})
	return err
}

// Patch applies the patch and returns the patched nodeStatus.
func (c *FakeNodeStatuses) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.NodeStatus, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(nodestatusesResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.NodeStatus{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.NodeStatus), err
}

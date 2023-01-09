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

// FakeClusterSettingses implements ClusterSettingsInterface
type FakeClusterSettingses struct {
	Fake *FakeGlobalTsmV1
}

var clustersettingsesResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "clustersettingses"}

var clustersettingsesKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "ClusterSettings"}

// Get takes name of the clusterSettings, and returns the corresponding clusterSettings object, and an error if there is any.
func (c *FakeClusterSettingses) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.ClusterSettings, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(clustersettingsesResource, name), &globaltsmtanzuvmwarecomv1.ClusterSettings{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ClusterSettings), err
}

// List takes label and field selectors, and returns the list of ClusterSettingses that match those selectors.
func (c *FakeClusterSettingses) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.ClusterSettingsList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(clustersettingsesResource, clustersettingsesKind, opts), &globaltsmtanzuvmwarecomv1.ClusterSettingsList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.ClusterSettingsList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.ClusterSettingsList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.ClusterSettingsList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested clusterSettingses.
func (c *FakeClusterSettingses) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(clustersettingsesResource, opts))
}

// Create takes the representation of a clusterSettings and creates it.  Returns the server's representation of the clusterSettings, and an error, if there is any.
func (c *FakeClusterSettingses) Create(ctx context.Context, clusterSettings *globaltsmtanzuvmwarecomv1.ClusterSettings, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.ClusterSettings, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(clustersettingsesResource, clusterSettings), &globaltsmtanzuvmwarecomv1.ClusterSettings{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ClusterSettings), err
}

// Update takes the representation of a clusterSettings and updates it. Returns the server's representation of the clusterSettings, and an error, if there is any.
func (c *FakeClusterSettingses) Update(ctx context.Context, clusterSettings *globaltsmtanzuvmwarecomv1.ClusterSettings, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.ClusterSettings, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(clustersettingsesResource, clusterSettings), &globaltsmtanzuvmwarecomv1.ClusterSettings{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ClusterSettings), err
}

// Delete takes name of the clusterSettings and deletes it. Returns an error if one occurs.
func (c *FakeClusterSettingses) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(clustersettingsesResource, name), &globaltsmtanzuvmwarecomv1.ClusterSettings{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeClusterSettingses) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(clustersettingsesResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.ClusterSettingsList{})
	return err
}

// Patch applies the patch and returns the patched clusterSettings.
func (c *FakeClusterSettingses) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.ClusterSettings, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(clustersettingsesResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.ClusterSettings{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ClusterSettings), err
}
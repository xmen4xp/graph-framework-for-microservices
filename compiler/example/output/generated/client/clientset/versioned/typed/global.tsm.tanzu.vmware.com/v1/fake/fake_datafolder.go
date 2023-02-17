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

// FakeDataFolders implements DataFolderInterface
type FakeDataFolders struct {
	Fake *FakeGlobalTsmV1
}

var datafoldersResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "datafolders"}

var datafoldersKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "DataFolder"}

// Get takes name of the dataFolder, and returns the corresponding dataFolder object, and an error if there is any.
func (c *FakeDataFolders) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.DataFolder, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(datafoldersResource, name), &globaltsmtanzuvmwarecomv1.DataFolder{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.DataFolder), err
}

// List takes label and field selectors, and returns the list of DataFolders that match those selectors.
func (c *FakeDataFolders) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.DataFolderList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(datafoldersResource, datafoldersKind, opts), &globaltsmtanzuvmwarecomv1.DataFolderList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.DataFolderList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.DataFolderList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.DataFolderList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested dataFolders.
func (c *FakeDataFolders) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(datafoldersResource, opts))
}

// Create takes the representation of a dataFolder and creates it.  Returns the server's representation of the dataFolder, and an error, if there is any.
func (c *FakeDataFolders) Create(ctx context.Context, dataFolder *globaltsmtanzuvmwarecomv1.DataFolder, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.DataFolder, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(datafoldersResource, dataFolder), &globaltsmtanzuvmwarecomv1.DataFolder{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.DataFolder), err
}

// Update takes the representation of a dataFolder and updates it. Returns the server's representation of the dataFolder, and an error, if there is any.
func (c *FakeDataFolders) Update(ctx context.Context, dataFolder *globaltsmtanzuvmwarecomv1.DataFolder, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.DataFolder, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(datafoldersResource, dataFolder), &globaltsmtanzuvmwarecomv1.DataFolder{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.DataFolder), err
}

// Delete takes name of the dataFolder and deletes it. Returns an error if one occurs.
func (c *FakeDataFolders) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(datafoldersResource, name), &globaltsmtanzuvmwarecomv1.DataFolder{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeDataFolders) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(datafoldersResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.DataFolderList{})
	return err
}

// Patch applies the patch and returns the patched dataFolder.
func (c *FakeDataFolders) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.DataFolder, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(datafoldersResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.DataFolder{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.DataFolder), err
}
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

// FakeApplicationInfos implements ApplicationInfoInterface
type FakeApplicationInfos struct {
	Fake *FakeGlobalTsmV1
}

var applicationinfosResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "applicationinfos"}

var applicationinfosKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "ApplicationInfo"}

// Get takes name of the applicationInfo, and returns the corresponding applicationInfo object, and an error if there is any.
func (c *FakeApplicationInfos) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.ApplicationInfo, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(applicationinfosResource, name), &globaltsmtanzuvmwarecomv1.ApplicationInfo{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ApplicationInfo), err
}

// List takes label and field selectors, and returns the list of ApplicationInfos that match those selectors.
func (c *FakeApplicationInfos) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.ApplicationInfoList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(applicationinfosResource, applicationinfosKind, opts), &globaltsmtanzuvmwarecomv1.ApplicationInfoList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.ApplicationInfoList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.ApplicationInfoList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.ApplicationInfoList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested applicationInfos.
func (c *FakeApplicationInfos) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(applicationinfosResource, opts))
}

// Create takes the representation of a applicationInfo and creates it.  Returns the server's representation of the applicationInfo, and an error, if there is any.
func (c *FakeApplicationInfos) Create(ctx context.Context, applicationInfo *globaltsmtanzuvmwarecomv1.ApplicationInfo, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.ApplicationInfo, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(applicationinfosResource, applicationInfo), &globaltsmtanzuvmwarecomv1.ApplicationInfo{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ApplicationInfo), err
}

// Update takes the representation of a applicationInfo and updates it. Returns the server's representation of the applicationInfo, and an error, if there is any.
func (c *FakeApplicationInfos) Update(ctx context.Context, applicationInfo *globaltsmtanzuvmwarecomv1.ApplicationInfo, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.ApplicationInfo, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(applicationinfosResource, applicationInfo), &globaltsmtanzuvmwarecomv1.ApplicationInfo{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ApplicationInfo), err
}

// Delete takes name of the applicationInfo and deletes it. Returns an error if one occurs.
func (c *FakeApplicationInfos) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(applicationinfosResource, name), &globaltsmtanzuvmwarecomv1.ApplicationInfo{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeApplicationInfos) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(applicationinfosResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.ApplicationInfoList{})
	return err
}

// Patch applies the patch and returns the patched applicationInfo.
func (c *FakeApplicationInfos) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.ApplicationInfo, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(applicationinfosResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.ApplicationInfo{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ApplicationInfo), err
}

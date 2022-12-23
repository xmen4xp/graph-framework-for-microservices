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

// FakeAppTemplates implements AppTemplateInterface
type FakeAppTemplates struct {
	Fake *FakeGlobalTsmV1
}

var apptemplatesResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "apptemplates"}

var apptemplatesKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "AppTemplate"}

// Get takes name of the appTemplate, and returns the corresponding appTemplate object, and an error if there is any.
func (c *FakeAppTemplates) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.AppTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(apptemplatesResource, name), &globaltsmtanzuvmwarecomv1.AppTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.AppTemplate), err
}

// List takes label and field selectors, and returns the list of AppTemplates that match those selectors.
func (c *FakeAppTemplates) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.AppTemplateList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(apptemplatesResource, apptemplatesKind, opts), &globaltsmtanzuvmwarecomv1.AppTemplateList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.AppTemplateList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.AppTemplateList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.AppTemplateList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested appTemplates.
func (c *FakeAppTemplates) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(apptemplatesResource, opts))
}

// Create takes the representation of a appTemplate and creates it.  Returns the server's representation of the appTemplate, and an error, if there is any.
func (c *FakeAppTemplates) Create(ctx context.Context, appTemplate *globaltsmtanzuvmwarecomv1.AppTemplate, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.AppTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(apptemplatesResource, appTemplate), &globaltsmtanzuvmwarecomv1.AppTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.AppTemplate), err
}

// Update takes the representation of a appTemplate and updates it. Returns the server's representation of the appTemplate, and an error, if there is any.
func (c *FakeAppTemplates) Update(ctx context.Context, appTemplate *globaltsmtanzuvmwarecomv1.AppTemplate, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.AppTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(apptemplatesResource, appTemplate), &globaltsmtanzuvmwarecomv1.AppTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.AppTemplate), err
}

// Delete takes name of the appTemplate and deletes it. Returns an error if one occurs.
func (c *FakeAppTemplates) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(apptemplatesResource, name), &globaltsmtanzuvmwarecomv1.AppTemplate{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAppTemplates) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(apptemplatesResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.AppTemplateList{})
	return err
}

// Patch applies the patch and returns the patched appTemplate.
func (c *FakeAppTemplates) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.AppTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(apptemplatesResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.AppTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.AppTemplate), err
}

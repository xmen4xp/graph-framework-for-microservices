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

	gnstsmtanzuvmwarecomv1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/crd_generated/apis/gns.tsm.tanzu.vmware.com/v1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeAdditionalGnsDatas implements AdditionalGnsDataInterface
type FakeAdditionalGnsDatas struct {
	Fake *FakeGnsTsmV1
}

var additionalgnsdatasResource = schema.GroupVersionResource{Group: "gns.tsm.tanzu.vmware.com", Version: "v1", Resource: "additionalgnsdatas"}

var additionalgnsdatasKind = schema.GroupVersionKind{Group: "gns.tsm.tanzu.vmware.com", Version: "v1", Kind: "AdditionalGnsData"}

// Get takes name of the additionalGnsData, and returns the corresponding additionalGnsData object, and an error if there is any.
func (c *FakeAdditionalGnsDatas) Get(ctx context.Context, name string, options v1.GetOptions) (result *gnstsmtanzuvmwarecomv1.AdditionalGnsData, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(additionalgnsdatasResource, name), &gnstsmtanzuvmwarecomv1.AdditionalGnsData{})
	if obj == nil {
		return nil, err
	}
	return obj.(*gnstsmtanzuvmwarecomv1.AdditionalGnsData), err
}

// List takes label and field selectors, and returns the list of AdditionalGnsDatas that match those selectors.
func (c *FakeAdditionalGnsDatas) List(ctx context.Context, opts v1.ListOptions) (result *gnstsmtanzuvmwarecomv1.AdditionalGnsDataList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(additionalgnsdatasResource, additionalgnsdatasKind, opts), &gnstsmtanzuvmwarecomv1.AdditionalGnsDataList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &gnstsmtanzuvmwarecomv1.AdditionalGnsDataList{ListMeta: obj.(*gnstsmtanzuvmwarecomv1.AdditionalGnsDataList).ListMeta}
	for _, item := range obj.(*gnstsmtanzuvmwarecomv1.AdditionalGnsDataList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested additionalGnsDatas.
func (c *FakeAdditionalGnsDatas) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(additionalgnsdatasResource, opts))
}

// Create takes the representation of a additionalGnsData and creates it.  Returns the server's representation of the additionalGnsData, and an error, if there is any.
func (c *FakeAdditionalGnsDatas) Create(ctx context.Context, additionalGnsData *gnstsmtanzuvmwarecomv1.AdditionalGnsData, opts v1.CreateOptions) (result *gnstsmtanzuvmwarecomv1.AdditionalGnsData, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(additionalgnsdatasResource, additionalGnsData), &gnstsmtanzuvmwarecomv1.AdditionalGnsData{})
	if obj == nil {
		return nil, err
	}
	return obj.(*gnstsmtanzuvmwarecomv1.AdditionalGnsData), err
}

// Update takes the representation of a additionalGnsData and updates it. Returns the server's representation of the additionalGnsData, and an error, if there is any.
func (c *FakeAdditionalGnsDatas) Update(ctx context.Context, additionalGnsData *gnstsmtanzuvmwarecomv1.AdditionalGnsData, opts v1.UpdateOptions) (result *gnstsmtanzuvmwarecomv1.AdditionalGnsData, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(additionalgnsdatasResource, additionalGnsData), &gnstsmtanzuvmwarecomv1.AdditionalGnsData{})
	if obj == nil {
		return nil, err
	}
	return obj.(*gnstsmtanzuvmwarecomv1.AdditionalGnsData), err
}

// Delete takes name of the additionalGnsData and deletes it. Returns an error if one occurs.
func (c *FakeAdditionalGnsDatas) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(additionalgnsdatasResource, name), &gnstsmtanzuvmwarecomv1.AdditionalGnsData{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAdditionalGnsDatas) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(additionalgnsdatasResource, listOpts)

	_, err := c.Fake.Invokes(action, &gnstsmtanzuvmwarecomv1.AdditionalGnsDataList{})
	return err
}

// Patch applies the patch and returns the patched additionalGnsData.
func (c *FakeAdditionalGnsDatas) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *gnstsmtanzuvmwarecomv1.AdditionalGnsData, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(additionalgnsdatasResource, name, pt, data, subresources...), &gnstsmtanzuvmwarecomv1.AdditionalGnsData{})
	if obj == nil {
		return nil, err
	}
	return obj.(*gnstsmtanzuvmwarecomv1.AdditionalGnsData), err
}

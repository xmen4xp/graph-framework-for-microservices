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

// FakeFeatureFlags implements FeatureFlagInterface
type FakeFeatureFlags struct {
	Fake *FakeGlobalTsmV1
}

var featureflagsResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "featureflags"}

var featureflagsKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "FeatureFlag"}

// Get takes name of the featureFlag, and returns the corresponding featureFlag object, and an error if there is any.
func (c *FakeFeatureFlags) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.FeatureFlag, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(featureflagsResource, name), &globaltsmtanzuvmwarecomv1.FeatureFlag{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.FeatureFlag), err
}

// List takes label and field selectors, and returns the list of FeatureFlags that match those selectors.
func (c *FakeFeatureFlags) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.FeatureFlagList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(featureflagsResource, featureflagsKind, opts), &globaltsmtanzuvmwarecomv1.FeatureFlagList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.FeatureFlagList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.FeatureFlagList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.FeatureFlagList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested featureFlags.
func (c *FakeFeatureFlags) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(featureflagsResource, opts))
}

// Create takes the representation of a featureFlag and creates it.  Returns the server's representation of the featureFlag, and an error, if there is any.
func (c *FakeFeatureFlags) Create(ctx context.Context, featureFlag *globaltsmtanzuvmwarecomv1.FeatureFlag, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.FeatureFlag, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(featureflagsResource, featureFlag), &globaltsmtanzuvmwarecomv1.FeatureFlag{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.FeatureFlag), err
}

// Update takes the representation of a featureFlag and updates it. Returns the server's representation of the featureFlag, and an error, if there is any.
func (c *FakeFeatureFlags) Update(ctx context.Context, featureFlag *globaltsmtanzuvmwarecomv1.FeatureFlag, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.FeatureFlag, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(featureflagsResource, featureFlag), &globaltsmtanzuvmwarecomv1.FeatureFlag{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.FeatureFlag), err
}

// Delete takes name of the featureFlag and deletes it. Returns an error if one occurs.
func (c *FakeFeatureFlags) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(featureflagsResource, name), &globaltsmtanzuvmwarecomv1.FeatureFlag{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeFeatureFlags) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(featureflagsResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.FeatureFlagList{})
	return err
}

// Patch applies the patch and returns the patched featureFlag.
func (c *FakeFeatureFlags) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.FeatureFlag, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(featureflagsResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.FeatureFlag{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.FeatureFlag), err
}

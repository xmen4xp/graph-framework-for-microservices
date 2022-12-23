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

// FakeDNSConfigFolders implements DNSConfigFolderInterface
type FakeDNSConfigFolders struct {
	Fake *FakeGlobalTsmV1
}

var dnsconfigfoldersResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "dnsconfigfolders"}

var dnsconfigfoldersKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "DNSConfigFolder"}

// Get takes name of the dNSConfigFolder, and returns the corresponding dNSConfigFolder object, and an error if there is any.
func (c *FakeDNSConfigFolders) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.DNSConfigFolder, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(dnsconfigfoldersResource, name), &globaltsmtanzuvmwarecomv1.DNSConfigFolder{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.DNSConfigFolder), err
}

// List takes label and field selectors, and returns the list of DNSConfigFolders that match those selectors.
func (c *FakeDNSConfigFolders) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.DNSConfigFolderList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(dnsconfigfoldersResource, dnsconfigfoldersKind, opts), &globaltsmtanzuvmwarecomv1.DNSConfigFolderList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.DNSConfigFolderList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.DNSConfigFolderList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.DNSConfigFolderList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested dNSConfigFolders.
func (c *FakeDNSConfigFolders) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(dnsconfigfoldersResource, opts))
}

// Create takes the representation of a dNSConfigFolder and creates it.  Returns the server's representation of the dNSConfigFolder, and an error, if there is any.
func (c *FakeDNSConfigFolders) Create(ctx context.Context, dNSConfigFolder *globaltsmtanzuvmwarecomv1.DNSConfigFolder, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.DNSConfigFolder, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(dnsconfigfoldersResource, dNSConfigFolder), &globaltsmtanzuvmwarecomv1.DNSConfigFolder{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.DNSConfigFolder), err
}

// Update takes the representation of a dNSConfigFolder and updates it. Returns the server's representation of the dNSConfigFolder, and an error, if there is any.
func (c *FakeDNSConfigFolders) Update(ctx context.Context, dNSConfigFolder *globaltsmtanzuvmwarecomv1.DNSConfigFolder, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.DNSConfigFolder, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(dnsconfigfoldersResource, dNSConfigFolder), &globaltsmtanzuvmwarecomv1.DNSConfigFolder{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.DNSConfigFolder), err
}

// Delete takes name of the dNSConfigFolder and deletes it. Returns an error if one occurs.
func (c *FakeDNSConfigFolders) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(dnsconfigfoldersResource, name), &globaltsmtanzuvmwarecomv1.DNSConfigFolder{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeDNSConfigFolders) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(dnsconfigfoldersResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.DNSConfigFolderList{})
	return err
}

// Patch applies the patch and returns the patched dNSConfigFolder.
func (c *FakeDNSConfigFolders) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.DNSConfigFolder, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(dnsconfigfoldersResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.DNSConfigFolder{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.DNSConfigFolder), err
}

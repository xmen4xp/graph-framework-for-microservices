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

// FakeCertificateAuthorityConfigNs implements CertificateAuthorityConfigNInterface
type FakeCertificateAuthorityConfigNs struct {
	Fake *FakeGlobalTsmV1
}

var certificateauthorityconfignsResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "certificateauthorityconfigns"}

var certificateauthorityconfignsKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "CertificateAuthorityConfigN"}

// Get takes name of the certificateAuthorityConfigN, and returns the corresponding certificateAuthorityConfigN object, and an error if there is any.
func (c *FakeCertificateAuthorityConfigNs) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.CertificateAuthorityConfigN, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(certificateauthorityconfignsResource, name), &globaltsmtanzuvmwarecomv1.CertificateAuthorityConfigN{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.CertificateAuthorityConfigN), err
}

// List takes label and field selectors, and returns the list of CertificateAuthorityConfigNs that match those selectors.
func (c *FakeCertificateAuthorityConfigNs) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.CertificateAuthorityConfigNList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(certificateauthorityconfignsResource, certificateauthorityconfignsKind, opts), &globaltsmtanzuvmwarecomv1.CertificateAuthorityConfigNList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.CertificateAuthorityConfigNList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.CertificateAuthorityConfigNList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.CertificateAuthorityConfigNList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested certificateAuthorityConfigNs.
func (c *FakeCertificateAuthorityConfigNs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(certificateauthorityconfignsResource, opts))
}

// Create takes the representation of a certificateAuthorityConfigN and creates it.  Returns the server's representation of the certificateAuthorityConfigN, and an error, if there is any.
func (c *FakeCertificateAuthorityConfigNs) Create(ctx context.Context, certificateAuthorityConfigN *globaltsmtanzuvmwarecomv1.CertificateAuthorityConfigN, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.CertificateAuthorityConfigN, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(certificateauthorityconfignsResource, certificateAuthorityConfigN), &globaltsmtanzuvmwarecomv1.CertificateAuthorityConfigN{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.CertificateAuthorityConfigN), err
}

// Update takes the representation of a certificateAuthorityConfigN and updates it. Returns the server's representation of the certificateAuthorityConfigN, and an error, if there is any.
func (c *FakeCertificateAuthorityConfigNs) Update(ctx context.Context, certificateAuthorityConfigN *globaltsmtanzuvmwarecomv1.CertificateAuthorityConfigN, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.CertificateAuthorityConfigN, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(certificateauthorityconfignsResource, certificateAuthorityConfigN), &globaltsmtanzuvmwarecomv1.CertificateAuthorityConfigN{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.CertificateAuthorityConfigN), err
}

// Delete takes name of the certificateAuthorityConfigN and deletes it. Returns an error if one occurs.
func (c *FakeCertificateAuthorityConfigNs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(certificateauthorityconfignsResource, name), &globaltsmtanzuvmwarecomv1.CertificateAuthorityConfigN{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeCertificateAuthorityConfigNs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(certificateauthorityconfignsResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.CertificateAuthorityConfigNList{})
	return err
}

// Patch applies the patch and returns the patched certificateAuthorityConfigN.
func (c *FakeCertificateAuthorityConfigNs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.CertificateAuthorityConfigN, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(certificateauthorityconfignsResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.CertificateAuthorityConfigN{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.CertificateAuthorityConfigN), err
}

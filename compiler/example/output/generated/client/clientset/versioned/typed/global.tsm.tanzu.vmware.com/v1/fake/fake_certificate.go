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

// FakeCertificates implements CertificateInterface
type FakeCertificates struct {
	Fake *FakeGlobalTsmV1
}

var certificatesResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "certificates"}

var certificatesKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "Certificate"}

// Get takes name of the certificate, and returns the corresponding certificate object, and an error if there is any.
func (c *FakeCertificates) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.Certificate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(certificatesResource, name), &globaltsmtanzuvmwarecomv1.Certificate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.Certificate), err
}

// List takes label and field selectors, and returns the list of Certificates that match those selectors.
func (c *FakeCertificates) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.CertificateList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(certificatesResource, certificatesKind, opts), &globaltsmtanzuvmwarecomv1.CertificateList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.CertificateList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.CertificateList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.CertificateList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested certificates.
func (c *FakeCertificates) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(certificatesResource, opts))
}

// Create takes the representation of a certificate and creates it.  Returns the server's representation of the certificate, and an error, if there is any.
func (c *FakeCertificates) Create(ctx context.Context, certificate *globaltsmtanzuvmwarecomv1.Certificate, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.Certificate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(certificatesResource, certificate), &globaltsmtanzuvmwarecomv1.Certificate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.Certificate), err
}

// Update takes the representation of a certificate and updates it. Returns the server's representation of the certificate, and an error, if there is any.
func (c *FakeCertificates) Update(ctx context.Context, certificate *globaltsmtanzuvmwarecomv1.Certificate, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.Certificate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(certificatesResource, certificate), &globaltsmtanzuvmwarecomv1.Certificate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.Certificate), err
}

// Delete takes name of the certificate and deletes it. Returns an error if one occurs.
func (c *FakeCertificates) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(certificatesResource, name), &globaltsmtanzuvmwarecomv1.Certificate{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeCertificates) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(certificatesResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.CertificateList{})
	return err
}

// Patch applies the patch and returns the patched certificate.
func (c *FakeCertificates) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.Certificate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(certificatesResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.Certificate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.Certificate), err
}
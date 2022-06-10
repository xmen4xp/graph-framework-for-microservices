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

	authenticationnexusorgv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/authentication.nexus.org/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeOIDCs implements OIDCInterface
type FakeOIDCs struct {
	Fake *FakeAuthenticationNexusV1
}

var oidcsResource = schema.GroupVersionResource{Group: "authentication.nexus.org", Version: "v1", Resource: "oidcs"}

var oidcsKind = schema.GroupVersionKind{Group: "authentication.nexus.org", Version: "v1", Kind: "OIDC"}

// Get takes name of the oIDC, and returns the corresponding oIDC object, and an error if there is any.
func (c *FakeOIDCs) Get(ctx context.Context, name string, options v1.GetOptions) (result *authenticationnexusorgv1.OIDC, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(oidcsResource, name), &authenticationnexusorgv1.OIDC{})
	if obj == nil {
		return nil, err
	}
	return obj.(*authenticationnexusorgv1.OIDC), err
}

// List takes label and field selectors, and returns the list of OIDCs that match those selectors.
func (c *FakeOIDCs) List(ctx context.Context, opts v1.ListOptions) (result *authenticationnexusorgv1.OIDCList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(oidcsResource, oidcsKind, opts), &authenticationnexusorgv1.OIDCList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &authenticationnexusorgv1.OIDCList{ListMeta: obj.(*authenticationnexusorgv1.OIDCList).ListMeta}
	for _, item := range obj.(*authenticationnexusorgv1.OIDCList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested oIDCs.
func (c *FakeOIDCs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(oidcsResource, opts))
}

// Create takes the representation of a oIDC and creates it.  Returns the server's representation of the oIDC, and an error, if there is any.
func (c *FakeOIDCs) Create(ctx context.Context, oIDC *authenticationnexusorgv1.OIDC, opts v1.CreateOptions) (result *authenticationnexusorgv1.OIDC, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(oidcsResource, oIDC), &authenticationnexusorgv1.OIDC{})
	if obj == nil {
		return nil, err
	}
	return obj.(*authenticationnexusorgv1.OIDC), err
}

// Update takes the representation of a oIDC and updates it. Returns the server's representation of the oIDC, and an error, if there is any.
func (c *FakeOIDCs) Update(ctx context.Context, oIDC *authenticationnexusorgv1.OIDC, opts v1.UpdateOptions) (result *authenticationnexusorgv1.OIDC, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(oidcsResource, oIDC), &authenticationnexusorgv1.OIDC{})
	if obj == nil {
		return nil, err
	}
	return obj.(*authenticationnexusorgv1.OIDC), err
}

// Delete takes name of the oIDC and deletes it. Returns an error if one occurs.
func (c *FakeOIDCs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(oidcsResource, name), &authenticationnexusorgv1.OIDC{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeOIDCs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(oidcsResource, listOpts)

	_, err := c.Fake.Invokes(action, &authenticationnexusorgv1.OIDCList{})
	return err
}

// Patch applies the patch and returns the patched oIDC.
func (c *FakeOIDCs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *authenticationnexusorgv1.OIDC, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(oidcsResource, name, pt, data, subresources...), &authenticationnexusorgv1.OIDC{})
	if obj == nil {
		return nil, err
	}
	return obj.(*authenticationnexusorgv1.OIDC), err
}

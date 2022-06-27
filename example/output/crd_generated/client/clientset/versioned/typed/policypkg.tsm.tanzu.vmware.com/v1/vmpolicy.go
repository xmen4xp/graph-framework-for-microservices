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

package v1

import (
	"context"
	"time"

	v1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/crd_generated/apis/policypkg.tsm.tanzu.vmware.com/v1"
	scheme "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/crd_generated/client/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// VMpoliciesGetter has a method to return a VMpolicyInterface.
// A group's client should implement this interface.
type VMpoliciesGetter interface {
	VMpolicies() VMpolicyInterface
}

// VMpolicyInterface has methods to work with VMpolicy resources.
type VMpolicyInterface interface {
	Create(ctx context.Context, vMpolicy *v1.VMpolicy, opts metav1.CreateOptions) (*v1.VMpolicy, error)
	Update(ctx context.Context, vMpolicy *v1.VMpolicy, opts metav1.UpdateOptions) (*v1.VMpolicy, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.VMpolicy, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.VMpolicyList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.VMpolicy, err error)
	VMpolicyExpansion
}

// vMpolicies implements VMpolicyInterface
type vMpolicies struct {
	client rest.Interface
}

// newVMpolicies returns a VMpolicies
func newVMpolicies(c *PolicypkgTsmV1Client) *vMpolicies {
	return &vMpolicies{
		client: c.RESTClient(),
	}
}

// Get takes name of the vMpolicy, and returns the corresponding vMpolicy object, and an error if there is any.
func (c *vMpolicies) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.VMpolicy, err error) {
	result = &v1.VMpolicy{}
	err = c.client.Get().
		Resource("vmpolicies").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of VMpolicies that match those selectors.
func (c *vMpolicies) List(ctx context.Context, opts metav1.ListOptions) (result *v1.VMpolicyList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.VMpolicyList{}
	err = c.client.Get().
		Resource("vmpolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested vMpolicies.
func (c *vMpolicies) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("vmpolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a vMpolicy and creates it.  Returns the server's representation of the vMpolicy, and an error, if there is any.
func (c *vMpolicies) Create(ctx context.Context, vMpolicy *v1.VMpolicy, opts metav1.CreateOptions) (result *v1.VMpolicy, err error) {
	result = &v1.VMpolicy{}
	err = c.client.Post().
		Resource("vmpolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(vMpolicy).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a vMpolicy and updates it. Returns the server's representation of the vMpolicy, and an error, if there is any.
func (c *vMpolicies) Update(ctx context.Context, vMpolicy *v1.VMpolicy, opts metav1.UpdateOptions) (result *v1.VMpolicy, err error) {
	result = &v1.VMpolicy{}
	err = c.client.Put().
		Resource("vmpolicies").
		Name(vMpolicy.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(vMpolicy).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the vMpolicy and deletes it. Returns an error if one occurs.
func (c *vMpolicies) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("vmpolicies").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *vMpolicies) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("vmpolicies").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched vMpolicy.
func (c *vMpolicies) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.VMpolicy, err error) {
	result = &v1.VMpolicy{}
	err = c.client.Patch(pt).
		Resource("vmpolicies").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
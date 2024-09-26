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
	v1 "/build/apis/gns.tsm.tanzu.vmware.com/v1"
	scheme "/build/client/clientset/versioned/scheme"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// IgnoreChildsGetter has a method to return a IgnoreChildInterface.
// A group's client should implement this interface.
type IgnoreChildsGetter interface {
	IgnoreChilds() IgnoreChildInterface
}

// IgnoreChildInterface has methods to work with IgnoreChild resources.
type IgnoreChildInterface interface {
	Create(ctx context.Context, ignoreChild *v1.IgnoreChild, opts metav1.CreateOptions) (*v1.IgnoreChild, error)
	Update(ctx context.Context, ignoreChild *v1.IgnoreChild, opts metav1.UpdateOptions) (*v1.IgnoreChild, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.IgnoreChild, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.IgnoreChildList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.IgnoreChild, err error)
	IgnoreChildExpansion
}

// ignoreChilds implements IgnoreChildInterface
type ignoreChilds struct {
	client rest.Interface
}

// newIgnoreChilds returns a IgnoreChilds
func newIgnoreChilds(c *GnsTsmV1Client) *ignoreChilds {
	return &ignoreChilds{
		client: c.RESTClient(),
	}
}

// Get takes name of the ignoreChild, and returns the corresponding ignoreChild object, and an error if there is any.
func (c *ignoreChilds) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.IgnoreChild, err error) {
	result = &v1.IgnoreChild{}
	err = c.client.Get().
		Resource("ignorechilds").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of IgnoreChilds that match those selectors.
func (c *ignoreChilds) List(ctx context.Context, opts metav1.ListOptions) (result *v1.IgnoreChildList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.IgnoreChildList{}
	err = c.client.Get().
		Resource("ignorechilds").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested ignoreChilds.
func (c *ignoreChilds) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("ignorechilds").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a ignoreChild and creates it.  Returns the server's representation of the ignoreChild, and an error, if there is any.
func (c *ignoreChilds) Create(ctx context.Context, ignoreChild *v1.IgnoreChild, opts metav1.CreateOptions) (result *v1.IgnoreChild, err error) {
	result = &v1.IgnoreChild{}
	err = c.client.Post().
		Resource("ignorechilds").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(ignoreChild).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a ignoreChild and updates it. Returns the server's representation of the ignoreChild, and an error, if there is any.
func (c *ignoreChilds) Update(ctx context.Context, ignoreChild *v1.IgnoreChild, opts metav1.UpdateOptions) (result *v1.IgnoreChild, err error) {
	result = &v1.IgnoreChild{}
	err = c.client.Put().
		Resource("ignorechilds").
		Name(ignoreChild.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(ignoreChild).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the ignoreChild and deletes it. Returns an error if one occurs.
func (c *ignoreChilds) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("ignorechilds").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *ignoreChilds) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("ignorechilds").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched ignoreChild.
func (c *ignoreChilds) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.IgnoreChild, err error) {
	result = &v1.IgnoreChild{}
	err = c.client.Patch(pt).
		Resource("ignorechilds").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

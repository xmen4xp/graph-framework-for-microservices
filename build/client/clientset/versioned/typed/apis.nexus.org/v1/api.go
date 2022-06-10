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

	v1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/apis.nexus.org/v1"
	scheme "golang-appnet.eng.vmware.com/nexus-sdk/api/build/client/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ApisGetter has a method to return a ApiInterface.
// A group's client should implement this interface.
type ApisGetter interface {
	Apis() ApiInterface
}

// ApiInterface has methods to work with Api resources.
type ApiInterface interface {
	Create(ctx context.Context, api *v1.Api, opts metav1.CreateOptions) (*v1.Api, error)
	Update(ctx context.Context, api *v1.Api, opts metav1.UpdateOptions) (*v1.Api, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Api, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.ApiList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Api, err error)
	ApiExpansion
}

// apis implements ApiInterface
type apis struct {
	client rest.Interface
}

// newApis returns a Apis
func newApis(c *ApisNexusV1Client) *apis {
	return &apis{
		client: c.RESTClient(),
	}
}

// Get takes name of the api, and returns the corresponding api object, and an error if there is any.
func (c *apis) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.Api, err error) {
	result = &v1.Api{}
	err = c.client.Get().
		Resource("apis").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Apis that match those selectors.
func (c *apis) List(ctx context.Context, opts metav1.ListOptions) (result *v1.ApiList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ApiList{}
	err = c.client.Get().
		Resource("apis").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested apis.
func (c *apis) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("apis").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a api and creates it.  Returns the server's representation of the api, and an error, if there is any.
func (c *apis) Create(ctx context.Context, api *v1.Api, opts metav1.CreateOptions) (result *v1.Api, err error) {
	result = &v1.Api{}
	err = c.client.Post().
		Resource("apis").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(api).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a api and updates it. Returns the server's representation of the api, and an error, if there is any.
func (c *apis) Update(ctx context.Context, api *v1.Api, opts metav1.UpdateOptions) (result *v1.Api, err error) {
	result = &v1.Api{}
	err = c.client.Put().
		Resource("apis").
		Name(api.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(api).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the api and deletes it. Returns an error if one occurs.
func (c *apis) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("apis").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *apis) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("apis").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched api.
func (c *apis) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Api, err error) {
	result = &v1.Api{}
	err = c.client.Patch(pt).
		Resource("apis").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

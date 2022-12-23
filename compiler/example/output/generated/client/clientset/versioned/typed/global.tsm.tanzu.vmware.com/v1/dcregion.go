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
	v1 "nexustempmodule/apis/global.tsm.tanzu.vmware.com/v1"
	scheme "nexustempmodule/client/clientset/versioned/scheme"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// DCRegionsGetter has a method to return a DCRegionInterface.
// A group's client should implement this interface.
type DCRegionsGetter interface {
	DCRegions() DCRegionInterface
}

// DCRegionInterface has methods to work with DCRegion resources.
type DCRegionInterface interface {
	Create(ctx context.Context, dCRegion *v1.DCRegion, opts metav1.CreateOptions) (*v1.DCRegion, error)
	Update(ctx context.Context, dCRegion *v1.DCRegion, opts metav1.UpdateOptions) (*v1.DCRegion, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.DCRegion, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.DCRegionList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.DCRegion, err error)
	DCRegionExpansion
}

// dCRegions implements DCRegionInterface
type dCRegions struct {
	client rest.Interface
}

// newDCRegions returns a DCRegions
func newDCRegions(c *GlobalTsmV1Client) *dCRegions {
	return &dCRegions{
		client: c.RESTClient(),
	}
}

// Get takes name of the dCRegion, and returns the corresponding dCRegion object, and an error if there is any.
func (c *dCRegions) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.DCRegion, err error) {
	result = &v1.DCRegion{}
	err = c.client.Get().
		Resource("dcregions").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of DCRegions that match those selectors.
func (c *dCRegions) List(ctx context.Context, opts metav1.ListOptions) (result *v1.DCRegionList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.DCRegionList{}
	err = c.client.Get().
		Resource("dcregions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested dCRegions.
func (c *dCRegions) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("dcregions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a dCRegion and creates it.  Returns the server's representation of the dCRegion, and an error, if there is any.
func (c *dCRegions) Create(ctx context.Context, dCRegion *v1.DCRegion, opts metav1.CreateOptions) (result *v1.DCRegion, err error) {
	result = &v1.DCRegion{}
	err = c.client.Post().
		Resource("dcregions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(dCRegion).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a dCRegion and updates it. Returns the server's representation of the dCRegion, and an error, if there is any.
func (c *dCRegions) Update(ctx context.Context, dCRegion *v1.DCRegion, opts metav1.UpdateOptions) (result *v1.DCRegion, err error) {
	result = &v1.DCRegion{}
	err = c.client.Put().
		Resource("dcregions").
		Name(dCRegion.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(dCRegion).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the dCRegion and deletes it. Returns an error if one occurs.
func (c *dCRegions) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("dcregions").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *dCRegions) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("dcregions").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched dCRegion.
func (c *dCRegions) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.DCRegion, err error) {
	result = &v1.DCRegion{}
	err = c.client.Patch(pt).
		Resource("dcregions").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

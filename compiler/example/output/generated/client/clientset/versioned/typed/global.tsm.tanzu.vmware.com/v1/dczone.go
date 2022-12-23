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

// DCZonesGetter has a method to return a DCZoneInterface.
// A group's client should implement this interface.
type DCZonesGetter interface {
	DCZones() DCZoneInterface
}

// DCZoneInterface has methods to work with DCZone resources.
type DCZoneInterface interface {
	Create(ctx context.Context, dCZone *v1.DCZone, opts metav1.CreateOptions) (*v1.DCZone, error)
	Update(ctx context.Context, dCZone *v1.DCZone, opts metav1.UpdateOptions) (*v1.DCZone, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.DCZone, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.DCZoneList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.DCZone, err error)
	DCZoneExpansion
}

// dCZones implements DCZoneInterface
type dCZones struct {
	client rest.Interface
}

// newDCZones returns a DCZones
func newDCZones(c *GlobalTsmV1Client) *dCZones {
	return &dCZones{
		client: c.RESTClient(),
	}
}

// Get takes name of the dCZone, and returns the corresponding dCZone object, and an error if there is any.
func (c *dCZones) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.DCZone, err error) {
	result = &v1.DCZone{}
	err = c.client.Get().
		Resource("dczones").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of DCZones that match those selectors.
func (c *dCZones) List(ctx context.Context, opts metav1.ListOptions) (result *v1.DCZoneList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.DCZoneList{}
	err = c.client.Get().
		Resource("dczones").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested dCZones.
func (c *dCZones) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("dczones").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a dCZone and creates it.  Returns the server's representation of the dCZone, and an error, if there is any.
func (c *dCZones) Create(ctx context.Context, dCZone *v1.DCZone, opts metav1.CreateOptions) (result *v1.DCZone, err error) {
	result = &v1.DCZone{}
	err = c.client.Post().
		Resource("dczones").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(dCZone).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a dCZone and updates it. Returns the server's representation of the dCZone, and an error, if there is any.
func (c *dCZones) Update(ctx context.Context, dCZone *v1.DCZone, opts metav1.UpdateOptions) (result *v1.DCZone, err error) {
	result = &v1.DCZone{}
	err = c.client.Put().
		Resource("dczones").
		Name(dCZone.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(dCZone).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the dCZone and deletes it. Returns an error if one occurs.
func (c *dCZones) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("dczones").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *dCZones) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("dczones").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched dCZone.
func (c *dCZones) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.DCZone, err error) {
	result = &v1.DCZone{}
	err = c.client.Patch(pt).
		Resource("dczones").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

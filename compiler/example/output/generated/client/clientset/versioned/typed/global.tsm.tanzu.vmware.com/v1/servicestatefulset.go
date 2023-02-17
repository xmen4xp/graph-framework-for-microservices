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

// ServiceStatefulSetsGetter has a method to return a ServiceStatefulSetInterface.
// A group's client should implement this interface.
type ServiceStatefulSetsGetter interface {
	ServiceStatefulSets() ServiceStatefulSetInterface
}

// ServiceStatefulSetInterface has methods to work with ServiceStatefulSet resources.
type ServiceStatefulSetInterface interface {
	Create(ctx context.Context, serviceStatefulSet *v1.ServiceStatefulSet, opts metav1.CreateOptions) (*v1.ServiceStatefulSet, error)
	Update(ctx context.Context, serviceStatefulSet *v1.ServiceStatefulSet, opts metav1.UpdateOptions) (*v1.ServiceStatefulSet, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.ServiceStatefulSet, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.ServiceStatefulSetList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ServiceStatefulSet, err error)
	ServiceStatefulSetExpansion
}

// serviceStatefulSets implements ServiceStatefulSetInterface
type serviceStatefulSets struct {
	client rest.Interface
}

// newServiceStatefulSets returns a ServiceStatefulSets
func newServiceStatefulSets(c *GlobalTsmV1Client) *serviceStatefulSets {
	return &serviceStatefulSets{
		client: c.RESTClient(),
	}
}

// Get takes name of the serviceStatefulSet, and returns the corresponding serviceStatefulSet object, and an error if there is any.
func (c *serviceStatefulSets) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.ServiceStatefulSet, err error) {
	result = &v1.ServiceStatefulSet{}
	err = c.client.Get().
		Resource("servicestatefulsets").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ServiceStatefulSets that match those selectors.
func (c *serviceStatefulSets) List(ctx context.Context, opts metav1.ListOptions) (result *v1.ServiceStatefulSetList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ServiceStatefulSetList{}
	err = c.client.Get().
		Resource("servicestatefulsets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested serviceStatefulSets.
func (c *serviceStatefulSets) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("servicestatefulsets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a serviceStatefulSet and creates it.  Returns the server's representation of the serviceStatefulSet, and an error, if there is any.
func (c *serviceStatefulSets) Create(ctx context.Context, serviceStatefulSet *v1.ServiceStatefulSet, opts metav1.CreateOptions) (result *v1.ServiceStatefulSet, err error) {
	result = &v1.ServiceStatefulSet{}
	err = c.client.Post().
		Resource("servicestatefulsets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(serviceStatefulSet).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a serviceStatefulSet and updates it. Returns the server's representation of the serviceStatefulSet, and an error, if there is any.
func (c *serviceStatefulSets) Update(ctx context.Context, serviceStatefulSet *v1.ServiceStatefulSet, opts metav1.UpdateOptions) (result *v1.ServiceStatefulSet, err error) {
	result = &v1.ServiceStatefulSet{}
	err = c.client.Put().
		Resource("servicestatefulsets").
		Name(serviceStatefulSet.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(serviceStatefulSet).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the serviceStatefulSet and deletes it. Returns an error if one occurs.
func (c *serviceStatefulSets) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("servicestatefulsets").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *serviceStatefulSets) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("servicestatefulsets").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched serviceStatefulSet.
func (c *serviceStatefulSets) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ServiceStatefulSet, err error) {
	result = &v1.ServiceStatefulSet{}
	err = c.client.Patch(pt).
		Resource("servicestatefulsets").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
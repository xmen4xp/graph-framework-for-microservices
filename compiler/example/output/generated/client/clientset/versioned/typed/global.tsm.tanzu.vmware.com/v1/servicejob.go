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

// ServiceJobsGetter has a method to return a ServiceJobInterface.
// A group's client should implement this interface.
type ServiceJobsGetter interface {
	ServiceJobs() ServiceJobInterface
}

// ServiceJobInterface has methods to work with ServiceJob resources.
type ServiceJobInterface interface {
	Create(ctx context.Context, serviceJob *v1.ServiceJob, opts metav1.CreateOptions) (*v1.ServiceJob, error)
	Update(ctx context.Context, serviceJob *v1.ServiceJob, opts metav1.UpdateOptions) (*v1.ServiceJob, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.ServiceJob, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.ServiceJobList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ServiceJob, err error)
	ServiceJobExpansion
}

// serviceJobs implements ServiceJobInterface
type serviceJobs struct {
	client rest.Interface
}

// newServiceJobs returns a ServiceJobs
func newServiceJobs(c *GlobalTsmV1Client) *serviceJobs {
	return &serviceJobs{
		client: c.RESTClient(),
	}
}

// Get takes name of the serviceJob, and returns the corresponding serviceJob object, and an error if there is any.
func (c *serviceJobs) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.ServiceJob, err error) {
	result = &v1.ServiceJob{}
	err = c.client.Get().
		Resource("servicejobs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ServiceJobs that match those selectors.
func (c *serviceJobs) List(ctx context.Context, opts metav1.ListOptions) (result *v1.ServiceJobList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ServiceJobList{}
	err = c.client.Get().
		Resource("servicejobs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested serviceJobs.
func (c *serviceJobs) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("servicejobs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a serviceJob and creates it.  Returns the server's representation of the serviceJob, and an error, if there is any.
func (c *serviceJobs) Create(ctx context.Context, serviceJob *v1.ServiceJob, opts metav1.CreateOptions) (result *v1.ServiceJob, err error) {
	result = &v1.ServiceJob{}
	err = c.client.Post().
		Resource("servicejobs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(serviceJob).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a serviceJob and updates it. Returns the server's representation of the serviceJob, and an error, if there is any.
func (c *serviceJobs) Update(ctx context.Context, serviceJob *v1.ServiceJob, opts metav1.UpdateOptions) (result *v1.ServiceJob, err error) {
	result = &v1.ServiceJob{}
	err = c.client.Put().
		Resource("servicejobs").
		Name(serviceJob.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(serviceJob).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the serviceJob and deletes it. Returns an error if one occurs.
func (c *serviceJobs) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("servicejobs").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *serviceJobs) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("servicejobs").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched serviceJob.
func (c *serviceJobs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ServiceJob, err error) {
	result = &v1.ServiceJob{}
	err = c.client.Patch(pt).
		Resource("servicejobs").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

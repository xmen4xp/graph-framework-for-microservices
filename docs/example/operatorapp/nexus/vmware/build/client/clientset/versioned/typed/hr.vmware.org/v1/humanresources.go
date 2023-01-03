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
	v1 "vmware/build/apis/hr.vmware.org/v1"
	scheme "vmware/build/client/clientset/versioned/scheme"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// HumanResourcesesGetter has a method to return a HumanResourcesInterface.
// A group's client should implement this interface.
type HumanResourcesesGetter interface {
	HumanResourceses() HumanResourcesInterface
}

// HumanResourcesInterface has methods to work with HumanResources resources.
type HumanResourcesInterface interface {
	Create(ctx context.Context, humanResources *v1.HumanResources, opts metav1.CreateOptions) (*v1.HumanResources, error)
	Update(ctx context.Context, humanResources *v1.HumanResources, opts metav1.UpdateOptions) (*v1.HumanResources, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.HumanResources, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.HumanResourcesList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.HumanResources, err error)
	HumanResourcesExpansion
}

// humanResourceses implements HumanResourcesInterface
type humanResourceses struct {
	client rest.Interface
}

// newHumanResourceses returns a HumanResourceses
func newHumanResourceses(c *HrVmwareV1Client) *humanResourceses {
	return &humanResourceses{
		client: c.RESTClient(),
	}
}

// Get takes name of the humanResources, and returns the corresponding humanResources object, and an error if there is any.
func (c *humanResourceses) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.HumanResources, err error) {
	result = &v1.HumanResources{}
	err = c.client.Get().
		Resource("humanresourceses").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of HumanResourceses that match those selectors.
func (c *humanResourceses) List(ctx context.Context, opts metav1.ListOptions) (result *v1.HumanResourcesList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.HumanResourcesList{}
	err = c.client.Get().
		Resource("humanresourceses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested humanResourceses.
func (c *humanResourceses) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("humanresourceses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a humanResources and creates it.  Returns the server's representation of the humanResources, and an error, if there is any.
func (c *humanResourceses) Create(ctx context.Context, humanResources *v1.HumanResources, opts metav1.CreateOptions) (result *v1.HumanResources, err error) {
	result = &v1.HumanResources{}
	err = c.client.Post().
		Resource("humanresourceses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(humanResources).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a humanResources and updates it. Returns the server's representation of the humanResources, and an error, if there is any.
func (c *humanResourceses) Update(ctx context.Context, humanResources *v1.HumanResources, opts metav1.UpdateOptions) (result *v1.HumanResources, err error) {
	result = &v1.HumanResources{}
	err = c.client.Put().
		Resource("humanresourceses").
		Name(humanResources.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(humanResources).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the humanResources and deletes it. Returns an error if one occurs.
func (c *humanResourceses) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("humanresourceses").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *humanResourceses) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("humanresourceses").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched humanResources.
func (c *humanResourceses) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.HumanResources, err error) {
	result = &v1.HumanResources{}
	err = c.client.Patch(pt).
		Resource("humanresourceses").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
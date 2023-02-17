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

// SloConfigsGetter has a method to return a SloConfigInterface.
// A group's client should implement this interface.
type SloConfigsGetter interface {
	SloConfigs() SloConfigInterface
}

// SloConfigInterface has methods to work with SloConfig resources.
type SloConfigInterface interface {
	Create(ctx context.Context, sloConfig *v1.SloConfig, opts metav1.CreateOptions) (*v1.SloConfig, error)
	Update(ctx context.Context, sloConfig *v1.SloConfig, opts metav1.UpdateOptions) (*v1.SloConfig, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.SloConfig, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.SloConfigList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.SloConfig, err error)
	SloConfigExpansion
}

// sloConfigs implements SloConfigInterface
type sloConfigs struct {
	client rest.Interface
}

// newSloConfigs returns a SloConfigs
func newSloConfigs(c *GlobalTsmV1Client) *sloConfigs {
	return &sloConfigs{
		client: c.RESTClient(),
	}
}

// Get takes name of the sloConfig, and returns the corresponding sloConfig object, and an error if there is any.
func (c *sloConfigs) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.SloConfig, err error) {
	result = &v1.SloConfig{}
	err = c.client.Get().
		Resource("sloconfigs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of SloConfigs that match those selectors.
func (c *sloConfigs) List(ctx context.Context, opts metav1.ListOptions) (result *v1.SloConfigList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.SloConfigList{}
	err = c.client.Get().
		Resource("sloconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested sloConfigs.
func (c *sloConfigs) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("sloconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a sloConfig and creates it.  Returns the server's representation of the sloConfig, and an error, if there is any.
func (c *sloConfigs) Create(ctx context.Context, sloConfig *v1.SloConfig, opts metav1.CreateOptions) (result *v1.SloConfig, err error) {
	result = &v1.SloConfig{}
	err = c.client.Post().
		Resource("sloconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(sloConfig).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a sloConfig and updates it. Returns the server's representation of the sloConfig, and an error, if there is any.
func (c *sloConfigs) Update(ctx context.Context, sloConfig *v1.SloConfig, opts metav1.UpdateOptions) (result *v1.SloConfig, err error) {
	result = &v1.SloConfig{}
	err = c.client.Put().
		Resource("sloconfigs").
		Name(sloConfig.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(sloConfig).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the sloConfig and deletes it. Returns an error if one occurs.
func (c *sloConfigs) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("sloconfigs").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *sloConfigs) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("sloconfigs").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched sloConfig.
func (c *sloConfigs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.SloConfig, err error) {
	result = &v1.SloConfig{}
	err = c.client.Patch(pt).
		Resource("sloconfigs").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
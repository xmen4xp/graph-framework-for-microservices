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

// UserDiscoveryRTsGetter has a method to return a UserDiscoveryRTInterface.
// A group's client should implement this interface.
type UserDiscoveryRTsGetter interface {
	UserDiscoveryRTs() UserDiscoveryRTInterface
}

// UserDiscoveryRTInterface has methods to work with UserDiscoveryRT resources.
type UserDiscoveryRTInterface interface {
	Create(ctx context.Context, userDiscoveryRT *v1.UserDiscoveryRT, opts metav1.CreateOptions) (*v1.UserDiscoveryRT, error)
	Update(ctx context.Context, userDiscoveryRT *v1.UserDiscoveryRT, opts metav1.UpdateOptions) (*v1.UserDiscoveryRT, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.UserDiscoveryRT, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.UserDiscoveryRTList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.UserDiscoveryRT, err error)
	UserDiscoveryRTExpansion
}

// userDiscoveryRTs implements UserDiscoveryRTInterface
type userDiscoveryRTs struct {
	client rest.Interface
}

// newUserDiscoveryRTs returns a UserDiscoveryRTs
func newUserDiscoveryRTs(c *GlobalTsmV1Client) *userDiscoveryRTs {
	return &userDiscoveryRTs{
		client: c.RESTClient(),
	}
}

// Get takes name of the userDiscoveryRT, and returns the corresponding userDiscoveryRT object, and an error if there is any.
func (c *userDiscoveryRTs) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.UserDiscoveryRT, err error) {
	result = &v1.UserDiscoveryRT{}
	err = c.client.Get().
		Resource("userdiscoveryrts").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of UserDiscoveryRTs that match those selectors.
func (c *userDiscoveryRTs) List(ctx context.Context, opts metav1.ListOptions) (result *v1.UserDiscoveryRTList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.UserDiscoveryRTList{}
	err = c.client.Get().
		Resource("userdiscoveryrts").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested userDiscoveryRTs.
func (c *userDiscoveryRTs) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("userdiscoveryrts").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a userDiscoveryRT and creates it.  Returns the server's representation of the userDiscoveryRT, and an error, if there is any.
func (c *userDiscoveryRTs) Create(ctx context.Context, userDiscoveryRT *v1.UserDiscoveryRT, opts metav1.CreateOptions) (result *v1.UserDiscoveryRT, err error) {
	result = &v1.UserDiscoveryRT{}
	err = c.client.Post().
		Resource("userdiscoveryrts").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(userDiscoveryRT).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a userDiscoveryRT and updates it. Returns the server's representation of the userDiscoveryRT, and an error, if there is any.
func (c *userDiscoveryRTs) Update(ctx context.Context, userDiscoveryRT *v1.UserDiscoveryRT, opts metav1.UpdateOptions) (result *v1.UserDiscoveryRT, err error) {
	result = &v1.UserDiscoveryRT{}
	err = c.client.Put().
		Resource("userdiscoveryrts").
		Name(userDiscoveryRT.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(userDiscoveryRT).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the userDiscoveryRT and deletes it. Returns an error if one occurs.
func (c *userDiscoveryRTs) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("userdiscoveryrts").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *userDiscoveryRTs) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("userdiscoveryrts").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched userDiscoveryRT.
func (c *userDiscoveryRTs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.UserDiscoveryRT, err error) {
	result = &v1.UserDiscoveryRT{}
	err = c.client.Patch(pt).
		Resource("userdiscoveryrts").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

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

// SLOFoldersGetter has a method to return a SLOFolderInterface.
// A group's client should implement this interface.
type SLOFoldersGetter interface {
	SLOFolders() SLOFolderInterface
}

// SLOFolderInterface has methods to work with SLOFolder resources.
type SLOFolderInterface interface {
	Create(ctx context.Context, sLOFolder *v1.SLOFolder, opts metav1.CreateOptions) (*v1.SLOFolder, error)
	Update(ctx context.Context, sLOFolder *v1.SLOFolder, opts metav1.UpdateOptions) (*v1.SLOFolder, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.SLOFolder, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.SLOFolderList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.SLOFolder, err error)
	SLOFolderExpansion
}

// sLOFolders implements SLOFolderInterface
type sLOFolders struct {
	client rest.Interface
}

// newSLOFolders returns a SLOFolders
func newSLOFolders(c *GlobalTsmV1Client) *sLOFolders {
	return &sLOFolders{
		client: c.RESTClient(),
	}
}

// Get takes name of the sLOFolder, and returns the corresponding sLOFolder object, and an error if there is any.
func (c *sLOFolders) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.SLOFolder, err error) {
	result = &v1.SLOFolder{}
	err = c.client.Get().
		Resource("slofolders").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of SLOFolders that match those selectors.
func (c *sLOFolders) List(ctx context.Context, opts metav1.ListOptions) (result *v1.SLOFolderList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.SLOFolderList{}
	err = c.client.Get().
		Resource("slofolders").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested sLOFolders.
func (c *sLOFolders) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("slofolders").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a sLOFolder and creates it.  Returns the server's representation of the sLOFolder, and an error, if there is any.
func (c *sLOFolders) Create(ctx context.Context, sLOFolder *v1.SLOFolder, opts metav1.CreateOptions) (result *v1.SLOFolder, err error) {
	result = &v1.SLOFolder{}
	err = c.client.Post().
		Resource("slofolders").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(sLOFolder).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a sLOFolder and updates it. Returns the server's representation of the sLOFolder, and an error, if there is any.
func (c *sLOFolders) Update(ctx context.Context, sLOFolder *v1.SLOFolder, opts metav1.UpdateOptions) (result *v1.SLOFolder, err error) {
	result = &v1.SLOFolder{}
	err = c.client.Put().
		Resource("slofolders").
		Name(sLOFolder.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(sLOFolder).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the sLOFolder and deletes it. Returns an error if one occurs.
func (c *sLOFolders) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("slofolders").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *sLOFolders) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("slofolders").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched sLOFolder.
func (c *sLOFolders) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.SLOFolder, err error) {
	result = &v1.SLOFolder{}
	err = c.client.Patch(pt).
		Resource("slofolders").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

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

// NodeFoldersGetter has a method to return a NodeFolderInterface.
// A group's client should implement this interface.
type NodeFoldersGetter interface {
	NodeFolders() NodeFolderInterface
}

// NodeFolderInterface has methods to work with NodeFolder resources.
type NodeFolderInterface interface {
	Create(ctx context.Context, nodeFolder *v1.NodeFolder, opts metav1.CreateOptions) (*v1.NodeFolder, error)
	Update(ctx context.Context, nodeFolder *v1.NodeFolder, opts metav1.UpdateOptions) (*v1.NodeFolder, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.NodeFolder, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.NodeFolderList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.NodeFolder, err error)
	NodeFolderExpansion
}

// nodeFolders implements NodeFolderInterface
type nodeFolders struct {
	client rest.Interface
}

// newNodeFolders returns a NodeFolders
func newNodeFolders(c *GlobalTsmV1Client) *nodeFolders {
	return &nodeFolders{
		client: c.RESTClient(),
	}
}

// Get takes name of the nodeFolder, and returns the corresponding nodeFolder object, and an error if there is any.
func (c *nodeFolders) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.NodeFolder, err error) {
	result = &v1.NodeFolder{}
	err = c.client.Get().
		Resource("nodefolders").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of NodeFolders that match those selectors.
func (c *nodeFolders) List(ctx context.Context, opts metav1.ListOptions) (result *v1.NodeFolderList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.NodeFolderList{}
	err = c.client.Get().
		Resource("nodefolders").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested nodeFolders.
func (c *nodeFolders) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("nodefolders").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a nodeFolder and creates it.  Returns the server's representation of the nodeFolder, and an error, if there is any.
func (c *nodeFolders) Create(ctx context.Context, nodeFolder *v1.NodeFolder, opts metav1.CreateOptions) (result *v1.NodeFolder, err error) {
	result = &v1.NodeFolder{}
	err = c.client.Post().
		Resource("nodefolders").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(nodeFolder).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a nodeFolder and updates it. Returns the server's representation of the nodeFolder, and an error, if there is any.
func (c *nodeFolders) Update(ctx context.Context, nodeFolder *v1.NodeFolder, opts metav1.UpdateOptions) (result *v1.NodeFolder, err error) {
	result = &v1.NodeFolder{}
	err = c.client.Put().
		Resource("nodefolders").
		Name(nodeFolder.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(nodeFolder).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the nodeFolder and deletes it. Returns an error if one occurs.
func (c *nodeFolders) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("nodefolders").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *nodeFolders) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("nodefolders").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched nodeFolder.
func (c *nodeFolders) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.NodeFolder, err error) {
	result = &v1.NodeFolder{}
	err = c.client.Patch(pt).
		Resource("nodefolders").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

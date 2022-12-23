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

// LocalRegistrationServiceClustersGetter has a method to return a LocalRegistrationServiceClusterInterface.
// A group's client should implement this interface.
type LocalRegistrationServiceClustersGetter interface {
	LocalRegistrationServiceClusters() LocalRegistrationServiceClusterInterface
}

// LocalRegistrationServiceClusterInterface has methods to work with LocalRegistrationServiceCluster resources.
type LocalRegistrationServiceClusterInterface interface {
	Create(ctx context.Context, localRegistrationServiceCluster *v1.LocalRegistrationServiceCluster, opts metav1.CreateOptions) (*v1.LocalRegistrationServiceCluster, error)
	Update(ctx context.Context, localRegistrationServiceCluster *v1.LocalRegistrationServiceCluster, opts metav1.UpdateOptions) (*v1.LocalRegistrationServiceCluster, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.LocalRegistrationServiceCluster, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.LocalRegistrationServiceClusterList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.LocalRegistrationServiceCluster, err error)
	LocalRegistrationServiceClusterExpansion
}

// localRegistrationServiceClusters implements LocalRegistrationServiceClusterInterface
type localRegistrationServiceClusters struct {
	client rest.Interface
}

// newLocalRegistrationServiceClusters returns a LocalRegistrationServiceClusters
func newLocalRegistrationServiceClusters(c *GlobalTsmV1Client) *localRegistrationServiceClusters {
	return &localRegistrationServiceClusters{
		client: c.RESTClient(),
	}
}

// Get takes name of the localRegistrationServiceCluster, and returns the corresponding localRegistrationServiceCluster object, and an error if there is any.
func (c *localRegistrationServiceClusters) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.LocalRegistrationServiceCluster, err error) {
	result = &v1.LocalRegistrationServiceCluster{}
	err = c.client.Get().
		Resource("localregistrationserviceclusters").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of LocalRegistrationServiceClusters that match those selectors.
func (c *localRegistrationServiceClusters) List(ctx context.Context, opts metav1.ListOptions) (result *v1.LocalRegistrationServiceClusterList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.LocalRegistrationServiceClusterList{}
	err = c.client.Get().
		Resource("localregistrationserviceclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested localRegistrationServiceClusters.
func (c *localRegistrationServiceClusters) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("localregistrationserviceclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a localRegistrationServiceCluster and creates it.  Returns the server's representation of the localRegistrationServiceCluster, and an error, if there is any.
func (c *localRegistrationServiceClusters) Create(ctx context.Context, localRegistrationServiceCluster *v1.LocalRegistrationServiceCluster, opts metav1.CreateOptions) (result *v1.LocalRegistrationServiceCluster, err error) {
	result = &v1.LocalRegistrationServiceCluster{}
	err = c.client.Post().
		Resource("localregistrationserviceclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(localRegistrationServiceCluster).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a localRegistrationServiceCluster and updates it. Returns the server's representation of the localRegistrationServiceCluster, and an error, if there is any.
func (c *localRegistrationServiceClusters) Update(ctx context.Context, localRegistrationServiceCluster *v1.LocalRegistrationServiceCluster, opts metav1.UpdateOptions) (result *v1.LocalRegistrationServiceCluster, err error) {
	result = &v1.LocalRegistrationServiceCluster{}
	err = c.client.Put().
		Resource("localregistrationserviceclusters").
		Name(localRegistrationServiceCluster.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(localRegistrationServiceCluster).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the localRegistrationServiceCluster and deletes it. Returns an error if one occurs.
func (c *localRegistrationServiceClusters) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("localregistrationserviceclusters").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *localRegistrationServiceClusters) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("localregistrationserviceclusters").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched localRegistrationServiceCluster.
func (c *localRegistrationServiceClusters) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.LocalRegistrationServiceCluster, err error) {
	result = &v1.LocalRegistrationServiceCluster{}
	err = c.client.Patch(pt).
		Resource("localregistrationserviceclusters").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

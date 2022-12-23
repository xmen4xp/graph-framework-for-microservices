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

// ProjectQueriesGetter has a method to return a ProjectQueryInterface.
// A group's client should implement this interface.
type ProjectQueriesGetter interface {
	ProjectQueries() ProjectQueryInterface
}

// ProjectQueryInterface has methods to work with ProjectQuery resources.
type ProjectQueryInterface interface {
	Create(ctx context.Context, projectQuery *v1.ProjectQuery, opts metav1.CreateOptions) (*v1.ProjectQuery, error)
	Update(ctx context.Context, projectQuery *v1.ProjectQuery, opts metav1.UpdateOptions) (*v1.ProjectQuery, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.ProjectQuery, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.ProjectQueryList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ProjectQuery, err error)
	ProjectQueryExpansion
}

// projectQueries implements ProjectQueryInterface
type projectQueries struct {
	client rest.Interface
}

// newProjectQueries returns a ProjectQueries
func newProjectQueries(c *GlobalTsmV1Client) *projectQueries {
	return &projectQueries{
		client: c.RESTClient(),
	}
}

// Get takes name of the projectQuery, and returns the corresponding projectQuery object, and an error if there is any.
func (c *projectQueries) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.ProjectQuery, err error) {
	result = &v1.ProjectQuery{}
	err = c.client.Get().
		Resource("projectqueries").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ProjectQueries that match those selectors.
func (c *projectQueries) List(ctx context.Context, opts metav1.ListOptions) (result *v1.ProjectQueryList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ProjectQueryList{}
	err = c.client.Get().
		Resource("projectqueries").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested projectQueries.
func (c *projectQueries) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("projectqueries").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a projectQuery and creates it.  Returns the server's representation of the projectQuery, and an error, if there is any.
func (c *projectQueries) Create(ctx context.Context, projectQuery *v1.ProjectQuery, opts metav1.CreateOptions) (result *v1.ProjectQuery, err error) {
	result = &v1.ProjectQuery{}
	err = c.client.Post().
		Resource("projectqueries").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(projectQuery).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a projectQuery and updates it. Returns the server's representation of the projectQuery, and an error, if there is any.
func (c *projectQueries) Update(ctx context.Context, projectQuery *v1.ProjectQuery, opts metav1.UpdateOptions) (result *v1.ProjectQuery, err error) {
	result = &v1.ProjectQuery{}
	err = c.client.Put().
		Resource("projectqueries").
		Name(projectQuery.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(projectQuery).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the projectQuery and deletes it. Returns an error if one occurs.
func (c *projectQueries) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("projectqueries").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *projectQueries) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("projectqueries").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched projectQuery.
func (c *projectQueries) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ProjectQuery, err error) {
	result = &v1.ProjectQuery{}
	err = c.client.Patch(pt).
		Resource("projectqueries").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

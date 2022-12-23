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

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "nexustempmodule/apis/global.tsm.tanzu.vmware.com/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ResourceGroupLister helps list ResourceGroups.
// All objects returned here must be treated as read-only.
type ResourceGroupLister interface {
	// List lists all ResourceGroups in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.ResourceGroup, err error)
	// Get retrieves the ResourceGroup from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.ResourceGroup, error)
	ResourceGroupListerExpansion
}

// resourceGroupLister implements the ResourceGroupLister interface.
type resourceGroupLister struct {
	indexer cache.Indexer
}

// NewResourceGroupLister returns a new ResourceGroupLister.
func NewResourceGroupLister(indexer cache.Indexer) ResourceGroupLister {
	return &resourceGroupLister{indexer: indexer}
}

// List lists all ResourceGroups in the indexer.
func (s *resourceGroupLister) List(selector labels.Selector) (ret []*v1.ResourceGroup, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ResourceGroup))
	})
	return ret, err
}

// Get retrieves the ResourceGroup from the index for a given name.
func (s *resourceGroupLister) Get(name string) (*v1.ResourceGroup, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("resourcegroup"), name)
	}
	return obj.(*v1.ResourceGroup), nil
}

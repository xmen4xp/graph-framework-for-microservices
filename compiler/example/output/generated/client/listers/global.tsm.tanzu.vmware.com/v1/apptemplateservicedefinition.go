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

// AppTemplateServiceDefinitionLister helps list AppTemplateServiceDefinitions.
// All objects returned here must be treated as read-only.
type AppTemplateServiceDefinitionLister interface {
	// List lists all AppTemplateServiceDefinitions in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.AppTemplateServiceDefinition, err error)
	// Get retrieves the AppTemplateServiceDefinition from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.AppTemplateServiceDefinition, error)
	AppTemplateServiceDefinitionListerExpansion
}

// appTemplateServiceDefinitionLister implements the AppTemplateServiceDefinitionLister interface.
type appTemplateServiceDefinitionLister struct {
	indexer cache.Indexer
}

// NewAppTemplateServiceDefinitionLister returns a new AppTemplateServiceDefinitionLister.
func NewAppTemplateServiceDefinitionLister(indexer cache.Indexer) AppTemplateServiceDefinitionLister {
	return &appTemplateServiceDefinitionLister{indexer: indexer}
}

// List lists all AppTemplateServiceDefinitions in the indexer.
func (s *appTemplateServiceDefinitionLister) List(selector labels.Selector) (ret []*v1.AppTemplateServiceDefinition, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.AppTemplateServiceDefinition))
	})
	return ret, err
}

// Get retrieves the AppTemplateServiceDefinition from the index for a given name.
func (s *appTemplateServiceDefinitionLister) Get(name string) (*v1.AppTemplateServiceDefinition, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("apptemplateservicedefinition"), name)
	}
	return obj.(*v1.AppTemplateServiceDefinition), nil
}
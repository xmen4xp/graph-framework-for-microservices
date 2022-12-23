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

// GnsRoutingRuleConfigLister helps list GnsRoutingRuleConfigs.
// All objects returned here must be treated as read-only.
type GnsRoutingRuleConfigLister interface {
	// List lists all GnsRoutingRuleConfigs in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.GnsRoutingRuleConfig, err error)
	// Get retrieves the GnsRoutingRuleConfig from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.GnsRoutingRuleConfig, error)
	GnsRoutingRuleConfigListerExpansion
}

// gnsRoutingRuleConfigLister implements the GnsRoutingRuleConfigLister interface.
type gnsRoutingRuleConfigLister struct {
	indexer cache.Indexer
}

// NewGnsRoutingRuleConfigLister returns a new GnsRoutingRuleConfigLister.
func NewGnsRoutingRuleConfigLister(indexer cache.Indexer) GnsRoutingRuleConfigLister {
	return &gnsRoutingRuleConfigLister{indexer: indexer}
}

// List lists all GnsRoutingRuleConfigs in the indexer.
func (s *gnsRoutingRuleConfigLister) List(selector labels.Selector) (ret []*v1.GnsRoutingRuleConfig, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.GnsRoutingRuleConfig))
	})
	return ret, err
}

// Get retrieves the GnsRoutingRuleConfig from the index for a given name.
func (s *gnsRoutingRuleConfigLister) Get(name string) (*v1.GnsRoutingRuleConfig, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("gnsroutingruleconfig"), name)
	}
	return obj.(*v1.GnsRoutingRuleConfig), nil
}

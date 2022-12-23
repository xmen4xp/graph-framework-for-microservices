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

// ProgressiveUpgradeLister helps list ProgressiveUpgrades.
// All objects returned here must be treated as read-only.
type ProgressiveUpgradeLister interface {
	// List lists all ProgressiveUpgrades in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.ProgressiveUpgrade, err error)
	// Get retrieves the ProgressiveUpgrade from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.ProgressiveUpgrade, error)
	ProgressiveUpgradeListerExpansion
}

// progressiveUpgradeLister implements the ProgressiveUpgradeLister interface.
type progressiveUpgradeLister struct {
	indexer cache.Indexer
}

// NewProgressiveUpgradeLister returns a new ProgressiveUpgradeLister.
func NewProgressiveUpgradeLister(indexer cache.Indexer) ProgressiveUpgradeLister {
	return &progressiveUpgradeLister{indexer: indexer}
}

// List lists all ProgressiveUpgrades in the indexer.
func (s *progressiveUpgradeLister) List(selector labels.Selector) (ret []*v1.ProgressiveUpgrade, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ProgressiveUpgrade))
	})
	return ret, err
}

// Get retrieves the ProgressiveUpgrade from the index for a given name.
func (s *progressiveUpgradeLister) Get(name string) (*v1.ProgressiveUpgrade, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("progressiveupgrade"), name)
	}
	return obj.(*v1.ProgressiveUpgrade), nil
}

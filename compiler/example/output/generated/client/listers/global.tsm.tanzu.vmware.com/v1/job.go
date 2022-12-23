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

// JobLister helps list Jobs.
// All objects returned here must be treated as read-only.
type JobLister interface {
	// List lists all Jobs in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Job, err error)
	// Get retrieves the Job from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.Job, error)
	JobListerExpansion
}

// jobLister implements the JobLister interface.
type jobLister struct {
	indexer cache.Indexer
}

// NewJobLister returns a new JobLister.
func NewJobLister(indexer cache.Indexer) JobLister {
	return &jobLister{indexer: indexer}
}

// List lists all Jobs in the indexer.
func (s *jobLister) List(selector labels.Selector) (ret []*v1.Job, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Job))
	})
	return ret, err
}

// Get retrieves the Job from the index for a given name.
func (s *jobLister) Get(name string) (*v1.Job, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("job"), name)
	}
	return obj.(*v1.Job), nil
}

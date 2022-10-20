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

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	"context"
	time "time"

	configtsmtanzuvmwarecomv1 "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/crd_generated/apis/config.tsm.tanzu.vmware.com/v1"
	versioned "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/crd_generated/client/clientset/versioned"
	internalinterfaces "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/crd_generated/client/informers/externalversions/internalinterfaces"
	v1 "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/crd_generated/client/listers/config.tsm.tanzu.vmware.com/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// DomainInformer provides access to a shared informer and lister for
// Domains.
type DomainInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.DomainLister
}

type domainInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewDomainInformer constructs a new informer for Domain type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewDomainInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredDomainInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredDomainInformer constructs a new informer for Domain type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredDomainInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ConfigTsmV1().Domains().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ConfigTsmV1().Domains().Watch(context.TODO(), options)
			},
		},
		&configtsmtanzuvmwarecomv1.Domain{},
		resyncPeriod,
		indexers,
	)
}

func (f *domainInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredDomainInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *domainInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&configtsmtanzuvmwarecomv1.Domain{}, f.defaultInformer)
}

func (f *domainInformer) Lister() v1.DomainLister {
	return v1.NewDomainLister(f.Informer().GetIndexer())
}
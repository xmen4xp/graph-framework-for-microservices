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
	managementvmwareorgv1 "vmware/build/apis/management.vmware.org/v1"
	versioned "vmware/build/client/clientset/versioned"
	internalinterfaces "vmware/build/client/informers/externalversions/internalinterfaces"
	v1 "vmware/build/client/listers/management.vmware.org/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// MgrInformer provides access to a shared informer and lister for
// Mgrs.
type MgrInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.MgrLister
}

type mgrInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewMgrInformer constructs a new informer for Mgr type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewMgrInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredMgrInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredMgrInformer constructs a new informer for Mgr type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredMgrInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ManagementVmwareV1().Mgrs().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ManagementVmwareV1().Mgrs().Watch(context.TODO(), options)
			},
		},
		&managementvmwareorgv1.Mgr{},
		resyncPeriod,
		indexers,
	)
}

func (f *mgrInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredMgrInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *mgrInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&managementvmwareorgv1.Mgr{}, f.defaultInformer)
}

func (f *mgrInformer) Lister() v1.MgrLister {
	return v1.NewMgrLister(f.Informer().GetIndexer())
}

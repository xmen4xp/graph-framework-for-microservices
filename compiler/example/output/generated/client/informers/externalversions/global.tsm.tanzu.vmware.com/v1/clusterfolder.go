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
	globaltsmtanzuvmwarecomv1 "nexustempmodule/apis/global.tsm.tanzu.vmware.com/v1"
	versioned "nexustempmodule/client/clientset/versioned"
	internalinterfaces "nexustempmodule/client/informers/externalversions/internalinterfaces"
	v1 "nexustempmodule/client/listers/global.tsm.tanzu.vmware.com/v1"
	time "time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ClusterFolderInformer provides access to a shared informer and lister for
// ClusterFolders.
type ClusterFolderInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.ClusterFolderLister
}

type clusterFolderInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewClusterFolderInformer constructs a new informer for ClusterFolder type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewClusterFolderInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredClusterFolderInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredClusterFolderInformer constructs a new informer for ClusterFolder type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredClusterFolderInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.GlobalTsmV1().ClusterFolders().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.GlobalTsmV1().ClusterFolders().Watch(context.TODO(), options)
			},
		},
		&globaltsmtanzuvmwarecomv1.ClusterFolder{},
		resyncPeriod,
		indexers,
	)
}

func (f *clusterFolderInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredClusterFolderInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *clusterFolderInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&globaltsmtanzuvmwarecomv1.ClusterFolder{}, f.defaultInformer)
}

func (f *clusterFolderInformer) Lister() v1.ClusterFolderLister {
	return v1.NewClusterFolderLister(f.Informer().GetIndexer())
}
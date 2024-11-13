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
	policypkgtsmtanzuvmwarecomv1 "/build/apis/policypkg.tsm.tanzu.vmware.com/v1"
	versioned "/build/client/clientset/versioned"
	internalinterfaces "/build/client/informers/externalversions/internalinterfaces"
	v1 "/build/client/listers/policypkg.tsm.tanzu.vmware.com/v1"
	time "time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// VMpolicyInformer provides access to a shared informer and lister for
// VMpolicies.
type VMpolicyInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.VMpolicyLister
}

type vMpolicyInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewVMpolicyInformer constructs a new informer for VMpolicy type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewVMpolicyInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredVMpolicyInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredVMpolicyInformer constructs a new informer for VMpolicy type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredVMpolicyInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.PolicypkgTsmV1().VMpolicies().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.PolicypkgTsmV1().VMpolicies().Watch(context.TODO(), options)
			},
		},
		&policypkgtsmtanzuvmwarecomv1.VMpolicy{},
		resyncPeriod,
		indexers,
	)
}

func (f *vMpolicyInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredVMpolicyInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *vMpolicyInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&policypkgtsmtanzuvmwarecomv1.VMpolicy{}, f.defaultInformer)
}

func (f *vMpolicyInformer) Lister() v1.VMpolicyLister {
	return v1.NewVMpolicyLister(f.Informer().GetIndexer())
}
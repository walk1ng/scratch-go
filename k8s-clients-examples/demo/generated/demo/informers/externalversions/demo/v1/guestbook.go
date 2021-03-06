/*


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
	time "time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	demov1 "walk1ng.io/demo/api/demo/v1"
	versioned "walk1ng.io/demo/generated/demo/clientset/versioned"
	internalinterfaces "walk1ng.io/demo/generated/demo/informers/externalversions/internalinterfaces"
	v1 "walk1ng.io/demo/generated/demo/listers/demo/v1"
)

// GuestBookInformer provides access to a shared informer and lister for
// GuestBooks.
type GuestBookInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.GuestBookLister
}

type guestBookInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewGuestBookInformer constructs a new informer for GuestBook type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewGuestBookInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredGuestBookInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredGuestBookInformer constructs a new informer for GuestBook type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredGuestBookInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.DemoV1().GuestBooks(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.DemoV1().GuestBooks(namespace).Watch(options)
			},
		},
		&demov1.GuestBook{},
		resyncPeriod,
		indexers,
	)
}

func (f *guestBookInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredGuestBookInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *guestBookInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&demov1.GuestBook{}, f.defaultInformer)
}

func (f *guestBookInformer) Lister() v1.GuestBookLister {
	return v1.NewGuestBookLister(f.Informer().GetIndexer())
}

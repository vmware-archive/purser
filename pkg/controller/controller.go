/*
 * Copyright (c) 2018 VMware Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package controller

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"

	groups_v1 "github.com/vmware/purser/pkg/apis/groups/v1"
	subscriber_v1 "github.com/vmware/purser/pkg/apis/subscriber/v1"

	apps_v1beta1 "k8s.io/api/apps/v1beta1"
	batch_v1 "k8s.io/api/batch/v1"
	api_v1 "k8s.io/api/core/v1"
	ext_v1beta1 "k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

// Kubeclient is kubernetes Clientset
var Kubeclient *kubernetes.Clientset

// Controller holds Kubernetes controller components
type Controller struct {
	clientset kubernetes.Interface
	queue     workqueue.RateLimitingInterface
	informer  cache.SharedIndexInformer
	conf      *Config
}

// Event indicate the informerEvent
type Event struct {
	key          string
	eventType    string
	resourceType string
	data         interface{}
	captureTime  meta_v1.Time
}

// Start runs the controller goroutine.
// nolint: gocyclo, interfacer
func Start(conf *Config) {
	Kubeclient = conf.Kubeclient

	if conf.Resource.Pod {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return Kubeclient.CoreV1().Pods(meta_v1.NamespaceAll).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return Kubeclient.CoreV1().Pods(meta_v1.NamespaceAll).Watch(options)
				},
			},
			&api_v1.Pod{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(Kubeclient, informer, "Pod")
		c.conf = conf
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.Node {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return Kubeclient.CoreV1().Nodes().List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return Kubeclient.CoreV1().Nodes().Watch(options)
				},
			},
			&api_v1.Node{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(Kubeclient, informer, "Node")
		c.conf = conf
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.PersistentVolume {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return Kubeclient.CoreV1().PersistentVolumes().List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return Kubeclient.CoreV1().PersistentVolumes().Watch(options)
				},
			},
			&api_v1.PersistentVolume{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(Kubeclient, informer, "PersistentVolume")
		c.conf = conf
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.PersistentVolumeClaim {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return Kubeclient.CoreV1().PersistentVolumeClaims(meta_v1.NamespaceAll).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return Kubeclient.CoreV1().PersistentVolumeClaims(meta_v1.NamespaceAll).Watch(options)
				},
			},
			&api_v1.PersistentVolumeClaim{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(Kubeclient, informer, "PersistentVolumeClaim")
		c.conf = conf
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.Service {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return Kubeclient.CoreV1().Services(meta_v1.NamespaceAll).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return Kubeclient.CoreV1().Services(meta_v1.NamespaceAll).Watch(options)
				},
			},
			&api_v1.Service{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(Kubeclient, informer, "Service")
		c.conf = conf
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.ReplicaSet {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return Kubeclient.ExtensionsV1beta1().ReplicaSets(meta_v1.NamespaceAll).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return Kubeclient.ExtensionsV1beta1().ReplicaSets(meta_v1.NamespaceAll).Watch(options)
				},
			},
			&ext_v1beta1.ReplicaSet{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(Kubeclient, informer, "ReplicaSet")
		c.conf = conf
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.DaemonSet {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return Kubeclient.ExtensionsV1beta1().DaemonSets(meta_v1.NamespaceAll).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return Kubeclient.ExtensionsV1beta1().DaemonSets(meta_v1.NamespaceAll).Watch(options)
				},
			},
			&ext_v1beta1.DaemonSet{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(Kubeclient, informer, "DaemonSet")
		c.conf = conf
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.Deployment {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return Kubeclient.AppsV1beta1().Deployments(meta_v1.NamespaceAll).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return Kubeclient.AppsV1beta1().Deployments(meta_v1.NamespaceAll).Watch(options)
				},
			},
			&apps_v1beta1.Deployment{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(Kubeclient, informer, "Deployment")
		c.conf = conf
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.StatefulSet {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return Kubeclient.AppsV1beta1().StatefulSets(meta_v1.NamespaceAll).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return Kubeclient.AppsV1beta1().StatefulSets(meta_v1.NamespaceAll).Watch(options)
				},
			},
			&apps_v1beta1.StatefulSet{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(Kubeclient, informer, "StatefulSet")
		c.conf = conf
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.Job {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return Kubeclient.BatchV1().Jobs(meta_v1.NamespaceAll).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return Kubeclient.BatchV1().Jobs(meta_v1.NamespaceAll).Watch(options)
				},
			},
			&batch_v1.Job{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(Kubeclient, informer, "Job")
		c.conf = conf
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.Namespace {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return Kubeclient.CoreV1().Namespaces().List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return Kubeclient.CoreV1().Namespaces().Watch(options)
				},
			},
			&api_v1.Namespace{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(Kubeclient, informer, "Namespace")
		c.conf = conf
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.Group {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return conf.Groupcrdclient.List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return conf.Groupcrdclient.Watch(options)
				},
			},
			&groups_v1.Group{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(Kubeclient, informer, "Group")
		c.conf = conf
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.Subscriber {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return conf.Subscriberclient.List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return conf.Subscriberclient.Watch(options)
				},
			},
			&subscriber_v1.Subscriber{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(Kubeclient, informer, "Subscriber")
		c.conf = conf
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}

func newResourceController(client kubernetes.Interface, informer cache.SharedIndexInformer, resourceType string) *Controller {
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	var newEvent Event
	var err error
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			newEvent.key, err = cache.MetaNamespaceKeyFunc(obj)
			newEvent.eventType = Create
			newEvent.resourceType = resourceType
			newEvent.captureTime = meta_v1.Now()
			log.Printf("Processing add to %v: %s", resourceType, newEvent.key)
			if err == nil {
				queue.Add(newEvent)
			}
		},
		// TODO: Fixme
		UpdateFunc: func(old, new interface{}) {
			/*newEvent.key, err = cache.MetaNamespaceKeyFunc(old)
			newEvent.eventType = "update"
			newEvent.resourceType = resourceType
			log.Printf("Processing update to %v: %s", resourceType, newEvent.key)
			if err == nil {
				queue.Add(newEvent)
			}*/
		},
		DeleteFunc: func(obj interface{}) {
			newEvent.key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			newEvent.eventType = Delete
			newEvent.resourceType = resourceType
			newEvent.data = obj
			newEvent.captureTime = meta_v1.Now()
			log.Printf("Processing delete to %v: %s", resourceType, newEvent.key)
			if err == nil {
				queue.Add(newEvent)
			}
		},
	})

	return &Controller{
		clientset: client,
		informer:  informer,
		queue:     queue,
	}
}

// Run initiates the controller
func (c *Controller) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}

	log.Println("Purser controller synced and ready")
	wait.Until(c.runWorker, time.Second, stopCh)
}

// HasSynced is required for the cache.Controller interface.
func (c *Controller) HasSynced() bool {
	return c.informer.HasSynced()
}

// LastSyncResourceVersion is required for the cache.Controller interface.
func (c *Controller) LastSyncResourceVersion() string {
	return c.informer.LastSyncResourceVersion()
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
		// continue looping
	}
}

func (c *Controller) processNextItem() bool {
	newEvent, quit := c.queue.Get()

	if quit {
		return false
	}
	defer c.queue.Done(newEvent)
	err := c.processItem(newEvent.(Event))
	if err == nil {
		c.queue.Forget(newEvent)
	} else {
		log.Printf("Error processing %s (giving up): %v", newEvent.(Event).key, err)
		c.queue.Forget(newEvent)
		utilruntime.HandleError(err)
	}

	return true
}

func (c *Controller) processItem(newEvent Event) error {
	obj, _, err := c.informer.GetIndexer().GetByKey(newEvent.key)
	if err != nil {
		return fmt.Errorf("error fetching object with key %s from store: %v", newEvent.key, err)
	}

	// process events based on its type
	switch newEvent.eventType {
	case Create:
		str, err := json.Marshal(obj)
		if err != nil {
			log.Errorf("Error marshalling object %s", obj)
		}
		payload := &Payload{Key: newEvent.key, EventType: newEvent.eventType, ResourceType: newEvent.resourceType,
			CloudType: "aws", Data: string(str), CaptureTime: newEvent.captureTime}
		c.conf.RingBuffer.Put(payload)
		return nil
	case Update:
		// TODO: Decide on what needs to be propagated.
		return nil
	case Delete:
		str, err := json.Marshal(newEvent.data)
		if err != nil {
			log.Errorf("Error marshalling object %s", newEvent.data)
		}
		payload := &Payload{Key: newEvent.key, EventType: newEvent.eventType, ResourceType: newEvent.resourceType,
			CloudType: "aws", Data: string(str), CaptureTime: newEvent.captureTime}
		c.conf.RingBuffer.Put(payload)
		return nil
	}
	return nil
}

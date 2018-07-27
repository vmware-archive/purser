package controller

import (
	api_v1 "k8s.io/api/core/v1"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"kuber-controller/client"
	"kuber-controller/config"
	"kuber-controller/utils"
	"os"
	"os/signal"
	"syscall"
	"time"
	log "github.com/Sirupsen/logrus"
	"kuber-controller/buffering"
	"kuber-controller/uploader"
	"sync"
	"fmt"
	"kuber-controller/metrics"
)

type Controller struct {
	clientset kubernetes.Interface
	queue     workqueue.RateLimitingInterface
	informer  cache.SharedIndexInformer
}

// Event indicate the informerEvent
type Event struct {
	key          string
	eventType    string
	namespace    string
	resourceType string
}

var serverStartTime time.Time
var groupcrdclient *client.GroupCrdClient
var ringBuffer *buffering.RingBuffer

func init() {
	ringBuffer = &buffering.RingBuffer{Size: buffering.BUFFER_SIZE, Mutex: &sync.Mutex{}}
}

func TestCrdFlow() {
	crdclient := GetApiExtensionClient()
	CreateCRDInstance(crdclient, "xyz", "namespace")
	ListCrdInstances(crdclient)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}

func Start(conf *config.Config) {
	// initialize client for api extension server
	groupcrdclient = GetApiExtensionClient()

	//go uploader.UploadData(ringBuffer)

	//var kubeClient kubernetes.Interface
	var kubeClient *kubernetes.Clientset
	_, err := rest.InClusterConfig()
	if err != nil {
		kubeClient = utils.GetClientOutOfCluster()
	} else {
		kubeClient = utils.GetClient()
	}

	if conf.Resource.Pod {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().Pods(meta_v1.NamespaceAll).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().Pods(meta_v1.NamespaceAll).Watch(options)
				},
			},
			&api_v1.Pod{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, informer, "pod")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.Node {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().Nodes().List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().Nodes().Watch(options)
				},
			},
			&api_v1.Node{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, informer, "node")
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
			newEvent.eventType = "create"
			newEvent.resourceType = resourceType
			log.Printf("Processing add to %v: %s", resourceType, newEvent.key)
			if err == nil {
				queue.Add(newEvent)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			newEvent.key, err = cache.MetaNamespaceKeyFunc(old)
			newEvent.eventType = "update"
			newEvent.resourceType = resourceType
			log.Printf("Processing update to %v: %s", resourceType, newEvent.key)
			if err == nil {
				queue.Add(newEvent)
			}
		},
		DeleteFunc: func(obj interface{}) {
			log.Printf("Delete object: %s\n", obj)
			newEvent.key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			newEvent.eventType = "delete"
			newEvent.resourceType = resourceType
			newEvent.namespace = utils.GetObjectMetaData(obj).Namespace
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

// Run starts the kubewatch controller
func (c *Controller) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	log.Println("Starting kubewatch controller")
	serverStartTime = time.Now().Local()

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	log.Println("Kubewatch controller synced and ready")
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
		// No error, reset the ratelimit counters
		c.queue.Forget(newEvent)
		/*} else if c.queue.NumRequeues(newEvent) < maxRetries {
		c.logger.Errorf("Error processing %s (will retry): %v", newEvent.(Event).key, err)
		c.queue.AddRateLimited(newEvent)*/
	} else {
		log.Printf("Error processing %s (giving up): %v", newEvent.(Event).key, err)
		c.queue.Forget(newEvent)
		utilruntime.HandleError(err)
	}

	return true
}

var count int16

func (c *Controller) processItem(newEvent Event) error {
	obj, _, err := c.informer.GetIndexer().GetByKey(newEvent.key)
	if err != nil {
		return fmt.Errorf("Error fetching object with key %s from store: %v", newEvent.key, err)
	}

	// process events based on its type
	switch newEvent.eventType {
	case "create":
		switch obj.(type) {
		case *api_v1.Pod:
			pod := obj.(*api_v1.Pod)
			payload := &uploader.Payload{Key: newEvent.key, EventType: newEvent.eventType, Namespace: newEvent.namespace,
				ResourceType: newEvent.resourceType, Data: &obj}
			ringBuffer.Put(payload)

			met := metrics.CalculatePodStatsFromContainers(pod)
			//metrics.PrintPodStats(pod, met)
			UpdateNamespaceGroupCrd(groupcrdclient, pod.Namespace, "namespace", pod.Name, met)
			//UpdateLabelGroupCrd(crdclient, met, pod)
			//UpdateCustomGroupCrd(crdclient, met, pod)

		case *api_v1.Node:
			count++
			log.Printf("Event type is node and count =  %d\n", count)
		}
		return nil
	case "update":
		// Decide on what needs to be propagated.
		return nil
	case "delete":
		payload := &uploader.Payload{Key: newEvent.key, EventType: newEvent.eventType, Namespace: newEvent.namespace,
			ResourceType: newEvent.resourceType, Data: &obj}
		ringBuffer.Put(payload)
		return nil
	}
	return nil
}

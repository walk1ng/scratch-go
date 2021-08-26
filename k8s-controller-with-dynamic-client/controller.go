package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

var (
	redisFailoverResource = schema.GroupVersionResource{
		Group:    "databases.spotahome.com",
		Version:  "v1",
		Resource: "redisfailovers",
	}

	secretResource = schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "secrets",
	}
)

const (
	maxTries = 3
)

/*
The controller will watch the redisfailover CR and
remove associated redis-auth secret while redisfailover CR was deleted
*/

type RedisFailoverController struct {
	informer cache.SharedIndexInformer
	client   dynamic.Interface
	queue    workqueue.RateLimitingInterface
	stopper  chan struct{}
}

func NewRedidFailoverController(dc dynamic.Interface) (*RedisFailoverController, error) {
	// dynamic informer factory
	dynamicInformerFactory := dynamicinformer.NewDynamicSharedInformerFactory(dc, time.Second*5)
	// dynamic informer for resource
	informer := dynamicInformerFactory.ForResource(redisFailoverResource).Informer()

	// queue
	q := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// informer event handler
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				q.Add(key)
			}
		},
	})

	return &RedisFailoverController{
		informer: informer,
		client:   dc,
		queue:    q,
		stopper:  make(chan struct{}),
	}, nil
}

func (controller *RedisFailoverController) Stop() {
	close(controller.stopper)
}

func (controller *RedisFailoverController) Run() {
	defer utilruntime.HandleCrash()

	defer controller.queue.ShutDown()

	// start the informer
	go controller.informer.Run(controller.stopper)

	// wait the cache completely synced
	if !cache.WaitForCacheSync(controller.stopper, controller.informer.HasSynced) {
		utilruntime.HandleError(errors.New("timeout of caches sync"))
		return
	} else {
		log.Println("caches synced")
	}

	// start workers, normally runWorker will infinity loop until some error happen(or run finished).
	// the wait.Until will restart runWorker after one second until the stopper channel cloased
	wait.Until(controller.runWorker, time.Second, controller.stopper)

}

func (controller *RedisFailoverController) runWorker() {
	log.Println("start worker")

	// infinity loop
	for {
		/*
			Get blocks until it can return an item to be processed. If shutdown = true,
			the caller should end their goroutine. You must call Done with item when you
			have finished processing it.
		*/
		key, shutdown := controller.queue.Get()
		if shutdown {
			return
		}

		// process the key "namespace/name"
		err := controller.processItem(key.(string))
		if err == nil {
			controller.queue.Forget(key)
		} else if controller.queue.NumRequeues(key) < maxTries {
			// requeue the key
			controller.queue.AddRateLimited(key)
		} else {
			controller.queue.Forget(key)
			utilruntime.HandleError(err)
		}

		// mark the key was done
		controller.queue.Done(key)

	}

}

func (controller *RedisFailoverController) processItem(redisFailover string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(redisFailover)
	if err != nil {
		utilruntime.HandleError(errors.WithMessage(err, "invalid resource key"))
		return nil
	}

	// mock to process the related redis auth secret
	secretName := fmt.Sprintf("redis-auth-%s", name)

	// client for secret
	client := controller.client.Resource(secretResource).Namespace(namespace)

	// get the related redis auth secret
	_, err = client.Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		if kerrors.IsNotFound(err) {
			log.Printf("secret %s/%s for redisfailover %s/%s not found\n", namespace, secretName, namespace, name)
			return err
		}

		utilruntime.HandleError(errors.WithMessagef(err, "failed to get secret %s/%s for redisfailover %s/%s: %v\n", namespace, secretName, namespace, name, err))

		return err
	}

	log.Printf("process secret %s/%s for redisfailover %s/%s\n", namespace, secretName, namespace, name)
	/*

	 */

	return nil
}

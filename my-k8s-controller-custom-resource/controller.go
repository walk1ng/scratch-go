package main

import (
	"fmt"
	samplecrdv1 "my-k8s-controller-custom-resource/pkg/apis/samplecrd/v1"
	clientset "my-k8s-controller-custom-resource/pkg/client/clientset/versioned"
	networkScheme "my-k8s-controller-custom-resource/pkg/client/clientset/versioned/scheme"
	informers "my-k8s-controller-custom-resource/pkg/client/informers/externalversions/samplecrd/v1"
	listers "my-k8s-controller-custom-resource/pkg/client/listers/samplecrd/v1"
	"time"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const controllerAgentName = "network-controller"

type Controller struct {
	// kubeclientsets
	kubeclientset kubernetes.Interface
	// network clientsets
	networkclientset clientset.Interface
	// network lister
	networkLister listers.NetworkLister

	networkSynced cache.InformerSynced

	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources
	// to the API server
	recorder record.EventRecorder
}

func NewController(
	kubeclientset kubernetes.Interface,
	networkclientset clientset.Interface,
	networkInformer informers.NetworkInformer) *Controller {

	utilruntime.Must(networkScheme.AddToScheme(scheme.Scheme))
	glog.V(4).Info("Creating event broadcaster")
	eventBroadCaster := record.NewBroadcaster()
	eventBroadCaster.StartLogging(glog.Infof)
	eventBroadCaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadCaster.NewRecorder(scheme.Scheme, corev1.EventSource{
		Component: controllerAgentName,
	})

	controller := &Controller{
		kubeclientset:    kubeclientset,
		networkclientset: networkclientset,
		networkLister:    networkInformer.Lister(),
		networkSynced:    networkInformer.Informer().HasSynced,
		workqueue:        workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Networks"),
		recorder:         recorder,
	}

	glog.Info("Setting up event handlers")
	// set up an event handler for when network resources change
	networkInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueNetwork,
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			oldNetwork := oldObj.(*samplecrdv1.Network)
			newNetwork := newObj.(*samplecrdv1.Network)
			if oldNetwork.ResourceVersion == newNetwork.ResourceVersion {
				return
			}
			glog.Info("network updates, enqueue the new one")
			controller.enqueueNetwork(newObj)
		},
		DeleteFunc: controller.enqueueNetworkForDelete,
	})

	return controller
}

func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	glog.Info("Starting Network control loop")

	glog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.networkSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("Starting workers")
	for i := 0; i < threadiness; i++ {
		wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("Started workers")
	<-stopCh
	glog.Info("Shutting down workers")

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {

	}
}

func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()
	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var key string
		var ok bool

		// convert key to string
		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj)
			utilruntime.HandleError(errors.New(fmt.Sprintf("expected string key but got %#v", obj)))
			return nil
		}

		// failed to sync
		if err := c.syncHandler(key); err != nil {
			return errors.WithMessagef(err, "error syncing %s: %s", key, err.Error())
		}

		// sync successfully
		c.workqueue.Forget(obj)
		glog.Info("Successfully sync key ", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	network, err := c.networkLister.Networks(namespace).Get(name)
	if err != nil {
		if kerrors.IsNotFound(err) {
			glog.Warningf("Network: %s/%s does not exist in local cache, will delete if from Neutron ...", name, namespace)
			glog.Infof("[Neutron] Deleting network: %s/%s ...", namespace, name)

			return nil
		}

		utilruntime.HandleError(fmt.Errorf("failed to list network by: %s/%s", namespace, name))
		return err
	}

	glog.Infof("[Neutron] Try to process network: %#v ...", network)

	// FIX ME: Do diff().
	//
	// actualNetwork, exists := neutron.Get(namespace, name)
	//
	// if !exists {
	// 	neutron.Create(namespace, name)
	// } else if !reflect.DeepEqual(actualNetwork, network) {
	// 	neutron.Update(namespace, name)
	// }

	// c.recorder.Event(network, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func (c *Controller) enqueueNetwork(obj interface{}) {
	var key string
	var err error

	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}

	c.workqueue.AddRateLimited(key)
}

func (c *Controller) enqueueNetworkForDelete(obj interface{}) {
	var key string
	var err error

	key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		utilruntime.HandleError(err)
		return
	}

	c.workqueue.AddRateLimited(key)
}

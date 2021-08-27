package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/user"
	"path"
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := newKubeConfig()
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	stopper := make(chan struct{})

	informerFactory := informers.NewSharedInformerFactory(clientset, time.Second*10)
	podInformer := informerFactory.Core().V1().Pods()

	podLister := podInformer.Lister()

	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// method 1, assert obj as metav1.Object
			// and use its method GetNamespace() and GetName()
			pod := obj.(metav1.Object)
			log.Printf("New Pod %s/%s added\n", pod.GetNamespace(), pod.GetName())
		},
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {

			// assert object to corev1.Pod
			oldPod := oldObj.(*corev1.Pod)
			newPod := newObj.(*corev1.Pod)

			// get the desired pod to use pod lister to fetch from local indexed store
			name, namespace := oldPod.GetName(), oldPod.GetNamespace()
			oldPod1, err := podLister.Pods(namespace).Get(name)
			if err != nil {
				utilruntime.HandleError(errors.WithMessagef(err, "pod lister to get pod %s/%s failed", namespace, name))
				return
			}

			if oldPod.GetResourceVersion() == newPod.GetResourceVersion() {
				log.Printf("Old Pod %s/%s updated due to caches synced", oldPod.GetNamespace(), oldPod.GetName())
				prettyPrint("LOCAL STORE", oldPod1)
				return
			}

			log.Printf("Old Pod %s/%s updated", oldPod.GetNamespace(), oldPod.GetName())
			prettyPrint("BEFORE", oldPod)
			prettyPrint("AFTER", newPod)
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err != nil {
				utilruntime.HandleError(err)
				return
			}

			// method 2: use cache.SplitMetaNamespaceKey to get namespace and name of object
			namespace, name, err := cache.SplitMetaNamespaceKey(key)
			if err != nil {
				utilruntime.HandleError(errors.WithMessage(err, fmt.Sprintf("invalid resource key: %s", key)))
				return
			}

			log.Printf("Old Pod %s/%s deleted\n", namespace, name)

		},
	})

	podInformer.Informer().Run(stopper)

}

func newKubeConfig() (*rest.Config, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	kubeconfig := path.Join(user.HomeDir, ".kube", "config")
	c, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		c, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func prettyPrint(title string, obj interface{}) {
	log.Printf(" =============  %s =============\n", title)
	data, _ := json.MarshalIndent(obj, "", "    ")
	log.Println(string(data))
	log.Printf(" =============  %s =============\n", title)
}

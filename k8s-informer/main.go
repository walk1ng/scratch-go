package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	watchNamespace = "dev"
)

func main() {
	// get rest config
	c, err := newKubeConfig()
	if err != nil {
		panic(err)
	}

	// get the clientset
	clientset, err := kubernetes.NewForConfig(c)
	if err != nil {
		panic(err)
	}

	// create informer factory
	informerFactory := informers.NewSharedInformerFactory(clientset, time.Second*10)
	// deployment informer
	deploymentInformer := informerFactory.Apps().V1().Deployments()
	// get informer
	informer := deploymentInformer.Informer()
	// get lister
	lister := deploymentInformer.Lister()

	// register event handler
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		// when deployemnt added
		AddFunc: func(obj interface{}) {
			dep := obj.(*appsv1.Deployment)
			if dep.Namespace == watchNamespace {
				log.Printf("%s added\n", dep.Name)
			}
		},
		// when deployment updated
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			old := oldObj.(*appsv1.Deployment)
			new := newObj.(*appsv1.Deployment)
			if old.Namespace == watchNamespace {
				log.Printf("deployment updated, old: %s, new: %s\n", old.Name, new.Name)
			}
		},
		// when deployment deleted
		DeleteFunc: func(obj interface{}) {
			dep := obj.(*appsv1.Deployment)
			if dep.Namespace == watchNamespace {
				log.Printf("%s deleted\n", dep.Name)
			}
		},
	})

	// channel for stop the informerfactory
	stopper := make(chan struct{})
	defer close(stopper)

	// start informer(factory)
	informerFactory.Start(stopper)
	// sync etcd to informer's local store
	informerFactory.WaitForCacheSync(stopper)

	// get all deployments from dev namespace
	deployments, err := lister.Deployments(watchNamespace).List(labels.Everything())
	if err != nil {
		fmt.Println("Failed to list all deployments in dev namespace by informer lister")
		panic(err)
	}

	for _, dep := range deployments {
		fmt.Printf("%s\n", dep.Name)
	}

	// hang
	<-stopper

}

func home() (string, error) {
	h, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return h, nil
}

func newKubeConfig() (*rest.Config, error) {
	h, err := home()
	if err != nil {
		return nil, err
	}

	kubeconfig := path.Join(h, ".kube", "config")

	c, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return c, nil
}

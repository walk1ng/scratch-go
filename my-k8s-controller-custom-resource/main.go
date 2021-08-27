package main

import (
	"flag"
	"os/user"
	"path"
	"time"

	"my-k8s-controller-custom-resource/pkg/client/clientset/versioned"
	informers "my-k8s-controller-custom-resource/pkg/client/informers/externalversions"

	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	masterURL  string
	kubeconfig string
)

func main() {

	stopper := make(chan struct{})

	config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		panic(err)
	}

	kubeclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	networkclient, err := versioned.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// create network informer
	networkInformerFactory := informers.NewSharedInformerFactory(networkclient, time.Second*5)
	networkInformer := networkInformerFactory.Samplecrd().V1().Networks()

	// create the controller
	controller := NewController(kubeclient, networkclient, networkInformer)

	// start the networkinformer
	go networkInformerFactory.Start(stopper)

	// run the controller
	if err := controller.Run(2, stopper); err != nil {
		glog.Fatalf("controller run failure with error: %v", err)
	}

}

func init() {

	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	homeDir := path.Join(u.HomeDir, ".kube", "config")

	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&kubeconfig, "kubeconfig", homeDir, "Path to a kubeconfig. Only required if out-of-cluster.")
}

package main

import (
	"flag"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	demov1 "walk1ng.io/demo/generated/demo/clientset/versioned/typed/demo/v1"
)

var client *demov1.DemoV1Client

func main() {
	kubeconfig := flag.String("kubeconfig", "/home/sysadmin/.kube/config", "path to the kube config")
	flag.Parse()

	var config *rest.Config
	var err error
	config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		config, _ = rest.InClusterConfig()
	}

	client, err = demov1.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	gblist, err := client.GuestBooks("default").List(v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmtStr := "%-30s%-30s%-30s%-30s\n"
	fmt.Printf(fmtStr, "Namespace", "Name", "Spec.Foo", "GVK")
	for _, gb := range gblist.Items {
		fmt.Printf(fmtStr, gb.Namespace, gb.Name, gb.Spec.Foo, gb.GetObjectKind().GroupVersionKind().String())
	}

}

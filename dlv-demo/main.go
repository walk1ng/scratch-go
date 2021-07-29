package main

import (
	"context"
	"fmt"
	"os/user"
	"path"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	user, _ := user.Current()
	config := path.Join(user.HomeDir, ".kube", "config")

	c, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		fmt.Println("rest.InClusterConfig")
		c, err = rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
	}

	clientset, err := kubernetes.NewForConfig(c)
	if err != nil {
		panic(err)
	}

	pods, err := clientset.CoreV1().Pods("default").List(context.Background(), v1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for i, pod := range pods.Items {
		fmt.Printf("pod %d: %s", i, pod.GetName())
	}

	time.Sleep(time.Second * 3600)
}

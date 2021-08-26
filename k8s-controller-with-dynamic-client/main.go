package main

import (
	"context"
	"log"
	"os/user"
	"path"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := newKubeConfig()
	if err != nil {
		log.Printf("failed to get config: %v", err)
		return
	}

	dc, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Printf("failed to get dynamic client: %v", err)
		return
	}

	ctl, err := NewRedidFailoverController(dc)
	if err != nil {
		log.Printf("failed to new redisfailover controller: %v", err)
		return
	}

	getAllRedisFailovers(dc)
	ctl.Run()

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

func getAllRedisFailovers(client dynamic.Interface) {
	l, err := client.Resource(redisFailoverResource).Namespace("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("failed to list redisfailover: %v\n", err)
		return
	}

	for i, rf := range l.Items {
		log.Printf("%d - %s/%s", i, rf.GetNamespace(), rf.GetName())
	}
}

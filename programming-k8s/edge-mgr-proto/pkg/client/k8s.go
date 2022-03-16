package client

import (
	k8sscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

var KubeClient *Client

func Init(config *rest.Config) (err error) {
	KubeClient, err = newK8sClient(config)
	return err
}

func newK8sClient(config *rest.Config) (*Client, error) {
	return newGenericClient(config, k8sscheme.Scheme)
}

package client

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Client struct {
	C runtimeclient.Client
}

func newGenericClient(config *rest.Config, scheme *runtime.Scheme) (*Client, error) {
	cli := &Client{}
	c, err := runtimeclient.New(config, runtimeclient.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, err
	}
	cli.C = c
	return cli, nil
}

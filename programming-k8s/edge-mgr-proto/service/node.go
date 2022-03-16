package service

import (
	"edge-mgr-proto/pkg/client"
	"edge-mgr-proto/pkg/informers"

	corev1 "k8s.io/api/core/v1"
)

type NodeService interface {
	Get(name string) (*corev1.Node, error)
}

type nodeService struct {
	Informers  informers.Informer
	KubeClient *client.Client
}

func newNodeService(informers informers.Informer, kubeClient *client.Client) *nodeService {
	return &nodeService{
		Informers:  informers,
		KubeClient: kubeClient,
	}
}

func (svc *nodeService) Get(name string) (*corev1.Node, error) {
	return svc.Informers.CoreV1().Nodes().Lister().Get(name)
}

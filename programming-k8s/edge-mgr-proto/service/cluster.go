package service

import (
	"edge-mgr-proto/pkg/client"
	"edge-mgr-proto/pkg/informers"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type ClusterService interface {
	GetClusterNodes() ([]*corev1.Node, error)
	GetNodeIPList() ([]string, error)
}

type clusterService struct {
	Informers  informers.Informer
	KubeClient *client.Client
}

func newClusterService(informers informers.Informer, kubeClient *client.Client) *clusterService {
	return &clusterService{
		Informers:  informers,
		KubeClient: kubeClient,
	}
}

func (svc *clusterService) GetClusterNodes() ([]*corev1.Node, error) {
	return svc.Informers.CoreV1().Nodes().Lister().List(labels.Nothing())
}

func (svc *clusterService) GetNodeIPList() ([]string, error) {
	nodes, err := svc.GetClusterNodes()
	if err != nil {
		return []string{}, err
	}

	nodeIPList := make([]string, len(nodes))
	for _, node := range nodes {
		for _, addr := range node.Status.Addresses {
			if addr.Type == corev1.NodeExternalIP {
				nodeIPList = append(nodeIPList, addr.Address)
			}
		}
	}
	return nodeIPList, nil
}

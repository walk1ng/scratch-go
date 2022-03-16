package service

import (
	"edge-mgr-proto/pkg/client"
	"edge-mgr-proto/pkg/informers"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type ClusterService interface {
	GetClusterNodes() ([]*corev1.Node, error)
	GetNodeList() ([]string, error)
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
	return svc.Informers.CoreV1().Nodes().Lister().List(labels.Everything())
}

func (svc *clusterService) GetNodeList() ([]string, error) {
	nodes, err := svc.GetClusterNodes()
	if err != nil {
		return []string{}, err
	}

	fmt.Println("len of nodes:", len(nodes))

	nodeList := make([]string, len(nodes))
	for i, node := range nodes {
		fmt.Println(node.Name)
		nodeList[i] = node.Name
	}
	fmt.Println(nodeList)
	return nodeList, nil
}

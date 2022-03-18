package service

import (
	"edge-mgr-proto/pkg/client"
	"edge-mgr-proto/pkg/informers"
	"edge-mgr-proto/types"
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

type ClusterService interface {
	Version(*corev1.Node) string
	Status() types.ClusterStatus
	PodCapacity(*corev1.Node) uint
	GetNodeList([]*corev1.Node) []string
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

func (svc *clusterService) Version(master *corev1.Node) string {
	return master.Status.NodeInfo.KubeletVersion
}

func (svc *clusterService) Status() types.ClusterStatus {
	// TODO
	return types.ClusterStatusReady
}

func (svc *clusterService) PodCapacity(master *corev1.Node) uint {
	return uint(master.Status.Capacity.Pods().Value())
}

func (svc *clusterService) GetNodeList(nodes []*corev1.Node) []string {
	fmt.Println("len of nodes:", len(nodes))

	nodeList := make([]string, len(nodes))
	for i, node := range nodes {
		fmt.Println(node.Name)
		nodeList[i] = node.Name
	}
	fmt.Println(nodeList)
	return nodeList
}

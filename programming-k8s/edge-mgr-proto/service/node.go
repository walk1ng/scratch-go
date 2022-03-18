package service

import (
	"edge-mgr-proto/pkg/client"
	"edge-mgr-proto/pkg/informers"
	"edge-mgr-proto/types"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
)

type NodeService interface {
	Get(name string) (*corev1.Node, error)
	Status(name string) types.NodeConditionStatus
	_list(selector labels.Selector) ([]*corev1.Node, error)
	ListAll() ([]*corev1.Node, error)
	MasterNodes() ([]*corev1.Node, error)
	WorkerNodes() ([]*corev1.Node, error)
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

func (svc *nodeService) Status(name string) types.NodeConditionStatus {
	node, err := svc.Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			return types.NodeStatusNotExist
		}
		return types.NodeStatusUnknown
	}

	for _, condition := range node.Status.Conditions {
		if condition.Type != corev1.NodeReady {
			continue
		}
		if condition.Status == corev1.ConditionTrue {
			return types.NodeStatusReady
		}
		return types.NodeStatusNotReady
	}

	return types.NodeStatusUnknown
}

func (svc *nodeService) _list(selector labels.Selector) ([]*corev1.Node, error) {
	return svc.Informers.CoreV1().Nodes().Lister().List(selector)
}

func (svc *nodeService) ListAll() ([]*corev1.Node, error) {
	return svc._list(labels.Everything())
}

func (svc *nodeService) MasterNodes() ([]*corev1.Node, error) {
	ret := make([]*corev1.Node, 0)
	allNodes, err := svc.ListAll()
	if err != nil {
		return ret, err
	}

	// filter the master node
	for _, node := range allNodes {
		// kubeadm standard label for master
		if _, ok := node.Labels["node-role.kubernetes.io/master"]; ok {
			ret = append(ret, node)
			continue
		}
		// rke standard label for master
		if v, ok := node.Labels["node-role.kubernetes.io/controlplane"]; v == "true" && ok {
			ret = append(ret, node)
			continue
		}
	}

	return ret, nil
}

func (svc *nodeService) WorkerNodes() ([]*corev1.Node, error) {
	ret := make([]*corev1.Node, 0)
	allNodes, err := svc.ListAll()
	if err != nil {
		return ret, err
	}

	// filter the worker node
	for _, node := range allNodes {
		// kubeadm standard label for master
		if _, ok := node.Labels["node-role.kubernetes.io/master"]; !ok {
			ret = append(ret, node)
			continue
		}
		// rke standard label for worker
		if v, ok := node.Labels["node-role.kubernetes.io/worker"]; v == "true" && ok {
			ret = append(ret, node)
			continue
		}
	}

	return ret, nil
}

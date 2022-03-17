package service

import (
	"edge-mgr-proto/pkg/client"
	"edge-mgr-proto/pkg/informers"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type PodService interface {
	Get(namespace, name string) (*corev1.Pod, error)
	List(selector labels.Selector) ([]*corev1.Pod, error)
	ListAll() ([]*corev1.Pod, error)
	ListByNamespace(namespace string, selector labels.Selector) ([]*corev1.Pod, error)
	ListByNode(nodeName string, selector labels.Selector) ([]*corev1.Pod, error)
}

type podService struct {
	Informers  informers.Informer
	KubeClient *client.Client
}

func newPodService(informers informers.Informer, kubeClient *client.Client) *podService {
	return &podService{
		Informers:  informers,
		KubeClient: kubeClient,
	}
}

func (svc *podService) Get(namespace, name string) (*corev1.Pod, error) {
	return svc.Informers.CoreV1().Pods().Lister().Pods(namespace).Get(name)
}

func (svc *podService) List(selector labels.Selector) ([]*corev1.Pod, error) {
	return svc.Informers.CoreV1().Pods().Lister().List(selector)
}

func (svc *podService) ListAll() ([]*corev1.Pod, error) {
	return svc.List(labels.Everything())
}

func (svc *podService) ListByNamespace(namespace string, selector labels.Selector) ([]*corev1.Pod, error) {
	return svc.Informers.CoreV1().Pods().Lister().Pods(namespace).List(selector)
}

func (svc *podService) ListByNode(nodeName string, selector labels.Selector) ([]*corev1.Pod, error) {
	ret := make([]*corev1.Pod, 0)
	allPods, err := svc.List(selector)
	if err != nil {
		return ret, err
	}

	// filter by nodename
	for _, pod := range allPods {
		if pod.Spec.NodeName == nodeName {
			ret = append(ret, pod)
		}
	}

	return ret, nil
}

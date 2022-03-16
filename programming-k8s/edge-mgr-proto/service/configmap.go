package service

import (
	"edge-mgr-proto/pkg/client"
	"edge-mgr-proto/pkg/informers"

	corev1 "k8s.io/api/core/v1"
)

type ConfigMapService interface {
	Get(namespace, name string) (*corev1.ConfigMap, error)
}

type configmapService struct {
	Informers  informers.Informer
	KubeClient *client.Client
}

func newConfigMapService(informers informers.Informer, kubeClient *client.Client) *configmapService {
	return &configmapService{
		Informers:  informers,
		KubeClient: kubeClient,
	}
}

func (svc *configmapService) Get(namespace, name string) (*corev1.ConfigMap, error) {
	return svc.Informers.CoreV1().ConfigMaps().Lister().ConfigMaps(namespace).Get(name)
}

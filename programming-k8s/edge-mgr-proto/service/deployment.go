package service

import (
	"context"
	"edge-mgr-proto/pkg/informers"
	"encoding/json"

	"edge-mgr-proto/pkg/client"

	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "k8s.io/api/apps/v1"
)

type DeploymentService interface {
	Create(namespace string, data []byte) error
	Get(namespace, name string) (*appsv1.Deployment, error)
}

type deploymentService struct {
	Informers  informers.Informer
	KubeClient *client.Client
}

func newDeploymentService(informers informers.Informer, kubeClient *client.Client) *deploymentService {
	return &deploymentService{
		Informers:  informers,
		KubeClient: kubeClient,
	}
}

func (svc *deploymentService) Create(namespace string, data []byte) error {
	deploy := appsv1.Deployment{}
	if err := json.Unmarshal(data, &deploy); err != nil {
		return err
	}
	deploy.SetNamespace(namespace)
	return svc.KubeClient.C.Create(context.Background(), &deploy, &runtimeclient.CreateOptions{})
}

func (svc *deploymentService) Get(namespace, name string) (*appsv1.Deployment, error) {
	return svc.Informers.Deployment().Lister().Deployments(namespace).Get(name)
}

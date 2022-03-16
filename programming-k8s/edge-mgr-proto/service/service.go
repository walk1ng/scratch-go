package service

import (
	"edge-mgr-proto/mq"
	"edge-mgr-proto/pkg/client"
	"edge-mgr-proto/pkg/informers"
)

var Svc Service

type Service interface {
	Deployment() DeploymentService
	ConfigMap() ConfigMapService
	Node() NodeService
	Cluster() ClusterService
}

type EdgeMgrService struct {
	WorkChan   *mq.WorkChannel
	Informers  informers.Informer
	KubeClient *client.Client
}

func Init(ch *mq.WorkChannel, informers informers.Informer, kubeClient *client.Client) error {
	Svc = newEdgeMgrService(ch, informers, kubeClient)
	return nil
}

func newEdgeMgrService(ch *mq.WorkChannel, informers informers.Informer, kubeClient *client.Client) *EdgeMgrService {
	return &EdgeMgrService{
		WorkChan:   ch,
		Informers:  informers,
		KubeClient: kubeClient,
	}
}

func (svc *EdgeMgrService) Deployment() DeploymentService {
	return newDeploymentService(svc.Informers, svc.KubeClient)
}

func (svc *EdgeMgrService) ConfigMap() ConfigMapService {
	return newConfigMapService(svc.Informers, svc.KubeClient)
}

func (svc *EdgeMgrService) Node() NodeService {
	return newNodeService(svc.Informers, svc.KubeClient)
}

func (svc *EdgeMgrService) Cluster() ClusterService {
	return newClusterService(svc.Informers, svc.KubeClient)
}

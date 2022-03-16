package service

import (
	"edge-mgr-proto/mq"
	"edge-mgr-proto/pkg/client"
	"edge-mgr-proto/pkg/informers"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

var Svc Service

type Service interface {
	// k8s
	Deployment() DeploymentService
	ConfigMap() ConfigMapService
	Node() NodeService
	Cluster() ClusterService
	// metrics
	Prometheus() PrometheusService
}

type EdgeMgrService struct {
	WorkChan      *mq.WorkChannel
	Informers     informers.Informer
	KubeClient    *client.Client
	PrometheusAPI v1.API
}

func Init(ch *mq.WorkChannel, informers informers.Informer, kubeClient *client.Client, promAPI v1.API) error {
	Svc = newEdgeMgrService(ch, informers, kubeClient, promAPI)
	return nil
}

func newEdgeMgrService(ch *mq.WorkChannel, informers informers.Informer, kubeClient *client.Client, promAPI v1.API) *EdgeMgrService {
	return &EdgeMgrService{
		WorkChan:      ch,
		Informers:     informers,
		KubeClient:    kubeClient,
		PrometheusAPI: promAPI,
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

func (svc *EdgeMgrService) Prometheus() PrometheusService {
	return newPrometheusService(svc.PrometheusAPI)
}

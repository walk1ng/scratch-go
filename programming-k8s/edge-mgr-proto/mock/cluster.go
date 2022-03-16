package mock

import (
	"edge-mgr-proto/models"
	"sync"
)

var K8sClusterMgr *K8sClusterManager

type K8sClusterDB []*models.Cluster

type K8sClusterManager struct {
	sync.Mutex
	DB K8sClusterDB
}

func newK8sClusterManager() *K8sClusterManager {
	mgr := &K8sClusterManager{}
	mgr.DB = make(K8sClusterDB, 10)
	return mgr
}

func InitClusterDB() error {
	K8sClusterMgr = newK8sClusterManager()
	return nil
}

package controller

import (
	"edge-mgr-proto/models"
	"edge-mgr-proto/service"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/labels"
)

func GetClusterOverview(c *gin.Context) {
	zap.L().Info("GetClusterOverview")
	// nodes in cluster
	nodes, err := service.Svc.Node().ListAll()
	if err != nil {
		zap.L().Error("failed to list all nodes", zap.Error(err))
		ResponseErrorWithMsg(c, CodeK8sGetFailure, err)
		return
	}
	nodeList := service.Svc.Cluster().GetNodeList(nodes)

	// masters in cluster
	masters, err := service.Svc.Node().MasterNodes()
	if err != nil {
		zap.L().Error("failed to list masters", zap.Error(err))
		ResponseErrorWithMsg(c, CodeK8sGetFailure, err)
		return
	}
	masterList := service.Svc.Cluster().GetNodeList(masters)

	// workers in cluster
	workers, err := service.Svc.Node().WorkerNodes()
	if err != nil {
		zap.L().Error("failed to list workers", zap.Error(err))
		ResponseErrorWithMsg(c, CodeK8sGetFailure, err)
		return
	}
	workerList := service.Svc.Cluster().GetNodeList(workers)

	pods, err := service.Svc.Pod().ListAll()
	if err != nil {
		zap.L().Error("failed to list all pods", zap.Error(err))
		ResponseErrorWithMsg(c, CodeK8sGetFailure, err)
		return
	}

	// metrics
	clusterCpuUsage := service.Svc.Prometheus().GetClusterCpuUsage(nodeList)
	clusterMemUsage := service.Svc.Prometheus().GetClusterMemoryUsage(nodeList)
	clusterDiskUsage := service.Svc.Prometheus().GetClusterDiskUsage(nodeList)

	overview := &models.ClusterOverviewResponse{
		KubernetesVersion: service.Svc.Cluster().Version(masters[0]),
		Status:            service.Svc.Cluster().Status(),
		PodCapacity:       service.Svc.Cluster().PodCapacity(masters[0]),
		PodCount:          uint(len(pods)),
		NodeCount:         uint(len(nodes)),
		Nodes:             nodeList,
		MasterCount:       uint(len(masters)),
		Masters:           masterList,
		WorkerCount:       uint(len(workers)),
		Workers:           workerList,
		CpuTotal:          clusterCpuUsage.Total,
		CpuUsed:           clusterCpuUsage.Used,
		MemoryBytesTotal:  clusterMemUsage.BytesTotal,
		MemoryBytesUsed:   clusterMemUsage.BytesUsed,
		DiskBytesTotal:    clusterDiskUsage.BytesTotal,
		DiskBytesUsed:     clusterDiskUsage.BytesUsed,
	}

	ResponseSuccess(c, overview)
}

func GetClusterNodesOverview(c *gin.Context) {
	zap.L().Info("GetClusterNodesOverview")
	var req models.ClusterNodesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("GetClusterNodesOverview", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParam, err)
		return
	}

	if len(req.Nodes) == 0 {
		zap.L().Error("GetClusterNodesOverview", zap.Error(errors.New("at least 1 node required")))
		ResponseErrorWithMsg(c, CodeInvalidParam, "at least 1 node required")
		return
	}

	zap.L().Info("GetClusterNodesOverview", zap.Strings("nodes", req.Nodes))

	resp := &models.ClusterNodesResponse{
		Count: 0,
		Nodes: make([]*models.ClusterNodeOverview, len(req.Nodes)),
	}

	for i, nodeName := range req.Nodes {
		node, err := service.Svc.Node().Get(nodeName)
		resp.Nodes[i] = &models.ClusterNodeOverview{
			Name:   nodeName,
			Status: service.Svc.Node().Status(nodeName),
		}
		if err != nil {
			zap.L().Error("GetClusterNodesOverview", zap.Error(err))
			continue
		}

		// node addresses
		resp.Nodes[i].Hostname, _ = service.Svc.Node().GetNodeHostname(nodeName)
		resp.Nodes[i].InternalIP, _ = service.Svc.Node().GetNodeInternalIP(nodeName)

		// pods
		pods, err := service.Svc.Pod().ListByNode(nodeName, labels.Everything())
		if err == nil {
			resp.Nodes[i].PodCount = uint(len(pods))
		}

		// metrics
		cpuUsage, _ := service.Svc.Prometheus().GetNodeCpuUsage(nodeName)
		memUsage, _ := service.Svc.Prometheus().GetNodeMemoryUsage(nodeName)
		diskUsage, _ := service.Svc.Prometheus().GetNodeDiskUsage(nodeName)
		diskIOUsage, _ := service.Svc.Prometheus().GetNodeDiskIOUsage(nodeName)
		resp.Nodes[i].CpuUsage = float64(cpuUsage)
		resp.Nodes[i].MemoryUsage = float64(memUsage)
		resp.Nodes[i].DiskUsage = float64(diskUsage)
		resp.Nodes[i].DiskIOUsage = float64(diskIOUsage)

		// node info
		resp.Nodes[i].NodeSystemInfo = node.Status.NodeInfo
	}

	resp.Count = uint(len(req.Nodes))

	ResponseSuccess(c, resp)
}

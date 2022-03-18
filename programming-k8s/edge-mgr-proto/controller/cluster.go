package controller

import (
	"edge-mgr-proto/models"
	"edge-mgr-proto/service"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/labels"
)

func GetClusterOverview(c *gin.Context) {
	// nodes in cluster
	nodes, err := service.Svc.Node().ListAll()
	if err != nil {
		ResponseErrorWithMsg(c, CodeK8sGetFailure, err)
		return
	}
	nodeList := service.Svc.Cluster().GetNodeList(nodes)

	// masters in cluster
	masters, err := service.Svc.Node().MasterNodes()
	if err != nil {
		ResponseErrorWithMsg(c, CodeK8sGetFailure, err)
		return
	}
	masterList := service.Svc.Cluster().GetNodeList(masters)

	// workers in cluster
	workers, err := service.Svc.Node().WorkerNodes()
	if err != nil {
		ResponseErrorWithMsg(c, CodeK8sGetFailure, err)
		return
	}
	workerList := service.Svc.Cluster().GetNodeList(workers)

	pods, err := service.Svc.Pod().ListAll()
	if err != nil {
		ResponseErrorWithMsg(c, CodeK8sGetFailure, err)
		return
	}

	// metrics
	clusterCpuUsage := service.Svc.Prometheus().GetClusterCpuUsage(nodeList)
	clusterMemUsage := service.Svc.Prometheus().GetClusterMemoryUsage(nodeList)
	// clusterDiskUsage := service.Svc.Prometheus().GetClusterDiskUsage(nodeList)

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
		CpuTotal:          float64(clusterCpuUsage["total"]),
		CpuUsed:           float64(clusterCpuUsage["used"]),
		MemoryBytesTotal:  float64(clusterMemUsage["total_bytes"]),
		MemoryBytesUsed:   float64(clusterMemUsage["used_bytes"]),
		// DiskBytesTotal:    float64(clusterDiskUsage["total_bytes"]),
		// DiskBytesUsed:     float64(clusterDiskUsage["used_bytes"]),
	}

	ResponseSuccess(c, overview)
}

func GetClusterNodesOverview(c *gin.Context) {
	var req models.ClusterNodesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseErrorWithMsg(c, CodeInvalidParam, err)
		return
	}

	if len(req.Nodes) == 0 {
		ResponseErrorWithMsg(c, CodeInvalidParam, "empty node list")
		return
	}

	resp := &models.ClusterNodesResponse{
		Count: 0,
		Nodes: make([]*models.ClusterNodeOverview, len(req.Nodes)),
	}

	for i, nodeName := range req.Nodes {
		_, err := service.Svc.Node().Get(nodeName)
		resp.Nodes[i] = &models.ClusterNodeOverview{
			Name:   nodeName,
			Status: service.Svc.Node().Status(nodeName),
		}
		if err != nil {
			continue
		}

		// pods
		pods, err := service.Svc.Pod().ListByNode(nodeName, labels.Everything())
		if err == nil {
			resp.Nodes[i].PodCount = uint(len(pods))
		}

		// metrics
		cpuUsage := service.Svc.Prometheus().GetNodeCpuUsage(nodeName)
		memUsage := service.Svc.Prometheus().GetNodeMemoryUsage(nodeName)
		diskUsage := service.Svc.Prometheus().GetNodeDiskUsage(nodeName)
		diskIOUsage := service.Svc.Prometheus().GetNodeDiskIOUsage(nodeName)
		resp.Nodes[i].CpuUsage = float64(cpuUsage)
		resp.Nodes[i].MemoryUsage = float64(memUsage)
		resp.Nodes[i].DiskUsage = float64(diskUsage)
		resp.Nodes[i].DiskIOUsage = float64(diskIOUsage)
	}

	resp.Count = uint(len(req.Nodes))

	ResponseSuccess(c, resp)
}

package controller

import (
	"edge-mgr-proto/models"
	"edge-mgr-proto/service"

	"github.com/gin-gonic/gin"
)

// func GetClusterNodes(c *gin.Context) {
// 	nodeList := service.Svc.Cluster().GetNodeList()
// 	if err != nil {
// 		ResponseErrorWithMsg(c, CodeK8sGetFailure, err)
// 		return
// 	}
// 	ResponseSuccess(c, service.Svc.Prometheus().GetClusterCpuUsage(nodeList))
// }

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
		Status:            string(service.Svc.Cluster().Status()),
		PodCapacity:       service.Svc.Cluster().PodCapacity(masters[0]),
		Pods:              uint(len(pods)),
		Nodes:             uint(len(nodes)),
		NodeList:          nodeList,
		Masters:           uint(len(masters)),
		MasterList:        masterList,
		Workers:           uint(len(workers)),
		WorkerList:        workerList,
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
	var req *models.ClusterNodesRequest
	if err := c.ShouldBindJSON(req); err != nil {
		ResponseErrorWithMsg(c, CodeInvalidParam, err)
		return
	}

	if len(req.Nodes) == 0 {
		ResponseErrorWithMsg(c, CodeInvalidParam, "empty node list")
		return
	}

	overview := &models.ClusterNodesResponse{}
	overview.Count = uint(len(req.Nodes))

	for _, node := range req.Nodes {

	}

}

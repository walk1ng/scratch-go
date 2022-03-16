package controller

import (
	"edge-mgr-proto/service"

	"github.com/gin-gonic/gin"
)

func GetClusterNodes(c *gin.Context) {
	nodeList, err := service.Svc.Cluster().GetNodeList()
	if err != nil {
		ResponseErrorWithMsg(c, CodeK8sGetFailure, err)
		return
	}
	ResponseSuccess(c, service.Svc.Prometheus().GetClusterCpuUsage(nodeList))
}

func GetClusterOverview(c *gin.Context) {

}

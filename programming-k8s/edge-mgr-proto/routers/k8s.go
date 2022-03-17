package routers

import (
	"edge-mgr-proto/controller"

	"github.com/gin-gonic/gin"
)

func initK8sRouter(g *gin.RouterGroup) {
	k8sGroup := g.Group("k8s")
	{
		// cluster
		k8sGroup.GET("/cluster/overview", controller.GetClusterOverview)
		k8sGroup.POST("/cluster/node", controller.GetClusterNodesOverview)

		// node
		k8sGroup.GET("/nodes", controller.GetNodes)
		k8sGroup.GET("/nodes/:name/details", controller.GetNodeDetails)

		// deployment
		k8sGroup.GET("/namespaces/:namespace/deployments/:name", controller.GetDeployment)
		k8sGroup.POST("/namespaces/:namespace/deployments", controller.CreateDeployment)

		// configmap
		k8sGroup.GET("/namespaces/:namespace/configmaps/:name", controller.GetConfigMap)
	}
}

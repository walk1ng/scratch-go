package controller

import (
	"edge-mgr-proto/models"
	"edge-mgr-proto/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
)

func GetNodes(c *gin.Context) {

}

func GetNodeDetails(c *gin.Context) {
	name := c.Params.ByName("name")
	node, err := service.Svc.Node().Get(name)
	if err != nil {
		zap.L().Error("failed to get node", zap.Error(err), zap.String("name", name))
		var code ResCode
		if errors.IsNotFound(err) {
			code = CodeK8sResNotExist
		} else {
			code = CodeK8sGetFailure
		}
		ResponseErrorWithMsg(c, code, err)
		return
	}

	detail := models.NodeDetails{}
	detail.Node = node
	detail.NodeInfo = node.Status.NodeInfo

	ResponseSuccess(c, detail)

}

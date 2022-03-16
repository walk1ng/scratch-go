package controller

import (
	"edge-mgr-proto/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetConfigMap(c *gin.Context) {
	namespace := c.Params.ByName("namespace")
	name := c.Params.ByName("name")

	cm, err := service.Svc.ConfigMap().Get(namespace, name)
	if err != nil {
		zap.L().Error("failed to get configmap", zap.Error(err), zap.String("namespace", namespace), zap.String("name", name))
		ResponseErrorWithMsg(c, CodeK8sGetFailure, err)
		return
	}

	ResponseSuccess(c, cm)
}

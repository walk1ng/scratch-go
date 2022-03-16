package controller

import (
	"edge-mgr-proto/mq"
	"edge-mgr-proto/service"
	"io/ioutil"

	"edge-mgr-proto/common"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetDeployment(c *gin.Context) {
	namespace := c.Params.ByName("namespace")
	name := c.Params.ByName("name")

	deploy, err := service.Svc.Deployment().Get(namespace, name)
	if err != nil {
		zap.L().Error("failed to get deployment", zap.Error(err), zap.String("namespace", namespace), zap.String("name", name))
		ResponseErrorWithMsg(c, CodeK8sGetFailure, err)
		return
	}

	ResponseSuccess(c, deploy)
}

func CreateDeployment(c *gin.Context) {
	namespace := c.Params.ByName("namespace")
	data, err := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	if err != nil {
		zap.L().Error("create deployment with invalid data", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParam, err)
		return
	}
	println(string(data))

	// mutex and enqueue
	common.WorkChan.Lock()
	defer common.WorkChan.Unlock()
	common.WorkChan.Queue <- mq.Message{
		Namespace: namespace,
		Verb:      mq.Create,
	}

	err = service.Svc.Deployment().Create(namespace, data)
	if err != nil {
		zap.L().Error("failed to create deployment", zap.Error(err), zap.String("namespace", namespace))
		ResponseErrorWithMsg(c, CodeK8sCreateFailure, err)
		return
	}

	ResponseSuccess(c, "ok")
}

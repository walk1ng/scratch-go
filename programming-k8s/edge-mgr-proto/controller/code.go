package controller

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeK8sResExist
	CodeK8sResNotExist
	CodeK8sGetFailure
	CodeK8sDeleteFailure
	CodeK8sCreateFailure
	CodeK8sUpdateFailure
	CodeServerBusy
	Code404
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:          "success",
	CodeInvalidParam:     "请求参数错误",
	CodeK8sResExist:      "K8s资源已存在",
	CodeK8sResNotExist:   "K8s资源不存在",
	CodeK8sGetFailure:    "获取K8s资源失败",
	CodeK8sDeleteFailure: "删除K8s资源失败",
	CodeK8sCreateFailure: "创建K8s资源失败",
	CodeK8sUpdateFailure: "更新K8s资源失败",
	CodeServerBusy:       "服务繁忙",
	Code404:              "api not found",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}

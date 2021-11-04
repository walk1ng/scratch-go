package e

const (
	ERR_BAD_REQUEST int = iota + 1000
	ERR_INTERNAL
	OK
)

var (
	myStatus = make(map[int]string)
)

func init() {
	myStatus[ERR_BAD_REQUEST] = "请求报文格式错误"
	myStatus[ERR_INTERNAL] = "服务内部错误"
	myStatus[OK] = "成功"

}

func GetStatusMessage(code int) string {
	v, ok := myStatus[code]
	if !ok {
		return ""
	}
	return v
}

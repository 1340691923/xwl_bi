//自定义请求辅助方法层
package request

//自定义业务异常
const (
	IdNullError      = 100002
	EmptyParmasError = 100003
	EmptyEventError  = 100004
)

var ErrorMap = map[int]string{
	IdNullError:      "id不能为空！",
	EmptyParmasError: "请求参数不能为空",
	EmptyEventError:  "事件名不能为空",
}

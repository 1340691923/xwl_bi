package report

// 内置异常
const (
	ServerErr     int = 10001
	AppParmasErr  int = 10002
	ReportTypeErr int = 10003
)

// 内置异常表 TOKEN_ERROR
var ERROR_TABLE = map[int]string{
	ServerErr:     "服务异常",
	AppParmasErr:  "appid错误或者appkey错误",
	ReportTypeErr: "上报类型错误",
}

//自定义响应 辅助方法层
package response

import (
	"io"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	fiber "github.com/gofiber/fiber/v2"

	. "github.com/1340691923/xwl_bi/platform-basic-libs/my_error"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

//自定义响应方法
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const (
	SUCCESS = 0
	ERROR   = 500
)

const (
	SearchSuccess       = "查询成功"
	DeleteSuccess       = "删除成功"
	OperateSuccess      = "操作成功"
	LogoutSuccess       = "注销成功"
	ChangeLayoutSuccess = "修改布局成功"
)

func (this *Response) JsonDealErr(err error) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b, _ := json.Marshal(this.DealErr(err))
	return util.BytesToStr(b)
}

//trace
func (this *Response) DealErr(err error) (errorTrace []string) {
	errorTrace = append(errorTrace, err.Error())
	if err != nil {
		for i := 1; ; i++ {
			pc, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			f := runtime.FuncForPC(pc)
			if f.Name() != "runtime.main" && f.Name() != "runtime.goexit" && !strings.Contains(file, "Response.go") {
				errStrings := "文件名:" + file + ",行数:" + strconv.Itoa(line) + ",函数名:" + f.Name()
				errorTrace = append(errorTrace, errStrings)
			}
		}
	}
	return errorTrace
}

//正确信息
func (this *Response) Success(ctx *fiber.Ctx, msg string, data interface{}) error {
	this.Msg = msg
	this.Data = data
	this.send(ctx, SUCCESS)
	return nil
}

//错误信息
func (this *Response) FastError(write io.Writer, err error) error {
	myErr := ErrorToErrorCode(err)

	this.Output(write, map[string]interface{}{
		"code": myErr.Code(),
		"msg":  myErr.Error(),
	})
	return nil
}

//错误信息
func (this *Response) Error(ctx *fiber.Ctx, err error) error {
	errorTrace := this.getTrace(err)

	myErr := ErrorToErrorCode(err)

	logs.Logger.Error("Error", zap.Strings("err", this.DealErr(myErr)))

	this.Msg = myErr.Error()
	this.Data = errorTrace
	this.send(ctx, myErr.Code())
	return nil
}

//输出
func (this *Response) send(ctx *fiber.Ctx, code int) error {
	this.Code = code
	var err error
	if this.Code != 0 {
		err = ctx.Status(http.StatusAccepted).JSON(this)
	} else {
		err = ctx.Status(http.StatusOK).JSON(this)
	}

	if err != nil {
		ctx.Status(http.StatusAccepted).JSON(map[string]interface{}{"msg": err, "code": 500})
	}
	return nil
}

//输出
func (this *Response) Output(write io.Writer, data map[string]interface{}) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b, _ := json.Marshal(data)
	write.Write(b)
	return nil
}

//得到trace信息
func (this *Response) getTrace(err error) []string {
	goEnv := os.Getenv("GO_ENV")
	var errorTrace []string
	if goEnv == "product" {
		errorTrace = this.DealErr(err)
	}
	return errorTrace
}

//处理异常（业务异常和默认异常）
func ErrorToErrorCode(err error) *MyError {
	if err == nil {
		return nil
	}
	errorCode, ok := err.(*MyError)

	if ok {
		return errorCode
	}
	return NewError(err.Error(), ERROR).(*MyError)
}

func (this *Response) ReturnValOrNull(value, empty interface{}) interface{} {
	var vValue = reflect.ValueOf(value)
	if value == nil || (vValue.Kind() == reflect.Slice && vValue.Len() == 0) {
		return empty
	}
	return value
}

func (this *Response) SliceReturnValOrNull(value []string, empty interface{}) interface{} {
	if value == nil || len(value) == 0 {
		return empty
	}
	return value
}

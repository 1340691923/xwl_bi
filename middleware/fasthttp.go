package middleware

import (
	"fmt"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"runtime"
	"strings"
	"time"
)

func Cors(handle fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if r := recover(); r != nil {
				//打印调用栈信息
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				stackInfo := fmt.Sprintf("%s", buf[:n])
				logs.Logger.Sugar().Errorf("panic stack info %s", stackInfo)
				logs.Logger.Sugar().Errorf("--->HaveLoginUserSign Error:", r)
				ctx.Response.SetBodyString(`{"code":500}`)
			}
		}()
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type") //header的类型
		handle(ctx)
	}
}

func FTimer(handle fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		startT := time.Now()
		handle(ctx)
		logs.Logger.Info("handle lost time", zap.String("time", time.Now().Sub(startT).String()))
	}
}

func WechatSpider(handle fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		UserAgent := strings.ToLower(util.Bytes2str(ctx.Request.Header.UserAgent()))

		if util.InstrArr(model.GlobConfig.Report.UserAgentBanList, UserAgent) {
			logs.Logger.Error("WechatSpider", zap.String("该UserAgent禁止访问接口！", ctx.Request.Header.String()))
			util.WriteJSON(ctx, map[string]interface{}{
				"code": 500,
				"msg":  "该UserAgent禁止访问接口！",
			})
			return
		}
		handle(ctx)
	}
}

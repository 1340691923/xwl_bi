/***
该项目用于采集用户行为数据 以及 埋点数据分析服务
*/
package main

import (
	"flag"
	"fmt"
	"github.com/1340691923/xwl_bi/application"
	"github.com/1340691923/xwl_bi/controller"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/middleware"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/sinker"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/buaazp/fasthttprouter"
	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"go.uber.org/zap"
	"log"
	"net/http/pprof"
	"time"
)

var (
	configFileDir  string
	configFileName string
	configFileExt  string
)

func init() {
	flag.StringVar(&configFileDir, "configFileDir", "config", "配置文件夹名")
	flag.StringVar(&configFileName, "configFileName", "config", "配置文件名")
	flag.StringVar(&configFileExt, "configFileExt", "json", "配置文件后缀")
	flag.Parse()
}

// By 肖文龙
func main() {
	app := application.NewApp(
		"report_server",
		application.WithConfigFileDir(configFileDir),
		application.WithConfigFileName(configFileName),
		application.WithConfigFileExt(configFileExt),
		application.RegisterInitFnObserver(application.InitLogs),
		application.RegisterInitFnObserver(application.InitKafkaSyncProduce),
		application.RegisterInitFnObserver(application.InitKafkaAsyncProduce),
		application.RegisterInitFnObserver(application.InitRedisPool),
		application.RegisterInitFnObserver(application.InitMysql),
		application.RegisterInitFnObserver(application.InitClickHouse),
		application.RegisterInitFnObserver(application.RefreshTableId),
	)

	err := app.InitConfig().NotifyInitFnObservers().Error()

	if err != nil {
		logs.Logger.Error("数据系统 初始化失败", zap.Error(err))
		panic(err)
	}

	defer app.Close()

	go func() {
		for {
			select {
			case err := <-db.KafkaASyncProducer.Errors():
				logs.Logger.Error(" db.KafkaASyncProducer.Errors", zap.Error(err))
				time.Sleep(time.Hour)
			default:

			}
		}
	}()
	go sinker.ClearDimsCacheByTimeBylocal(time.Second * 20)

	router := fasthttprouter.New()

	router.GET("/debug/pprof/", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Index))
	router.GET("/debug/pprof/cmdline", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Cmdline))
	router.GET("/debug/pprof/profile", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Profile))
	router.GET("/debug/pprof/symbol", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Symbol))
	router.GET("/debug/pprof/trace", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Trace))
	router.GET("/debug/pprof/allocs", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Handler("allocs").ServeHTTP))
	router.GET("/debug/pprof/block", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Handler("block").ServeHTTP))
	router.GET("/debug/pprof/goroutine", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Handler("goroutine").ServeHTTP))
	router.GET("/debug/pprof/heap", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Handler("heap").ServeHTTP))
	router.GET("/debug/pprof/mutex", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Handler("mutex").ServeHTTP))
	router.GET("/debug/pprof/threadcreate", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Handler("threadcreate").ServeHTTP))
	router.POST("/test", func(ctx *fasthttp.RequestCtx) {
		ctx.WriteString(`{"code":0}`)
	})

	router.GET("/GetWordParse", controller.GetWordParse)

	//写着写着变成了flutter的嵌套语法哈哈哈
	router.POST(
		"/sync_json/:typ/:appid/:appkey/:eventName/:debug",
		middleware.Cors(
			middleware.WechatSpider(
				controller.ReportController{}.ReportAction,
			),
		),
	)

	server := &fasthttp.Server{
		Handler: router.Handler,
	}
	if model.GlobConfig.Report.ReadTimeout != 0 {
		server.ReadTimeout = time.Duration(model.GlobConfig.Report.ReadTimeout) * time.Second
	}

	if model.GlobConfig.Report.WriteTimeout != 0 {
		server.WriteTimeout = time.Duration(model.GlobConfig.Report.WriteTimeout) * time.Second
	}

	if model.GlobConfig.Report.MaxConnsPerIP != 0 {
		server.MaxConnsPerIP = model.GlobConfig.Report.MaxConnsPerIP
	}

	if model.GlobConfig.Report.MaxRequestsPerConn != 0 {
		server.MaxRequestsPerConn = model.GlobConfig.Report.MaxRequestsPerConn
	}

	if model.GlobConfig.Report.IdleTimeout != 0 {
		server.IdleTimeout = time.Duration(model.GlobConfig.Report.IdleTimeout) * time.Second
	}

	go func() {
		port := fmt.Sprintf(":%v", model.GlobConfig.Report.ReportPort)
		logs.Logger.Sugar().Infof("service start", port)
		log.Println(fmt.Sprintf("上报服务启动成功 ,性能检测入口为: http://127.0.0.1:%v", model.GlobConfig.Report.ReportPort))
		if err = server.ListenAndServe(port); err != nil {
			logs.Logger.Error("service err", zap.Error(err))
			log.Panic(err)
		}
	}()

	app.WaitForExitSign(func() {
		logs.Logger.Sugar().Infof("数据上报服务停止中...")
		if err = server.Shutdown(); err != nil {
			logs.Logger.Sugar().Infof("数据上报服务停止失败 err", zap.Error(err))
		} else {
			logs.Logger.Sugar().Infof("数据上报服务停止成功...")
		}
	})

}

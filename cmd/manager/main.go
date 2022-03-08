package main

import (
	"flag"
	"github.com/1340691923/xwl_bi/application"
	"github.com/1340691923/xwl_bi/engine/logs"
	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
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

// 管理后台
// By 肖文龙
func main() {

	app := application.NewApp(
		"manager",
		application.WithConfigFileDir(configFileDir),
		application.WithConfigFileName(configFileName),
		application.WithConfigFileExt(configFileExt),
		application.RegisterInitFnObserver(application.InitLogs),
		application.RegisterInitFnObserver(application.InitMysql),
		application.RegisterInitFnObserver(application.InitTask),
		application.RegisterInitFnObserver(application.InitRbac),
		application.RegisterInitFnObserver(application.InitOpenWinBrowser),
		application.RegisterInitFnObserver(application.InitClickHouse),
		application.RegisterInitFnObserver(application.InitRedisPool),
		application.RegisterInitFnObserver(application.InitDebugSarama),
	)

	err := app.
		InitConfig().
		NotifyInitFnObservers().
		Error()

	if err != nil {
		logs.Logger.Error("BI 初始化失败", zap.String("err.Error()", err.Error()))
		panic(err)
	}

	defer app.Close()

	app.RunManager()

	app.WaitForExitSign(func() {
		logs.Logger.Sugar().Infof("BI 服务停止成功...")
	})

}

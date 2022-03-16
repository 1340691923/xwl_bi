package main

import (
	"flag"
	"fmt"
	"github.com/1340691923/xwl_bi/application"
	"github.com/1340691923/xwl_bi/cmd/init_app/ck"
	"github.com/1340691923/xwl_bi/cmd/init_app/kafka"
	"github.com/1340691923/xwl_bi/cmd/init_app/mysql"
	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/go-sql-driver/mysql"
	"log"
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

// 初始化程序
// By 肖文龙
func main() {

	app := application.NewApp(
		"init_app",
		application.WithConfigFileDir(configFileDir),
		application.WithConfigFileName(configFileName),
		application.WithConfigFileExt(configFileExt),
		application.RegisterInitFnObserver(application.InitLogs),
		application.RegisterInitFnObserver(application.InitMysql),
		application.RegisterInitFnObserver(application.InitClickHouse),
	)

	err := app.InitConfig().
		NotifyInitFnObservers().
		Error()

	if err != nil {
		log.Println(fmt.Sprintf("初始化失败%s", err.Error()))
		panic(err)
	}

	defer app.Close()

	kafka.Init()
	ck.Init()
	mysql.Init()
	log.Println("数据已全部初始化完毕！")
}

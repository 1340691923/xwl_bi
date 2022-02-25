package main

import (
	"flag"
	"fmt"
	"github.com/1340691923/xwl_bi/application"
	"log"
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

func init(){
	app := application.NewApp(
		"init_app",
		application.WithConfigFileDir(configFileDir),
		application.WithConfigFileName(configFileName),
		application.WithConfigFileExt(configFileExt),
		application.RegisterInitFnObserver(application.InitLogs),
		application.RegisterInitFnObserver(application.InitRedisPool),
	)

	err := app.InitConfig().
		NotifyInitFnObservers().
		Error()

	if err != nil {
		log.Println(fmt.Sprintf("初始化失败%s",err.Error()))
		panic(err)
	}
}

func main(){
	chExit := make(chan int,1)
	go func() {
		time.Sleep(1 * time.Second)
		chExit <- 1
		time.Sleep(5 * time.Second)
		close(chExit)
	}()

	for  {
		select {
		case v, ok := <-chExit:
			if !ok {
				fmt.Println("close channel 2", v)
				goto EXIT2
			}

			fmt.Println("ch2 val =", v)
		}
	}

	EXIT2:
	fmt.Println("exit testSelectFor2")
	select {
	default:

	}
}
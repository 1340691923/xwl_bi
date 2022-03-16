package application

import (
	"fmt"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"github.com/1340691923/xwl_bi/router"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/viper"
)

type InitFnObserver func() (deferFunc func(), err error)

// Options方法
type NewAppOptions func(app *App)

// App 结构体 启动应用基本配置
type App struct {
	configFileDir,
	configFileName,
	configFileExt,
	appName string
	cmdName         string
	InitFnObservers []InitFnObserver
	err             error
	deferFuncs      []func()
}

// 设置配置文件格式   例如:json,conf 等等
func RegisterInitFnObserver(fn InitFnObserver) NewAppOptions {
	return func(app *App) {
		app.InitFnObservers = append(app.InitFnObservers, fn)
	}
}

// 设置配置文件格式   例如:json,conf 等等
func WithConfigFileExt(configFileExt string) NewAppOptions {
	return func(app *App) {
		app.configFileExt = configFileExt
	}
}

// 设置配置文件目录
func WithConfigFileDir(configFileDir string) NewAppOptions {
	return func(app *App) {
		app.configFileDir = configFileDir
	}
}

// 设置配置文件名
func WithConfigFileName(configFileName string) NewAppOptions {
	return func(app *App) {
		app.configFileName = configFileName
	}
}

// 设置应用名
func WithCmdName(cmdName string) NewAppOptions {
	return func(app *App) {
		app.cmdName = cmdName
	}
}

// App 构造方法
func NewApp(cmdName string, opts ...NewAppOptions) *App {
	app := &App{
		configFileDir:  "config",
		configFileName: "config.json",
		appName:        "铸龙-BI",
		cmdName:        cmdName,
	}
	for _, opt := range opts {
		opt(app)
	}
	return app
}

// 初始化配置
func (this *App) InitConfig() *App {
	config := viper.New()
	config.AddConfigPath(this.configFileDir)
	config.SetConfigName(this.configFileName)
	config.SetConfigType(this.configFileExt)
	if err := config.ReadInConfig(); err != nil {
		log.Println("GlobConfig err", err)
		this.err = err
		return this
	}

	if err := config.Unmarshal(&model.GlobConfig); err != nil {
		this.err = err
		return this
	}

	model.CmdName = this.cmdName

	return this
}

func (this *App) NotifyInitFnObservers() *App {
	this.deferFuncs = []func(){}

	for _, fnObserver := range this.InitFnObservers {
		var fn func()
		fn, this.err = fnObserver()
		if this.err != nil {
			return this
		}
		this.deferFuncs = append(this.deferFuncs, fn)
	}
	return this
}

//关闭app
func (this *App) Close() {
	for _, fn := range this.deferFuncs {
		fn()
	}
}

// 获取配置文件夹
func (this *App) getConfigDir() string {
	return filepath.Join(
		util.GetCurrentDirectory(),
		this.configFileDir,
	)
}

// 是否有异常
func (this *App) Error() (err error) {
	return this.err
}

func (this *App) RunManager() {
	appServer := router.Init()

	go func() {
		if err := appServer.Listen(fmt.Sprintf(":%v", model.GlobConfig.Manager.Port)); err != nil {
			logs.Logger.Error("BI 后台服务 http服务启动失败:", zap.String("err.Error()", err.Error()))
			log.Panic(err)
		}
	}()
}

func (this *App) WaitForExitSign(exitFunc ...func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	<-c
	for index := range exitFunc {
		exitFunc[index]()
	}
}

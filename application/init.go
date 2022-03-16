package application

import (
	"fmt"
	"github.com/1340691923/xwl_bi/controller"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/rbac"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/report"
	"github.com/1340691923/xwl_bi/platform-basic-libs/sinker"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"github.com/gofiber/websocket/v2"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"io"
	"log"
	"time"

	"path/filepath"
	"strconv"
)

// 初始化日志
func InitLogs() (fn func(), err error) {
	logger := logs.NewLog(
		logs.WithLogPath(filepath.Join(model.GlobConfig.Comm.Log.LogDir, model.CmdName)),
		logs.WithStorageDays(model.GlobConfig.Comm.Log.StorageDays),
	)
	logs.Logger, err = logger.InitLog()
	if err != nil {
		return
	}
	log.Println(fmt.Sprintf("日志组件初始化成功！日志所在目录：%v，保存天数为：%v", model.GlobConfig.Comm.Log.LogDir, model.GlobConfig.Comm.Log.StorageDays))
	fn = func() {}
	return
}

// 初始化mysql连接
func InitMysql() (fn func(), err error) {
	config := model.GlobConfig.Comm.Mysql
	dbSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		config.Username,
		config.Pwd,
		config.IP,
		config.Port,
		config.DbName)
	db.Sqlx, err = db.NewSQLX(
		"mysql",
		dbSource,
		config.MaxOpenConns,
		config.MaxIdleConns,
	)
	if err != nil {
		return
	}
	log.Println(fmt.Sprintf("Mysql组件初始化成功！连接：%v，最大打开连接数：%v，最大等待连接数:%v",
		dbSource,
		config.MaxOpenConns,
		config.MaxIdleConns,
	))
	fn = func() {}
	return
}

// 初始化mysql连接
func InitClickHouse() (fn func(), err error) {
	config := model.GlobConfig.Comm.ClickHouse
	dbSource := fmt.Sprintf(
		"tcp://%s:%s?username=%s&password=%s&database=%s&compress=true",
		config.IP,
		config.Port,
		config.Username,
		config.Pwd,
		config.DbName,
	)
	db.ClickHouseSqlx, err = db.NewSQLX(
		"clickhouse",
		dbSource,
		config.MaxOpenConns,
		config.MaxIdleConns,
	)

	if err != nil {
		return
	}
	log.Println(fmt.Sprintf("ClickHouse组件初始化成功！连接：%v，最大打开连接数：%v，最大等待连接数:%v",
		dbSource,
		config.MaxOpenConns,
		config.MaxIdleConns,
	))
	fn = func() {}
	return
}

// 初始化redis
func InitRedisPool() (fn func(), err error) {
	config := model.GlobConfig.Comm.Redis

	db.RedisPool = db.NewRedisPool(config.Addr, config.Passwd, config.Db, config.MaxIdle, config.MaxActive)

	log.Println(fmt.Sprintf("Redis组件初始化成功！连接：%v，DB：%v，密码:%v MaxIdle:%v MaxActive:%v",
		config.Addr,
		config.Db,
		config.Passwd,
		config.MaxIdle,
		config.MaxActive,
	))
	fn = func() {}
	return
}

// 初始化项目启动任务
func InitTask() (fn func(), err error) {
	fn = func() {}
	return
}

// 初始化项目启动任务
func InitRbac() (fn func(), err error) {
	config := model.GlobConfig.Comm.Mysql
	dbSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		config.Username,
		config.Pwd,
		config.IP,
		config.Port,
		config.DbName)
	err = rbac.Run("mysql", dbSource)
	if err != nil {
		return
	}
	log.Println(fmt.Sprintf("Rbac组件初始化成功！连接：%v",
		dbSource,
	))
	return
}

// 掉起浏览器
func InitOpenWinBrowser() (fn func(), err error) {
	config := model.GlobConfig
	if !config.Manager.DeBug {
		port := ":" + strconv.Itoa(int(config.Manager.Port))
		uri := fmt.Sprintf("%s%s", "http://127.0.0.1", port)
		util.OpenWinBrowser(uri)
		log.Println(fmt.Sprintf("将打开浏览器！地址为：%v",
			uri,
		))
	}
	fn = func() {}
	return
}

//初始化kafka异步生产者
func InitKafkaAsyncProduce() (fn func(), err error) {
	config := model.GlobConfig.Comm.Kafka
	conn, err := db.NewKafkaAsyncProduce(config.Addresses, config.Username, config.Password)
	if err != nil {
		return
	}
	db.KafkaASyncProducer = conn
	fn = func() {
		log.Println("KafkaASyncProducer 关闭了")
		db.KafkaASyncProducer.Close()
	}
	return
}

//初始化kafka同步步生产者
func InitKafkaSyncProduce() (fn func(), err error) {
	config := model.GlobConfig.Comm.Kafka
	conn, err := db.NewKafkaSyncProduce(config.Addresses, config.Username, config.Password)
	if err != nil {
		return
	}
	db.KafkaSyncProducer = conn

	fn = func() {
		log.Println("KafkaSyncProducer 关闭了")
		db.KafkaSyncProducer.Close()
	}

	return
}

//
func InitDebugSarama() (fn func(), err error) {
	debugSarama := sinker.NewKafkaSarama()
	err = debugSarama.Init(model.GlobConfig.Comm.Kafka, model.GlobConfig.Comm.Kafka.DebugDataTopicName, model.GlobConfig.Comm.Kafka.DebugDataGroup, func(msg model.InputMessage, markFn func()) {

		distinctId := gjson.GetBytes(msg.Value, "distinct_id").String()

		managerMap, ok := controller.ConnUUidMap.Load(distinctId)

		if ok {
			managerMap.(*controller.ManagerConnMap).Conns.Range(func(key, value interface{}) bool {

				if err := value.(*websocket.Conn).WriteJSON(map[string]interface{}{
					"code": 1,
					"data": util.Bytes2str(msg.Value),
				}); err != nil {
					if err == io.EOF {
						logs.Logger.Error("客户端已经断开WsSocket!", zap.Error(err))
					} else if err.Error() == "use of closed network connection" {
						logs.Logger.Error("服务端已经断开WsSocket!", zap.Error(err))
					} else {
						logs.Logger.Error("socket err!", zap.Error(err))
					}
					managerMap.(*controller.ManagerConnMap).DeleteConn(key.(string))
					controller.ConnUUidMap.Store(distinctId, managerMap)
				}
				return true
			})
		}

	}, func() {

	})
	if err != nil {
		return
	}

	go debugSarama.Run()

	log.Println(fmt.Sprintf("Debug数据消费者启动！"))

	fn = func() {
		debugSarama.Stop()
	}
	return
}

func RefreshTableId() (fn func(), err error) {
	fn = func() {
		go report.RefreshTableIdMap(5 * time.Minute)
	}
	return
}

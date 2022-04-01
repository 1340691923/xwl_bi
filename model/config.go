//应用启动引擎层
package model

import (
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"os"
)

var GlobConfig Config

var CmdName string

//全局配置结构体
type Config struct {
	Manager ManagerConfig `json:"manager"`
	Report  ReportConfig  `json:"report"`
	Sinker  SinkerConfig  `json:"sinker"`
	Comm    struct {
		Log        LogConfig        `json:"log"`
		Mysql      MysqlConfig      `json:"mysql"`
		ClickHouse ClickHouseConfig `json:"clickhouse"`
		Kafka      KafkaCfg         `json:"kafka"`
		Redis      RedisConfig      `json:"redis"`
	} `json:"comm"`
}

type ManagerConfig struct {
	Port              uint16 `json:"port"`              //铸龙分析系统http启动端口
	CkQueryLimit      int    `json:"ckQueryLimit"`      //clickhouse 查询限流器阈值
	CkQueryExpiration int    `json:"ckQueryExpiration"` //clickhouse 查询限流器阈值
	JwtSecret         string `json:"jwtSecret"`
	DeBug             bool   `json:"deBug"`
}

type SinkerConfig struct {
	ReportAcceptStatus  BatchConfig `json:"reportAcceptStatus"`
	ReportData2CK       BatchConfig `json:"reportData2CK"`
	RealTimeWarehousing BatchConfig `json:"realTimeWarehousing"`
	PprofHttpPort       uint16      `json:"pprofHttpPort"`
}

type RedisConfig struct {
	Addr      string `json:"addr"`
	Passwd    string `json:"passwd"`
	Db        int    `json:"db"`
	MaxIdle   int    `json:"maxIdle"`
	MaxActive int    `json:"maxActive"`
}

type ClickHouseConfig struct {
	Username             string `json:"username"`
	Pwd                  string `json:"pwd"`
	IP                   string `json:"ip"`
	Port                 string `json:"port"`
	DbName               string `json:"dbName"`
	MaxOpenConns         int    `json:"maxOpenConns"`
	MaxIdleConns         int    `json:"maxIdleConns"`
	MacrosShardKeyName   string `json:"macrosShardKeyName"`
	MacrosReplicaKeyName string `json:"macrosReplicaKeyName"`
	ClusterName          string `json:"clusterName"`
}

type MysqlConfig struct {
	Username     string `json:"username"`
	Pwd          string `json:"pwd"`
	IP           string `json:"ip"`
	Port         string `json:"port"`
	DbName       string `json:"dbName"`
	MaxOpenConns int    `json:"maxOpenConns"`
	MaxIdleConns int    `json:"maxIdleConns"`
}

type ReportConfig struct {
	ReportPort         uint16   `json:"reportPort"` //上报程序启动端口
	ReadTimeout        int      `json:"readTimeout"`
	WriteTimeout       int      `json:"writeTimeout"`
	MaxConnsPerIP      int      `json:"maxConnsPerIP"`
	MaxRequestsPerConn int      `json:"maxRequestsPerConn"`
	IdleTimeout        int      `json:"idleTimeout"`
	UserAgentBanList   []string `json:"userAgentBanList"`
}

type LogConfig struct {
	StorageDays int    `json:"storageDays"` //日志保留天数
	LogDir      string `json:"logDir"`      //日志保留文件夹地址
}

func (this *Config) GetCkQueryLimit() int {
	if this.Manager.CkQueryLimit == 0 {
		return 10
	}
	return this.Manager.CkQueryLimit
}

func (this *Config) GetCkQueryExpiration() int {
	if this.Manager.CkQueryExpiration == 0 {
		return 2
	}
	return this.Manager.CkQueryExpiration
}

func (this *Config) GetKafkaCfgProducerType() string {
	if this.Comm.Kafka.ProducerType == "" {
		return "sync"
	}
	return this.Comm.Kafka.ProducerType
}

type KafkaCfg struct {
	NumPartitions      int32    `json:"numPartitions"`
	Addresses          []string `json:"addresses"`
	Username           string   `json:"username"`
	Password           string   `json:"password"`
	ReportTopicName    string   `json:"reportTopicName"`
	ConsumerGroupName  string   `json:"consumerGroupName"`
	RealTimeDataGroup  string   `json:"realTimeDataGroup"`
	ReportData2CKGroup string   `json:"reportData2CKGroup"`
	DebugDataTopicName string   `json:"debugDataTopicName"`
	DebugDataGroup     string   `json:"debugDataGroup"`
	ProducerType       string   `json:"producer_type"`
}

type BatchConfig struct {
	BufferSize    int `json:"bufferSize"`
	FlushInterval int `json:"flushInterval"`
}

//下载配置文件
func DownloadConfigFile(fname string) (err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var config Config
	filePtr, err := os.Create(fname)
	if err != nil {
		return errors.New(fmt.Sprintf("创建配置文件异常:%s", err.Error()))
	}
	defer filePtr.Close()
	// 带JSON缩进格式写文件
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return errors.New(fmt.Sprintf("创建配置文件异常:%s", err.Error()))
	}
	_, err = filePtr.Write(data)
	return
}

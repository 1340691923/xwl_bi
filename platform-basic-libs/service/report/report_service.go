package report

import (
	"bytes"
	"encoding/json"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/my_error"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/myapp"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"github.com/Shopify/sarama"
	"github.com/garyburd/redigo/redis"
	"go.uber.org/zap"
	"sync"
	"time"
)

type ReportService struct {
}

var tableIdMap sync.Map

func RefreshTableIdMap(t time.Duration) {
	for {
		time.Sleep(t)
		tableIdMap.Range(func(key, value interface{}) bool {
			tableIdMap.Delete(key)
			return true
		})
	}
}

func (this *ReportService) GetTableid(appid, appkey string) (table string, err error) {
	buff := new(bytes.Buffer)
	buff.WriteString(appid)
	buff.WriteString("_xwl_")
	buff.WriteString(appkey)
	key := buff.String()

	if val, found := tableIdMap.Load(key); found {
		table = val.(string)
		return
	}

	conn := db.RedisPool.Get()
	defer conn.Close()

	if table, err = myapp.GetAppidToTableid(conn, key); err != nil {
		if err == redis.ErrNil {
			err = my_error.NewBusiness(ERROR_TABLE, AppParmasErr)
			return
		} else {
			logs.Logger.Error("GetTableid", zap.Error(err))
			err = my_error.NewBusiness(ERROR_TABLE, ServerErr)
			return
		}
	}

	tableIdMap.Store(key, table)
	return
}

const (
	DebugToDB    = "1"
	DebugNotToDB = "2"
)

func (this *ReportService) IsDebugUser(debug, xwlDistinctId, tableId string) bool {
	debugArr := []string{DebugToDB, DebugNotToDB}

	if !util.InstrArr(debugArr, debug) {
		return false
	}
	Hash := "DebugDeviceID_" + tableId

	conn := db.RedisPool.Get()
	defer conn.Close()

	i, debugErr := redis.Int(conn.Do("SISMEMBER", Hash, xwlDistinctId))
	if util.FilterRedisNilErr(debugErr) {
		logs.Logger.Error("debugErr.err", zap.Error(debugErr))
		return false
	}
	if i < 1 {
		return false
	}

	return true
}

func (this *ReportService) CantInflowOfKakfa(debug, xwlDistinctId, tableId string) bool {
	debugArr := []string{"1", "2"}

	if !util.InstrArr(debugArr, debug) {
		return false
	}
	Hash := "DebugDeviceID_" + tableId

	conn := db.RedisPool.Get()
	defer conn.Close()

	i, debugErr := redis.Int(conn.Do("SISMEMBER", Hash, xwlDistinctId))
	if util.FilterRedisNilErr(debugErr) {
		logs.Logger.Error("debugErr.err", zap.Error(debugErr))
		return false
	}
	if i < 1 {
		return false
	}

	return true
}

func sendMsg(msg *sarama.ProducerMessage) (err error) {
	switch model.GlobConfig.GetKafkaCfgProducerType() {
	case "async":
		db.KafkaASyncProducer.Input() <- msg
	case "sync":
		_, _, err = db.KafkaSyncProducer.SendMessage(msg)
	}
	return
}

func (this *ReportService) InflowOfDebugData(data map[string]interface{}, eventName string) (err error) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = model.GlobConfig.Comm.Kafka.DebugDataTopicName
	sendData, _ := json.Marshal(data)
	msg.Value = sarama.ByteEncoder(sendData)
	msg.Timestamp = time.Now()

	return sendMsg(msg)

}

package report

import (
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/my_error"
	model2 "github.com/1340691923/xwl_bi/platform-basic-libs/sinker/model"
	"sync"

	"github.com/Shopify/sarama"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type ReportInterface interface {
	NewReportType(appid, tableId, debug, timeNow, eventName, ip string, body []byte)
	GetkafkaData() model.KafkaData
	InflowOfKakfa() (err error)
	Put()
}

type UserReport struct {
	kafkaData model.KafkaData
}

var userPool = sync.Pool{
	New: func() interface{} {
		return &UserReport{}
	},
}

var eventPool = sync.Pool{
	New: func() interface{} {
		return &EventReport{}
	},
}

func (this *UserReport) NewReportType(appid, tableId, debug, timeNow, eventName, ip string, body []byte) {
	this.kafkaData.APPID = appid
	this.kafkaData.TableId = tableId
	this.kafkaData.Debug = debug
	this.kafkaData.ReqData = body
	this.kafkaData.Ip = ip
	this.kafkaData.ReportTime = timeNow
	this.kafkaData.ReportType = model.UserReportType
	this.kafkaData.EventName = "用户属性"
}

func (this *UserReport) GetkafkaData() model.KafkaData {
	return this.kafkaData
}

func (this *UserReport) InflowOfKakfa() (err error) {

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	msg := &sarama.ProducerMessage{}
	msg.Topic = model.GlobConfig.Comm.Kafka.ReportTopicName
	sendData, _ := json.Marshal(this.kafkaData)
	msg.Value = sarama.ByteEncoder(sendData)
	msg.Timestamp = time.Now()

	return sendMsg(msg)
}

func (this *UserReport) Put() {
	userPool.Put(this)
}

type EventReport struct {
	kafkaData model.KafkaData
}

func (this *EventReport) NewReportType(appid, tableId, debug, timeNow, eventName, ip string, body []byte) {
	this.kafkaData.APPID = appid
	this.kafkaData.TableId = tableId
	this.kafkaData.Debug = debug
	this.kafkaData.ReqData = body
	this.kafkaData.ReportTime = timeNow
	this.kafkaData.ReportType = model.EventReportType
	this.kafkaData.EventName = eventName
	this.kafkaData.Ip = ip
}

func (this *EventReport) InflowOfKakfa() (err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	msg := &sarama.ProducerMessage{}
	msg.Topic = model.GlobConfig.Comm.Kafka.ReportTopicName
	sendData, _ := json.Marshal(this.kafkaData)
	msg.Value = sarama.ByteEncoder(sendData)
	msg.Timestamp = time.Now()

	return sendMsg(msg)
}

func (this *EventReport) GetkafkaData() model.KafkaData {
	return this.kafkaData
}

func (this *EventReport) Put() {
	eventPool.Put(this)
}

func NewEventReport() ReportInterface {
	return eventPool.Get().(*EventReport)
}

func NewUserReport() ReportInterface {
	return userPool.Get().(*UserReport)
}

var duckMap = map[string]func() ReportInterface{
	model2.ReportUserProperties:  NewUserReport,
	model2.ReportEventProperties: NewEventReport,
}

func GetReportDuck(typ string) (reportInterface ReportInterface, err error) {
	var ok bool
	var fn func() ReportInterface
	if fn, ok = duckMap[typ]; !ok {
		err = my_error.NewBusiness(ERROR_TABLE, ReportTypeErr)
		return
	}
	reportInterface = fn()
	return
}

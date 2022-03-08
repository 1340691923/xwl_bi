package report

import (
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/my_error"
	model2 "github.com/1340691923/xwl_bi/platform-basic-libs/sinker/model"
	"sync"

	"github.com/Shopify/sarama"
	"time"
)

type ReportInterface interface {
	NewReportType(data *ReportTypeData)
	GetkafkaData() model.KafkaData
	InflowOfKakfa(marshaler func(v interface{}) ([]byte, error)) (err error)
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

type ReportTypeData struct {
	Appid string
	TableId string
	Debug string
	TimeNow string
	EventName string
	Ip string
	Body []byte
}

func (this *UserReport) NewReportType(data *ReportTypeData) {
	this.kafkaData.APPID = data.Appid
	this.kafkaData.TableId = data.TableId
	this.kafkaData.Debug = data.Debug
	this.kafkaData.ReqData = data.Body
	this.kafkaData.Ip = data.Ip
	this.kafkaData.ReportTime = data.TimeNow
	this.kafkaData.ReportType = model.UserReportType
	this.kafkaData.EventName = "用户属性"
}

func (this *UserReport) GetkafkaData() model.KafkaData {
	return this.kafkaData
}

func (this *UserReport) InflowOfKakfa(marshaler func(v interface{}) ([]byte, error)) (err error) {

	msg := &sarama.ProducerMessage{}
	msg.Topic = model.GlobConfig.Comm.Kafka.ReportTopicName
	sendData, _ := marshaler(this.kafkaData)
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

func (this *EventReport) NewReportType(data *ReportTypeData) {
	this.kafkaData.APPID = data.Appid
	this.kafkaData.TableId = data.TableId
	this.kafkaData.Debug = data.Debug
	this.kafkaData.ReqData = data.Body
	this.kafkaData.ReportTime = data.TimeNow
	this.kafkaData.ReportType = model.EventReportType
	this.kafkaData.EventName = data.EventName
	this.kafkaData.Ip = data.Ip
}

func (this *EventReport) InflowOfKakfa(marshaler func(v interface{}) ([]byte, error)) (err error) {

	msg := &sarama.ProducerMessage{}
	msg.Topic = model.GlobConfig.Comm.Kafka.ReportTopicName
	sendData, _ := marshaler(this.kafkaData)

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

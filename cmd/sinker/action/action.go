package action

import (
	"bytes"
	"errors"
	"fmt"
	"go.uber.org/zap"

	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/consumer_data"
	"github.com/1340691923/xwl_bi/platform-basic-libs/sinker"
	parser "github.com/1340691923/xwl_bi/platform-basic-libs/sinker/parse"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"github.com/valyala/fastjson"
	"strconv"
	"strings"
	"sync"
)

const (
	PresetAttribute  = 1
	CustomAttribute  = 2
	IsUserAttribute  = 1
	IsEventAttribute = 2
)

var MetaAttrRelationSet sync.Map
var AttributeMap sync.Map
var MetaEventMap sync.Map

var metaAttrRelationChan = make(chan map[string]interface{}, 1000)
var attributeChan = make(chan map[string]interface{}, 1000)
var metaEventChan = make(chan map[string]interface{}, 1000)

func MysqlConsumer() {
	for {
		select {
		case m := <-metaAttrRelationChan:
			if _, err := db.SqlBuilder.Insert("meta_attr_relation").SetMap(m).RunWith(db.Sqlx).Exec(); err != nil && !strings.Contains(err.Error(), "1062") {
				logs.Logger.Error("meta_attr_relation insert", zap.Error(err))
			}
		case m := <-attributeChan:
			if _, err := db.SqlBuilder.Insert("attribute").SetMap(m).RunWith(db.Sqlx).Exec(); err != nil && !strings.Contains(err.Error(), "1062") {
				logs.Logger.Error("attribute insert", zap.Error(err))
			}
		case m := <-metaEventChan:
			_, err := db.SqlBuilder.Insert("meta_event").SetMap(m).RunWith(db.Sqlx).Exec()
			if err != nil && !strings.Contains(err.Error(), "1062") {
				logs.Logger.Error("metaEvent insert", zap.Error(err))
			}
		default:

		}
	}
}

func AddRealTimeData(kafkaData model.KafkaData, data string, realTimeWarehousing *consumer_data.RealTimeWarehousing) (err error) {

	clientReportData := consumer_data.ClientReportData{
		Data:    data,
		TableId: kafkaData.TableId,
		Date:    util.Str2Time(kafkaData.ReportTime, util.TimeFormat).Format(util.TimeFormatDay4),
	}
	err = clientReportData.CreateIndex()
	if err != nil {
		logs.Logger.Error(" clientReportData.CreateIndex", zap.Error(err))
	}
	bulkIndexRequest := clientReportData.GetReportData()
	err = realTimeWarehousing.Add(bulkIndexRequest)
	return err
}

func AddMetaEvent(kafkaData model.KafkaData) (err error) {
	if kafkaData.ReportType == model.EventReportType {
		redisConn := db.RedisPool.Get()
		defer redisConn.Close()

		b := bytes.Buffer{}
		b.WriteString(kafkaData.TableId)
		b.WriteString("_")
		b.WriteString(kafkaData.EventName)
		bStr := b.String()

		_, found := MetaEventMap.Load(bStr)

		if !found {
			m := map[string]interface{}{
				"appid":      kafkaData.TableId,
				"event_name": kafkaData.EventName,
			}
			metaEventChan <- m
			MetaEventMap.Store(bStr, struct{}{})
		}
	}
	return nil
}

func AddTableColumn(kafkaData model.KafkaData, failFunc func(data consumer_data.ReportAcceptStatusData), tableName string, ReqDataObject *parser.FastjsonMetric) (err error) {

	dims, err := sinker.GetDims(model.GlobConfig.Comm.ClickHouse.DbName, tableName, nil, db.ClickHouseSqlx)
	if err != nil {
		logs.Logger.Error("sinker.GetDims", zap.Error(err))
		return
	}


	obj, err := ReqDataObject.GetParseObject().Object()
	if err != nil {
		logs.Logger.Error("ReqDataObject.GetParseObject().Object()", zap.Error(err))
		return
	}

	knownKeys := []string{}
	newKeys := new(sync.Map)
	var foundNewKey bool

	tableId, _ := strconv.Atoi(kafkaData.TableId)

	GetReportTypeErr := kafkaData.GetReportTypeErr()

	for _, column := range dims {
		knownKeys = append(knownKeys, column.Name)

		if obj.Get(column.Name) != nil {
			reportType := parser.FjDetectType(obj.Get(column.Name))
			if reportType != column.Type {
				if !(reportType == parser.Int && column.Type == parser.Float) && !(reportType == parser.Float && column.Type == parser.Int) {
					errorReason := fmt.Sprintf("%s的类型错误，正确类型为%v，上报类型为%v(%v)", column.Name, parser.TypeRemarkMap[column.Type], parser.TypeRemarkMap[reportType], obj.Get(column.Name).String())
					failFunc(consumer_data.ReportAcceptStatusData{
						PartDate:       kafkaData.ReportTime,
						TableId:        tableId,
						ReportType:     GetReportTypeErr,
						DataName:       kafkaData.EventName,
						ErrorReason:    errorReason,
						ErrorHandling:  "丢弃数据",
						ReportData:     util.Bytes2str(kafkaData.ReqData),
						XwlKafkaOffset: kafkaData.Offset,
						Status:         consumer_data.FailStatus,
					})
					return errors.New(errorReason)
				}
			}
		}
	}


	b := bytes.Buffer{}

	obj.Visit(func(key []byte, v *fastjson.Value) {

		columnName := util.Bytes2str(key)

		func() {

			b.Reset()
			b.WriteString(kafkaData.TableId)
			b.WriteString("_")
			b.WriteString(kafkaData.EventName)
			b.WriteString("_")
			b.WriteString(columnName)
			bStr := b.String()

			_, found := MetaAttrRelationSet.Load(bStr)

			if !found {
				m := map[string]interface{}{
					"event_name": kafkaData.EventName,
					"event_attr": columnName,
					"app_id":     kafkaData.TableId,
				}
				metaAttrRelationChan <- m
				MetaAttrRelationSet.Store(bStr, struct{}{})
			}

		}()

		if !util.InstrArr(knownKeys, columnName) {
			foundNewKey = true
			newKeys.Store(columnName, parser.FjDetectType(obj.Get(columnName)))
		}

		func() {

			b.Reset()
			b.WriteString(kafkaData.TableId)
			b.WriteString("_xwl_")
			b.WriteString(strconv.Itoa(kafkaData.ReportType))
			b.WriteString("_")
			b.WriteString(columnName)
			AttributeMapkey := b.String()

			_, found := AttributeMap.Load(AttributeMapkey)

			if !found {
				var (
					attributeType, attributeSource int
				)
				if _, ok := parser.SysColumn[columnName]; ok {
					attributeType = PresetAttribute
				} else {
					attributeType = CustomAttribute
				}

				switch kafkaData.ReportType {
				case model.UserReportType:
					attributeSource = IsUserAttribute
				case model.EventReportType:
					attributeSource = IsEventAttribute
				}
				m := map[string]interface{}{
					"attribute_name":   columnName,
					"data_type":        parser.FjDetectType(obj.Get(columnName)),
					"attribute_type":   attributeType,
					"attribute_source": attributeSource,
					"app_id":           kafkaData.TableId,
				}
				attributeChan <- m

				AttributeMap.Store(AttributeMapkey, struct{}{})
			}

		}()

	})

	if foundNewKey {
		dims, err = sinker.ChangeSchema(newKeys, model.GlobConfig.Comm.ClickHouse.DbName, tableName, dims)
		if err != nil {
			logs.Logger.Error("err", zap.Error(err))
		}
		func() {
			redisConn := db.RedisPool.Get()
			defer redisConn.Close()
			dimsCachekey := sinker.GetDimsCachekey(model.GlobConfig.Comm.ClickHouse.DbName, tableName)
			_, err = redisConn.Do("del", dimsCachekey)
			if err != nil {
				logs.Logger.Error("err", zap.Error(err))
			}
		}()
	}

	consumer_data.TableColumnMap.Store(tableName, dims)
	return
}

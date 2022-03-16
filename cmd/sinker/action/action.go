package action

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/consumer_data"
	"github.com/1340691923/xwl_bi/platform-basic-libs/sinker"
	parser "github.com/1340691923/xwl_bi/platform-basic-libs/sinker/parse"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"github.com/valyala/fastjson"
	"go.uber.org/zap"
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

var metaAttrRelationChan = make(chan metaAttrRelationModel, 10000)
var attributeChan = make(chan attributeModel, 10000)
var metaEventChan = make(chan metaEventModel, 10000)

type metaAttrRelationModel struct {
	EventName string
	EventAttr string
	AppId     string
}

type attributeModel struct {
	AttributeName    string
	DataType         int
	AttributeType    int
	attribute_source int
	App_id           string
}

type metaEventModel struct {
	EventName string
	AppId     string
}

func MysqlConsumer() {
	for {
		select {
		case m := <-metaAttrRelationChan:
			if _, err := db.Sqlx.Exec(`insert into  meta_attr_relation(app_id,event_name,event_attr) values (?,?,?);`,
				m.AppId, m.EventName, m.EventAttr); err != nil && !strings.Contains(err.Error(), "1062") {
				logs.Logger.Sugar().Errorf("meta_attr_relation insert", m, err)
			}
		case m := <-attributeChan:
			if _, err := db.Sqlx.Exec(`insert into  attribute(app_id,attribute_source,attribute_type,data_type,attribute_name) values (?,?,?,?,?);`,
				m.App_id, m.attribute_source, m.AttributeType, m.DataType, m.AttributeName); err != nil && !strings.Contains(err.Error(), "1062") {
				logs.Logger.Sugar().Errorf("attribute insert", m, err)
			}
		case m := <-metaEventChan:
			_, err := db.Sqlx.Exec(`insert into  meta_event(appid,event_name) values (?,?);`, m.AppId, m.EventName)
			if err != nil && !strings.Contains(err.Error(), "1062") {
				logs.Logger.Sugar().Errorf("metaEvent insert", m, err)
			}
		default:

		}
	}
}

func AddMetaEvent(kafkaData model.KafkaData) (err error) {
	if kafkaData.ReportType == model.EventReportType {
		b := bytes.Buffer{}
		b.WriteString(kafkaData.TableId)
		b.WriteString("_")
		b.WriteString(kafkaData.EventName)
		bStr := b.String()

		_, found := MetaEventMap.Load(bStr)

		if !found {
			metaEventChan <- metaEventModel{
				EventName: kafkaData.EventName,
				AppId:     kafkaData.TableId,
			}
			MetaEventMap.Store(bStr, struct{}{})
		}
	}
	return nil
}

func AddTableColumn(kafkaData model.KafkaData, failFunc func(data consumer_data.ReportAcceptStatusData), tableName string, ReqDataObject *parser.FastjsonMetric) (err error) {

	dims, err := sinker.GetDims(model.GlobConfig.Comm.ClickHouse.DbName, tableName, nil, db.ClickHouseSqlx, false)
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

		columnName := string(key)

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
				metaAttrRelationChan <- metaAttrRelationModel{
					EventName: kafkaData.EventName,
					EventAttr: columnName,
					AppId:     kafkaData.TableId,
				}
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

				attributeChan <- attributeModel{
					AttributeName:    columnName,
					DataType:         parser.FjDetectType(obj.Get(columnName)),
					AttributeType:    attributeType,
					attribute_source: attributeSource,
					App_id:           kafkaData.TableId,
				}

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
			_, err = redisConn.Do("unlink", dimsCachekey)
			if err != nil {
				_, err = redisConn.Do("del", dimsCachekey)
				if err != nil {
					logs.Logger.Error("err", zap.Error(err))
				}
			}
			sinker.ClearDimsCacheByKey(dimsCachekey)

		}()
	}

	consumer_data.TableColumnMap.Store(tableName, dims)
	return
}

package analysis

import (
	"encoding/json"
	"fmt"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/platform-basic-libs/request"
	"github.com/1340691923/xwl_bi/platform-basic-libs/response"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/analysis/utils"
	parser "github.com/1340691923/xwl_bi/platform-basic-libs/sinker/parse"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type BehaviorAnalysisService struct {
}

func (this *BehaviorAnalysisService) GetConfigs(appid int) (eventNameList []response.MetaEventListRes, attributeMap map[int][]response.AttributeRes, err error) {

	attributeMap = map[int][]response.AttributeRes{}

	if err := db.Sqlx.Select(&eventNameList, "select event_name,show_name from meta_event where appid = ?", appid); err != nil {
		return eventNameList, attributeMap, err
	}

	var attributeRes []response.AttributeRes
	if err := db.Sqlx.Select(&attributeRes, "select attribute_name,show_name,data_type,attribute_type,attribute_source from attribute where app_id = ? and (status = 1 or attribute_type = 1) and attribute_name not in ('xwl_part_date','xwl_kafka_offset','xwl_part_event')  order by attribute_type asc", appid); err != nil {
		return eventNameList, attributeMap, err
	}

	for k, v := range attributeRes {
		if _, ok := parser.TypeRemarkMap[v.DataType]; ok {
			attributeRes[k].DataTypeFormat = parser.TypeRemarkMap[v.DataType]
		}
		if _, ok := parser.SysColumn[v.AttributeName]; ok {
			attributeRes[k].ShowName = parser.SysColumn[v.AttributeName]
		}
		if _, ok := attributeMap[v.AttributeSource]; ok {
			attributeMap[v.AttributeSource] = append(attributeMap[v.AttributeSource], attributeRes[k])
		} else {
			attributeMap[v.AttributeSource] = []response.AttributeRes{attributeRes[k]}
		}
	}

	return
}

type AttributeName struct {
	AttributeName string            `json:"attribute_name" db:"attribute_name"`
	ShowName      string            `json:"show_name" db:"show_name"`
	DataType      int               `json:"data_type" db:"data_type"`
	Analysis      map[string]string `json:"analysis" db:"-"`
}

func (this *BehaviorAnalysisService) LoadPropQuotas(reqData request.LoadPropQuotasReq) (attributeNameList []AttributeName, err error) {

	if err := db.Sqlx.Select(&attributeNameList,
		"select attribute_name,show_name,data_type from attribute "+
			"  where app_id = ?  and (status = 1 or attribute_type = 1) and attribute_source = 2 "+
			"and attribute_name not in ('xwl_part_event','xwl_part_date') and attribute_name in "+
			"(select event_attr from meta_attr_relation where  app_id = ? and event_name = ?)",
		reqData.Appid,
		reqData.Appid,
		reqData.EventName,
	); err != nil {
		return nil, err
	}

	for index, v := range attributeNameList {
		if v.ShowName == "" {
			attributeNameList[index].ShowName = v.AttributeName
		}
		if _, ok := parser.SysColumn[v.AttributeName]; ok {
			attributeNameList[index].ShowName = parser.SysColumn[v.AttributeName]
		}

		switch v.DataType {
		case parser.Int:
			fallthrough
		case parser.Float:
			attributeNameList[index].Analysis = utils.IntPropQuotas
		case parser.String:
			attributeNameList[index].Analysis = utils.StringPropQuotas
		}
	}
	return attributeNameList, nil
}

type ValueStruct struct {
	Value interface{} `json:"value" db:"value"`
}

func (this *BehaviorAnalysisService) GetValues(appid string, table string, col string, reqData []byte) (values []ValueStruct, err error) {

	cache := NewCache(time.Minute*2, fmt.Sprintf("%s_%s_%s_%s", "GetValues", appid, table, col), reqData)

	resData, redisErr := cache.LoadData()

	if util.FilterRedisNilErr(redisErr) {
		return values, err
	}
	if len(resData) > 0 {
		var json = jsoniter.ConfigCompatibleWithStandardLibrary

		err := json.Unmarshal(resData, &values)

		if err != nil {
			return values, err
		}

		return values, err
	}

	tableName := ""
	switch table {
	case "1":
		tableName = "xwl_user" + appid
	case "2":
		tableName = "xwl_event" + appid
	}

	SQL := "select DISTINCT " + col + "  as value from " + tableName + " where  isNotNull(" + col + ") ;"

	err = db.ClickHouseSqlx.Select(&values, SQL)
	if err != nil {
		return values, err
	}

	for index, v := range values {
		switch v.Value.(type) {
		case time.Time:
			values[index].Value = v.Value.(time.Time).Format(util.TimeFormat)
		default:
			break
		}
	}

	resB, err := json.Marshal(values)
	if err != nil {
		return values, err
	}
	cache.SetData(resB)

	return values, err
}

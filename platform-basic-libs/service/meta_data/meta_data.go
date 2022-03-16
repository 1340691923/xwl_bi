package meta_data

import (
	"bytes"
	"fmt"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/platform-basic-libs/request"
	"github.com/1340691923/xwl_bi/platform-basic-libs/response"
	parser "github.com/1340691923/xwl_bi/platform-basic-libs/sinker/parse"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"github.com/garyburd/redigo/redis"
	"strings"
	"time"
)

type MetaDataService struct {
	Appid string `json:"appid"`
}

type YesterdayCountRes struct {
	XwlPartEvent string `db:"xwl_part_event"`
	Count        int    `db:"count"`
}

func (this *MetaDataService) UpdateYesterdayCount() (err error) {
	Hash := "EventYesterdayCountCache"
	redisConn := db.RedisPool.Get()
	defer redisConn.Close()

	lastUpdateDay, redisErr := redis.Int64(redisConn.Do("hget", Hash, this.Appid))

	if util.FilterRedisNilErr(redisErr) {
		return redisErr
	}
	timeNow := time.Now().Unix()

	if util.IsSameDay(lastUpdateDay, timeNow) {
		return nil
	}

	var datas []YesterdayCountRes

	err = db.ClickHouseSqlx.Select(&datas, "SELECT xwl_part_event,count() as count from xwl_event"+this.Appid+" xe where toDateTime(formatDateTime(xwl_part_date,'%Y-%m-%d')) = yesterday() group by xwl_part_event ")

	if err != nil {
		return err
	}

	_, err = db.Sqlx.Exec("update meta_event set yesterday_count = 0  where  appid= ?;", this.Appid)
	if err != nil {
		return err
	}

	if len(datas) > 0 {
		sqlStr := this.createSql(datas)

		_, err = db.Sqlx.Exec(sqlStr, this.Appid)
		if err != nil {
			return err
		}
	}

	_, err = redisConn.Do("hset", Hash, this.Appid, timeNow)
	if err != nil {
		return err
	}
	return nil
}

func (this *MetaDataService) createSql(datas []YesterdayCountRes) string {
	sqlStart := `UPDATE meta_event SET yesterday_count = CASE event_name `
	sqlCaseWhen := ""
	eventNameArr := []string{}
	for _, data := range datas {
		eventNameArr = append(eventNameArr, fmt.Sprintf("'%v'", data.XwlPartEvent))
		sqlCaseWhen = sqlCaseWhen + fmt.Sprintf("  WHEN '%v' THEN %v  ", data.XwlPartEvent, data.Count)
	}
	sqlEnd := " END WHERE event_name IN (" + strings.Join(eventNameArr, ",") + ") and appid= ? "
	return sqlStart + sqlCaseWhen + sqlEnd
}

func (this *MetaDataService) GetEventNameShowMap() (mapStr string, err error) {
	var eventNameList []response.MetaEventListRes

	if err := db.Sqlx.Select(&eventNameList, "select event_name,show_name from meta_event where appid = ?", this.Appid); err != nil {
		return "", err
	}
	buff := bytes.Buffer{}
	buff.WriteString("map(")
	for index, eventNameObj := range eventNameList {
		buff.WriteString("'")
		buff.WriteString(eventNameObj.EventName)
		buff.WriteString("'")
		buff.WriteString(",")
		buff.WriteString("'")
		if strings.TrimSpace(eventNameObj.ShowName) == "" {
			buff.WriteString(eventNameObj.EventName)
		} else {
			buff.WriteString(eventNameObj.ShowName)
		}
		buff.WriteString("'")

		if index != len(eventNameList)-1 {
			buff.WriteString(",")
		}
	}
	buff.WriteString(") ")
	return buff.String(), err
}

func (this *MetaDataService) MetaEventList() (res []response.MetaEventListRes, err error) {
	if err := db.Sqlx.Select(&res, "select event_name,show_name,yesterday_count from meta_event where appid = ?", this.Appid); err != nil {
		return res, err
	}
	return res, err
}

func (this *MetaDataService) MetaEventListByAttr(attr string) (res []response.MetaEventListRes, err error) {

	eventNameList := []string{}

	if err := db.Sqlx.Select(&eventNameList, "select event_name from meta_attr_relation where app_id = ? and event_attr = ? ", this.Appid, attr); err != nil {
		return res, err
	}

	var resTmp []response.MetaEventListRes
	if err := db.Sqlx.Select(&resTmp, "select event_name,show_name,yesterday_count from meta_event where appid = ?", this.Appid); err != nil {
		return res, err
	}

	for _, v := range resTmp {
		if util.InstrArr(eventNameList, v.EventName) {
			res = append(res, v)
		}
	}
	return res, err
}

func (this *MetaDataService) AttrManager(typ int) (res []response.AttributeRes, err error) {
	if err := db.Sqlx.Select(&res, "select attribute_name,show_name,data_type,attribute_type,status from attribute where app_id = ? and attribute_source =?", this.Appid, typ); err != nil {
		return res, err
	}
	for k, v := range res {
		if _, ok := parser.TypeRemarkMap[v.DataType]; ok {
			res[k].DataTypeFormat = parser.TypeRemarkMap[v.DataType]
		}
		if _, ok := parser.SysColumn[v.AttributeName]; ok {
			res[k].ShowName = parser.SysColumn[v.AttributeName]
			res[k].AttributeType = 1
		}
	}
	return res, err
}

func (this *MetaDataService) UpdateAttrInvisible(reqData request.UpdateAttrInvisibleReq) (err error) {
	appid := reqData.Appid
	attributeSource := reqData.AttributeSource
	attributeName := reqData.AttributeName
	status := reqData.Status
	if _, err := db.Sqlx.Exec("update attribute set status = ? where   app_id = ? and attribute_source = ? and attribute_name = ?;", status, appid, attributeSource, attributeName); err != nil {
		return err
	}

	return nil
}

func (this *MetaDataService) AttrManagerByMeta(reqData request.AttrManagerByMetaReq) (res []response.AttributeRes, err error) {
	appid := reqData.Appid
	typ := reqData.Typ

	eventName := reqData.EventName
	eventAttrList := []string{}

	if err := db.Sqlx.Select(&eventAttrList, "select event_attr from meta_attr_relation where  app_id = ? and event_name = ? ", appid, eventName); err != nil {
		return nil, err
	}

	var resTmp []response.AttributeRes

	if err := db.Sqlx.Select(&resTmp, "select attribute_name,show_name,data_type,attribute_type from attribute where app_id = ? and attribute_source =?", appid, typ); err != nil {
		return nil, err
	}

	for k, v := range resTmp {
		if _, ok := parser.TypeRemarkMap[v.DataType]; ok {
			resTmp[k].DataTypeFormat = parser.TypeRemarkMap[v.DataType]
		}
		if _, ok := parser.SysColumn[v.AttributeName]; ok {
			resTmp[k].ShowName = parser.SysColumn[v.AttributeName]
		}
		if util.InstrArr(eventAttrList, v.AttributeName) {

			res = append(res, resTmp[k])
		}
	}
	return res, err
}

func (this *MetaDataService) UpdateAttrShowName(reqData request.UpdateAttrShowNameReq) (err error) {
	appid := reqData.Appid
	attributeName := reqData.AttributeName
	attributeSource := reqData.Typ
	showName := reqData.ShowName
	if _, err := db.Sqlx.Exec("update attribute set show_name = ? where   app_id = ? and attribute_source = ? and attribute_name = ?;", showName, appid, attributeSource, attributeName); err != nil {
		return err
	}
	return nil
}

func (this *MetaDataService) UpdateEventShowName(reqData request.UpdateShowNameReq) (err error) {
	appid := reqData.Appid
	eventName := reqData.EventName
	showName := reqData.ShowName
	if _, err := db.Sqlx.Exec("update meta_event set show_name = ? where event_name = ? and appid = ? ;", showName, eventName, appid); err != nil {
		return err
	}
	return nil
}

func (this *MetaDataService) GetCalcuSymbolData(reqData request.GetCalcuSymbolDataReq) (res []response.AttributeRes, err error) {
	appid := reqData.Appid
	eventName := reqData.EventName

	eventAttrList := []string{}

	if err := db.Sqlx.Select(&eventAttrList, "select event_attr from meta_attr_relation where app_id = ? and event_name = ?  ", appid, eventName); err != nil {
		return nil, err
	}

	var resTmp []response.AttributeRes

	if err := db.Sqlx.Select(&resTmp, "select attribute_name,show_name,data_type,attribute_type from attribute where app_id = ? and (status = 1 or attribute_type = 1)  and attribute_source =2", appid); err != nil {
		return nil, err
	}

	for k, v := range resTmp {
		if _, ok := parser.TypeRemarkMap[v.DataType]; ok {
			resTmp[k].DataTypeFormat = parser.TypeRemarkMap[v.DataType]
		}
		if _, ok := parser.SysColumn[v.AttributeName]; ok {
			resTmp[k].ShowName = parser.SysColumn[v.AttributeName]
		}
		if util.InstrArr(eventAttrList, v.AttributeName) {

			res = append(res, resTmp[k])
		}
	}

	return res, nil

}

type EventNameAndTheAttr struct {
	EventNameDesc string `json:"event_name_desc" db:"event_name_desc"`
	EventName     string `json:"event_name" db:"event_name"`
	AttributeName string `json:"attribute_name" db:"attribute_name"`
	AttributeDesc string `json:"attribute_desc" db:"attribute_desc"`
	DataType      string `json:"data_type" db:"data_type"`
	AttributeType string `json:"attribute_type" db:"attribute_type"`
}

func (this *MetaDataService) GetAnalyseSelectOptions(appid int) (eventNameAndTheAttrList []EventNameAndTheAttr, err error) {

	getEventNameAndTheAttrSql := `
SELECT e.show_name as event_name_desc,
	e.event_name,
	a.attribute_name,
	a.show_name as attribute_desc,
	a.data_type,
	a.attribute_type
	
FROM


	(select * from meta_attr_relation where app_id = ?) mu
	INNER JOIN (select * from meta_event where appid = ?  )  e ON e.event_name = mu.event_name
	INNER JOIN (select * from attribute where app_id = ? and (status = 1 or attribute_type = 1) and attribute_source = 2 and attribute_name not in ('xwl_part_event','xwl_part_date','xwl_kafka_offset') ) a ON a.attribute_name = mu.event_attr `

	err = db.Sqlx.Select(&eventNameAndTheAttrList, getEventNameAndTheAttrSql, appid, appid, appid)
	if err != nil {
		return
	}
	return
}

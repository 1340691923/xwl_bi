package analysis

import (
	"fmt"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/platform-basic-libs/my_error"
	"github.com/1340691923/xwl_bi/platform-basic-libs/request"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/meta_data"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"time"
)

type UserEventDetail struct {
	req             request.UserEventDetailReq
	eventNameMapStr string
}

func (this *UserEventDetail) GetList() (interface{}, error) {
	SQL, args, err := this.GetExecSql()
	if err != nil {
		return nil, err
	}
	logs.Logger.Sugar().Infof("sql", SQL, args)
	rows, err := db.ClickHouseSqlx.Query(SQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	columnLength := len(columns)
	cache := make([]interface{}, columnLength)
	for index, _ := range cache {
		var a interface{}
		cache[index] = &a
	}
	var list []map[string]interface{}

	for rows.Next() {
		err := rows.Scan(cache...)
		if err != nil {
			return nil, err
		}
		item := make(map[string]interface{})
		for i, data := range cache {
			item[columns[i]] = *data.(*interface{})
		}
		list = append(list, item)
	}

	for index := range list {
		obj := list[index]

		for k, v := range obj {

			switch v.(type) {
			case time.Time:
				obj[k] = v.(time.Time).Format(util.TimeFormat)
			}
		}
		list[index] = obj
	}

	res := map[string][]map[string]interface{}{}

	tmp := []string{}

	for index := range list {
		if _, ok := res[list[index]["date_year"].(string)]; !ok {
			m := []map[string]interface{}{}
			m = append(m, list[index])
			res[list[index]["date_year"].(string)] = m
			tmp = append(tmp, list[index]["date_year"].(string))
		} else {
			m := res[list[index]["date_year"].(string)]
			m = append(m, list[index])
			res[list[index]["date_year"].(string)] = m
		}
	}

	resList := []map[string][]map[string]interface{}{}

	for index := range tmp {
		val := res[tmp[index]]
		m := map[string][]map[string]interface{}{tmp[index]: val}
		resList = append(resList, m)
	}

	return map[string]interface{}{"list": resList}, nil
}

func (this *UserEventDetail) GetExecSql() (SQL string, allArgs []interface{}, err error) {

	eventWhereSql := ""

	if this.req.EventName != "" {
		eventWhereSql = fmt.Sprintf(` and xwl_part_event = '%s' `, this.req.EventName)
	}

	SQL = `
		with ` + this.eventNameMapStr + `  as eventMap    
		select formatDateTime(xwl_part_date,'%Y年%m月%d日') as date_year,
               formatDateTime(xwl_part_date,'%H点%M分%S秒') as date_t,
               mapValues(eventMap)[indexOf(mapKeys(eventMap), xwl_part_event)] as xwl_part_event_desc,*  from xwl_event` + strconv.Itoa(this.req.Appid) + `
			prewhere xwl_distinct_id = ? and xwl_part_date >= toDateTime(?) and xwl_part_date <= toDateTime(?) and xwl_part_event in (?) ` + eventWhereSql + `     
			order by xwl_part_date ` + this.req.OrderBy + ` limit ?,?    `
	allArgs = append(allArgs, this.req.UserID, this.req.Date[0]+" 00:00:00", this.req.Date[1]+" 23:59:59", this.req.EventNames, db.CreatePage(uint64(this.req.Page), uint64(this.req.PageSize)), this.req.PageSize)
	return
}

func NewUserEventDetailList(reqData []byte) (Ianalysis, error) {
	obj := &UserEventDetail{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(reqData, &obj.req)
	if err != nil {
		return nil, err
	}
	metaDataService := meta_data.MetaDataService{Appid: strconv.Itoa(obj.req.Appid)}
	mapStr, err := metaDataService.GetEventNameShowMap()
	if err != nil {
		return nil, err
	}
	obj.eventNameMapStr = mapStr

	if obj.req.Page == 0 {
		obj.req.Page = 1
	}

	if obj.req.PageSize == 0 {
		obj.req.PageSize = 30
	}
	if len(obj.req.EventNames) == 0 {
		return nil, my_error.NewBusiness(ERROR_TABLE, EventNameEmptyError)
	}

	return obj, nil
}

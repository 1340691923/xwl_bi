package analysis

import (
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/platform-basic-libs/my_error"
	"github.com/1340691923/xwl_bi/platform-basic-libs/request"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/meta_data"
	jsoniter "github.com/json-iterator/go"
	"strconv"
)

type UserEventCount struct {
	req             request.UserEventCountReq
	eventNameMapStr string
}

func (this *UserEventCount) GetList() (interface{}, error) {

	type EventPie struct {
		XwlPartEvent string      `json:"name" db:"xwl_part_event"`
		EventScale   interface{} `json:"value" db:"event_scale"`
	}

	var EventPieList []EventPie

	sql1 := ` 
 			with ` + this.eventNameMapStr + `  as eventMap 
			select mapValues(eventMap)[indexOf(mapKeys(eventMap), xwl_part_event)] as xwl_part_event,round(count(*),2) as event_scale from xwl_event` + strconv.Itoa(this.req.Appid) + ` 
			prewhere xwl_distinct_id = ? and xwl_part_date >= toDateTime(?) and xwl_part_date <= toDateTime(?) and xwl_part_event in (?)  group by xwl_part_event `

	err := db.ClickHouseSqlx.Select(&EventPieList, sql1, this.req.UserID, this.req.Date[0]+" 00:00:00", this.req.Date[1]+" 23:59:59", this.req.EventNames)

	logs.Logger.Sugar().Infof("sql", sql1, this.req.UserID, this.req.Date[0], this.req.Date[1], this.req.EventNames)

	if err != nil {
		return nil, err
	}

	type EventLine struct {
		DateGroup string `json:"date_group" db:"date_group"`
		Count     int    `json:"count" db:"count"`
	}

	var EventLineList []EventLine

	groupSql, groupCol := this.GetGroupDateSql()

	sql2 := `select ` + groupCol + ` ,count(*) as count from xwl_event` + strconv.Itoa(this.req.Appid) +
		` prewhere xwl_distinct_id = ? and xwl_part_date >= toDateTime(?) and xwl_part_date <= toDateTime(?) and xwl_part_event in (?) ` + ` group by ` + groupSql + ` order by ` + groupSql

	err = db.ClickHouseSqlx.Select(&EventLineList, sql2, this.req.UserID, this.req.Date[0]+" 00:00:00", this.req.Date[1]+" 23:59:59", this.req.EventNames)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"EventPieList": EventPieList, "EventLineList": EventLineList}, nil
}

func (this *UserEventCount) GetExecSql() (SQL string, allArgs []interface{}, err error) {
	return
}

func (this *UserEventCount) GetGroupDateSql() (groupSQL string, groupCol string) {

	switch this.req.WindowTimeFormat {
	case ByDay:
		return "  date_group ", " formatDateTime(xwl_part_date,'%Y年%m月%d日') as date_group "
	case ByHour:
		return "  date_group ", " formatDateTime(xwl_part_date,'%Y年%m月%d日 %H点') as date_group "
	case ByMinute:
		return "  date_group ", " formatDateTime(xwl_part_date,'%Y年%m月%d日 %H点%M分') as date_group "
	case ByWeek:
		return "  date_group ", " formatDateTime(xwl_part_date,'%Y年%m月 星期%u')  as date_group "
	case Monthly:
		return "  date_group ", " formatDateTime(xwl_part_date,'%Y年%m月') as date_group"
	}

	return
}

func NewUserEventCountList(reqData []byte) (Ianalysis, error) {
	obj := &UserEventCount{}
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
	if len(obj.req.Date) < 2 {
		return nil, my_error.NewBusiness(ERROR_TABLE, TimeError)
	}
	if len(obj.req.EventNames) == 0 {
		return nil, my_error.NewBusiness(ERROR_TABLE, EventNameEmptyError)
	}

	return obj, nil
}

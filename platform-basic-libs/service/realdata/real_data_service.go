package realdata

import (
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/platform-basic-libs/response"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/consumer_data"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"strconv"
)

type RealDataService struct {
}

//错误数据列表
func (this RealDataService) FailDataList(minutes int, appid int) (failDataResList []response.FailDataRes, err error) {

	err = db.ClickHouseSqlx.Select(&failDataResList, `
			select  
			toStartOfInterval(a.part_date, INTERVAL `+strconv.Itoa(minutes)+`  minute) as interval_date,
			formatDateTime(interval_date,'%Y-%m-%d') as year ,formatDateTime(interval_date,'%H:%M') as start_minute, formatDateTime(addMinutes(interval_date, ?),'%H:%M') as end_minute,
			count(report_data) as count,a.error_reason,a.error_handling,report_type 
			from (select * from xwl_acceptance_status prewhere table_id = ? and status = ? order by part_date desc limit 1000 ) a
			group by interval_date,a.error_reason,a.error_handling,report_type
			order by interval_date desc;
	`, minutes, appid, consumer_data.FailStatus)

	return
}

func (this RealDataService) FailDataDesc(appid, startTime, endTime, errorReason, errorHandling, reportType string) (data string, err error) {
	err = db.ClickHouseSqlx.Get(&data, `
			select report_data from xwl_acceptance_status prewhere
			table_id = `+appid+`
			and part_date >= '`+startTime+`'
			and part_date <= '`+endTime+`'
			and error_reason = '`+errorReason+`'
			and error_handling = '`+errorHandling+`'
			and status = `+strconv.Itoa(consumer_data.FailStatus)+`
			and report_type = '`+reportType+`' LIMIT  1
	`)
	return
}

func (this RealDataService) ReportCount(appid string, startTime string, endTime string) (res []response.ReportCountRes, err error) {
	type count struct {
		DataName string `db:"data_name"`
		Count    int    `db:"count"`
	}

	type ShowNameTmp struct {
		EventName string `db:"event_name"`
		ShowName  string `db:"show_name"`
	}

	var allCountArr []count
	var failCountArr []count
	var succCountArr []count
	var showNameTmpArr []ShowNameTmp
	mysqlErr := db.Sqlx.
		Select(&showNameTmpArr, "select event_name,show_name from meta_event where appid = ?", appid)
	if util.FilterMysqlNilErr(mysqlErr) {
		return nil, mysqlErr
	}
	err = db.ClickHouseSqlx.Select(&allCountArr, `select data_name,count() as count from xwl_acceptance_status xas prewhere table_id = `+appid+` and  part_date >= '`+startTime+`'  and part_date <= '`+endTime+`' group by data_name`)
	if err != nil {
		return nil, err
	}

	err = db.ClickHouseSqlx.Select(&failCountArr, `select data_name,count() as count from xwl_acceptance_status xas prewhere   status = 0 and	table_id = `+appid+` and  part_date >= '`+startTime+`'  and part_date <= '`+endTime+`' group by data_name`)
	if err != nil {
		return nil, err
	}
	err = db.ClickHouseSqlx.Select(&succCountArr, `select data_name,count() as count from xwl_acceptance_status xas prewhere  status = 1 and table_id = `+appid+` and  part_date >= '`+startTime+`'  and part_date <= '`+endTime+`' group by data_name`)
	if err != nil {
		return nil, err
	}

	resMap := map[string]response.ReportCountRes{}

	for _, data := range allCountArr {
		var tmp response.ReportCountRes
		tmp.ReceivedCount = data.Count
		tmp.DataName = data.DataName
		resMap[data.DataName] = tmp
	}

	for _, data := range failCountArr {
		if _, found := resMap[data.DataName]; found {
			tmp := resMap[data.DataName]
			tmp.FailCount = data.Count
			resMap[data.DataName] = tmp
		}
	}

	for _, data := range succCountArr {
		if _, found := resMap[data.DataName]; found {
			tmp := resMap[data.DataName]
			tmp.SuccCount = data.Count
			resMap[data.DataName] = tmp
		}
	}

	for _, data := range showNameTmpArr {
		if _, found := resMap[data.EventName]; found {
			tmp := resMap[data.EventName]
			tmp.ShowName = data.ShowName
			resMap[data.EventName] = tmp
		}
	}

	for _, data := range resMap {
		res = append(res, data)
	}
	return res, nil
}

func (this RealDataService) EventFailDesc(appid, startTime, endTime, dataName string) (res []response.EventFailDescRes, err error) {
	err = db.ClickHouseSqlx.Select(&res, `select error_reason,count() as count,any(report_data) as report_data from xwl_acceptance_status prewhere 
			table_id = `+appid+`
			and part_date >= '`+startTime+`'
			and part_date <= '`+endTime+`'
			and data_name = '`+dataName+`'
			and status = `+strconv.Itoa(consumer_data.FailStatus)+`
			group by  error_reason`)
	if err != nil {
		return nil, err
	}
	return res, nil
}

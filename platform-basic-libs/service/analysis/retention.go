package analysis

import (
	"fmt"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/platform-basic-libs/my_error"
	"github.com/1340691923/xwl_bi/platform-basic-libs/request"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/analysis/utils"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type Retention struct {
	sql  string
	args []interface{}
	req  request.RetentionReqData
}

func (this *Retention) GetList() (interface{}, error) {

	sqls, args, err := this.GetExecSql()

	if err != nil {
		return nil, err
	}

	logs.Logger.Sugar().Infof("sql", sqls, args, err)

	type Res struct {
		Dates string     `json:"dates" db:"dates"`
		Value []uint64   `json:"value" db:"value"`
		UI    [][]string `json:"ui" db:"ui"`
	}

	var res []Res

	err = db.ClickHouseSqlx.Select(&res, sqls, args...)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"alldata": res}, nil
}

func (this *Retention) getSqlByDate(t time.Time) (SQL string, allArgs []interface{}, err error) {

	var tmp = func(index int) (firstDayEventNameSql string, args []interface{}, err error) {

		firstDayEventNameSql = `xwl_part_event ='` + this.req.ZhibiaoArr[index].EventName + `' and  toYYYYMMDD(xwl_part_date) = '` + t.Format(util.TimeFormatDay) + `' `
		var sql = ""
		if len(this.req.ZhibiaoArr[index].Relation.Filts) > 0 {
			firstDayEventNameSql = firstDayEventNameSql + " and "

			sql, args, _, err = utils.GetWhereSql(this.req.ZhibiaoArr[index].Relation)

			if err == nil {
				firstDayEventNameSql = firstDayEventNameSql + sql
			}
		}
		return
	}

	firstDayEventNameSql, args, err := tmp(0)

	if err != nil {
		return
	}

	allArgs = append(allArgs, args...)

	firstDayEventName2Sql, args, err := tmp(1)

	if err != nil {
		return
	}

	allArgs = append(allArgs, args...)

	retentionSql := firstDayEventNameSql + `,` + firstDayEventName2Sql

	sumArr := make([]string, this.req.WindowTime)
	uiArr := make([]string, this.req.WindowTime)
	var sql string

	retentionPartDate := t

	for i := 0; i < this.req.WindowTime; i++ {

		retentionPartDate = retentionPartDate.AddDate(0, 0, 1)

		retentionSql = retentionSql + ","

		sumArr[i] = fmt.Sprintf("sum(r[%s])", strconv.Itoa(i+3))
		uiArr[i] = fmt.Sprintf("groupUniqArray(if(r[%s]=1,xwl_distinct_id,null))", strconv.Itoa(i+3))

		retentionSql = retentionSql + ` xwl_part_event ='` + this.req.ZhibiaoArr[1].EventName + `' and  toYYYYMMDD(xwl_part_date) = '` + retentionPartDate.Format(util.TimeFormatDay) + `' `

		if len(this.req.ZhibiaoArr[1].Relation.Filts) > 0 {
			retentionSql = retentionSql + " and "
			sql, args, _, err = utils.GetWhereSql(this.req.ZhibiaoArr[1].Relation)

			if err != nil {
				logs.Logger.Error("err", zap.Error(err))
				return
			}

			allArgs = append(allArgs, args...)
			retentionSql = retentionSql + sql
		}
	}

	var userFilterSql string
	var userFilterArgs []interface{}

	if len(this.req.WhereFilterByUser.Filts) > 0 {
		var colArr []string
		var sql string
		sql, userFilterArgs, colArr, err = utils.GetWhereSql(this.req.WhereFilterByUser)
		if err != nil {
			return SQL, allArgs, err
		}
		userFilterSql = `and xwl_distinct_id in ( select xwl_distinct_id from ` + utils.GetUserTableView(this.req.Appid, colArr) + ` where ` + sql + ")"
	}

	whereFilterSql, whereFilterArgs, _, err := utils.GetWhereSql(this.req.WhereFilter)

	whereFilterSql = whereFilterSql + this.sql

	whereFilterArgs = append(whereFilterArgs, this.args...)

	parteventWhereSql := " xwl_part_event in (?,?) "

	allArgs = append(allArgs, this.req.ZhibiaoArr[0].EventName, this.req.ZhibiaoArr[1].EventName)

	allArgs = append(allArgs, whereFilterArgs...)

	allArgs = append(allArgs, userFilterArgs...)

	SQL = `
			SELECT
				'` + t.Format(util.TimeFormatDay2) + `' AS dates,
				array(sum(r[1]),sum(r[2]),` + strings.Join(sumArr, ",") + `) as value,
				array(groupUniqArray(if(r[1]=1,xwl_distinct_id,null)),groupUniqArray(if(r[2]=1,xwl_distinct_id,null)), ` + strings.Join(uiArr, ",") + `) as ui
				FROM
				(
					SELECT
   					 xwl_distinct_id,
   				 retention(` + retentionSql + `) AS r
				FROM xwl_event` + strconv.Itoa(this.req.Appid) + `  
				prewhere xwl_part_date >= toDateTime('` + t.Format(util.TimeFormat) + `') and xwl_part_date <= toDateTime('` + t.AddDate(0, 0, this.req.WindowTime+1).Format(util.TimeFormat) + `') and ` + parteventWhereSql + `  and ` + whereFilterSql + ` ` + userFilterSql + `
				
				GROUP BY xwl_distinct_id
			) limit 1000
       `
	return
}

func (this *Retention) parseReqDate() []time.Time {

	if len(this.req.Date) < 2 {
		return nil
	}

	startTimeFormat := this.req.Date[0]
	endTimeFormat := this.req.Date[1]

	if startTimeFormat == endTimeFormat {
		t := make([]time.Time, 1)
		t[0] = util.Str2Time(startTimeFormat, util.TimeFormatDay2)
		return t
	}
	startT := util.Str2Time(startTimeFormat, util.TimeFormatDay2)
	endT := util.Str2Time(endTimeFormat, util.TimeFormatDay2)

	t := []time.Time{}
	for ; startT.Before(endT.AddDate(0, 0, 1)); startT = startT.AddDate(0, 0, 1) {
		t = append(t, startT)
	}

	return t
}

func (this *Retention) GetExecSql() (SQL string, allArgs []interface{}, err error) {
	sqlArr := []string{}
	for _, date := range this.parseReqDate() {
		sql, args, err := this.getSqlByDate(date)
		if err != nil {
			return "", nil, err
		}
		sqlArr = append(sqlArr, sql)
		allArgs = append(allArgs, args...)
	}

	return strings.Join(sqlArr, "union all"), allArgs, err
}

func NewRetention(reqData []byte) (Ianalysis, error) {
	obj := &Retention{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(reqData, &obj.req)
	if err != nil {
		return nil, err
	}
	if len(obj.req.Date) < 2 {
		return nil, my_error.NewBusiness(ERROR_TABLE, TimeError)
	}
	if len(obj.req.ZhibiaoArr) != 2 {
		return nil, my_error.NewBusiness(ERROR_TABLE, ZhiBiaoNumError)
	}

	obj.sql, obj.args, err = utils.GetUserGroupSqlAndArgs(obj.req.UserGroup, obj.req.Appid)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

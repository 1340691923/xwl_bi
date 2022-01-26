package analysis

import (
	"database/sql"
	"github.com/1340691923/xwl_bi/engine/logs"

	"fmt"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/platform-basic-libs/my_error"
	"github.com/1340691923/xwl_bi/platform-basic-libs/request"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/analysis/utils"
	jsoniter "github.com/json-iterator/go"
	"strconv"
)

type Funnel struct {
	sql  string
	args []interface{}
	req  request.FunnelReqData
}

func (this *Funnel) GetExecSql() (SQL string, allArgs []interface{}, err error) {

	startTime := this.req.Date[0] + " 00:00:00"
	endTime := this.req.Date[1] + " 23:59:59"

	windowSql := ""

	for _, zhibiao := range this.req.ZhibiaoArr {
		windowSql = windowSql + ","

		windowSql = windowSql + fmt.Sprintf(" xwl_part_event = '%v' ", zhibiao.EventName)

		if len(zhibiao.Relation.Filts) > 0 {
			windowSql = windowSql + " and "
			sql, args, _, err := utils.GetWhereSql(zhibiao.Relation)

			if err != nil {
				return SQL, allArgs, err
			}

			allArgs = append(allArgs, args...)

			windowSql = windowSql + sql
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

	allArgs = append(allArgs, whereFilterArgs...)

	allArgs = append(allArgs, this.args...)

	allArgs = append(allArgs, userFilterArgs...)

	if err != nil {
		logs.Logger.Sugar().Errorf("req.WhereFilter", this.req.WhereFilter)
		return
	}

	SQL = `SELECT '总体' as groupkey,level_index,count(1) as count,groupUniqArray(xwl_distinct_id) as ui  FROM
			(
				SELECT  xwl_distinct_id,
					arrayWithConstant(windowFunnel_level, 1) levels, 
					arrayJoin(arrayEnumerate( levels )) level_index
				  FROM (
					SELECT
					  xwl_distinct_id,
					  windowFunnel(` + strconv.Itoa(this.req.WindowTime) + `)(
						xwl_part_date
						` + windowSql + `
					  ) AS windowFunnel_level
					FROM  xwl_event` + strconv.Itoa(this.req.Appid) + `   
					WHERE xwl_part_date >= toDateTime('` + startTime + `') and xwl_part_date <= toDateTime('` + endTime + `') and ` + whereFilterSql + ` ` + userFilterSql + `
					GROUP BY xwl_distinct_id
				)
			)
			group by level_index
			ORDER BY level_index limit 1000
	`

	if len(this.req.GroupBy) > 0 {
		groupSql := `SELECT toString(groupkey) as groupkey,level_index,count(1) as count ,groupUniqArray(xwl_distinct_id) as ui FROM
			(
				SELECT  xwl_distinct_id, groupkey,
					arrayWithConstant(windowFunnel_level, 1) levels, 
					arrayJoin(arrayEnumerate( levels )) level_index
				  FROM (
					SELECT
					  xwl_distinct_id, ` + this.req.GroupBy[0] + ` groupkey,
					  windowFunnel(` + strconv.Itoa(this.req.WindowTime) + `)(
						xwl_part_date
						` + windowSql + ` 
					  ) AS windowFunnel_level
					  FROM xwl_event` + strconv.Itoa(this.req.Appid) + ` 
					  WHERE xwl_part_date >= toDateTime('` + startTime + `') and xwl_part_date <= toDateTime('` + endTime + `') and ` + whereFilterSql + `  ` + userFilterSql + `
					
					GROUP BY xwl_distinct_id,groupkey
				)
			)
			group by groupkey,level_index
			ORDER BY groupkey,level_index limit 1000
	`
		SQL = fmt.Sprintf("%s UNION ALL %s", SQL, groupSql)
		allArgs = append(allArgs, allArgs...)
	}

	return SQL, allArgs, err
}

type FunnelRes struct {
	LevelIndex int      `json:"level_index" db:"level_index"`
	Count      int      `json:"count" db:"count"`
	UI         []string `json:"ui" db:"ui"`
}

type FunnelGroupRes struct {
	Groupkey   sql.NullString `json:"groupkey" db:"groupkey"`
	LevelIndex int            `json:"level_index" db:"level_index"`
	UI         []string       `json:"ui" db:"ui"`
	Count      int            `json:"count" db:"count"`
}

func (this *Funnel) GetList() (interface{}, error) {

	sql, args, err := this.GetExecSql()

	if err != nil {
		return nil, err
	}

	logs.Logger.Sugar().Infof("SQL", sql, args)

	var funnelGroupResList []FunnelGroupRes

	if err := db.ClickHouseSqlx.Select(&funnelGroupResList, sql, args...); err != nil {
		return nil, err
	}

	groupData := map[string][]FunnelRes{}

	for _, v := range funnelGroupResList {
		if _, ok := groupData[v.Groupkey.String]; !ok {
			groupData[v.Groupkey.String] = []FunnelRes{}
		}
		groupData[v.Groupkey.String] = append(groupData[v.Groupkey.String], FunnelRes{
			LevelIndex: v.LevelIndex,
			Count:      v.Count,
			UI:         v.UI,
		})
	}

	return map[string]interface{}{"groupData": groupData}, nil
}

func NewFunnel(reqData []byte) (Ianalysis, error) {
	obj := &Funnel{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(reqData, &obj.req)
	if err != nil {
		return nil, err
	}

	if len(obj.req.Date) < 2 {
		return nil, my_error.NewBusiness(ERROR_TABLE, TimeError)
	}
	if len(obj.req.ZhibiaoArr) > 30 {
		return nil, my_error.NewBusiness(ERROR_TABLE, ZhiBiaoNumError)
	}
	var T int
	switch obj.req.WindowTimeFormat {
	case "天":
		T = 60 * 60 * 24
	case "小时":
		T = 60 * 60
	case "分钟":
		T = 60
	case "秒":
		T = 1
	default:
		return nil, my_error.NewBusiness(ERROR_TABLE, TimeError)
	}
	obj.req.WindowTime = obj.req.WindowTime * T
	obj.sql, obj.args, err = utils.GetUserGroupSqlAndArgs(obj.req.UserGroup, obj.req.Appid)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

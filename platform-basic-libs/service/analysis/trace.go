package analysis

import (
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/platform-basic-libs/my_error"
	"github.com/1340691923/xwl_bi/platform-basic-libs/request"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/analysis/utils"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/meta_data"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"strings"
)

type Trace struct {
	sql             string
	args            []interface{}
	eventNameMapStr string
	req             request.TraceReqData
	sqlTyp          int
}

const ChartSql int = 1

const TableSql int = 2

func (this *Trace) setSqlTyp(typ int) {
	this.sqlTyp = typ
}

func (this *Trace) GetList() (interface{}, error) {

	type Res1 struct {
		Source string `json:"source"`
		Target string `json:"target"`
		Value  int    `json:"value"`
	}

	type Res2 struct {
		Name string `json:"name"`
	}

	this.setSqlTyp(TableSql)

	SQL, args, err := this.GetExecSql()
	if err != nil {
		return nil, err
	}

	logs.Logger.Sugar().Infof("sql", SQL, args, err)

	type TableRes struct {
		Trace     string   `json:"trace" db:"trace"`
		UserCount uint64   `json:"user_count" db:"user_count"`
		Ui        []string `json:"ui" db:"ui"`
	}

	var tableRes []TableRes

	err = db.ClickHouseSqlx.Select(&tableRes, SQL, args...)
	if err != nil {
		return nil, err
	}

	this.setSqlTyp(ChartSql)

	SQL, args, err = this.GetExecSql()

	if err != nil {
		return nil, err
	}

	logs.Logger.Sugar().Infof("sql", SQL, args, err)

	type ChartsRes struct {
		Event        []string `json:"event" db:"trace2"`
		SumUserCount uint64   `json:"sum_user_count" db:"sum_user_count"`
	}

	var chartsRes []ChartsRes

	err = db.ClickHouseSqlx.Select(&chartsRes, SQL, args...)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"tableRes": tableRes, "chartRes": chartsRes}, nil
}

func (this *Trace) GetTableSql() (SQL string, allArgs []interface{}, err error) {
	startTime := this.req.Date[0] + " 00:00:00"
	endTime := this.req.Date[1] + " 23:59:59"

	windowSql, allArgs, err := getZhibiaoFilterSqlArgs(this.req.ZhibiaoArr)
	if err != nil {
		return
	}
	userFilterSql, userFilterArgs, err := getUserfilterSqlArgs(this.req.WhereFilterByUser, this.req.Appid)
	if err != nil {
		return
	}
	whereFilterSql, whereFilterArgs, _, err := utils.GetWhereSql(this.req.WhereFilter)
	if err != nil {
		return
	}

	if len(this.req.EventNames) > 1 {
		whereFilterSql = whereFilterSql + " and xwl_part_event in (?) "
		whereFilterArgs = append(whereFilterArgs, this.req.EventNames)
	}

	whereFilterSql = whereFilterSql + this.sql

	allArgs = append(allArgs, whereFilterArgs...)

	allArgs = append(allArgs, this.args...)

	allArgs = append(allArgs, userFilterArgs...)

	SQL = `
		  select trace ,user_count,ui from  (SELECT
			 result_chain as trace,
			  uniqCombined(user_id) AS user_count,groupUniqArray(user_id) ui
			  FROM (
					select
					  xwl_distinct_id as user_id,
                      ` + this.eventNameMapStr + ` as eventMap,
					  arrayStringConcat(
					   arrayCompact(    
						arrayMap(
						  b -> tupleElement(b, 1),
						  arraySort(  
							y -> tupleElement(y, 2),
							arrayFilter(
							  (x, y) -> y - x.2 <= ` + strconv.Itoa(this.req.WindowTime) + ` and  y - x.2  >=  0,
							  arrayMap(
								(x, y) -> (x, y),
								groupArray(mapValues(eventMap)[indexOf(mapKeys(eventMap), xwl_part_event)]),
								groupArray(toUnixTimestamp(xwl_part_date))
							  ),  
							  arrayWithConstant(
								length(groupArray(toUnixTimestamp(xwl_part_date))),
								maxIf(toUnixTimestamp(xwl_part_date) ` + windowSql + `)
							  )
							)
						  )
						)
					   ),
					  '->'
					 ) result_chain
					from
					  xwl_event` + strconv.Itoa(this.req.Appid) + `
					prewhere
			   xwl_part_date >= toDateTime('` + startTime + `')
				AND xwl_part_date <= toDateTime('` + endTime + `') and ` + whereFilterSql + ` ` + userFilterSql + `
					
					group by
					  xwl_distinct_id
					  HAVING notEmpty(result_chain)
			  ) tab 
			  GROUP BY result_chain
			  ORDER BY user_count DESC) limit 1000
	`

	return
}

func (this *Trace) GetChartSql() (SQL string, allArgs []interface{}, err error) {
	startTime := this.req.Date[0] + " 00:00:00"
	endTime := this.req.Date[1] + " 23:59:59"

	windowSql, allArgs, err := getZhibiaoFilterSqlArgs(this.req.ZhibiaoArr)
	if err != nil {
		return
	}
	userFilterSql, userFilterArgs, err := getUserfilterSqlArgs(this.req.WhereFilterByUser, this.req.Appid)
	if err != nil {
		return
	}
	whereFilterSql, whereFilterArgs, _, err := utils.GetWhereSql(this.req.WhereFilter)
	if err != nil {
		return
	}

	if len(this.req.EventNames) > 1 {
		whereFilterSql = whereFilterSql + " and xwl_part_event in (?) "
		whereFilterArgs = append(whereFilterArgs, this.req.EventNames)
	}

	whereFilterSql = whereFilterSql + this.sql

	allArgs = append(allArgs, whereFilterArgs...)

	allArgs = append(allArgs, this.args...)

	allArgs = append(allArgs, userFilterArgs...)

	SQL = `
		
		select   splitByString('->',arrayJoin(arrayMap(	(x, y) -> concat(concat( x,'->'),y),arraySlice(trace, 1, length(trace) - 1),arraySlice(trace, 2, length(trace)) )))    as trace2,sum(1)  as sum_user_count
		
		from (
		select if(length(trace33)==1,arrayPushBack(trace33,str_link),arrayPushBack(arraySlice(trace33, 1, length(trace33) - 1),str_link)) as trace  from 
		(
		
		with  ` + this.eventNameMapStr + ` as eventMap 
		
		select result_chain as trace33,concat(mapValues(eventMap)[indexOf(mapKeys(eventMap), '` + this.req.ZhibiaoArr[0].EventName + `')],'(最终路径)') as str_link  from 
		
		(
		select
          					
          					 arrayMap( 
                                                  b -> tuple(b.3,b.1),
                                                  arraySort(  
                                                        y -> tupleElement(y, 2),
                                                        arrayFilter(
                                                   (x, y) -> y - x.2 <=  ` + strconv.Itoa(this.req.WindowTime) + ` and  y - x.2 >= 0,
                                                          arrayMap(
                                                                (x, y,z) -> (x, y,z),
                                                                groupArray(event.1),
                                                                groupArray(toUnixTimestamp(part_date)),
                                                                groupArray(event.2)
                                                          ), 
                                                          arrayWithConstant(
                                                                length(groupArray(toUnixTimestamp(part_date))),
                                                            maxIf(toUnixTimestamp(part_date) ` + strings.ReplaceAll(windowSql, "xwl_part_event", "event.1") + `  )
                                                          )
                                                        )
                                                  )
                                                )
                                                 as  tmpArr,
                                                 
                                                 arrayMap(
                                                  b -> b.1,
                                                 arrayFilter((x, y) -> ((x.2) != y), tmpArr, arrayPushFront(arrayPopBack(tmpArr.2), ''))
                                                 )
                                                 
                                               as result_chain
                                        from
                               (
                               with  ` + this.eventNameMapStr + ` as eventMap 
                              select  
              
              arrayJoin(
               
                
              arrayMap(
                          (x, y,z,q) -> (x, y,z,q),
                             groupArray(mapValues(eventMap)[indexOf(mapKeys(eventMap), xwl_part_event)]),
                             groupArray(xwl_part_date),
                             groupArray(xwl_distinct_id),
                            arrayMap((x, y) -> CONCAT(x, concat('(',concat(toString(y),')'))),groupArray(mapValues(eventMap)[indexOf(mapKeys(eventMap), xwl_part_event)]),arrayEnumerateUniq(groupArray(mapValues(eventMap)[indexOf(mapKeys(eventMap), xwl_part_event)])))
                        )
              ) x,
               (x.1,x.4) as event,x.2 as part_date,x.3 as ui
               
               from   xwl_event` + strconv.Itoa(this.req.Appid) + `
						prewhere
			   xwl_part_date >= toDateTime('` + startTime + `')
				AND xwl_part_date <= toDateTime('` + endTime + `') and ` + whereFilterSql + ` ` + userFilterSql + `
                                
                               )
                                          
                                        group by
                                          ui
                                          HAVING notEmpty(tmpArr) 

			))) group by trace2

	`

	/*SQL = `select  splitByString('->', arrayJoin(arrayMap(	(x, y) -> concat(concat( x,'->'),y),arraySlice(trace, 1, length(trace) - 1),arraySlice(trace, 2, length(trace)) ))) as trace2,sum(1) as sum_user_count from (

			select if(length(trace33)==1,arrayPushBack(trace33,str_link),arrayPushBack(arraySlice(trace33, 1, length(trace33) - 1),str_link)) as trace from (
			select  arrayMap((x, y) -> CONCAT(x, concat('(',concat(toString(y),')'))),result_chain,arrayEnumerateUniq(result_chain)) as trace33,concat(mapValues(eventMap)[indexOf(mapKeys(eventMap), '`+this.req.ZhibiaoArr[0].EventName+`')],'(最终路径)') as str_link  from (

				select
		  xwl_distinct_id as user_id, ` + this.eventNameMapStr + ` as eventMap,

						   arrayCompact(
							arrayMap(
							  b -> tupleElement(b, 1),
							  arraySort(
								y -> tupleElement(y, 2),
								arrayFilter(
							   (x, y) -> y - x.2 <= ` + strconv.Itoa(this.req.WindowTime) + ` and  y - x.2 >= 0,
								  arrayMap(
									(x, y) -> (x, y),
						 			groupArray(mapValues(eventMap)[indexOf(mapKeys(eventMap), xwl_part_event)]),
									groupArray(toUnixTimestamp(xwl_part_date))
								  ),
								  arrayWithConstant(
									length(groupArray(toUnixTimestamp(xwl_part_date))),
								    maxIf(toUnixTimestamp(xwl_part_date)` + windowSql + ` )
								  )
								)
							  )
							)
						   )

						  result_chain
						from
				       xwl_event` + strconv.Itoa(this.req.Appid) + `
							prewhere
				   xwl_part_date >= toDateTime('` + startTime + `')
					AND xwl_part_date <= toDateTime('` + endTime + `') and ` + whereFilterSql + ` ` + userFilterSql + `
						group by
						  xwl_distinct_id
						  HAVING notEmpty(result_chain)

			)))

			group by trace2

	    `*/

	return
}

func (this *Trace) GetExecSql() (SQL string, allArgs []interface{}, err error) {
	if this.sqlTyp == ChartSql {
		return this.GetChartSql()
	}

	return this.GetTableSql()
}

func NewTrace(reqData []byte) (Ianalysis, error) {
	obj := &Trace{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(reqData, &obj.req)
	if err != nil {
		return nil, err
	}
	if len(obj.req.Date) < 2 {
		return nil, my_error.NewBusiness(ERROR_TABLE, TimeError)
	}
	if len(obj.req.ZhibiaoArr) != 1 {
		return nil, my_error.NewBusiness(ERROR_TABLE, ZhiBiaoNumError)
	}
	obj.sqlTyp = TableSql
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
	metaDataService := meta_data.MetaDataService{Appid: strconv.Itoa(obj.req.Appid)}
	mapStr, err := metaDataService.GetEventNameShowMap()
	if err != nil {
		return nil, err
	}
	obj.eventNameMapStr = mapStr

	obj.sql, obj.args, err = utils.GetUserGroupSqlAndArgs(obj.req.UserGroup, obj.req.Appid)

	if err != nil {
		return nil, err
	}

	obj.req.EventNames = append(obj.req.EventNames, obj.req.ZhibiaoArr[0].EventName)

	return obj, nil
}

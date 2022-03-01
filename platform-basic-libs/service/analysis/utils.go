package analysis

import (
	"fmt"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/platform-basic-libs/request"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/analysis/utils"
)

func getUserfilterSqlArgs(analysisFilter request.AnalysisFilter, appid int) (userFilterSql string, userFilterArgs []interface{}, err error) {
	if len(analysisFilter.Filts) > 0 {
		var colArr []string
		var sql string
		sql, userFilterArgs, colArr, err = utils.GetWhereSql(analysisFilter)
		if err != nil {
			return
		}
		userFilterSql = ` and xwl_distinct_id in ( select xwl_distinct_id from ` + utils.GetUserTableView(appid, colArr) + ` where ` + sql + ") "
	}
	return
}

func getZhibiaoFilterSqlArgs(zhibiaoArr []request.Zhibiao) (windowSql string, allArgs []interface{}, err error) {

	for _, zhibiao := range zhibiaoArr {
		windowSql = windowSql + ","

		windowSql = windowSql + fmt.Sprintf(" xwl_part_event = '%v' ", zhibiao.EventName)

		if len(zhibiao.Relation.Filts) > 0 {
			windowSql = windowSql + " and "
			sql, args, _, err := utils.GetWhereSql(zhibiao.Relation)

			if err != nil {
				logs.Logger.Sugar().Errorf("zhibiao.Relation", zhibiao)
				return windowSql, allArgs, err
			}

			allArgs = append(allArgs, args...)

			windowSql = windowSql + sql
		}
	}

	return windowSql, allArgs, err
}

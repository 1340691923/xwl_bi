package analysis

import (
	"fmt"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/platform-basic-libs/my_error"
	"github.com/1340691923/xwl_bi/platform-basic-libs/request"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/analysis/utils"
	jsoniter "github.com/json-iterator/go"
)

type UserAttr struct {
	req  request.UserAttrReqData
	sql  string
	args []interface{}
}

func (this *UserAttr) GetList() (interface{}, error) {

	SQL, args, err := this.GetExecSql()
	if err != nil {
		return nil, err
	}

	logs.Logger.Sugar().Infof("sql", SQL, args, err)

	type TableRes struct {
		GroupKey string      `json:"name" db:"groupkey"`
		Amount   interface{} `json:"value" db:"amount"`
	}

	var tableRes []TableRes

	err = db.ClickHouseSqlx.Select(&tableRes, SQL, args...)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"tableRes": tableRes}, nil
}

func (this *UserAttr) UserCountSql() string {
	zhibiaoArr := this.req.ZhibiaoArr
	return utils.CountTypMap[zhibiaoArr[1]](zhibiaoArr[0])
}

func (this *UserAttr) getGroupClo() string {
	if len(this.req.GroupBy) > 0 {
		return fmt.Sprintf("%s as groupkey,", this.req.GroupBy[0])
	}
	return `'总体' as groupkey,`
}

func (this *UserAttr) getGroupSql() string {
	if len(this.req.GroupBy) > 0 {
		return fmt.Sprintf("group by %s", this.req.GroupBy[0])
	}
	return ``
}

func (this *UserAttr) GetExecSql() (SQL string, allArgs []interface{}, err error) {

	whereSql, allArgs, colArr, err := utils.GetWhereSql(this.req.WhereFilterByUser)

	whereSql = whereSql + this.sql
	allArgs = append(allArgs, this.args...)

	if err != nil {
		return
	}
	if len(this.req.GroupBy) > 0 {
		colArr = append(colArr, this.req.GroupBy[0])
	}
	if this.req.ZhibiaoArr[0] != utils.Default {
		colArr = append(colArr, this.req.ZhibiaoArr[0])
	}

	SQL = `select ` + this.getGroupClo() + `cast(coalesce(` + this.UserCountSql() + `, 0) as double) as amount from   ` + utils.GetUserTableView(this.req.Appid, colArr) + ` where ` + whereSql + this.getGroupSql() + " limit 1000 "

	return
}

func NewUserAttr(reqData []byte) (Ianalysis, error) {
	obj := &UserAttr{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(reqData, &obj.req)
	if err != nil {
		return nil, err
	}

	if len(obj.req.ZhibiaoArr) != 2 {
		return nil, my_error.NewBusiness(ERROR_TABLE, ZhiBiaoNumError)
	}

	if len(obj.req.GroupBy) > 1 {
		return nil, my_error.NewBusiness(ERROR_TABLE, GroupNumError)
	}

	obj.sql, obj.args, err = utils.GetUserGroupSqlAndArgs(obj.req.UserGroup, obj.req.Appid)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

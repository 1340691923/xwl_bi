package analysis

import (
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/platform-basic-libs/my_error"
	"github.com/1340691923/xwl_bi/platform-basic-libs/request"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/analysis/utils"
	parser "github.com/1340691923/xwl_bi/platform-basic-libs/sinker/parse"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type UserList struct {
	req     request.UserListReqData
	propMap map[string]string
}

func (this *UserList) GetList() (interface{}, error) {
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
	delete(this.propMap, "xwl_update_time")
	return map[string]interface{}{"alldata": list, "propMap": this.propMap}, nil

}

func (this *UserList) GetExecSql() (SQL string, allArgs []interface{}, err error) {
	allArgs = append(allArgs, this.req.UI)

	type Attribute struct {
		AttributeName string `db:"attribute_name" json:"attribute_name"` //属性名
		ShowName      string `db:"show_name" json:"show_name"`
	}
	var attributes []Attribute
	err = db.Sqlx.Select(&attributes, "select attribute_name,show_name from attribute where app_id = ? and attribute_source =? and (status = 1 or attribute_type = 1)", this.req.Appid, 1)
	if err != nil {
		return
	}

	fields := make([]string, len(attributes))

	for index, attribute := range attributes {

		fields[index] = attribute.AttributeName
		if _, ok := parser.SysColumn[attribute.AttributeName]; ok {
			this.propMap[attribute.AttributeName] = parser.SysColumn[attribute.AttributeName]
		} else {
			if attribute.ShowName == "" {
				this.propMap[attribute.AttributeName] = attribute.AttributeName
			} else {
				this.propMap[attribute.AttributeName] = attribute.ShowName
			}
		}
	}

	SQL = `select * from ` + utils.GetUserTableView(this.req.Appid, fields) + ` where xwl_distinct_id in (?)`
	return
}

func NewUserList(reqData []byte) (Ianalysis, error) {
	obj := &UserList{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(reqData, &obj.req)
	if err != nil {
		return nil, err
	}
	if len(obj.req.UI) <= 0 {
		return nil, my_error.NewBusiness(ERROR_TABLE, UIEmptyError)
	}
	obj.propMap = map[string]string{}
	return obj, nil
}

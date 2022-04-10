package utils

import (
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"github.com/Masterminds/squirrel"
	"strings"
)

func GetUserGroupSqlAndArgs(ids []int, appid int) (SQL string, Args []interface{}, err error) {
	if len(ids) == 0 {
		return " and ( 1 = 1 ) ", nil, err
	}

	sql, args, err := db.
		SqlBuilder.
		Select("user_list").
		From("user_group").
		Where(db.Eq{"appid": appid, "id": ids}).
		ToSql()

	var userGroupList []model.UserGroup

	err = db.Sqlx.Select(&userGroupList, sql, args...)

	if err != nil {
		return "", nil, err
	}

	or := squirrel.Or{}

	for index := range userGroupList {
		idStr, err := util.GzipUnCompress(userGroupList[index].UserList)
		if err != nil {
			return "", nil, err
		}
		id := strings.Split(idStr, ",")
		or = append(or, db.Eq{"xwl_distinct_id": [][]string{id}})
	}

	SQL, Args, err = or.ToSql()
	SQL = " and " + SQL
	return SQL, Args, err
}

package model

import (
	"github.com/1340691923/xwl_bi/engine/db"
)

type Curd interface {
	ProcessSqlWhere(sqlA db.SelectBuilder) db.SelectBuilder
	TableName() string
	ProcessSqlInsert(sqlA db.InsertBuilder) db.InsertBuilder
	ProcessSqlUpdate(id int, sqlA db.UpdateBuilder) db.UpdateBuilder
}

func SearchList(curd Curd, page, limit int, columns string, list interface{}) (err error) {
	sqlA := db.
		SqlBuilder.
		Select(columns).
		From(curd.TableName())
	sqlA = curd.ProcessSqlWhere(sqlA)
	sql, args, err := sqlA.
		Limit(uint64(limit)).
		Offset(db.CreatePage(uint64(page), uint64(limit))).
		OrderBy("id desc").
		ToSql()
	err = db.Sqlx.Select(list, sql, args...)
	return
}

func Count(curd Curd) (count int, err error) {
	sqlA := db.SqlBuilder.
		Select("count(*)").
		From(curd.TableName())
	sqlA = curd.ProcessSqlWhere(sqlA)
	sql, args, err := sqlA.ToSql()
	err = db.Sqlx.Get(&count, sql, args...)
	return
}

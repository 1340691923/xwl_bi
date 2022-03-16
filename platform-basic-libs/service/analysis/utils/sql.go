package utils

import (
	"errors"
	"fmt"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/platform-basic-libs/request"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"github.com/Masterminds/squirrel"
	"strconv"
	"strings"
)

type sqlI interface {
	ToSql() (string, []interface{}, error)
	Append(sqlizer squirrel.Sqlizer)
}

type Or struct {
	or db.Or
}

func (this *Or) Append(sqlizer squirrel.Sqlizer) {
	this.or = append(this.or, sqlizer)
}

func (this *Or) ToSql() (string, []interface{}, error) {
	return this.or.ToSql()
}

type And struct {
	and db.And
}

func (this *And) Append(sqlizer squirrel.Sqlizer) {
	this.and = append(this.and, sqlizer)
}

func (this *And) ToSql() (string, []interface{}, error) {
	return this.and.ToSql()
}

const COMPOUND = "COMPOUND"
const SIMPLE = "SIMPLE"
const AND = "且"
const OR = "或"

var noValueSymbolArr = []string{"isNotNull", "isNull"}
var rangeSymbolArr = []string{"range"}
var rangeTimeSymbolArr = []string{"rangeTime"}

func GetWhereSql(anlysisFilter request.AnalysisFilter) (SQL string, Args []interface{}, Cols []string, err error) {
	var arrP sqlI
	colArr := []string{}
	switch anlysisFilter.Relation {
	case AND:
		arrP = &And{}
	case OR:
		arrP = &Or{}
	default:
		return "", nil, nil, errors.New("错误的连接类型:" + anlysisFilter.Relation)
	}

	for _, v := range anlysisFilter.Filts {
		if v.FilterType == SIMPLE {
			colArr = append(colArr, v.ColumnName)
			arrP.Append(getExpr(v.ColumnName, v.Comparator, v.Ftv))
		} else {
			var arrC sqlI
			switch v.Relation {
			case AND:
				arrC = &And{}
			case OR:
				arrC = &Or{}
			default:
				return "", nil, nil, errors.New("错误的连接类型")
			}

			for _, v2 := range v.Filts {
				colArr = append(colArr, v2.ColumnName)
				arrC.Append(getExpr(v2.ColumnName, v2.Comparator, v2.Ftv))
			}
			arrP.Append(arrC)
		}
	}
	sql, args, err := arrP.ToSql()

	return sql, args, colArr, err
}

func getArgMax(col string) string {
	return fmt.Sprintf(" argMax(%s, %s) %s ", col, ReplacingMergeTreeKey, col)
}

const ReplacingMergeTreeKey = "xwl_update_time"

var SpecialCloArr = []string{"xwl_distinct_id", "xwl_update_time"}

func GetUserTableView(tableId int, fields []string) string {

	colArr := []string{}

	for i, field := range fields {
		if util.InstrArr(SpecialCloArr, field) {
			fields = append(fields[:i], fields[i:]...)
			continue
		}
		colArr = append(colArr, getArgMax(field))
	}

	if len(colArr) > 0 {
		return " (select xwl_distinct_id," + strings.Join(colArr, ",") + " from xwl_user" + strconv.Itoa(tableId) + " xu group by xwl_distinct_id) "
	}

	return " (select xwl_distinct_id from xwl_user" + strconv.Itoa(tableId) + " xu group by xwl_distinct_id) "
}

func getExpr(columnName, comparator string, ftv interface{}) squirrel.Sqlizer {

	if util.InstrArr(noValueSymbolArr, comparator) {
		return squirrel.Expr(fmt.Sprintf("%v(%v)", comparator, columnName))
	}
	if util.InstrArr(rangeSymbolArr, comparator) {
		return squirrel.Expr(fmt.Sprintf(" ( %v >= ? and %v <= ? ) ", columnName, columnName), ftv.([]interface{})[0], ftv.([]interface{})[1])
	}
	if util.InstrArr(rangeTimeSymbolArr, comparator) {
		if len(ftv.([]interface{})) != 2 {
			return squirrel.Expr(" 1 = 1 ")
		}
		return squirrel.Expr(fmt.Sprintf(" ( %v >= toDateTime(?) and %v <= toDateTime(?) ) ", columnName, columnName), ftv.([]interface{})[0], ftv.([]interface{})[1])
	}
	if comparator == "match" {
		return squirrel.Expr(fmt.Sprintf("match(%v,?) = 1", columnName), ftv)
	}
	if comparator == "notmatch" {
		return squirrel.Expr(fmt.Sprintf("match(%v,?) = 0", columnName), ftv)
	}

	if comparator == "=" {
		return db.Eq{columnName: ftv}
	}
	if comparator == "!=" {
		return db.NotEq{columnName: ftv}
	}
	return squirrel.Expr(fmt.Sprintf("%v %v ?", columnName, comparator), ftv)
}

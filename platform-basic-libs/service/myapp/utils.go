package myapp

import (
	"fmt"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/jwt"
)

func GetAppidsByToken(token string) (list []model.App, err error) {
	c, _ := jwt.ParseToken(token)

	selectBuilder := db.SqlBuilder.
		Select("id,app_name").
		From("app")

	if c.UserID != 1 {
		selectBuilder = selectBuilder.Where(fmt.Sprintf("FIND_IN_SET(%v,app_manager)", c.UserID)).Where(db.Eq{"is_close": 0})
	}

	sql, args, err := selectBuilder.ToSql()

	if err != nil {
		return
	}

	err = db.Sqlx.Select(&list, sql, args...)
	if err != nil {
		return
	}

	return
}

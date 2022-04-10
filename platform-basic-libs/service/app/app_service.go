package app

import (
	"errors"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/analysis/utils"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/myapp"
	"github.com/1340691923/xwl_bi/platform-basic-libs/sinker"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"strconv"
)

type AppService struct {
}

func (this *AppService) UpdateManager(app model.App, managerUid int32) (err error) {
	_, err = db.
		SqlBuilder.
		Update("app").
		SetMap(map[string]interface{}{
		"app_manager": app.AppManager,
		"update_by":   managerUid}).
		Where(db.Eq{"app_id": app.AppId}).
		RunWith(db.Sqlx).Exec()
	if err != nil {
		return
	}
	return
}

func (this *AppService) Create(app model.App, managerUid int32) (err error) {
	roolbackFn := func(id int64) {
		db.SqlBuilder.
			Delete("app").
			Where(db.Eq{"id": id}).
			RunWith(db.Sqlx).
			Exec()
	}

	app.AppId = util.GetUUid()
	app.AppKey = util.MD5HexHash(util.Str2bytes(util.GetUUid()))
	if app.SaveMonth < 1 {
		app.SaveMonth = 1
	}

	rows, err := db.
		SqlBuilder.
		Insert("app").
		SetMap(map[string]interface{}{
		"app_name":    app.AppName,
		"descibe":     app.Descibe,
		"app_id":      app.AppId,
		"app_key":     app.AppKey,
		"update_by":   managerUid,
		"create_by":   managerUid,
		"app_manager": managerUid,
		"save_mouth":  app.SaveMonth,
	}).RunWith(db.Sqlx).Exec()

	if err != nil {
		return
	}

	tableId, err := rows.LastInsertId()

	if err != nil {
		return
	}

	eventTableName := "xwl_event" + strconv.Itoa(int(tableId))

	_, err = db.ClickHouseSqlx.Exec(
		`  CREATE TABLE ` + eventTableName + ` ` + sinker.GetClusterSql() + ` (
			xwl_part_date DateTime DEFAULT now(),
			xwl_part_event String,
			xwl_account_id String,
			xwl_distinct_id String,
			xwl_lib_version String,
			xwl_os String,
		
			xwl_screen_width Int64,
			xwl_screen_height Int64,
			xwl_device_id String,
			xwl_network_type String,
			xwl_device_model String,
			xwl_ip  String,
			xwl_city String,
			xwl_province String,
			xwl_lib String,
			xwl_scene String,
			xwl_manufacturer String,
			xwl_os_version String,
			xwl_kafka_offset Int64,
			xwl_kafka_partition Int64
			)ENGINE = ` + sinker.GetMergeTree(eventTableName) + `
			PARTITION BY (toYYYYMM(xwl_part_date))
			ORDER BY (toYYYYMM(xwl_part_date),xwl_part_event) TTL xwl_part_date + toIntervalMonth(` + strconv.Itoa(app.SaveMonth) + `) SETTINGS index_granularity = 8192;`)
	if err != nil {
		roolbackFn(tableId)
		return
	}

	userTableName := "xwl_user" + strconv.Itoa(int(tableId))

	//由于clickhouse修改性能上的缺失，所以本项目修改数据采用增量覆盖的方式
	//因为ddl会修改，所以不采用进行物化视图建立的方案
	_, err = db.ClickHouseSqlx.Exec(
		`CREATE table ` + userTableName + ` ` + sinker.GetClusterSql() + ` (
				   xwl_account_id String, 
				   xwl_distinct_id String,
				   xwl_reg_time Nullable(DateTime),
				   xwl_update_time DateTime,
				   xwl_kafka_offset    Int64,
				   xwl_kafka_partition Int64, 
                   	xwl_ip  String,
					xwl_city String,
					xwl_province String
				)
				ENGINE = ` + sinker.GetReplacingMergeTree(userTableName, utils.ReplacingMergeTreeKey) + `
				ORDER BY xwl_distinct_id
				SETTINGS index_granularity = 8192;`)
	if err != nil {
		roolbackFn(tableId)
		return
	}

	conn := db.RedisPool.Get()
	defer conn.Close()
	if _, err = conn.Do("hset", "AppidToTableid", app.AppId+"_xwl_"+app.AppKey, tableId); err != nil {
		roolbackFn(tableId)
		return
	}
	return
}

func (this AppService) ChangeStatus(app model.App, managerUid int32) (err error) {
	if !util.InArr([]int{1, 0}, *app.IsClose) {
		return errors.New("无效操作")
	}
	_, err = db.
		SqlBuilder.
		Update("app").
		SetMap(map[string]interface{}{
			"is_close":  *app.IsClose,
			"update_by": managerUid}).
		Where(db.Eq{"app_id": app.AppId}).
		RunWith(db.Sqlx).
		Exec()
	if err != nil {
		return err
	}

	switch *app.IsClose {
	case 1:
		myapp.DeleteAppidToTableid(app.AppId, app.AppKey)
	case 0:
		if err = myapp.SetAppidToTableid(app.AppId, app.AppKey, app.Id); err != nil {
			return err
		}
	default:
		return errors.New("无效操作")
	}
	return nil
}

func (this *AppService) List(managerUid int32, app model.App) (list []model.App, count int, err error) {

	selectBuilder := db.
		SqlBuilder.
		Select("*").
		From("app")
	selectBuilder2 := db.
		SqlBuilder.
		Select("count(*)").
		From("app")

	if managerUid != 1 {
		selectBuilder = selectBuilder.Where(db.Eq{"create_by": managerUid})
		selectBuilder2 = selectBuilder2.Where(db.Eq{"create_by": managerUid})
	}

	if app.AppName != "" {
		selectBuilder = selectBuilder.Where(db.Like{"app_name": db.CreateLike(app.AppName)})
		selectBuilder2 = selectBuilder2.Where(db.Like{"app_name": db.CreateLike(app.AppName)})
	}
	if app.IsClose != nil {
		selectBuilder = selectBuilder.Where(db.Eq{"is_close": &app.IsClose})
		selectBuilder2 = selectBuilder2.Where(db.Eq{"is_close": &app.IsClose})
	}

	sql, args, err := selectBuilder.Limit(app.Limit).Offset(db.CreatePage(app.Page, app.Limit)).ToSql()

	if err != nil {
		return
	}

	err = db.
		Sqlx.
		Select(&list, sql, args...)
	if err != nil {
		return
	}

	err = selectBuilder2.RunWith(db.Sqlx).QueryRow().Scan(&count)

	if err != nil {
		return
	}

	return
}

func (this *AppService) ResetAppkey(managerUid int32, app model.App) (err error) {

	app.AppKey = util.MD5HexHash(util.Str2bytes(util.GetUUid()))

	_, err = db.
		SqlBuilder.
		Update("app").
		SetMap(map[string]interface{}{
		"app_key":   app.AppKey,
		"update_by": managerUid,
	}).Where(db.Eq{"app_id": app.AppId}).RunWith(db.Sqlx).Exec()
	if err != nil {
		return
	}

	if err = myapp.DeleteAppidToTableid(app.AppId, app.AppKey); err != nil {
		return
	}

	if err = myapp.SetAppidToTableid(app.AppId, app.AppKey, app.Id); err != nil {
		return
	}
	return
}

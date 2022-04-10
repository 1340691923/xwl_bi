package debug_data

import (
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/model"
)

type DebugData struct {
}

func (this *DebugData) AddDebugDeviceID(appid, deviceID, remark string, managerUid int32) (err error) {

	_, err = db.
		SqlBuilder.
		Insert("debug_device").
		SetMap(map[string]interface{}{
		"remark":    remark,
		"device_id": deviceID,
		"appid":     appid,
		"create_by": managerUid,
	}).RunWith(db.Sqlx).Exec()

	if err != nil {
		return
	}

	Hash := "DebugDeviceID_" + appid
	conn := db.RedisPool.Get()
	defer conn.Close()
	_, err = conn.Do("sadd", Hash, deviceID)

	if err != nil {
		return
	}

	return
}

func (this *DebugData) DelDebugDeviceID(appid, deviceID string, managerUid int32) (err error) {

	_, err = db.SqlBuilder.
		Delete("debug_device").
		Where(db.Eq{"device_id": deviceID, "appid": appid, "create_by": managerUid}).
		RunWith(db.Sqlx).
		Exec()
	if err != nil {
		return
	}
	Hash := "DebugDeviceID_" + appid
	conn := db.RedisPool.Get()
	defer conn.Close()
	_, err = conn.Do("srem", Hash, deviceID)

	if err != nil {
		return
	}

	return
}

func (this *DebugData) DebugDeviceIDList(appid int, managerUid int32) (res []model.DebugDevice, err error) {

	sql, args, err := db.SqlBuilder.
		Select("*").
		From("debug_device").
		Where(db.Eq{"appid": appid, "create_by": managerUid}).
		ToSql()
	if err != nil {
		return
	}

	err = db.Sqlx.Select(&res, sql, args...)
	if err != nil {
		return
	}

	return
}

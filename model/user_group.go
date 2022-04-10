package model

import (
	"errors"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
)

type UserGroup struct {
	Id           int      `db:"id" json:"id"`
	GroupName    string   `db:"group_name" json:"group_name"`
	GroupRemark  string   `db:"group_remark" json:"group_remark"`
	CreateBy     int      `db:"create_by" json:"create_by"`
	UserCount    int      `db:"user_count" json:"user_count"`
	UserList     []byte   `db:"user_list" json:"-"`
	UserListData []string `db:"-" json:"user_list"`
	CreateTime   string   `db:"create_time" json:"create_time"`
	UpdateTime   string   `db:"update_time" json:"update_time"`
}

func (this *UserGroup) Insert(managerUid int32, appid int, userCount int, userList []byte) (err error) {
	_, err = db.
		SqlBuilder.
		Insert("user_group").
		SetMap(map[string]interface{}{
		"group_name":   this.GroupName,
		"group_remark": this.GroupRemark,
		"create_by":    managerUid,
		"user_count":   userCount,
		"appid":        appid,
		"user_list":    userList,
	}).RunWith(db.Sqlx).Exec()
	if err != nil {
		if util.IsMysqlRepeatError(err) {
			return errors.New("分群名重复，请重新填写")
		}
		return err
	}
	return nil
}

func (this *UserGroup) ModifyUserGroup(managerUid int32, appid int) (err error) {
	_, err = db.SqlBuilder.
		Update("user_group").
		SetMap(map[string]interface{}{
			"group_name":   this.GroupName,
			"group_remark": this.GroupRemark,
		}).
		Where(
			db.Eq{
				"create_by": managerUid,
				"id":        this.Id,
				"appid":     appid,
			}).RunWith(db.Sqlx).Exec()
	if err != nil {
		if util.IsMysqlRepeatError(err) {
			return errors.New("分群名重复，请重新填写")
		}
		return err
	}
	return nil
}

func (this *UserGroup) DeleteUserGroupById(managerUid int32, appid int) (err error) {
	_, err = db.SqlBuilder.
		Delete("user_group").
		Where(db.Eq{"create_by": managerUid, "id": this.Id, "appid": appid}).
		RunWith(db.Sqlx).
		Exec()
	return
}

func (this *UserGroup) List(managerUid int32, appid int) (list []UserGroup, err error) {
	SQL, args, err := db.SqlBuilder.
		Select("id,group_name,group_remark,create_by,user_count,user_list,create_time,update_time").
		From("user_group").
		Where(db.Eq{"create_by": managerUid, "appid": appid}).
		ToSql()
	if err != nil {
		return nil, err
	}

	if err := db.Sqlx.Select(&list, SQL, args...); err != nil {
		return nil, err
	}
	return list, err
}

func (this *UserGroup) GetSelectOptions(managerUid int32, appid int) (list []UserGroup, err error) {
	SQL, args, err := db.SqlBuilder.
		Select("id,group_name").
		From("user_group").
		Where(
			db.Eq{
				"create_by": managerUid,
				"appid":     appid,
			}).
		ToSql()
	if err != nil {
		return
	}

	err = db.Sqlx.Select(&list, SQL, args...)
	if err != nil {
		return
	}

	return

}

package model

import (
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"go.uber.org/zap"
	"strings"
)

// GmUserModel BI用户
type GmUserModel struct {
	ID            int32  `json:"id" db:"id"`
	RoleId        int32  `json:"role_id" db:"role_id"`
	Username      string `json:"username" db:"username"`
	Password      string `json:"password" db:"password"`
	Realname      string `json:"realname" db:"realname"`
	CreateTime    string `db:"create_time" json:"create_time"`
	UpdateTime    string `db:"update_time" json:"update_time"`
	LastLoginTime string `db:"last_login_time" json:"last_login_time"`

	IsDel int8 `db:"is_del" json:"is_del"` //是否禁止该账号
}

//密码进行md5混淆
func (this GmUserModel) GetPassword() string {
	return util.MD5HexHash(util.Str2bytes(this.Password))
}

//是否存在该用户
func (this GmUserModel) Exsit() (b bool) {
	var count int
	err := db.Sqlx.Get(&count, "select count(*) from gm_user where username = ? and role_id = ? limit 1;", this.Username, this.RoleId)
	if err != nil || count == 0 {
		logs.Logger.Error("err", zap.Error(err))
		return false
	}
	return true
}

//登录
func (this GmUserModel) GetUserByUP() (gmUser GmUserModel, err error) {
	err = db.Sqlx.Get(&gmUser, "select id,username,password,role_id,realname,is_del from gm_user where username = ? and password = ? limit 1;", this.Username, this.GetPassword())
	return
}

//通过id查询用户信息
func (this GmUserModel) GetUserById() (gmUser GmUserModel, err error) {
	err = db.Sqlx.Get(&gmUser, "select id,username,password,role_id,realname from gm_user where id = ?;", this.ID)
	return
}

//新增用户
func (this GmUserModel) Insert() (id int64, err error) {
	rlt, err := db.Sqlx.Exec("insert into gm_user(username,password,role_id,realname)values(?,?,?,?)", this.Username, this.GetPassword(), this.RoleId, this.Realname)
	if err != nil {
		return
	}
	id, _ = rlt.LastInsertId()
	return
}

// Update
func (this GmUserModel) Update() (err error) {
	if strings.TrimSpace(this.Password) == ""{
		_, err = db.Sqlx.Exec("update gm_user set username = ?,role_id=?,realname=? where id = ? ;", this.Username, this.RoleId, this.Realname, this.ID)
		return
	}
	_, err = db.Sqlx.Exec("update gm_user set username = ?,password=?,role_id=?,realname=? where id = ? ;", this.Username, this.GetPassword(), this.RoleId, this.Realname, this.ID)
	return
}

// Update
func (this GmUserModel) UpdatePassById() (err error) {
	_, err = db.Sqlx.Exec("update gm_user set password=? where id = ? ;", this.GetPassword(), this.ID)
	return
}

// Select
func (this GmUserModel) Select(appid string) (gmUser []GmUserModel, err error) {
	if appid == "" {
		err = db.Sqlx.Select(&gmUser, "select * from gm_user;")
	} else {
		err = db.Sqlx.Select(&gmUser, "select * from gm_user where FIND_IN_SET(id,(select app_manager from app where id = "+appid+"));")
	}

	return
}

// Delete
func (this GmUserModel) Delete() (err error) {
	_, err = db.Sqlx.Exec("delete from gm_user where id = ? ;", this.ID)
	if err != nil {
		return
	}

	_, err = db.Sqlx.Exec(`UPDATE app SET app_manager = TRIM(BOTH ',' FROM REPLACE(CONCAT(',', app_manager, ','), concat(',',?,','), ','))
						 WHERE FIND_IN_SET(?, app_manager) `, this.ID, this.ID)
	if err != nil {
		return
	}

	return
}

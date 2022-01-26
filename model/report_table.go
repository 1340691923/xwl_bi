package model

import "github.com/1340691923/xwl_bi/engine/db"

type ReportTable struct {
	Id         int    `db:"id" json:"id"`
	Appid      int    `db:"appid" json:"appid"`
	UserId     int    `db:"user_id" json:"user_id"`
	Name       string `db:"name" json:"name"`
	RtType     int8   `db:"rt_type" json:"rt_type"`
	Data       string `db:"data" json:"data"`
	CreateTime string `db:"create_time" json:"create_time"`
	UpdateTime string `db:"update_time" json:"update_time"`
	Remark     string `db:"remark" json:"remark"`
}

func (this *ReportTable) InsertOrUpdate() (err error) {
	sql := `insert into report_table(appid,user_id,name,rt_type,data,remark)values(?,?,?,?,?,?) on duplicate key update data=?,remark=?`
	_, err = db.Sqlx.Exec(sql, this.Appid, this.UserId, this.Name, this.RtType, this.Data, this.Remark, this.Data, this.Remark)
	return
}

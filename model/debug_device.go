package model

type DebugDevice struct {
	Id         int    `db:"id" json:"id"`
	Appid      int    `db:"appid" json:"appid"`
	DeviceId   string `db:"device_id" json:"device_id"`
	Remark     string `db:"remark" json:"remark"`
	CreateBy   int    `db:"create_by" json:"create_by"`
	CreateTime string `db:"create_time" json:"create_time"`
}

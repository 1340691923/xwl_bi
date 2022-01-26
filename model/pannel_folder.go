package model

type PannelFolder struct {
	Id         int    `db:"id" json:"id"`
	FolderName string `db:"folder_name" json:"folder_name"`
	FolderType int8   `db:"folder_type" json:"folder_type"` //0为自己创建的
	CreateBy   int    `db:"create_by" json:"create_by"`
	CreateTime string `db:"create_time" json:"create_time"`
	UpdateTime string `db:"update_time" json:"update_time"`
	Appid      int    `db:"appid" json:"appid"`
}

type Pannel struct {
	Id           int    `db:"id" json:"id"`
	FolderId     int    `db:"folder_id" json:"folder_id"`
	PannelName   string `db:"pannel_name" json:"pannel_name"`
	Managers     string `db:"managers" json:"managers"`
	CreateBy     int    `db:"create_by" json:"create_by"`
	CreateTime   int64  `db:"create_time" json:"create_time"`
	UpdateTime   int64  `db:"update_time" json:"update_time"`
	ReportTables string `db:"report_tables" json:"report_tables"`
}

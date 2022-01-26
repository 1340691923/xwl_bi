package response

type FailDataRes struct {
	IntervalDate  string `json:"interval_date" db:"interval_date"`
	Year          string `json:"year" db:"year"`
	StartMinute   string `json:"start_minute" db:"start_minute"`
	EndMinute     string `json:"end_minute" db:"end_minute"`
	Count         int    `json:"count" db:"count"`
	ErrorReason   string `json:"error_reason" db:"error_reason"`
	ErrorHandling string `json:"error_handling" db:"error_handling"`
	ReportType    string `json:"report_type" db:"report_type"`
}

type ReportCountRes struct {
	DataName      string `json:"data_name"`
	ShowName      string `json:"show_name"`
	ReceivedCount int    `json:"received_count"`
	SuccCount     int    `json:"succ_count"`
	FailCount     int    `json:"fail_count"`
}

type EventFailDescRes struct {
	ErrorReason string `json:"error_reason" db:"error_reason"`
	Count       int    `json:"count" db:"count"`
	ReportData  string `json:"report_data" db:"report_data"`
}

type MetaEventListRes struct {
	EventName      string `json:"event_name" db:"event_name"`
	ShowName       string `json:"show_name" db:"show_name"`
	YesterdayCount string `json:"yesterday_count" db:"yesterday_count"`
}

type AttributeRes struct {
	AttributeName   string `db:"attribute_name" json:"attribute_name"` //属性名
	ShowName        string `db:"show_name" json:"show_name"`           //显示名
	DataType        int    `db:"data_type" json:"data_type"`           //数据类型
	AttributeType   int8   `db:"attribute_type" json:"attribute_type"` //默认为1 （1为预置属性，2为自定义属性）
	DataTypeFormat  string `db:"-" json:"data_type_format"`            //数据类型
	AttributeSource int    `json:"attribute_source" db:"attribute_source"`
	Status          int    `json:"status"`
}

type AttrCalcuSymbolData struct {
	AttributeName  string `db:"attribute_name" json:"attribute_name"` //属性名
	ShowName       string `db:"show_name" json:"show_name"`           //显示名
	DataType       int    `db:"data_type" json:"data_type"`           //数据类型
	AttributeType  int8   `db:"attribute_type" json:"attribute_type"` //默认为1 （1为预置属性，2为自定义属性）
	DataTypeFormat string `db:"-" json:"data_type_format"`            //数据类型
}

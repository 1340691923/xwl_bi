package model

import (
	"bytes"
)

const (
	UserReportType  = 1
	EventReportType = 2
	Debug           = 1
	Release         = 0
)

type KafkaData struct {
	APPID           string `json:"appid"`
	DistinctId      string `json:"distinct_id"`
	TableId         string `json:"table_id"`
	Ip              string `json:"ip"`
	ReportType      int    `json:"report_type"`
	Debug           string `json:"debug"`
	ReqData         []byte `json:"req_data"`
	ReportTime      string `json:"report_time"`
	ConsumptionTime string `json:"consumption_time"`
	EventName       string `json:"event_name"`
	Offset          int64  `json:"offset"`
}

func (this *KafkaData) GetTableName() (tableName string) {
	buff := bytes.Buffer{}

	switch this.ReportType {
	case UserReportType:
		buff.WriteString("xwl_user")
	case EventReportType:
		buff.WriteString("xwl_event")
	}

	buff.WriteString(this.TableId)
	tableName = buff.String()
	return
}

func (this *KafkaData) GetReportTypeErr() (eventType string) {

	switch this.ReportType {
	case UserReportType:
		eventType = "用户属性类型不合法"
	case EventReportType:
		eventType = "事件属性类型不合法"
	}
	return eventType
}

type EventData struct {
	AccountId       string                 `json:"xwl_account_id"`
	XwlDistinctId   string                 `json:"xwl_distinct_id"`
	XwlUserId       string                 `json:"xwl_user_id"`
	XwlPartEvent    string                 `json:"xwl_part_event"`
	XwlPartDate     string                 `json:"xwl_part_date"`
	XwlMpPlatform   string                 `json:"xwl_mp_platform"`
	XwlLibVersion   string                 `json:"xwl_lib_version"`
	XwlOs           string                 `json:"xwl_os"`
	XwlScreenWidth  string                 `json:"xwl_screen_width"`
	XwlCountryCode  string                 `json:"xwl_country_code"`
	XwlScreenHeight string                 `json:"xwl_screen_height"`
	XwlDeviceId     string                 `json:"xwl_device_id"`
	XwlNetworkType  string                 `json:"xwl_network_type"`
	XwlDeviceModel  string                 `json:"xwl_device_model"`
	XwlCity         string                 `json:"xwl_city"`
	XwlProvince     string                 `json:"xwl_province"`
	XwlLib          string                 `json:"xwl_lib"`
	XwlScene        string                 `json:"xwl_scene"`
	XwlManufacturer string                 `json:"xwl_manufacturer"`
	XwlOsVersion    string                 `json:"xwl_os_version"`
	XwlKafkaOffset  int64                  `json:"xwl_kafka_offset"`
	XwlIp           string                 `json:"xwl_ip"`
	TableId         string                 `json:"table_id"`
	Properties      map[string]interface{} `json:"properties"`
}

package parser

import (
	"fmt"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/platform-basic-libs/sinker/model"
	"strings"
)

const (
	TypeUnknown = iota
	Int
	Float
	String
	DateTime
	ElasticDateTime
	IntArray
	FloatArray
	StringArray
	DateTimeArray
)

var TypeRemarkMap = map[int]string{
	TypeUnknown:     "未知类型",
	Int:             "数字类型",
	Float:           "浮点数类型",
	String:          "字符串类型",
	DateTime:        "时间类型",
	ElasticDateTime: "时间类型",
	IntArray:        "数字数组类型",
	FloatArray:      "浮点数数组类型",
	StringArray:     "字符串数组类型",
	DateTimeArray:   "时间数组类型",
}

/**系统字段*/
var SysColumn = map[string]string{
	"xwl_account_id":      "账户ID",
	"xwl_distinct_id":     "访客ID",
	"xwl_reg_time":        "用户注册时间",
	"xwl_server_time":     "服务端入库时间",
	"xwl_update_time":     "用户信息最后修改时间",
	"xwl_kafka_offset":    "kafka偏移量",
	"xwl_part_date":       "事件入库时间",
	"xwl_part_event":      "事件名",
	"xwl_lib_version":     "SDK版本",
	"xwl_os":              "用户操作系统",
	"xwl_ip":              "事件生成时用户所在IP",
	"xwl_screen_width":    "用户屏幕宽度",
	"xwl_screen_height":   "用户屏幕高度",
	"xwl_device_id":       "用户设备ID",
	"xwl_network_type":    "用户网络类型",
	"xwl_device_model":    "用户机型",
	"xwl_city":            "用户所在城市",
	"xwl_province":        "用户所在省份",
	"xwl_lib":             "SDK语言",
	"xwl_scene":           "场景值",
	"xwl_manufacturer":    "设备提供商",
	"xwl_os_version":      "用户操作系统版本",
	"xwl_cpu":             "cpu",
	"xwl_client_time":     "客户端上报时间",
	"xwl_browser_version": "浏览器版本号",
	"xwl_browser":         "浏览器类型",
	"xwl_kafka_partition": "kafka分区",
}

type TypeInfo struct {
	Type     int
	Nullable bool
}

var (
	typeInfo map[string]TypeInfo
)

func GetValueByType(metric *FastjsonMetric, cwt *model.ColumnWithType) (val interface{}) {
	name := cwt.SourceName
	switch cwt.Type {
	case Int:
		val = metric.GetInt(name, cwt.Nullable)
	case Float:
		val = metric.GetFloat(name, cwt.Nullable)
	case String:
		val = metric.GetString(name, cwt.Nullable)
	case DateTime:
		val = metric.GetDateTime(name, cwt.Nullable)
	case ElasticDateTime:
		val = metric.GetElasticDateTime(name, cwt.Nullable)
	case IntArray:
		val = metric.GetArray(name, Int)
	case FloatArray:
		val = metric.GetArray(name, Float)
	case StringArray:
		val = metric.GetArray(name, String)
	case DateTimeArray:
		val = metric.GetArray(name, DateTime)
	default:
		logs.Logger.Sugar().Errorf("未知TYPE:%s", cwt.Type)
	}
	return
}

func init() {
	primTypeInfo := make(map[string]TypeInfo)
	typeInfo = make(map[string]TypeInfo)
	for _, t := range []string{"UInt8", "UInt16", "UInt32", "UInt64", "Int8", "Int16", "Int32", "Int64"} {
		primTypeInfo[t] = TypeInfo{Type: Int, Nullable: false}
	}
	for _, t := range []string{"Float32", "Float64"} {
		primTypeInfo[t] = TypeInfo{Type: Float, Nullable: false}
	}
	for _, t := range []string{"String"} {
		primTypeInfo[t] = TypeInfo{Type: String, Nullable: false}
	}
	for _, t := range []string{"Date", "DateTime"} {
		primTypeInfo[t] = TypeInfo{Type: DateTime, Nullable: false}
	}
	primTypeInfo["ElasticDateTime"] = TypeInfo{Type: ElasticDateTime, Nullable: false}
	for k, v := range primTypeInfo {
		typeInfo[k] = v
		nullK := fmt.Sprintf("Nullable(%s)", k)
		typeInfo[nullK] = TypeInfo{Type: v.Type, Nullable: true}
		arrK := fmt.Sprintf("Array(%s)", k)
		switch v.Type {
		case Int:
			typeInfo[arrK] = TypeInfo{Type: IntArray, Nullable: false}
		case Float:
			typeInfo[arrK] = TypeInfo{Type: FloatArray, Nullable: false}
		case String:
			typeInfo[arrK] = TypeInfo{Type: StringArray, Nullable: false}
		case DateTime:
			typeInfo[arrK] = TypeInfo{Type: DateTimeArray, Nullable: false}
		}
	}
}

func WhichType(typ string) (dataType int, nullable bool) {
	ti, ok := typeInfo[typ]
	if ok {
		dataType, nullable = ti.Type, ti.Nullable
		return
	}
	nullable = strings.HasPrefix(typ, "Nullable(")
	if nullable {
		typ = typ[len("Nullable(") : len(typ)-1]
	}
	if strings.HasPrefix(typ, "DateTime64") {
		dataType = DateTime
	} else if strings.HasPrefix(typ, "Array(DateTime64") {
		dataType = DateTimeArray
		nullable = false
	} else if strings.HasPrefix(typ, "Decimal") {
		dataType = Float
	} else if strings.HasPrefix(typ, "Array(Decimal") {
		dataType = FloatArray
		nullable = false
	} else if strings.HasPrefix(typ, "FixedString") {
		dataType = String
	} else if strings.HasPrefix(typ, "Array(FixedString") {
		dataType = StringArray
		nullable = false
	} else {
		logs.Logger.Sugar().Errorf("未知TYPE:%s", typ)
	}
	typeInfo[typ] = TypeInfo{Type: dataType, Nullable: nullable}
	return
}

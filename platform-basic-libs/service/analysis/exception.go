package analysis

// 内置异常
const (
	TimeError           int = 60001
	ZhiBiaoNumError     int = 60002
	GroupNumError       int = 60003
	GroupEmptyError     int = 60004
	UIEmptyError        int = 60005
	EventNameEmptyError int = 60006
)

// 内置异常表
var ERROR_TABLE = map[int]string{
	TimeError:           "筛选时间异常",
	ZhiBiaoNumError:     "筛选指标个数异常",
	GroupNumError:       "筛选分组个数异常",
	GroupEmptyError:     "筛选分组不能为空字段",
	UIEmptyError:        "用户id列表不能为空",
	EventNameEmptyError: "事件名不能为空",
}

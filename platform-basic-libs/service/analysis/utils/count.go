package utils

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
)

type CountType = string

const SPLIT = "$$$xwl$$$"

const (
	UserNum        = "1"
	AllSum         = "2"
	AvgCount       = "3"
	AvgSumByUser   = "4"
	MiddleCount    = "5"
	MaxCount       = "6"
	MinCount       = "7"
	DistincCount   = "8"
	MiddleCount5   = "9"
	MiddleCount10  = "10"
	MiddleCount20  = "11"
	MiddleCount25  = "12"
	MiddleCount30  = "13"
	MiddleCount40  = "14"
	MiddleCount60  = "15"
	MiddleCount70  = "16"
	MiddleCount75  = "17"
	MiddleCount80  = "18"
	MiddleCount90  = "19"
	MiddleCount95  = "20"
	MiddleCount99  = "21"
	AllCount       = "A1"
	ClickUserNum   = "A2"
	AvgCountByUser = "A3"
)

func allCount(col string) string {

	if col == Default {
		return toString("count()")
	}

	return toString(fmt.Sprintf("count(%s)", col))
}

func clickUserNum(col string) string {
	if col == Default {
		return toString(Round(NaN2Zero("COUNT(DISTINCT xwl_distinct_id)")))
	}

	return toString(Round(NaN2Zero(fmt.Sprintf("COUNT(DISTINCT %s)", col))))
}

var IntPropQuotas = map[string]string{
	AllSum:        "总和",
	AvgCount:      "均值",
	AvgSumByUser:  "人均值",
	MaxCount:      "最大值",
	MinCount:      "最小值",
	DistincCount:  "去重数",
	MiddleCount:   "中位数",
	MiddleCount5:  "5分位数",
	MiddleCount10: "10分位数",
	MiddleCount20: "20分位数",
	MiddleCount25: "25分位数",
	MiddleCount30: "30分位数",
	MiddleCount40: "40分位数",
	MiddleCount60: "60分位数",
	MiddleCount70: "70分位数",
	MiddleCount75: "75分位数",
	MiddleCount80: "80分位数",
	MiddleCount90: "90分位数",
	MiddleCount95: "95分位数",
	MiddleCount99: "99分位数",
}

var StringPropQuotas = map[string]string{
	DistincCount: "去重数",
}

const Default = "默认"

type getCountCol = func(col string) string

var autoAddId int32

func NaN2Zero(fn string) string {

	autoAddId = atomic.AddInt32(&autoAddId, 1)

	uuid := "col_" + strconv.Itoa(int(autoAddId))
	return fmt.Sprintf("if(isInfinite(if(isNaN(%s as %s),0,%s)),0,%s)", fn, uuid, uuid, uuid)
}

func toString(fn string) string {
	return fmt.Sprintf("toString(%s)", fn)
}

func Round(fn string) string {
	return fmt.Sprintf("round(%s,2)", fn)
}

func Divide(fn string) string {
	return fmt.Sprintf("divide(%s)", fn)
}

func ToFloat32OrZero(fn string) string {
	return fmt.Sprintf("toFloat64OrZero(CAST(%s,'String'))", fn)
}

var CountTypMap = map[CountType]getCountCol{
	UserNum: func(col string) string {

		if col == Default {
			return toString("count()")
		}

		return toString(fmt.Sprintf("count(%s)", col))
	},
	AllSum: func(col string) string {
		return toString(Round(NaN2Zero(fmt.Sprintf("sum(%s)", col))))
	},
	AvgCount: func(col string) string {
		return toString(Round(NaN2Zero(fmt.Sprintf(" avg(%s) ", col))))
	},
	AvgSumByUser: func(col string) string {
		return toString(Round(NaN2Zero(fmt.Sprintf("SUM(%s)/COUNT(DISTINCT xwl_distinct_id)", col))))
	},
	MiddleCount: func(col string) string {
		return toString(NaN2Zero(fmt.Sprintf("quantile(%s)", col)))
	},
	MaxCount: func(col string) string {
		return toString(Round(NaN2Zero(fmt.Sprintf("max(%s)", col))))
	},
	MinCount: func(col string) string {
		return toString(Round(NaN2Zero(fmt.Sprintf("min(%s)", col))))
	},
	DistincCount: func(col string) string {
		return toString(Round(NaN2Zero(fmt.Sprintf("count(DISTINCT %s)", col))))
	},
	AllCount:     allCount,
	ClickUserNum: clickUserNum,
	AvgCountByUser: func(col string) string {
		if col == Default {
			return toString(Round(NaN2Zero(fmt.Sprintf("%s/%s", ToFloat32OrZero(allCount(col)), ToFloat32OrZero(clickUserNum(col))))))
		}

		arr := strings.Split(col, SPLIT)

		return toString(Round(NaN2Zero(fmt.Sprintf("%s/%s", ToFloat32OrZero(allCount(arr[0])), ToFloat32OrZero(clickUserNum(arr[1]))))))
	},
	MiddleCount5: func(col string) string {
		return toString(getQuantile(0.05, col))
	},
	MiddleCount10: func(col string) string {
		return toString(getQuantile(0.1, col))
	},
	MiddleCount20: func(col string) string {
		return toString(getQuantile(0.2, col))
	},
	MiddleCount25: func(col string) string {
		return toString(getQuantile(0.25, col))
	},
	MiddleCount30: func(col string) string {
		return toString(getQuantile(0.3, col))
	},
	MiddleCount40: func(col string) string {
		return toString(getQuantile(0.4, col))
	},
	MiddleCount60: func(col string) string {
		return toString(getQuantile(0.6, col))
	},
	MiddleCount70: func(col string) string {
		return toString(getQuantile(0.7, col))
	},
	MiddleCount75: func(col string) string {
		return toString(getQuantile(0.75, col))
	},
	MiddleCount80: func(col string) string {
		return toString(getQuantile(0.8, col))
	},
	MiddleCount90: func(col string) string {
		return toString(getQuantile(0.9, col))
	},
	MiddleCount95: func(col string) string {
		return toString(getQuantile(0.95, col))
	},
	MiddleCount99: func(col string) string {
		return toString(getQuantile(0.99, col))
	},
}

func getQuantile(sum float64, col string) string {
	return NaN2Zero(fmt.Sprintf(" quantile(%v)(%s) ", sum, col))
}

const (
	TwoDecimalPlaces = "1"
	Percentage       = "2"
	Rounding         = "3"
)

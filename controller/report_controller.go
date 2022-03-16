package controller

import (
	"errors"
	"fmt"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/my_error"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/report"
	"github.com/1340691923/xwl_bi/platform-basic-libs/sinker"
	parser "github.com/1340691923/xwl_bi/platform-basic-libs/sinker/parse"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"math"
	"strings"
	"time"
)

type ReportController struct {
	BaseController
}

//上报接口
func (this ReportController) ReportAction(ctx *fasthttp.RequestCtx) {

	if strings.ToUpper(util.Bytes2str(ctx.Method())) == "OPTIONS" {
		return
	}

	var (
		typ       = ctx.UserValue("typ").(string)
		appid     = ctx.UserValue("appid").(string)
		appkey    = ctx.UserValue("appkey").(string)
		debug     = ctx.UserValue("debug").(string)
		eventName = ctx.UserValue("eventName").(string)
		body      = ctx.Request.Body()
	)
	if strings.TrimSpace(eventName) == "" {
		this.FastError(ctx, errors.New("事件名 不能为空"))
		return
	}

	if strings.TrimSpace(appid) == "" {
		this.FastError(ctx, errors.New("appid 不能为空"))
		return
	}

	reportService := report.ReportService{}

	tableId, err := reportService.GetTableid(appid, appkey)
	if err != nil {
		this.FastError(ctx, err)
		return
	}

	duck, err := report.GetReportDuck(typ)

	if err != nil {
		this.FastError(ctx, err)
		return
	}

	defer duck.Put()

	gjsonArr := gjson.GetManyBytes(body, "xwl_distinct_id", "xwl_ip", "xwl_part_date")

	xwlDistinctId := gjsonArr[0].String()
	xwlIp := gjsonArr[1].String()
	xwlPartDate := gjsonArr[2].String()
	if xwlDistinctId == "" {
		this.FastError(ctx, errors.New("xwl_distinct_id 不能为空"))
		return
	}

	if xwlIp == "" {
		xwlIp = util.CtxClientIP(ctx)
	}

	if xwlPartDate == "" {
		xwlPartDate = time.Now().Format(util.TimeFormat)
	}

	duck.NewReportType(appid, tableId, debug, xwlPartDate, eventName, xwlIp, ctx.PostBody())

	if reportService.IsDebugUser(debug, xwlDistinctId, tableId) {
		kafkaData := duck.GetkafkaData()

		pp := parser.FastjsonParser{}

		metric, debugErr := pp.Parse(kafkaData.ReqData)

		if debugErr != nil {
			logs.Logger.Error("parser.ParseKafkaData ", zap.Error(err))
			this.FastError(ctx, errors.New("服务异常"))
			return
		}

		dims, err := sinker.GetDims(model.GlobConfig.Comm.ClickHouse.DbName, kafkaData.GetTableName(), []string{}, db.ClickHouseSqlx, true)
		if err != nil {
			logs.Logger.Error("sinker.GetDims", zap.Error(err))
			this.FastError(ctx, errors.New("服务异常"))
			return
		}
		obj := metric.GetParseObject()
		m := map[string]interface{}{
			"data_name":   kafkaData.EventName,
			"report_data": util.Bytes2str(ctx.PostBody()),
			"report_time": kafkaData.ReportTime,
			"appid":       kafkaData.TableId,
			"distinct_id": xwlDistinctId,
		}
		haveFailAttr := false

		var eventType = kafkaData.GetReportTypeErr()

		for _, column := range dims {
			if obj.Get(column.Name) != nil {
				reportType := parser.FjDetectType(obj.Get(column.Name))
				if reportType != column.Type {
					if !(reportType == parser.Int && column.Type == parser.Float) && !(reportType == parser.Float && column.Type == parser.Int) {
						errorReason := fmt.Sprintf("%s的类型错误，正确类型为%v，上报类型为%v(%v)", column.Name, parser.TypeRemarkMap[column.Type], parser.TypeRemarkMap[reportType], obj.Get(column.Name).String())
						haveFailAttr = true
						m["error_reason"] = errorReason
						m["data_judge"] = eventType
					}
				}
			}
		}

		xwlUpdateTime := gjson.GetBytes(body, "xwl_update_time").String()
		clinetT := util.Str2Time(xwlUpdateTime, util.TimeFormat)
		serverT := util.Str2Time(kafkaData.ReportTime, util.TimeFormat)
		if math.Abs(serverT.Sub(clinetT).Minutes()) > 10 {
			m["error_reason"] = "客户端上报时间误差大于十分钟"
			m["data_judge"] = eventType
		}

		if !haveFailAttr {
			m["data_judge"] = "数据检验通过"
		}

		err = reportService.InflowOfDebugData(m, eventName)

		if err != nil {
			logs.Logger.Error("reportService.InflowOfDebugData", zap.Error(err))
			this.FastError(ctx, errors.New("服务异常"))
			return
		}

		if haveFailAttr {
			logs.Logger.Error("reportService.InflowOfDebugData", zap.String("error_reason", m["error_reason"].(string)))
			this.FastError(ctx, my_error.NewError(m["error_reason"].(string), 10006))
			return
		}
		if debug == report.DebugNotToDB {
			this.Output(ctx, map[string]interface{}{
				"code": 0,
				"msg":  "上报成功（数据不入库）",
			})
			return
		}
	}

	err = duck.InflowOfKakfa()
	if err != nil {
		this.FastError(ctx, err)
		return
	}

	ctx.WriteString(`{"code":0,"msg":"上报成功"}`)
	return
}

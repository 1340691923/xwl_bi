package controller

import (
	"github.com/1340691923/xwl_bi/platform-basic-libs/request"
	"github.com/1340691923/xwl_bi/platform-basic-libs/response"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/analysis"
	parser "github.com/1340691923/xwl_bi/platform-basic-libs/sinker/parse"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type BehaviorAnalysisController struct {
	BaseController
}

//获取分析面板初始化数据
func (this BehaviorAnalysisController) GetConfigs(ctx *fiber.Ctx) error {

	type ReqData struct {
		Appid int32 `json:"appid"`
	}
	var reqData ReqData
	err := ctx.BodyParser(&reqData)
	if err != nil {
		return this.Error(ctx, err)
	}

	behaviorAnalysisService := analysis.BehaviorAnalysisService{}

	eventNameList, attributeMap, err := behaviorAnalysisService.GetConfigs(int(reqData.Appid))

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, map[string]interface{}{"event_name_list": eventNameList, "attributeMap": attributeMap, "sys_col": parser.SysColumn})

}

//根据事件名查找指标的运算函数
func (this BehaviorAnalysisController) LoadPropQuotas(ctx *fiber.Ctx) error {

	var reqData request.LoadPropQuotasReq

	err := ctx.BodyParser(&reqData)
	if err != nil {
		return this.Error(ctx, err)
	}

	if err := this.CheckParameter([]request.CheckConfigStruct{
		{
			request.EmptyEventError,
			"event_name",
		},
	}, ctx); err != nil {
		return this.Error(ctx, err)
	}

	behaviorAnalysisService := analysis.BehaviorAnalysisService{}

	attributeNameList, err := behaviorAnalysisService.LoadPropQuotas(reqData)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, attributeNameList)

}

//获取上报字段所有的值
func (this BehaviorAnalysisController) GetValues(ctx *fiber.Ctx) error {

	type ReqData struct {
		Appid int32  `json:"appid"`
		Table string `json:"table"`
		Col   string `json:"col"`
	}
	var reqData ReqData
	err := ctx.BodyParser(&reqData)
	if err != nil {
		return this.Error(ctx, err)
	}

	appid := strconv.Itoa(int(reqData.Appid))
	table := reqData.Table
	col := reqData.Col

	behaviorAnalysisService := analysis.BehaviorAnalysisService{}

	values, err := behaviorAnalysisService.GetValues(appid, table, col, ctx.Body())

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, values)
}

//事件分析查询
func (this BehaviorAnalysisController) EventList(ctx *fiber.Ctx) error {

	i, err := analysis.NewAnalysisByCommand(analysis.EventComand, ctx.Body())

	if err != nil {
		return this.Error(ctx, err)
	}

	res, err := analysis.GetAnalysisRes(i)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, res)
}

//漏斗分析查询
func (this BehaviorAnalysisController) FunnelList(ctx *fiber.Ctx) error {

	i, err := analysis.NewAnalysisByCommand(analysis.FunnelComand, ctx.Body())

	if err != nil {
		return this.Error(ctx, err)
	}

	res, err := analysis.GetAnalysisRes(i)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, res)
}

//留存分析查询
func (this BehaviorAnalysisController) RetentionList(ctx *fiber.Ctx) error {

	i, err := analysis.NewAnalysisByCommand(analysis.RetentionComand, ctx.Body())

	if err != nil {
		return this.Error(ctx, err)
	}

	res, err := analysis.GetAnalysisRes(i)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, res)
}

//用户属性分析查询
func (this BehaviorAnalysisController) UserAttrList(ctx *fiber.Ctx) error {

	i, err := analysis.NewAnalysisByCommand(analysis.UserAttrCommand, ctx.Body())

	if err != nil {
		return this.Error(ctx, err)
	}

	res, err := analysis.GetAnalysisRes(i)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, res)
}

//智能路径分析查询
func (this BehaviorAnalysisController) TraceList(ctx *fiber.Ctx) error {

	i, err := analysis.NewAnalysisByCommand(analysis.TraceComand, ctx.Body())

	if err != nil {
		return this.Error(ctx, err)
	}

	res, err := analysis.GetAnalysisRes(i)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, res)
}

//用户列表查询
func (this BehaviorAnalysisController) UserList(ctx *fiber.Ctx) error {

	i, err := analysis.NewAnalysisByCommand(analysis.UserListCommand, ctx.Body())

	if err != nil {
		return this.Error(ctx, err)
	}

	res, err := analysis.GetAnalysisRes(i)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, res)
}

//查询用户访问过的事件详情
func (this BehaviorAnalysisController) UserEventDetailList(ctx *fiber.Ctx) error {

	i, err := analysis.NewAnalysisByCommand(analysis.UserEventDetailListCommand, ctx.Body())

	if err != nil {
		return this.Error(ctx, err)
	}

	res, err := analysis.GetAnalysisRes(i)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, res)
}

//查询用户访问过的事件统计
func (this BehaviorAnalysisController) UserEventCountList(ctx *fiber.Ctx) error {

	i, err := analysis.NewAnalysisByCommand(analysis.UserEventCountCommand, ctx.Body())

	if err != nil {
		return this.Error(ctx, err)
	}

	res, err := analysis.GetAnalysisRes(i)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, res)
}

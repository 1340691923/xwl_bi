package controller

import (
	"fmt"
	"github.com/1340691923/xwl_bi/platform-basic-libs/request"
	"github.com/1340691923/xwl_bi/platform-basic-libs/response"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/analysis"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/meta_data"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type MetaDataController struct {
	BaseController
}

//元事件列表
func (this MetaDataController) MetaEventList(ctx *fiber.Ctx) error {

	type ReqData struct {
		Appid int `json:"appid"`
	}

	var reqData ReqData

	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}

	appid := strconv.Itoa(reqData.Appid)

	metaData := meta_data.MetaDataService{Appid: appid}
	if err := metaData.UpdateYesterdayCount(); err != nil {
		return this.Error(ctx, err)
	}
	res, err := metaData.MetaEventList()
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, map[string]interface{}{"list": res})
}

//元事件列表（通过属性查找）
func (this MetaDataController) MetaEventListByAttr(ctx *fiber.Ctx) error {

	type ReqData struct {
		Appid int    `json:"appid"`
		Attr  string `json:"attr"`
	}

	var reqData ReqData

	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}

	appid := strconv.Itoa(reqData.Appid)
	attr := reqData.Attr

	metaData := meta_data.MetaDataService{Appid: appid}

	res, err := metaData.MetaEventListByAttr(attr)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, map[string]interface{}{"list": res})
}

//查看上报属性列表
func (this MetaDataController) AttrManager(ctx *fiber.Ctx) error {

	var reqData request.AttrManagerReq

	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}

	appid := reqData.Appid
	typ := reqData.Typ

	metaData := meta_data.MetaDataService{Appid: strconv.Itoa(appid)}

	res, err := metaData.AttrManager(typ)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, map[string]interface{}{"list": res})
}

//修改属性是否可见
func (this MetaDataController) UpdateAttrInvisible(ctx *fiber.Ctx) error {

	var reqData request.UpdateAttrInvisibleReq

	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}

	metaData := meta_data.MetaDataService{Appid: strconv.Itoa(reqData.Appid)}

	err := metaData.UpdateAttrInvisible(reqData)

	if err != nil {
		return this.Error(ctx, err)
	}

	if err := analysis.ClearCacheByAppid(fmt.Sprintf("%s_%s_%s_*", "GetValues", strconv.Itoa(reqData.Appid), reqData.AttributeSource)); err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

//查看上报属性列表（通过元事件）
func (this MetaDataController) AttrManagerByMeta(ctx *fiber.Ctx) error {

	var reqData request.AttrManagerByMetaReq

	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}

	metaData := meta_data.MetaDataService{Appid: strconv.Itoa(reqData.Appid)}

	res, err := metaData.AttrManagerByMeta(reqData)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, map[string]interface{}{"list": res})

}

//修改属性显示名
func (this MetaDataController) UpdateAttrShowName(ctx *fiber.Ctx) error {

	var reqData request.UpdateAttrShowNameReq

	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}

	metaData := meta_data.MetaDataService{Appid: strconv.Itoa(reqData.Appid)}

	err := metaData.UpdateAttrShowName(reqData)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

//更新事件显示名
func (this MetaDataController) UpdateShowName(ctx *fiber.Ctx) error {
	var reqData request.UpdateShowNameReq

	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}

	metaData := meta_data.MetaDataService{Appid: strconv.Itoa(reqData.Appid)}

	err := metaData.UpdateEventShowName(reqData)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

func (this MetaDataController) GetCalcuSymbolData(ctx *fiber.Ctx) error {

	var reqData request.GetCalcuSymbolDataReq

	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}

	metaData := meta_data.MetaDataService{Appid: strconv.Itoa(reqData.Appid)}

	res, err := metaData.GetCalcuSymbolData(reqData)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, map[string]interface{}{"list": res})
}

//获取分析数据时的下拉选配置
func (this MetaDataController) GetAnalyseSelectOptions(ctx *fiber.Ctx) error {

	var reqData request.GetAnalyseSelectOptionsReq

	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}

	metaData := meta_data.MetaDataService{Appid: strconv.Itoa(reqData.Appid)}

	res, err := metaData.GetAnalyseSelectOptions(reqData.Appid)

	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, res)
}

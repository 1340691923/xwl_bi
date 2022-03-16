package controller

import (
	"errors"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/jwt"
	"github.com/1340691923/xwl_bi/platform-basic-libs/request"
	"github.com/1340691923/xwl_bi/platform-basic-libs/response"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/pannel"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type PanelController struct {
	BaseController
}

//查看自己的已存报表
func (this PanelController) ReportTableList(ctx *fiber.Ctx) error {
	var req model.ReportTable
	if err := ctx.BodyParser(&req); err != nil {
		return this.Error(ctx, err)
	}

	c, _ := jwt.ParseToken(this.GetToken(ctx))

	pannelService := pannel.PannelService{}
	list, err := pannelService.ReportTableList(req.Appid, req.RtType, c.UserID)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, list)
}

//删除自己的已存报表
func (this PanelController) DeleteReportTableByID(ctx *fiber.Ctx) error {
	var reportTable model.ReportTable
	if err := ctx.BodyParser(&reportTable); err != nil {
		return this.Error(ctx, err)
	}

	c, _ := jwt.ParseToken(this.GetToken(ctx))

	pannelService := pannel.PannelService{}
	err := pannelService.DeleteReportTableByID(reportTable, c.UserID)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

//新增/修改自己的已存报表
func (this PanelController) AddReportTable(ctx *fiber.Ctx) error {
	var reportTable model.ReportTable
	if err := ctx.BodyParser(&reportTable); err != nil {
		return this.Error(ctx, err)
	}

	if reportTable.Name == "" {
		return this.Error(ctx, errors.New("报表名不能为空"))
	}

	c, _ := jwt.ParseToken(this.GetToken(ctx))
	reportTable.UserId = int(c.UserID)

	if reportTable.Remark == "" {
		reportTable.Remark = reportTable.Name
	}

	if err := reportTable.InsertOrUpdate(); err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

//检测报表名称是否重名
func (this PanelController) FindNameCount(ctx *fiber.Ctx) error {
	var req request.FindNameCount
	if err := ctx.BodyParser(&req); err != nil {
		return this.Error(ctx, err)
	}
	if req.Name == "" {
		return this.Error(ctx, errors.New("报表名不能为空"))
	}
	c, _ := jwt.ParseToken(this.GetToken(ctx))
	pannelService := pannel.PannelService{}
	count, err := pannelService.FindNameCount(req, c.UserID)
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, count)
}

//通过ID查看报表信息
func (this PanelController) FindRtById(ctx *fiber.Ctx) error {
	var req request.FindRtById
	if err := ctx.BodyParser(&req); err != nil {
		return this.Error(ctx, err)
	}

	c, _ := jwt.ParseToken(this.GetToken(ctx))

	pannelService := pannel.PannelService{}
	res, err := pannelService.FindRtById(req.Id, c.UserID)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, res)
}

//查看面板信息
func (this PanelController) GetPannelList(ctx *fiber.Ctx) error {
	var req request.GetPannelList
	if err := ctx.BodyParser(&req); err != nil {
		return this.Error(ctx, err)
	}

	c, _ := jwt.ParseToken(this.GetToken(ctx))
	pannelService := pannel.PannelService{}
	res, err := pannelService.GetPannelList(c.UserID, req.Appid)
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, res)
}

//新建文件夹
func (this PanelController) NewDir(ctx *fiber.Ctx) error {
	var reqData request.NewDir
	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}

	c, _ := jwt.ParseToken(this.GetToken(ctx))

	reqData.CreateBy = int(c.UserID)

	pannelService := pannel.PannelService{}
	err := pannelService.AddDir(reqData)

	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, nil)
}

//新建面板
func (this PanelController) NewPannel(ctx *fiber.Ctx) error {
	var reqData request.NewPannel
	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}

	c, _ := jwt.ParseToken(this.GetToken(ctx))

	if reqData.PannelName == "" {
		return this.Error(ctx, errors.New("面板名称不能为空"))
	}

	if reqData.FolderId == 0 {
		return this.Error(ctx, errors.New("文件夹不能为空"))
	}

	pannelService := pannel.PannelService{}
	err := pannelService.AddPannel(reqData, c.UserID)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

//修改面板名称
func (this PanelController) Rename(ctx *fiber.Ctx) error {
	var reqData model.Pannel
	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}

	c, _ := jwt.ParseToken(this.GetToken(ctx))

	if reqData.PannelName == "" {
		return this.Error(ctx, errors.New("面板名称不能为空"))
	}

	pannelService := pannel.PannelService{}
	err := pannelService.Rename(reqData.PannelName, reqData.Id, c.UserID)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

//迁移面板到指定文件夹
func (this PanelController) MovePannel2Dir(ctx *fiber.Ctx) error {
	var reqData model.Pannel
	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}

	c, _ := jwt.ParseToken(this.GetToken(ctx))

	pannelService := pannel.PannelService{}
	err := pannelService.MovePannel2Dir(reqData.Id, reqData.FolderId, c.UserID)

	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, nil)
}

//删除面板
func (this PanelController) DeletePannel(ctx *fiber.Ctx) error {
	var reqData model.Pannel
	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}

	c, _ := jwt.ParseToken(this.GetToken(ctx))

	pannelService := pannel.PannelService{}
	err := pannelService.DeletePannel(reqData.Id, c.UserID)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

//删除文件夹
func (this PanelController) DeleteDir(ctx *fiber.Ctx) error {
	var reqData model.PannelFolder
	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}

	c, _ := jwt.ParseToken(this.GetToken(ctx))

	reqData.CreateBy = int(c.UserID)

	pannelService := pannel.PannelService{}
	err := pannelService.DeleteDir(reqData)

	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, nil)
}

//复制面板
func (this PanelController) CopyPannel(ctx *fiber.Ctx) error {
	var reqData model.Pannel
	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}
	reportTablesArr := strings.Split(reqData.ReportTables, ",")

	if len(reportTablesArr) == 0 {
		return this.Error(ctx, errors.New("该面板的报表数量为0"))
	}

	c, _ := jwt.ParseToken(this.GetToken(ctx))

	reqData.CreateBy = int(c.UserID)

	pannelService := pannel.PannelService{}
	err := pannelService.CopyPannel(reqData, c.UserID)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

//修改面板的报表排序
func (this PanelController) UpdatePannelRt(ctx *fiber.Ctx) error {
	var reqData model.Pannel
	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}

	c, _ := jwt.ParseToken(this.GetToken(ctx))

	reqData.CreateBy = int(c.UserID)

	pannelService := pannel.PannelService{}

	err := pannelService.UpdatePannelRt(reqData.ReportTables, reqData.Id, c.UserID)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.ChangeLayoutSuccess, nil)
}

//分享面板给其他成员
func (this PanelController) UpdatePannelManager(ctx *fiber.Ctx) error {
	var reqData model.Pannel
	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}

	c, _ := jwt.ParseToken(this.GetToken(ctx))

	pannelService := pannel.PannelService{}

	err := pannelService.UpdatePannelManager(reqData.Managers, reqData.Id, c.UserID)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

//获取面板的报表相关属性
func (this PanelController) RtListByAppid(ctx *fiber.Ctx) error {
	var reqData model.ReportTable
	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}

	pannelService := pannel.PannelService{}

	res, err := pannelService.RtListByAppid(reqData.Appid)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, res)
}

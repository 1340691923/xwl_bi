package router

import (
	. "github.com/1340691923/xwl_bi/controller"
	"github.com/1340691923/xwl_bi/middleware"
	"github.com/1340691923/xwl_bi/platform-basic-libs/api_config"
	"github.com/gofiber/fiber/v2"
)

func runPannel(app *fiber.App) {
	c := api_config.NewApiRouterConfig()
	const AbsolutePath = "/api/pannel"
	appG := app.Group(AbsolutePath).Use(middleware.FilterAppid)
	{
		c.MountApi(api_config.MountApiBasePramas{Remark: "检测报表名称是否重名", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), PanelController{}.FindNameCount)

		c.MountApi(api_config.MountApiBasePramas{Remark: "获取面板的报表相关属性", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), PanelController{}.RtListByAppid)

		appG = appG.Use(middleware.OperaterLog)

		c.MountApi(api_config.MountApiBasePramas{Remark: "查看自己的已存报表", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), PanelController{}.ReportTableList)

		c.MountApi(api_config.MountApiBasePramas{Remark: "删除自己的已存报表", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), PanelController{}.DeleteReportTableByID)

		c.MountApi(api_config.MountApiBasePramas{Remark: "新增/修改自己的已存报表", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), PanelController{}.AddReportTable)

		c.MountApi(api_config.MountApiBasePramas{Remark: "通过ID查看报表信息", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), PanelController{}.FindRtById)

		c.MountApi(api_config.MountApiBasePramas{Remark: "查看面板信息", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), PanelController{}.GetPannelList)

		c.MountApi(api_config.MountApiBasePramas{Remark: "新建文件夹", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), PanelController{}.NewDir)

		c.MountApi(api_config.MountApiBasePramas{Remark: "新建面板", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), PanelController{}.NewPannel)

		c.MountApi(api_config.MountApiBasePramas{Remark: "修改面板名称", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), PanelController{}.Rename)

		c.MountApi(api_config.MountApiBasePramas{Remark: "迁移面板到指定文件夹", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), PanelController{}.MovePannel2Dir)

		c.MountApi(api_config.MountApiBasePramas{Remark: "删除面板", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), PanelController{}.DeletePannel)

		c.MountApi(api_config.MountApiBasePramas{Remark: "删除文件夹", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), PanelController{}.DeleteDir)

		c.MountApi(api_config.MountApiBasePramas{Remark: "复制面板", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), PanelController{}.CopyPannel)

		c.MountApi(api_config.MountApiBasePramas{Remark: "修改面板的报表排序", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), PanelController{}.UpdatePannelRt)

		c.MountApi(api_config.MountApiBasePramas{Remark: "分享面板给其他成员", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), PanelController{}.UpdatePannelManager)

	}
}

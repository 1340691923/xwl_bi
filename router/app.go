package router

import (
	. "github.com/1340691923/xwl_bi/controller"
	"github.com/1340691923/xwl_bi/middleware"
	"github.com/1340691923/xwl_bi/platform-basic-libs/api_config"
	"github.com/gofiber/fiber/v2"
)

func runApp(app *fiber.App) {
	apiRouterConfig := api_config.NewApiRouterConfig()
	const AbsolutePath = "/api/app"
	appG := app.Group(AbsolutePath)
	{
		apiRouterConfig.MountApi(api_config.MountApiBasePramas{Remark: "查询应用列表", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), AppController{}.List)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{Remark: "获取应用下拉选", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), AppController{}.Config)

		appG = appG.Use(middleware.OperaterLog)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{Remark: "创建应用", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), AppController{}.Create)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{Remark: "重置秘钥", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), AppController{}.ResetAppkey)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{Remark: "修改应用成员", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), AppController{}.UpdateManager)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{Remark: "修改应用状态", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), AppController{}.StatusAction)
	}
}

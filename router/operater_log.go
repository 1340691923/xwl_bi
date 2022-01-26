package router

import (
	. "github.com/1340691923/xwl_bi/controller"
	"github.com/1340691923/xwl_bi/platform-basic-libs/api_config"
	"github.com/gofiber/fiber/v2"
)

func runOperaterLog(app *fiber.App) {
	apiRouterConfig := api_config.NewApiRouterConfig()
	const AbsolutePath = "/api/operater_log"
	appG := app.Group(AbsolutePath)
	{
		apiRouterConfig.MountApi(api_config.MountApiBasePramas{Remark: "查看后台操作日志", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), GmOperaterController{}.ListAction)
	}
}

package router

import (
	. "github.com/1340691923/xwl_bi/controller"
	"github.com/1340691923/xwl_bi/middleware"
	"github.com/1340691923/xwl_bi/platform-basic-libs/api_config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
)

func runRealData(app *fiber.App) {
	apiRouterConfig := api_config.NewApiRouterConfig()
	const AbsolutePath = "/api/realdata"
	appG := app.Group(AbsolutePath).Use(middleware.FilterAppid)
	{
		appG := appG.Use(middleware.OperaterLog)

		apiRouterConfig.MountApi(
			api_config.MountApiBasePramas{
				Remark:       "实时数据列表",
				AbsolutePath: AbsolutePath,
			}, appG.(*fiber.Group).Use(limiter.New(limiter.Config{
				Max:        60,
				Expiration: 2 * time.Second,
			})).(*fiber.Group), RealDataController{}.List)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{Remark: "查看错误数据列表", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), RealDataController{}.FailDataList)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{Remark: "查看错误数据(抽样示例)", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), RealDataController{}.FailDataDesc)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{Remark: "查看上报统计列表", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), RealDataController{}.ReportCount)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{Remark: "事件错误信息查看", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), RealDataController{}.EventFailDesc)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{Remark: "添加测试设备", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), RealDataController{}.AddDebugDeviceID)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{Remark: "查看测试设备列表", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), RealDataController{}.DebugDeviceIDList)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{Remark: "删除测试设备", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), RealDataController{}.DelDebugDeviceID)

	}

}

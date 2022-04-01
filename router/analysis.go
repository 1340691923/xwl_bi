package router

import (
	. "github.com/1340691923/xwl_bi/controller"
	"github.com/1340691923/xwl_bi/middleware"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/api_config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
)

func runAnalysis(app *fiber.App) {
	c := api_config.NewApiRouterConfig()
	const AbsolutePath = "/api/analysis"
	appG := app.Group(AbsolutePath).Use(middleware.FilterAppid)
	{
		c.MountApi(api_config.MountApiBasePramas{Remark: "获取分析面板初始化数据", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), BehaviorAnalysisController{}.GetConfigs)
		c.MountApi(api_config.MountApiBasePramas{Remark: "根据事件名查找指标的运算函数", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), BehaviorAnalysisController{}.LoadPropQuotas)
		c.MountApi(api_config.MountApiBasePramas{Remark: "获取上报字段所有的值", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), BehaviorAnalysisController{}.GetValues)
		appG = appG.Use(middleware.OperaterLog)

		appG = appG.
			Use(limiter.New(
				limiter.Config{
					Max:        model.GlobConfig.GetCkQueryLimit(),
					Expiration: time.Duration(model.GlobConfig.GetCkQueryExpiration()) * time.Second,
				}))

		c.MountApi(api_config.MountApiBasePramas{Remark: "事件分析查询", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), BehaviorAnalysisController{}.EventList)
		c.MountApi(api_config.MountApiBasePramas{Remark: "漏斗分析查询", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), BehaviorAnalysisController{}.FunnelList)
		c.MountApi(api_config.MountApiBasePramas{Remark: "留存分析查询", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), BehaviorAnalysisController{}.RetentionList)

		c.MountApi(api_config.MountApiBasePramas{Remark: "用户属性分析查询", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), BehaviorAnalysisController{}.UserAttrList)
		c.MountApi(api_config.MountApiBasePramas{Remark: "用户列表查询", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), BehaviorAnalysisController{}.UserList)
		c.MountApi(api_config.MountApiBasePramas{Remark: "查询用户访问过的事件详情", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), BehaviorAnalysisController{}.UserEventDetailList)
		c.MountApi(api_config.MountApiBasePramas{Remark: "查询用户访问过的事件统计", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), BehaviorAnalysisController{}.UserEventCountList)

		/*f := appG.
			Use(limiter.New(
				limiter.Config{
					Max:        2,
					Expiration: 2 * time.Second,
				}))*/

		c.MountApi(api_config.MountApiBasePramas{Remark: "智能路径分析查询", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), BehaviorAnalysisController{}.TraceList)
	}
}

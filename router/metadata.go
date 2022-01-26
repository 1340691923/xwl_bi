package router

import (
	. "github.com/1340691923/xwl_bi/controller"
	"github.com/1340691923/xwl_bi/middleware"
	"github.com/1340691923/xwl_bi/platform-basic-libs/api_config"
	"github.com/gofiber/fiber/v2"
)

func runMetaData(app *fiber.App) {
	apiRouterConfig := api_config.NewApiRouterConfig()
	const AbsolutePath = "/api/metadata"
	appG := app.Group(AbsolutePath).Use(middleware.FilterAppid)
	{
		apiRouterConfig.MountApi(api_config.MountApiBasePramas{Remark: "获取分析数据时的下拉选配置", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), MetaDataController{}.GetAnalyseSelectOptions)

		appG = appG.Use(middleware.OperaterLog)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{Remark: "元事件列表（通过属性查找）", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), MetaDataController{}.MetaEventListByAttr)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{Remark: "元事件列表", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), MetaDataController{}.MetaEventList)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{Remark: "修改属性是否可见", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), MetaDataController{}.UpdateAttrInvisible)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{Remark: "更新事件显示名", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), MetaDataController{}.UpdateShowName)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{Remark: "查看上报属性列表", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), MetaDataController{}.AttrManager)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{Remark: "查看上报属性列表（通过元事件）", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), MetaDataController{}.AttrManagerByMeta)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{Remark: "修改属性显示名", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), MetaDataController{}.UpdateAttrShowName)

	}

}

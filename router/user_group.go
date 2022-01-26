package router

import (
	. "github.com/1340691923/xwl_bi/controller"
	"github.com/1340691923/xwl_bi/middleware"
	"github.com/1340691923/xwl_bi/platform-basic-libs/api_config"
	"github.com/gofiber/fiber/v2"
)

func runUserGroup(app *fiber.App) {
	c := api_config.NewApiRouterConfig()
	const AbsolutePath = "/api/user_group"
	appG := app.Group(AbsolutePath).Use(middleware.FilterAppid)
	{

		appG = appG.Use(middleware.OperaterLog)

		c.MountApi(api_config.MountApiBasePramas{Remark: "新增用户分群", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), UserGroupController{}.AddUserGroup)

		c.MountApi(api_config.MountApiBasePramas{Remark: "修改用户分群", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), UserGroupController{}.ModifyUserGroup)

		c.MountApi(api_config.MountApiBasePramas{Remark: "删除用户分群", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), UserGroupController{}.DeleteUserGroup)

		c.MountApi(api_config.MountApiBasePramas{Remark: "用户分群列表", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), UserGroupController{}.UserGroupList)

		c.MountApi(api_config.MountApiBasePramas{Remark: "用户分群下拉选", AbsolutePath: AbsolutePath}, appG.(*fiber.Group), UserGroupController{}.UserGroupSelect)
	}
}

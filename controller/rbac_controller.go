package controller

import (
	"github.com/1340691923/xwl_bi/platform-basic-libs/api_config"
	"github.com/1340691923/xwl_bi/platform-basic-libs/response"
	fiber "github.com/gofiber/fiber/v2"
)

//接口访问权限管理	直接放缓存
type RbacController struct {
	BaseController
}

//获取接口路由信息
func (this RbacController) UrlConfig(ctx *fiber.Ctx) error {
	apiRouterConfig := api_config.NewApiRouterConfig()
	return this.Success(ctx, response.SearchSuccess, apiRouterConfig.GetRouterConfigs())
}

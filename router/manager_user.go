package router

import (
	. "github.com/1340691923/xwl_bi/controller"
	"github.com/1340691923/xwl_bi/middleware"
	api_config "github.com/1340691923/xwl_bi/platform-basic-libs/api_config"

	. "github.com/gofiber/fiber/v2"
)

// BI用户 路由
func runGmUser(app *App) {
	apiRouterConfig := api_config.NewApiRouterConfig()
	const AbsolutePath = "/api/gm_user"
	gmUser := app.Group(AbsolutePath)
	{
		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "查询用户详细信息",
			AbsolutePath: AbsolutePath,
			RelativePath: "info",
		}, gmUser.(*Group), ManagerUserController{}.UserInfo)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "GM角色列表",
			AbsolutePath: AbsolutePath,
			RelativePath: "roles",
		}, gmUser.(*Group), ManagerRoleController{}.RolesAction)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "获取接口路由信息",
			AbsolutePath: AbsolutePath,
			RelativePath: "UrlConfig",
		}, gmUser.(*Group), RbacController{}.UrlConfig)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "退出登录",
			AbsolutePath: AbsolutePath,
			RelativePath: "logout",
		}, gmUser.(*Group), ManagerUserController{}.LogoutAction)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "查询用户列表",
			AbsolutePath: AbsolutePath,
			RelativePath: "userlist",
		}, gmUser.(*Group), ManagerUserController{}.UserListAction)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "角色下拉选",
			AbsolutePath: AbsolutePath,
			RelativePath: "roleOption",
		}, gmUser.(*Group), ManagerRoleController{}.RoleOptionAction)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "通过ID获取用户信息",
			AbsolutePath: AbsolutePath,
			RelativePath: "getUserById",
		}, gmUser.(*Group), ManagerUserController{}.GetUserByIdAction)

		gmUser = gmUser.Use(middleware.OperaterLog)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "修改自己的密码",
			AbsolutePath: AbsolutePath,
			RelativePath: "ModifyPwd",
		}, gmUser.(*Group), ManagerUserController{}.ModifyPwd)

		gmUser = gmUser.Use(middleware.OperaterLog)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "修改GM角色",
			AbsolutePath: AbsolutePath,
			RelativePath: "role/update",
		}, gmUser.(*Group), ManagerRoleController{}.RolesUpdateAction)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "新增GM角色",
			AbsolutePath: AbsolutePath,
			RelativePath: "role/add",
		}, gmUser.(*Group), ManagerRoleController{}.RolesAddAction)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "删除GM角色",
			AbsolutePath: AbsolutePath,
			RelativePath: "role/delete",
		}, gmUser.(*Group), ManagerRoleController{}.RolesDelAction)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "修改用户信息",
			AbsolutePath: AbsolutePath,
			RelativePath: "UpdateUser",
		}, gmUser.(*Group), ManagerUserController{}.UserUpdateAction)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "新增用户信息",
			AbsolutePath: AbsolutePath,
			RelativePath: "InsertUser",
		}, gmUser.(*Group), ManagerUserController{}.UserAddAction)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "删除用户",
			AbsolutePath: AbsolutePath,
			RelativePath: "DelUser",
		}, gmUser.(*Group), ManagerUserController{}.DeleteUserAction)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "解封/封禁用户账号",
			AbsolutePath: AbsolutePath,
			RelativePath: "ban",
		}, gmUser.(*Group), ManagerUserController{}.UserBanAction)
	}
}

package controller

import (
	"errors"
	"github.com/1340691923/xwl_bi/platform-basic-libs/jwt"
	"strconv"

	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/rbac"
	"github.com/1340691923/xwl_bi/platform-basic-libs/request"
	"github.com/1340691923/xwl_bi/platform-basic-libs/response"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/gm_role"
	. "github.com/gofiber/fiber/v2"
)

// GM角色控制器
type ManagerRoleController struct {
	BaseController
}

//获取所有的GM 角色
func (this ManagerRoleController) RolesAction(ctx *Ctx) error {
	var service gm_role.GmRoleService
	roles, err := service.Select()
	if err != nil {
		return this.Error(ctx, err)
	}
	list, err := service.GetRoles(roles)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, list)
}

//新增GM角色
func (this ManagerRoleController) RolesAddAction(ctx *Ctx) error {

	var model2 request.GmRoleModel

	err := ctx.BodyParser(&model2)
	if err != nil {
		return this.Error(ctx, err)
	}
	var roleModel model.GmRoleModel
	roleModel.ID = model2.ID
	roleModel.RoleName = model2.RoleName
	roleModel.RoleList = model2.RoleList
	roleModel.Description = model2.Description
	roleModel.ID = model2.ID
	id, err := roleModel.Insert()

	if err != nil {
		return this.Error(ctx, err)
	}

	go func() {
		for _, api := range model2.Api {
			_, err = rbac.Enforcer.AddPolicySafe(strconv.Itoa(int(id)), api, "*")
			if err != nil {
				logs.Logger.Sugar().Errorf("err:%s", err.Error())
				return
			}
		}
	}()

	return this.Success(ctx, response.OperateSuccess, map[string]interface{}{"id": id})
}

// 修改GM角色
func (this ManagerRoleController) RolesUpdateAction(ctx *Ctx) error {
	var model2 request.GmRoleModel
	err := ctx.BodyParser(&model2)
	if err != nil {
		return this.Error(ctx, err)
	}

	token := this.GetToken(ctx)
	claims, err := jwt.ParseToken(token)
	if err != nil {
		return this.Error(ctx, err)
	}

	if model2.ID == 1 && claims.RoleId != 1 {
		return this.Error(ctx, errors.New("您无权修改该角色!"))
	}

	var roleModel model.GmRoleModel
	roleModel.ID = model2.ID
	roleModel.RoleName = model2.RoleName
	roleModel.RoleList = model2.RoleList
	roleModel.Description = model2.Description
	err = roleModel.Update()

	rbac.Enforcer.RemoveFilteredPolicy(0, strconv.Itoa(model2.ID)) //先全清掉

	go func() {
		for _, api := range model2.Api {
			_, err = rbac.Enforcer.AddPolicySafe(strconv.Itoa(model2.ID), api, "*")
			if err != nil {
				logs.Logger.Sugar().Errorf("err:%s", err.Error())
				return
			}
		}
	}()
	return this.Success(ctx, response.OperateSuccess, nil)
}

// 删除GM角色
func (this ManagerRoleController) RolesDelAction(ctx *Ctx) error {

	var reqData request.RolesDelReq

	err := ctx.BodyParser(&reqData)
	if err != nil {
		return this.Error(ctx, err)
	}

	id := reqData.Id

	claims, err := jwt.ParseToken(this.GetToken(ctx))
	if err != nil {
		return this.Error(ctx, err)
	}

	if id == 1 && claims.RoleId != 1 {
		return this.Error(ctx, errors.New("您无权修改该角色!"))
	}

	var service gm_role.GmRoleService
	err = service.Delete(id)
	if err != nil {
		return this.Error(ctx, err)
	}
	rbac.Enforcer.RemoveFilteredPolicy(0, strconv.Itoa(id)) //先全清掉

	return this.Success(ctx, response.OperateSuccess, nil)
}

// 获取Gm角色下拉选
func (this ManagerRoleController) RoleOptionAction(ctx *Ctx) error {

	var model model.GmRoleModel

	list, err := model.Select()
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, list)
}

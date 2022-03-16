package controller

import (
	"errors"
	"github.com/1340691923/xwl_bi/platform-basic-libs/jwt"
	"github.com/1340691923/xwl_bi/platform-basic-libs/request"
	"github.com/1340691923/xwl_bi/platform-basic-libs/response"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/user_group"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type UserGroupController struct {
	BaseController
}

//新增用户分群
func (this UserGroupController) AddUserGroup(ctx *fiber.Ctx) error {
	var addUserGroup request.AddUserGroup
	if err := ctx.BodyParser(&addUserGroup); err != nil {
		return this.Error(ctx, err)
	}

	if strings.TrimSpace(addUserGroup.Name) == "" {
		return this.Error(ctx, errors.New("用户分群名不能为空"))
	}
	if len(addUserGroup.Ids) == 0 {
		return this.Error(ctx, errors.New("需分群用户ID列表不能为空"))
	}
	c, _ := jwt.ParseToken(this.GetToken(ctx))

	userGroupService := user_group.UserGroupService{
		ManagerID: c.UserID,
		Appid:     addUserGroup.Appid,
	}

	err := userGroupService.AddUserGroup(len(addUserGroup.Ids), addUserGroup.Ids, addUserGroup.Remark, addUserGroup.Name)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

//修改用户分群
func (this UserGroupController) ModifyUserGroup(ctx *fiber.Ctx) error {
	var modifyUserGroup request.ModifyUserGroup
	if err := ctx.BodyParser(&modifyUserGroup); err != nil {
		return this.Error(ctx, err)
	}

	if strings.TrimSpace(modifyUserGroup.Name) == "" {
		return this.Error(ctx, errors.New("用户分群名不能为空"))
	}

	if modifyUserGroup.Id == 0 {
		return this.Error(ctx, errors.New("用户分群ID不能为空"))
	}

	c, _ := jwt.ParseToken(this.GetToken(ctx))

	userGroupService := user_group.UserGroupService{
		ManagerID: c.UserID,
		Appid:     modifyUserGroup.Appid,
	}

	err := userGroupService.ModifyUserGroup(modifyUserGroup.Id, modifyUserGroup.Remark, modifyUserGroup.Name)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

//删除用户分群
func (this UserGroupController) DeleteUserGroup(ctx *fiber.Ctx) error {
	var deleteUserGroup request.DeleteUserGroup
	if err := ctx.BodyParser(&deleteUserGroup); err != nil {
		return this.Error(ctx, err)
	}

	if deleteUserGroup.Id == 0 {
		return this.Error(ctx, errors.New("用户分群ID不能为空"))
	}

	c, _ := jwt.ParseToken(this.GetToken(ctx))

	userGroupService := user_group.UserGroupService{
		ManagerID: c.UserID,
		Appid:     deleteUserGroup.Appid,
	}

	err := userGroupService.DeleteUserGroup(deleteUserGroup.Id)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

//用户分群列表
func (this UserGroupController) UserGroupList(ctx *fiber.Ctx) error {
	var userGroupList request.UserGroupList
	if err := ctx.BodyParser(&userGroupList); err != nil {
		return this.Error(ctx, err)
	}

	c, _ := jwt.ParseToken(this.GetToken(ctx))

	userGroupService := user_group.UserGroupService{
		ManagerID: c.UserID,
		Appid:     userGroupList.Appid,
	}

	list, err := userGroupService.UserGroupList()

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, list)
}

//用户分群下拉选
func (this UserGroupController) UserGroupSelect(ctx *fiber.Ctx) error {
	var userGroupList request.UserGroupList
	if err := ctx.BodyParser(&userGroupList); err != nil {
		return this.Error(ctx, err)
	}

	c, _ := jwt.ParseToken(this.GetToken(ctx))

	userGroupService := user_group.UserGroupService{
		ManagerID: c.UserID,
		Appid:     userGroupList.Appid,
	}

	list, err := userGroupService.Options()

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, list)
}

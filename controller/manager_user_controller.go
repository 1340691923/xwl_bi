package controller

import (
	"errors"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/jwt"
	"github.com/1340691923/xwl_bi/platform-basic-libs/request"
	"github.com/1340691923/xwl_bi/platform-basic-libs/response"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/gm_user"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	. "github.com/gofiber/fiber/v2"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

//BI用户控制器
type ManagerUserController struct {
	BaseController
}

// 登录
func (this ManagerUserController) Login(ctx *Ctx) error {

	type ReqData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var reqData ReqData

	err := ctx.BodyParser(&reqData)

	if err != nil {
		return this.Error(ctx, err)
	}

	username := reqData.Username
	password := reqData.Password
	var gmUserService gm_user.GmUserService
	token, err := gmUserService.CheckLogin(username, password)
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, "登陆成功", map[string]interface{}{"token": token})
}

//修改自己的密码
func (this ManagerUserController) ModifyPwd(ctx *Ctx) error {
	type ReqData struct {
		Password string `json:"password"`
	}

	var reqData ReqData

	err := ctx.BodyParser(&reqData)

	if err != nil {
		return this.Error(ctx, err)
	}

	if err := this.CheckParameter([]request.CheckConfigStruct{
		{
			request.EmptyParmasError,
			"password",
		},
	}, ctx); err != nil {
		return this.Error(ctx, err)
	}

	token := this.GetToken(ctx)
	claims, err := jwt.ParseToken(token)
	if err != nil {
		return this.Error(ctx, err)
	}
	gmUserModel := model.GmUserModel{}
	gmUserModel.ID = claims.UserID
	gmUserModel.Password = reqData.Password
	err = gmUserModel.UpdatePassById()
	if err != nil {
		return this.Error(ctx, err)
	}

	util.TokenBucket.LoadOrStore(token, claims.ExpiresAt)

	return this.Success(ctx, response.OperateSuccess, nil)
}

// 用户详细信息
func (this ManagerUserController) UserInfo(ctx *Ctx) error {
	var gmUserService gm_user.GmUserService
	token := this.GetToken(ctx)
	claims, err := jwt.ParseToken(token)
	if err != nil {
		return this.Error(ctx, err)
	}
	info, err := gmUserService.GetRoleInfo(claims.RoleId)
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, "登陆成功", map[string]interface{}{"roles": []string{"admin"}, "introduction": info.Description, "name": claims.RealName, "list": info.RoleList, "avatar": ""})
}

//退出登录
func (this ManagerUserController) LogoutAction(ctx *Ctx) error {
	token := this.GetToken(ctx)
	var claims *jwt.Claims
	claims, err := jwt.ParseToken(token)
	if err != nil {
		logs.Logger.Error("LogoutAction err", zap.Error(err))
		return this.Success(ctx, response.LogoutSuccess, nil)
	}
	util.TokenBucket.LoadOrStore(token, claims.ExpiresAt)

	return this.Success(ctx, response.LogoutSuccess, nil)
}

//BI用户列表
func (this ManagerUserController) UserListAction(ctx *Ctx) error {

	appid := gjson.GetBytes(ctx.Body(), "appid").String()

	var userModel model.GmUserModel
	list, err := userModel.Select(appid)
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, list)
}

//删除BI用户
func (this ManagerUserController) DeleteUserAction(ctx *Ctx) error {

	var reqData request.DeleteUserReq

	err := ctx.BodyParser(&reqData)

	if err != nil {
		return this.Error(ctx, err)
	}

	var userModel model.GmUserModel
	userModel.ID = reqData.Id
	if userModel.ID == 1 {
		return this.Error(ctx, errors.New("您无权删除该用户!"))
	}

	err = userModel.Delete()
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.DeleteSuccess, nil)
}

//用ID获取用户信息
func (this ManagerUserController) GetUserByIdAction(ctx *Ctx) error {
	var reqData request.GetUserByIdReq

	err := ctx.BodyParser(&reqData)

	if err != nil {
		return this.Error(ctx, err)
	}

	var userModel model.GmUserModel
	userModel.ID = reqData.Id
	gmUser, err := userModel.GetUserById()
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, gmUser)
}

// 修改BI用户信息
func (this ManagerUserController) UserUpdateAction(ctx *Ctx) error {

	var reqData request.UserUpdateReq

	err := ctx.BodyParser(&reqData)

	if err != nil {
		return this.Error(ctx, err)
	}

	var userModel model.GmUserModel
	var id = reqData.Id

	userModel.ID = int32(id)
	userModel.Realname = reqData.Realname
	userModel.RoleId = reqData.RoleId
	userModel.Password = reqData.Password
	userModel.Username = reqData.Username

	err = userModel.Update()
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, nil)
}

//新增BI用户
func (this ManagerUserController) UserAddAction(ctx *Ctx) error {

	var reqData request.UserAddReq

	err := ctx.BodyParser(&reqData)

	if err != nil {
		return this.Error(ctx, err)
	}

	var userModel model.GmUserModel

	userModel.Realname = reqData.Realname
	userModel.RoleId = reqData.RoleId
	userModel.Password = reqData.Password
	userModel.Username = reqData.Username

	id, err := userModel.Insert()
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, id)
}

//禁用或者解封账号
func (this ManagerUserController) UserBanAction(ctx *Ctx) error {

	var reqData request.UserBanReq

	err := ctx.BodyParser(&reqData)

	if err != nil {
		return this.Error(ctx, err)
	}

	id := reqData.Id
	ban_typ := reqData.Typ

	if id == 1 || id == 0 {
		return this.Error(ctx, errors.New("您无权操作该用户!"))
	}

	_, err = db.
		SqlBuilder.
		Update("gm_user").
		SetMap(map[string]interface{}{"is_del": ban_typ}).
		Where(db.Eq{"id": id}).
		RunWith(db.Sqlx).
		Exec()

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, id)
}

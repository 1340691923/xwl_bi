package controller

import (
	"errors"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/jwt"
	"github.com/1340691923/xwl_bi/platform-basic-libs/response"
	app2 "github.com/1340691923/xwl_bi/platform-basic-libs/service/app"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/myapp"
	"github.com/gofiber/fiber/v2"
)

type AppController struct {
	BaseController
}

//创建应用
func (this AppController) Create(ctx *fiber.Ctx) error {

	var app model.App
	err := ctx.BodyParser(&app)
	if err != nil {
		return this.Error(ctx, err)
	}

	if app.AppName == "" {
		return this.Error(ctx, errors.New("应用名不能为空"))
	}

	c, _ := jwt.ParseToken(this.GetToken(ctx))

	appService := app2.AppService{}

	err = appService.Create(app, c.UserID)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

//修改应用管理员
func (this AppController) UpdateManager(ctx *fiber.Ctx) error {
	var app model.App
	err := ctx.BodyParser(&app)
	if err != nil {
		return this.Error(ctx, err)
	}

	if app.AppId == "" {
		return this.Error(ctx, errors.New("应用ID不能为空"))
	}

	c, _ := jwt.ParseToken(this.GetToken(ctx))

	appService := app2.AppService{}

	err = appService.UpdateManager(app, c.UserID)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

//得到该用户的应用下拉选
func (this AppController) Config(ctx *fiber.Ctx) error {

	list, err := myapp.GetAppidsByToken(this.GetToken(ctx))

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, map[string]interface{}{
		"list": list,
	})
}

//重置应用的appkey
func (this AppController) ResetAppkey(ctx *fiber.Ctx) error {
	var app model.App
	err := ctx.BodyParser(&app)
	if err != nil {
		return this.Error(ctx, err)
	}

	c, _ := jwt.ParseToken(this.GetToken(ctx))

	appService := app2.AppService{}

	err = appService.ResetAppkey(c.UserID, app)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

//修改应用的状态
func (this AppController) StatusAction(ctx *fiber.Ctx) error {
	var app model.App
	err := ctx.BodyParser(&app)
	if err != nil {
		return this.Error(ctx, err)
	}
	c, _ := jwt.ParseToken(this.GetToken(ctx))

	appService := app2.AppService{}

	err = appService.ChangeStatus(app, c.UserID)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

func (this AppController) List(ctx *fiber.Ctx) error {
	var app model.App
	err := ctx.BodyParser(&app)
	if err != nil {
		return this.Error(ctx, err)
	}

	c, _ := jwt.ParseToken(this.GetToken(ctx))

	appService := app2.AppService{}

	list, count, err := appService.List(c.UserID, app)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, map[string]interface{}{
		"list":  list,
		"count": count,
	})
}

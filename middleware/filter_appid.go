package middleware

import (
	"errors"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/myapp"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"github.com/gofiber/fiber/v2"
	"github.com/tidwall/gjson"
)

func FilterAppid(ctx *fiber.Ctx) error {

	appid := gjson.GetBytes(ctx.Body(), "appid").Int()

	if appid == 0 {
		return res.Error(ctx, errors.New("请先在左上角选择您的应用"))
	}

	list, err := myapp.GetAppidsByToken(util.GetToken(ctx))

	if err != nil {
		return res.Error(ctx, err)
	}
	haveApp := false
	for _, v := range list {
		if v.Id == int(appid) {
			haveApp = true
			break
		}
	}

	if !haveApp {
		return res.Error(ctx, errors.New("您不属于该应用成员"))
	}

	err = ctx.Next()

	return err

}

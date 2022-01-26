// 中间件层
package middleware

import (
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/jwt"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

//fiber 没有将fasthttp的
func OperaterLog(ctx *fiber.Ctx) error {

	var err error
	token := util.GetToken(ctx)
	var claims *jwt.Claims
	claims, err = jwt.ParseToken(token)
	if err != nil {
		logs.Logger.Error("OperaterLog jwt err", zap.Error(err))
		return err
	}

	gmOperaterLog := model.GmOperaterLog{
		OperaterName:   claims.Username,
		OperaterId:     int(claims.UserID),
		OperaterAction: ctx.Path(),
		Method:         ctx.Method(),
		Body:           ctx.Body(),
		OperaterRoleId: int(claims.RoleId),
	}

	err = gmOperaterLog.Insert()

	if err != nil {
		logs.Logger.Error("OperaterLog", zap.Error(err))
	}

	return ctx.Next()

}

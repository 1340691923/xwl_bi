package controller

import (
	"github.com/1340691923/xwl_bi/platform-basic-libs/request"
	"github.com/1340691923/xwl_bi/platform-basic-libs/response"
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/gm_operater_log"
	"github.com/gofiber/fiber/v2"
)

type GmOperaterController struct {
	BaseController
}

//查看后台操作日志
func (this GmOperaterController) ListAction(ctx *fiber.Ctx) error {

	var reqData request.GmOperaterLogList

	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}

	gmOperaterLogService := gm_operater_log.GmOperaterLogService{}

	list, count, err := gmOperaterLogService.List(reqData)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, map[string]interface{}{"list": list, "count": count})
}

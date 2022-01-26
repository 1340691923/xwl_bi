package middleware

import (
	"github.com/1340691923/xwl_bi/platform-basic-libs/service/analysis"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

//防止狂点设置3秒的缓存
func AnalysisCache(ctx *fiber.Ctx) error {
	cache := analysis.NewCache(time.Second*3, ctx.Path(), ctx.Body())

	resData, redisErr := cache.LoadData()
	if util.FilterRedisNilErr(redisErr) {
		return res.Error(ctx, redisErr)
	}

	if len(resData) > 0 {
		return ctx.Send(resData)
	}

	err := ctx.Next()
	if err != nil {
		return res.Error(ctx, err)
	}
	if ctx.Response().StatusCode() == http.StatusOK {
		err = cache.SetData(ctx.Response().Body())
		if err != nil {
			return res.Error(ctx, err)
		}
	}

	return err

}

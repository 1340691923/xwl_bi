package gm_operater_log

import (
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/request"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"go.uber.org/zap"
)

type GmOperaterLogService struct {
}

func (this *GmOperaterLogService) List(reqData request.GmOperaterLogList) (list []model.GmOperaterLog, count int, err error) {
	if reqData.Page <= 0 {
		reqData.Page = 1
	}
	if reqData.Limit <= 0 {
		reqData.Limit = 10
	}

	page := reqData.Page
	limit := reqData.Limit

	operater_action := reqData.OperaterAction

	gmOperaterModel := &model.GmOperaterLog{
		OperaterRoleId: reqData.RoleId,
		OperaterId:     reqData.UserId,
		OperaterAction: operater_action,
		FilterDate:     reqData.Date,
	}
	listP := &list
	err = model.SearchList(gmOperaterModel, page, limit, "*", listP)
	if err != nil {
		return
	}
	count, err = model.Count(gmOperaterModel)
	if err != nil {
		return
	}

	for index := range list {

		body, err := util.GzipUnCompress(list[index].Body)
		if err != nil {
			logs.Logger.Error("err", zap.Error(err))
			continue
		}

		list[index].BodyStr = body
	}
	return
}

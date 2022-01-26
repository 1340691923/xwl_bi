package analysis

import (
	"errors"
	"fmt"
)

type Ianalysis interface {
	GetList() (interface{}, error)
	GetExecSql() (SQL string, allArgs []interface{}, err error) //后续开发  查看执行sql功能
}

type Command int

const (
	FunnelComand               Command = 1
	RetentionComand            Command = 2
	TraceComand                Command = 3
	EventComand                Command = 4
	UserAttrCommand            Command = 5
	UserListCommand            Command = 6
	UserEventDetailListCommand Command = 7
	UserEventCountCommand      Command = 8
)

var commandMap = map[Command]func(reqData []byte) (Ianalysis, error){
	FunnelComand:               NewFunnel,
	RetentionComand:            NewRetention,
	TraceComand:                NewTrace,
	EventComand:                NewEvent,
	UserAttrCommand:            NewUserAttr,
	UserListCommand:            NewUserList,
	UserEventDetailListCommand: NewUserEventDetailList,
	UserEventCountCommand:      NewUserEventCountList,
}

func NewAnalysisByCommand(command Command, reqData []byte) (i Ianalysis, err error) {
	var fn func(reqData []byte) (Ianalysis, error)
	var found bool
	if fn, found = commandMap[command]; !found {
		return nil, errors.New(fmt.Sprintf("没有找到该命令:%v", command))
	}

	return fn(reqData)
}

func GetAnalysisRes(i Ianalysis) (res interface{}, err error) {
	return i.GetList()
}

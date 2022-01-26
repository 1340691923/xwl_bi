package myapp

import (
	"fmt"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/garyburd/redigo/redis"
	"go.uber.org/zap"
)

const AppidToTableidHash = "AppidToTableid"

func SetAppidToTableid(appid, appkey string, tableID int) (err error) {
	conn := db.RedisPool.Get()
	defer conn.Close()
	_, err = conn.Do("hset", AppidToTableidHash, fmt.Sprintf("%s_xwl_%s", appid, appkey), tableID)
	if err != nil {
		logs.Logger.Error("SetAppidToTableid err", zap.Error(err))
	}
	return
}

func GetAppidToTableid(conn redis.Conn, key string) (tableID string, err error) {
	tableID, err = redis.String(conn.Do("hget", AppidToTableidHash, key))
	return
}

func DeleteAppidToTableid(appid, appkey string) (err error) {
	conn := db.RedisPool.Get()
	defer conn.Close()
	_, err = conn.Do("hdel", AppidToTableidHash, fmt.Sprintf("%s_xwl_%s", appid, appkey))
	if err != nil {
		logs.Logger.Error("DeleteAppidToTableid", zap.Error(err))
	}
	return
}

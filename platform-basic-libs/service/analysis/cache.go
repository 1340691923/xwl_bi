package analysis

import (
	"fmt"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"github.com/garyburd/redigo/redis"
	"go.uber.org/zap"
	"time"
)

type Cache struct {
	overTime     int    `json:"over_time"`
	analysisType string `json:"cache_key"`
	reqData      []byte `json:"req_data"`
}

func ClearCacheByAppid(key string) (err error) {
	conn := db.RedisPool.Get()
	defer conn.Close()
	_, err = conn.Do("unlink", key)
	if err != nil {
		_, err = conn.Do("del", key)
		if err != nil {
			logs.Logger.Error("err", zap.Error(err))
		}
	}
	return
}

func NewCache(overTime time.Duration, analysisType string, reqData []byte) *Cache {
	return &Cache{overTime: int(overTime.Seconds()), analysisType: analysisType, reqData: reqData}
}

func (this *Cache) getKey() string {
	return fmt.Sprintf("%s_%s", this.analysisType, util.MD5HexHash(this.reqData))
}

func (this *Cache) LoadData() (b []byte, err error) {
	conn := db.RedisPool.Get()
	defer conn.Close()
	b, err = redis.Bytes(conn.Do("get", this.getKey()))
	return
}

func (this *Cache) SetData(b []byte) (err error) {
	conn := db.RedisPool.Get()
	defer conn.Close()
	_, err = conn.Do("SETEX", this.getKey(), this.overTime, b)
	return
}

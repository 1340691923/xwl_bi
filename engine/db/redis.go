package db

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"strings"
)

var RedisPool *redis.Pool

// NewRedisPool 新建一个Redis连接池 URL优先
func NewRedisPool(addr, passwd string, db, maxIdle, MaxActive int) *redis.Pool {
	b := strings.HasPrefix(addr, "redis://")
	var dialFunc func() (redis.Conn, error)
	switch {
	case b && passwd == "":
		dialFunc = func() (redis.Conn, error) {
			return redis.DialURL(addr, redis.DialDatabase(db))
		}
	case b && passwd != "":
		dialFunc = func() (redis.Conn, error) {
			return redis.DialURL(addr, redis.DialDatabase(db), redis.DialPassword(passwd))
		}
	case !b && passwd == "":
		dialFunc = func() (redis.Conn, error) {
			return redis.Dial("tcp", addr, redis.DialDatabase(db))
		}
	case !b && passwd != "":
		dialFunc = func() (redis.Conn, error) {
			return redis.Dial("tcp", addr, redis.DialDatabase(db), redis.DialPassword(passwd))
		}
	}

	return &redis.Pool{
		MaxIdle:   maxIdle,
		MaxActive: MaxActive,
		Dial: func() (redis.Conn, error) {
			c, err := dialFunc()
			if err != nil {
				log.Println(fmt.Errorf("redis 连接失败:%v", err))
			}
			return c, err
		},
	}
}

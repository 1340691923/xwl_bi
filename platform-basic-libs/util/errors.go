package util

import (
	"database/sql"
	"strings"

	"github.com/garyburd/redigo/redis"
)

func FilterMysqlNilErr(err error) bool {
	if err != nil && err != sql.ErrNoRows {
		return true
	}
	return false
}

func IsMysqlRepeatError(err error) bool {
	if err != nil && strings.Contains(err.Error(), "Error 1062") {
		return true
	}
	return false
}

func FilterRedisNilErr(err error) bool {
	if err != nil && err != redis.ErrNil {
		return true
	}
	return false
}

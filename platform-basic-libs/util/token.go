package util

import (
	"log"
	"strconv"
	"sync"
	"github.com/sony/sonyflake"
)

var TokenBucket sync.Map

func GetUUid()string{
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		log.Println("err", err)
	}
	return strconv.Itoa(int(id))
}
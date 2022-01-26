package model

import "time"

type InputMessage struct {
	Topic     string
	Partition int
	Key       []byte
	Value     []byte
	Offset    int64
	Timestamp *time.Time
}

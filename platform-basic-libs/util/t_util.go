package util

import (
	"go.uber.org/zap"
	"time"
)

type Tutil struct {
	startT   time.Time
	tag      string
	logFloag bool
	logger   *zap.Logger
}

func NewTutil(tagP string, logFloag bool, logger *zap.Logger) *Tutil {
	return &Tutil{startT: time.Now(), tag: tagP, logFloag: logFloag, logger: logger}
}

func (this *Tutil) EndT(tagC string, haveNext ...bool) {
	if len(haveNext) == 0 {
		this.startT = time.Now()
	} else {
		if this.logFloag {
			this.logger.Sugar().Infof("%s(%s):lost time:%v", this.tag, tagC, time.Now().Sub(this.startT).String())
		}
	}
}

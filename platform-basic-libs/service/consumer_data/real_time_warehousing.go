package consumer_data

import (
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"go.uber.org/zap"
	"sync"
	"time"
)

type RealTimeWarehousingData struct {
	Appid      int64
	EventName  string
	CreateTime string
	Data       []byte
}

type RealTimeWarehousing struct {
	buffer        []*RealTimeWarehousingData
	bufferMutex   *sync.RWMutex
	batchSize     int
	flushInterval int
}

func NewRealTimeWarehousing(config model.BatchConfig) *RealTimeWarehousing {
	logs.Logger.Info("NewRealTimeWarehousing", zap.Int("batchSize", config.BufferSize), zap.Int("flushInterval", config.FlushInterval))
	realTimeWarehousing := &RealTimeWarehousing{
		buffer:        make([]*RealTimeWarehousingData, 0, config.BufferSize),
		bufferMutex:   new(sync.RWMutex),
		batchSize:     config.BufferSize,
		flushInterval: config.FlushInterval,
	}

	if config.FlushInterval > 0 {
		realTimeWarehousing.RegularFlushing()
	}

	return realTimeWarehousing
}

func (this *RealTimeWarehousing) Flush() (err error) {
	this.bufferMutex.Lock()
	if len(this.buffer) > 0 {
		startNow := time.Now()

		tx, err := db.ClickHouseSqlx.Begin()
		if err != nil {
			return err
		}

		stmt, err := tx.Prepare("INSERT INTO xwl_real_time_warehousing (table_id,event_name,create_time, report_data) VALUES (?,?,?)")
		if err != nil {
			return err
		}

		for _, buffer := range this.buffer {
			if _, err := stmt.Exec(
				buffer.Appid,
				buffer.EventName,
				buffer.CreateTime,
				util.Bytes2str(buffer.Data),
			); err != nil {
				stmt.Close()
				return err
			}
		}

		if err := tx.Commit(); err != nil {
			logs.Logger.Error("入库数据状态出现错误", zap.Error(err))
		} else {
			lostTime := time.Now().Sub(startNow).String()
			len := len(this.buffer)
			if len > 0 {
				logs.Logger.Info("入库实时数据成功", zap.String("所花时间", lostTime), zap.Int("数据长度为", len))
			}
		}
		stmt.Close()

		this.buffer = make([]*RealTimeWarehousingData, 0, this.batchSize)
	}
	this.bufferMutex.Unlock()
	return nil
}

func (this *RealTimeWarehousing) Add(data *RealTimeWarehousingData) (err error) {
	this.bufferMutex.Lock()
	this.buffer = append(this.buffer, data)
	this.bufferMutex.Unlock()

	if this.getBufferLength() >= this.batchSize {
		err := this.Flush()
		return err
	}

	return nil
}

func (this *RealTimeWarehousing) getBufferLength() int {
	this.bufferMutex.RLock()
	defer this.bufferMutex.RUnlock()
	return len(this.buffer)
}

func (this *RealTimeWarehousing) FlushAll() error {
	for this.getBufferLength() > 0 {
		if err := this.Flush(); err != nil {
			return err
		}
	}
	return nil
}

func (this *RealTimeWarehousing) RegularFlushing() {
	go func() {
		ticker := time.NewTicker(time.Duration(this.flushInterval) * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			if err := this.Flush(); err != nil {
				logs.Logger.Error("RealTimeWarehousing RegularFlushing", zap.Error(err))
			}
		}
	}()
}

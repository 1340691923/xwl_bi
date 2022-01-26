package consumer_data

import (
	"context"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	jsoniter "github.com/json-iterator/go"
	"github.com/olivere/elastic"
	"go.uber.org/zap"
	"sync"
	"time"
)

type RealTimeWarehousing struct {
	buffer        []*elastic.BulkIndexRequest
	bufferMutex   *sync.RWMutex
	batchSize     int
	flushInterval int
}

func NewRealTimeWarehousing(batchSize, flushInterval int) *RealTimeWarehousing {
	realTimeWarehousing := &RealTimeWarehousing{
		buffer:        make([]*elastic.BulkIndexRequest, 0, batchSize),
		bufferMutex:   new(sync.RWMutex),
		batchSize:     batchSize,
		flushInterval: flushInterval,
	}

	if flushInterval > 0 {
		realTimeWarehousing.RegularFlushing()
	}

	return realTimeWarehousing
}

func (this *RealTimeWarehousing) Flush() (err error) {
	this.bufferMutex.Lock()
	if len(this.buffer) > 0 {
		startNow := time.Now()
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		bulkRequest := db.EsClient.Bulk()

		for _, buffer := range this.buffer {
			bulkRequest.Add(buffer)
		}
		res, err := bulkRequest.Do(context.Background())

		if err != nil {
			logs.Logger.Error("ES出现错误，休息10秒钟继续", zap.Error(err))
			time.Sleep(time.Second * 10)
			this.Flush()
		} else {
			if res.Errors {
				resStr, _ := json.MarshalToString(res)
				logs.Logger.Error("ES出现错误", zap.String("res", resStr))
			} else {
				lostTime := time.Now().Sub(startNow).String()
				len := len(this.buffer)
				if len > 0 {
					logs.Logger.Info("ES入库成功", zap.String("所花时间", lostTime), zap.Int("数据长度为", len))
				}

			}
		}

		this.buffer = make([]*elastic.BulkIndexRequest, 0, this.batchSize)
	}
	this.bufferMutex.Unlock()
	return nil
}

func (this *RealTimeWarehousing) Add(data *elastic.BulkIndexRequest) (err error) {
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

package consumer_data

import (
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/model"
	"go.uber.org/zap"
	"sync"
	"time"
)

type ReportAcceptStatusData struct {
	PartDate       string
	ReportType     string
	DataName       string
	ErrorReason    string
	ErrorHandling  string
	ReportData     string
	XwlKafkaOffset int64
	TableId        int
	Status         int
}

type ReportAcceptStatus struct {
	buffer        []*ReportAcceptStatusData
	bufferMutex   *sync.RWMutex
	batchSize     int
	flushInterval int
}

const (
	FailStatus    = 0
	SuccessStatus = 1
)

func NewReportAcceptStatus(config model.BatchConfig) *ReportAcceptStatus {
	logs.Logger.Info("NewReportAcceptStatus", zap.Int("batchSize", config.BufferSize), zap.Int("flushInterval", config.FlushInterval))
	reportAcceptStatus := &ReportAcceptStatus{
		buffer:        make([]*ReportAcceptStatusData, 0, config.BufferSize),
		bufferMutex:   new(sync.RWMutex),
		batchSize:     config.BufferSize,
		flushInterval: config.FlushInterval,
	}

	if config.FlushInterval > 0 {
		reportAcceptStatus.RegularFlushing()
	}

	return reportAcceptStatus
}

func (this *ReportAcceptStatus) Flush() (err error) {

	this.bufferMutex.Lock()
	if len(this.buffer) == 0 {
		this.bufferMutex.Unlock()
		return nil
	}

	startNow := time.Now()

	tx, err := db.ClickHouseSqlx.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO xwl_acceptance_status (status,part_date, table_id,report_type, data_name, error_reason, error_handling, report_data, xwl_kafka_offset) VALUES (?,?,?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	for _, buffer := range this.buffer {
		if _, err := stmt.Exec(
			buffer.Status,
			buffer.PartDate,
			buffer.TableId,
			buffer.ReportType,
			buffer.DataName,
			buffer.ErrorReason,
			buffer.ErrorHandling,
			buffer.ReportData,
			buffer.XwlKafkaOffset,
		); err != nil {
			stmt.Close()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		logs.Logger.Error("入库数据状态出现错误", zap.Error(err))
	} else {
		len := len(this.buffer)
		if len > 0 {
			logs.Logger.Info("入库数据状态成功", zap.String("所花时间", time.Now().Sub(startNow).String()), zap.Int("数据长度为", len))
		}
	}
	stmt.Close()
	this.buffer = make([]*ReportAcceptStatusData, 0, this.batchSize)
	this.bufferMutex.Unlock()
	return nil
}

func (this *ReportAcceptStatus) Add(data *ReportAcceptStatusData) (err error) {
	this.bufferMutex.Lock()
	this.buffer = append(this.buffer, data)
	this.bufferMutex.Unlock()

	if this.getBufferLength() >= this.batchSize {
		err := this.Flush()
		return err
	}

	return nil
}

func (this *ReportAcceptStatus) getBufferLength() int {
	this.bufferMutex.RLock()
	defer this.bufferMutex.RUnlock()
	return len(this.buffer)
}

func (this *ReportAcceptStatus) FlushAll() error {
	for this.getBufferLength() > 0 {
		if err := this.Flush(); err != nil {
			return err
		}
	}
	return nil
}

func (this *ReportAcceptStatus) RegularFlushing() {
	go func() {
		ticker := time.NewTicker(time.Duration(this.flushInterval) * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			this.Flush()
		}
	}()
}

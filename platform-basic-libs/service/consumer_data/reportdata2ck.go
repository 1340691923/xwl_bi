package consumer_data

import (
	"bytes"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/model"
	model2 "github.com/1340691923/xwl_bi/platform-basic-libs/sinker/model"
	parser "github.com/1340691923/xwl_bi/platform-basic-libs/sinker/parse"
	"go.uber.org/zap"
	"strings"
	"sync"
	"time"
)

var TableColumnMap sync.Map

type ReportData2CK struct {
	buffer        []FastjsonMetricData
	bufferMutex   *sync.RWMutex
	batchSize     int
	flushInterval int
}

type FastjsonMetricData struct {
	FastjsonMetric *parser.FastjsonMetric
	TableName      string
}

func NewReportData2CK(config model.BatchConfig) *ReportData2CK {
	logs.Logger.Info("NewReportData2CK", zap.Int("batchSize", config.BufferSize), zap.Int("flushInterval", config.FlushInterval))
	reportData2CK := &ReportData2CK{
		buffer:        make([]FastjsonMetricData, 0, config.BufferSize),
		bufferMutex:   new(sync.RWMutex),
		batchSize:     config.BufferSize,
		flushInterval: config.FlushInterval,
	}
	if config.FlushInterval > 0 {
		reportData2CK.RegularFlushing()
	}

	return reportData2CK
}

func (this *ReportData2CK) Flush() (err error) {
	this.bufferMutex.Lock()
	if len(this.buffer) == 0 {
		this.bufferMutex.Unlock()
		return nil
	}
	startNow := time.Now()

	rowsMap := map[string][][]interface{}{}

	for _, obj := range this.buffer {

		rowArr := []interface{}{}
		rows := [][]interface{}{}
		if _, haveKey := rowsMap[obj.TableName]; haveKey {
			rows = rowsMap[obj.TableName]
		} else {
			rowsMap[obj.TableName] = rows
		}
		dims, _ := TableColumnMap.Load(obj.TableName)
		for _, dim := range dims.([]*model2.ColumnWithType) {

			val := parser.GetValueByType(obj.FastjsonMetric, dim)
			rowArr = append(rowArr, val)
		}

		rows = append(rows, rowArr)
		rowsMap[obj.TableName] = rows

	}

	bytesbuffer := bytes.Buffer{}

	TableColumnMap.Range(func(tableName, value interface{}) bool {

		if _, haveKey := rowsMap[tableName.(string)]; haveKey {

			seriesDims := value.([]*model2.ColumnWithType)
			serDimsQuoted := make([]string, len(seriesDims))
			params := make([]string, len(seriesDims))

			for i, serDim := range seriesDims {
				serDimsQuoted[i] = "`" + serDim.Name + "`"
				params[i] = "?"
			}

			bytesbuffer.WriteString("INSERT INTO ")
			bytesbuffer.WriteString(tableName.(string))
			bytesbuffer.WriteString(" (")
			bytesbuffer.WriteString(strings.Join(serDimsQuoted, ","))
			bytesbuffer.WriteString(") ")
			bytesbuffer.WriteString("VALUES (")
			bytesbuffer.WriteString(strings.Join(params, ","))
			bytesbuffer.WriteString(")")
			insertSql := bytesbuffer.String()
			bytesbuffer.Reset()

			tx, err := db.ClickHouseSqlx.Begin()
			if err != nil {
				logs.Logger.Error("CK入库失败", zap.Error(err))
				return false
			}

			stmt, err := tx.Prepare(insertSql)
			if err != nil {
				logs.Logger.Error("CK入库失败", zap.Error(err))
				return false
			}
			defer stmt.Close()
			haveFail := false

			for _, row := range rowsMap[tableName.(string)] {

				if _, err := stmt.Exec(row...); err != nil {
					logs.Logger.Error("CK入库失败", zap.Error(err))
					haveFail = true
				}
			}
			if !haveFail {
				if err := tx.Commit(); err != nil {

					logs.Logger.Error("CK入库失败", zap.Error(err))
					return false
				} else {
					len := len(this.buffer)
					logs.Logger.Info("CK入库成功，", zap.String("所花时间", time.Now().Sub(startNow).String()), zap.Int("数据长度为", len))
				}
			}
		}

		return true
	})

	this.buffer = make([]FastjsonMetricData, 0, this.batchSize)
	this.bufferMutex.Unlock()
	return nil
}

func (this *ReportData2CK) Add(data FastjsonMetricData) (err error) {
	this.bufferMutex.Lock()

	this.buffer = append(this.buffer, data)
	this.bufferMutex.Unlock()

	if this.getBufferLength() >= this.batchSize {
		err := this.Flush()
		return err
	}

	return nil
}

func (this *ReportData2CK) getBufferLength() int {
	this.bufferMutex.RLock()
	defer this.bufferMutex.RUnlock()
	return len(this.buffer)
}

func (this *ReportData2CK) FlushAll() error {
	for this.getBufferLength() > 0 {
		if err := this.Flush(); err != nil {
			return err
		}
	}
	return nil
}

func (this *ReportData2CK) RegularFlushing() {
	go func() {
		ticker := time.NewTicker(time.Duration(this.flushInterval) * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			this.Flush()
		}
	}()
}

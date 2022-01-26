package consumer_data

import (
	"bytes"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	model2 "github.com/1340691923/xwl_bi/platform-basic-libs/sinker/model"
	parser "github.com/1340691923/xwl_bi/platform-basic-libs/sinker/parse"
	"go.uber.org/zap"
	"log"
	"strings"
	"sync"
	"time"
)

var TableColumnMap sync.Map

type ReportData2CK struct {
	buffer        []map[string]*parser.FastjsonMetric
	bufferMutex   *sync.RWMutex
	batchSize     int
	flushInterval int
}

func NewReportData2CK(batchSize int, flushInterval int) *ReportData2CK {
	logs.Logger.Info("NewReportData2CK", zap.Int("batchSize", batchSize), zap.Int("flushInterval", flushInterval))
	reportData2CK := &ReportData2CK{
		buffer:        make([]map[string]*parser.FastjsonMetric, 0, batchSize),
		bufferMutex:   new(sync.RWMutex),
		batchSize:     batchSize,
		flushInterval: flushInterval,
	}
	if flushInterval > 0 {
		reportData2CK.RegularFlushing()
	}

	return reportData2CK
}

func (this *ReportData2CK) Flush() (err error) {
	this.bufferMutex.Lock()

	startNow := time.Now()

	rowsMap := map[string][][]interface{}{}
	for _, data := range this.buffer {

		for tableName, metric := range data {

			rows := [][]interface{}{}

			if _, haveKey := rowsMap[tableName]; haveKey {
				rows = rowsMap[tableName]
			} else {
				rowsMap[tableName] = rows
			}
			v, _ := TableColumnMap.Load(tableName)
			dims := v.([]*model2.ColumnWithType)
			var rowArr []interface{}
			for _, dim := range dims {
				val := parser.GetValueByType(metric, dim)
				rowArr = append(rowArr, val)
			}
			rows = append(rows, rowArr)
			rowsMap[tableName] = rows
		}
	}

	buffer := bytes.Buffer{}

	TableColumnMap.Range(func(key, value interface{}) bool {

		tableName := key.(string)

		if _, haveKey := rowsMap[tableName]; haveKey {

			seriesDims := value.([]*model2.ColumnWithType)
			serDimsQuoted := make([]string, len(seriesDims))
			params := make([]string, len(seriesDims))
			for i, serDim := range seriesDims {
				serDimsQuoted[i] = "`" + serDim.Name + "`"
				params[i] = "?"
			}

			buffer.WriteString("INSERT INTO ")
			buffer.WriteString(tableName)
			buffer.WriteString(" (")
			buffer.WriteString(strings.Join(serDimsQuoted, ","))
			buffer.WriteString(") ")
			buffer.WriteString("VALUES (")
			buffer.WriteString(strings.Join(params, ","))
			buffer.WriteString(")")

			insertSql := buffer.String()
			buffer.Reset()
			tx, err := db.ClickHouseSqlx.Begin()
			if err != nil {
				logs.Logger.Error("CK入库失败", zap.Error(err))
				return false
			}
			log.Println("insertSql",insertSql)
			stmt, err := tx.Prepare(insertSql)
			if err != nil {
				logs.Logger.Error("CK入库失败", zap.Error(err))
				return false
			}
			defer stmt.Close()
			haveFail := false
			for _, row := range rowsMap[tableName] {
				log.Println("row",row)
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
					lostTime := time.Now().Sub(startNow).String()
					len := len(this.buffer)
					logs.Logger.Info("CK入库成功，", zap.String("所花时间", lostTime), zap.Int("数据长度为", len))
				}
			}
		}

		return true
	})

	this.buffer = make([]map[string]*parser.FastjsonMetric, 0, this.batchSize)
	this.bufferMutex.Unlock()
	return nil
}

func (this *ReportData2CK) Add(data map[string]*parser.FastjsonMetric) (err error) {
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

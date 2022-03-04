package consumer_data

import (
	"bytes"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	model2 "github.com/1340691923/xwl_bi/platform-basic-libs/sinker/model"
	parser "github.com/1340691923/xwl_bi/platform-basic-libs/sinker/parse"
	"go.uber.org/zap"
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
	if len(this.buffer)==0{
		this.bufferMutex.Unlock()
		return nil
	}
	startNow := time.Now()

	rowsMap := map[string][][]interface{}{}


	for bufferIndex := range this.buffer {
		for tableName := range this.buffer[bufferIndex] {
			rowArr := []interface{}{}
			rows := [][]interface{}{}
			if _, haveKey := rowsMap[tableName]; haveKey {
				rows = rowsMap[tableName]
			} else {
				rowsMap[tableName] = rows
			}
			dims, _ := TableColumnMap.Load(tableName)
			for _, dim := range dims.([]*model2.ColumnWithType) {
				val := parser.GetValueByType(this.buffer[bufferIndex][tableName], dim)
				logs.Logger.Sugar().Errorf("dim.SourceName",dim.SourceName,val)
				rowArr = append(rowArr, val)
			}

			rows = append(rows, rowArr)
			rowsMap[tableName] = rows
		}
	}

	bytesbuffer:=bytes.Buffer{}

	TableColumnMap.Range(func(tableName, value interface{}) bool {

		if _, haveKey := rowsMap[tableName.(string)]; haveKey {

			seriesDims := value.([]*model2.ColumnWithType)
			serDimsQuoted := make([]string, len(seriesDims))
			params := make([]string, len(seriesDims))

			for i, serDim := range seriesDims {
				serDimsQuoted[i] ="`"+serDim.Name+"`"
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
				logs.Logger.Sugar().Errorf("bytesbuffer.String()",insertSql)

				if _, err := stmt.Exec(row...); err != nil {
					logs.Logger.Error("CK入库失败", zap.Error(err))
					haveFail = true
				}
				logs.Logger.Sugar().Errorf("bytesbuffer.String() args",row...)
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

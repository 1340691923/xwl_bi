package sinker

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/model"
	model2 "github.com/1340691923/xwl_bi/platform-basic-libs/sinker/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/sinker/parse"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	ErrTblNotExist       = errors.Errorf("table doesn't exist")
	selectSQLTemplate    = `select name, type, default_kind from system.columns where database = '%s' and table = '%s'`
	lowCardinalityRegexp = regexp.MustCompile(`LowCardinality\((.+)\)`)
)

const DimsHash = "dimsHash_"

func GetDimsCachekey(database, table string) string {
	b := bytes.Buffer{}
	b.WriteString(DimsHash)
	b.WriteString(database)
	b.WriteString("_")
	b.WriteString(table)
	dimsCachekey := b.String()
	return dimsCachekey
}

var dimsCacheMap sync.Map

func ClearDimsCacheByTime(clearTime time.Duration) {

	for {
		time.Sleep(clearTime)
		dimsCacheMap.Range(func(key, value interface{}) bool {
			ClearDimsCacheByRedis(key.(string))
			dimsCacheMap.Delete(key)
			return true
		})

	}
}

func ClearDimsCacheByTimeBylocal(clearTime time.Duration) {

	for {
		time.Sleep(clearTime)
		dimsCacheMap.Range(func(key, value interface{}) bool {
			ClearDimsCacheByRedis(key.(string))
			dimsCacheMap.Delete(key)
			return true
		})

	}
}

func ClearDimsCacheByRedis(key string) {
	redisConn := db.RedisPool.Get()
	defer redisConn.Close()

	_, err := redisConn.Do("unlink", key)
	if err != nil {
		_, err = redisConn.Do("del", key)
		if err != nil {
			logs.Logger.Error("err", zap.Error(err))
		}
	}
}

func ClearDimsCacheByKey(key string) {
	dimsCacheMap.Delete(key)
}

func GetDims(database, table string, excludedColumns []string, conn *sqlx.DB, onlyRedis bool) (dims []*model2.ColumnWithType, err error) {

	dimsCachekey := GetDimsCachekey(database, table)
	if !onlyRedis {
		cache, load := dimsCacheMap.Load(dimsCachekey)
		if load {
			return cache.([]*model2.ColumnWithType), nil
		}
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	redisConn := db.RedisPool.Get()
	defer redisConn.Close()
	dimsBytes, redisErr := redis.Bytes(redisConn.Do("get", dimsCachekey))

	if redisErr == nil && len(dimsBytes) != 0 {
		dimsCache, err := util.GzipUnCompressByte(dimsBytes)
		if err == nil {
			jsonErr := json.Unmarshal(dimsCache, &dims)
			if jsonErr == nil {
				dimsCacheMap.Store(dimsCachekey, dims)
				return dims, err
			}
			logs.Logger.Error("jsonErr", zap.Error(jsonErr))
		} else {
			logs.Logger.Error("GzipUnCompressByte Err", zap.Error(err))
		}

	} else {
		logs.Logger.Error("redisErr", zap.Error(redisErr))
	}

	var rs *sql.Rows
	if rs, err = conn.Query(fmt.Sprintf(selectSQLTemplate, database, table)); err != nil {
		err = errors.Wrapf(err, "")
		return dims, err
	}
	defer rs.Close()

	var name, typ, defaultKind string
	for rs.Next() {
		if err = rs.Scan(&name, &typ, &defaultKind); err != nil {
			err = errors.Wrapf(err, "")
			return dims, err
		}
		typ = lowCardinalityRegexp.ReplaceAllString(typ, "$1")
		if !util.InstrArr(excludedColumns, name) && defaultKind != "MATERIALIZED" {
			tp, nullable := parser.WhichType(typ)
			dims = append(dims, &model2.ColumnWithType{Name: name, Type: tp, Nullable: nullable, SourceName: GetSourceName(name)})
		}
	}
	if len(dims) == 0 {
		err = errors.Wrapf(ErrTblNotExist, "%s.%s", database, table)
		return dims, err
	}
	dimsCacheMap.Store(dimsCachekey, dims)

	res, _ := json.Marshal(dims)
	s, err := util.GzipCompressByte(res)
	if err != nil {
		return dims, err
	}
	_, err = redisConn.Do("SETEX", dimsCachekey, 60*60*6, s)

	return dims, err
}

func GetSourceName(name string) (sourcename string) {
	sourcename = strings.Replace(name, ".", "\\.", -1)
	return
}

func ChangeSchema(newKeys *sync.Map, dbname, table string, dims []*model2.ColumnWithType) ([]*model2.ColumnWithType, error) {
	var queries []string
	var err error
	newKeys.Range(func(key, value interface{}) bool {

		strKey, _ := key.(string)
		intVal := value.(int)
		var strVal string
		switch intVal {
		case parser.Int:
			strVal = "Float64"
		case parser.Float:
			strVal = "Float64"
		case parser.String:
			strVal = "String"
		case parser.DateTime:
			strVal = "Nullable(DateTime)"
		case parser.IntArray:
			strVal = "Array(Int64)"
		case parser.FloatArray:
			strVal = "Array(Float64)"
		case parser.StringArray:
			strVal = "Array(String)"
		case parser.DateTimeArray:
			strVal = "Array(DateTime)"
		default:
			err = errors.Errorf("BUG: unsupported column type %s", strVal)
			return false
		}
		query := fmt.Sprintf("ALTER TABLE %s.%s %s ADD COLUMN IF NOT EXISTS `%s` %s", dbname, table, GetClusterSql(), strKey, strVal)
		queries = append(queries, query)
		tp, nullable := parser.WhichType(strVal)
		dims = append(dims, &model2.ColumnWithType{
			Name:       strKey,
			Type:       tp,
			Nullable:   nullable,
			SourceName: GetSourceName(strKey),
		})

		return true
	})

	//sort.Strings(queries)

	for _, query := range queries {
		logs.Logger.Info(fmt.Sprintf("executing sql=> %s", query), zap.String("table", table))
		if _, err = db.ClickHouseSqlx.Exec(query); err != nil {
			err = errors.Wrapf(err, query)
			return dims, err
		}
	}

	return dims, nil
}

func GetClusterSql() string {
	if model.GlobConfig.Comm.ClickHouse.ClusterName == "" {
		return " "
	}
	b := bytes.Buffer{}
	b.WriteString(" on cluster ")
	b.WriteString(model.GlobConfig.Comm.ClickHouse.ClusterName)
	b.WriteString(" ")
	clusterSql := b.String()
	return clusterSql
}

func GetMergeTree(tableName string) string {
	if model.GlobConfig.Comm.ClickHouse.ClusterName == "" {
		return "MergeTree"
	}
	return `ReplicatedMergeTree('/clickhouse/` + model.GlobConfig.Comm.ClickHouse.DbName + `/tables/{` + model.GlobConfig.Comm.ClickHouse.MacrosShardKeyName + `}/` + tableName + `', '{` + model.GlobConfig.Comm.ClickHouse.MacrosReplicaKeyName + `}')`
}

func GetReplacingMergeTree(tableName, ext string) string {
	if model.GlobConfig.Comm.ClickHouse.ClusterName == "" {
		return "ReplacingMergeTree"
	}
	return `ReplicatedReplacingMergeTree('/clickhouse/` + model.GlobConfig.Comm.ClickHouse.DbName + `/tables/{` + model.GlobConfig.Comm.ClickHouse.MacrosShardKeyName + `}/` + tableName + `', '{` + model.GlobConfig.Comm.ClickHouse.MacrosReplicaKeyName + `}',` + ext + `)`
}

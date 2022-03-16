package ck

import (
	"fmt"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/sinker"
	"log"
)

//初始化clickhouse 表数据
func Init() {
	var err error

	_, err = db.ClickHouseSqlx.Exec(` create database if not exists ` + model.GlobConfig.Comm.ClickHouse.DbName + ` ` + sinker.GetClusterSql())

	if err != nil {
		log.Println(fmt.Sprintf("clickhouse 建库 "+model.GlobConfig.Comm.ClickHouse.DbName+" 失败:%s", err.Error()))
		panic(err)
	}

	_, err = db.ClickHouseSqlx.Exec(`DROP TABLE IF EXISTS xwl_acceptance_status` + sinker.GetClusterSql() + `;`)

	if err != nil {
		log.Println(fmt.Sprintf("clickhouse 删除表 xwl_acceptance_status 失败:%s", err.Error()))
		panic(err)
	}

	_, err = db.ClickHouseSqlx.Exec(`
		
		CREATE TABLE xwl_acceptance_status ` + sinker.GetClusterSql() + `
		(
		
			table_id Int64,
		
			part_date DateTime DEFAULT now(),
		
			data_name String,
		
			error_reason String,
		
			error_handling String,
		
			report_type String,
		
			report_data String,
		
			xwl_kafka_offset Int64,
		
			part_event String,
		
			status Int32
		)
		ENGINE = ` + sinker.GetMergeTree("xwl_acceptance_status") + `
		PARTITION BY (toYYYYMMDD(part_date))
		ORDER BY (toYYYYMMDD(part_date),
		 table_id,
		 data_name,
		 error_reason,
		 error_handling,
		 report_type,
		 status)
		TTL part_date + toIntervalMonth(3)
		SETTINGS index_granularity = 8192;
`)
	if err != nil {
		log.Println(fmt.Sprintf("clickhouse 建表 xwl_acceptance_status 失败:%s", err.Error()))
		panic(err)
	}

	_, err = db.ClickHouseSqlx.Exec(`DROP TABLE IF EXISTS xwl_real_time_warehousing` + sinker.GetClusterSql() + `;`)

	if err != nil {
		log.Println(fmt.Sprintf("clickhouse 删除表 xwl_real_time_warehousing 失败:%s", err.Error()))
		panic(err)
	}

	_, err = db.ClickHouseSqlx.Exec(`
		
		CREATE TABLE xwl_real_time_warehousing ` + sinker.GetClusterSql() + `
		(
		
			table_id Int64,
		
			create_time DateTime DEFAULT now(),
		
			event_name String,
		
			report_data String
		)
		ENGINE = ` + sinker.GetMergeTree("xwl_real_time_warehousing") + ` 
		PARTITION BY (toYYYYMMDD(create_time))
		ORDER BY (toYYYYMMDD(create_time),
		 table_id,
		 event_name)
		TTL create_time + toIntervalMonth(3)
		SETTINGS index_granularity = 8192;
`)
	if err != nil {
		log.Println(fmt.Sprintf("clickhouse 建表 xwl_real_time_warehousing 失败:%s", err.Error()))
		panic(err)
	}

	log.Println("初始化CK数据完成！")
}

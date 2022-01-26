package mysql

import (
	_ "embed"
	"fmt"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"log"
	"strings"
)

//go:embed bi.sql
var SqlByte []byte

//初始化mysql数据
func Init() {
	var err error

	_, err = db.Sqlx.Exec(` create database if not exists ` + model.GlobConfig.Comm.Mysql.DbName)

	if err != nil {
		log.Println(fmt.Sprintf("mysql 执行建库语句失败:%s", err.Error()))
		panic(err)
	}

	execSqlArr := strings.Split(util.Bytes2str(SqlByte), ";")

	for _, execSql := range execSqlArr {
		_, err = db.Sqlx.Exec(execSql)
		if err != nil {
			log.Println(fmt.Sprintf("mysql 执行建表语句sql:%v失败:%s", execSql, err.Error()))
			panic(err)
		}
	}

	log.Println("初始化mysql数据完成！")
}

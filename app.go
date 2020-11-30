package goBase

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"time"
)

type APP struct {
	Conf   *Config
	DbCons map[string]*sql.DB
}

// GetMysqlCon 获取数据库连接
func (a *APP) GetMysqlCon(name string) *sql.DB{
	if con,ok := a.DbCons[name];ok{
		return con
	}
	conf,err := a.Conf.getDBConfByName(name)
	if err != nil{
		panic(err)
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%d)/%v?charset=%v",conf.User,conf.Pass,conf.Host,conf.Port,conf.DbName,conf.Charset))
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(conf.MaxOpenCons)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	a.DbCons[name] = db
	return db
}
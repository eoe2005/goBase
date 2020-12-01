package goBase

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// APP 应用
type APP struct {
	AppConfig *AppConfig
	DbCons    map[string]*sql.DB
	RedisCons map[string]*GRedis
	Aes       *Secure
}

// GetMysqlCon 获取数据库连接
func (a *APP) GetMysqlCon(name string) *sql.DB {
	if con, ok := a.DbCons[name]; ok {
		return con
	}
	conf, err := a.AppConfig.getDBConfByName(name)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%d)/%v?charset=%v", conf.User, conf.Pass, conf.Host, conf.Port, conf.DbName, conf.Charset))
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(conf.MaxOpenCons)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	a.DbCons[name] = db
	return db
}

// RandString 生成字符串
func (a *APP) RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

// ServerDefaultHandle 配置默认的路由信息
func (a *APP) ServerDefaultHandle(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if h, ok := a.AppConfig.APIRouters[path]; ok {
		if handle, o := h.(*ActionHandle); o {
			handle.Execute(a, w, r)
			return
		}
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("{\"code\" : 404,\"msg\":\"接口不存在\",\"data\":\"\"}"))
}

// Run 程序运行
func (a *APP) Run() {

	http.HandleFunc("/", a.ServerDefaultHandle)

	http.ListenAndServe(fmt.Sprintf(":%v", a.AppConfig.Port), nil)
}

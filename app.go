package goBase

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

// APP 应用
type APP struct {
	AppConfig *AppConfig
	DbCons    map[string]*sql.DB
	Aes       *Secure
}

// GetRedis 获取Redis连接
func (a *APP) GetRedis(name string) *redis.Client {
	conf, err := a.AppConfig.getRedisConfByName(name)
	if err != nil {
		panic(err)
	}
	return redis.NewClient(&redis.Options{
		Password: conf.Auth,
		Addr:     conf.Host,
		DB:       conf.DB,
	})
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
	req := &GReq{
		App:     a,
		W:       w,
		R:       r,
		Context: context.Background(),
	}
	path := r.URL.Path
	isLogin := false
	h, ok := a.AppConfig.Routers[path]
	if !ok {
		h, ok = a.AppConfig.RoutersLogined[path]
		isLogin = true
	}
	if ok {
		if isLogin {
			uid := req.GetUID()
			if uid < 1 {
				req.Fail(301, "账号没有登录")
				return
			}
		}
		t := reflect.TypeOf(h).Kind()
		if t == reflect.Struct {
			m := reflect.ValueOf(h).MethodByName("Handle")
			args := []reflect.Value{reflect.ValueOf(req)}
			m.Call(args)
			return
		}
		if t == reflect.Func {
			handle, _ := h.(func(*GReq))
			handle(req)
			return
		} else {
			req.Fail(404, "接口不存在")
			return
		}
		return
	}
	req.Fail(404, "接口不存在")

}

// Run 程序运行
func (a *APP) Run() {

	http.HandleFunc("/", a.ServerDefaultHandle)

	http.ListenAndServe(fmt.Sprintf(":%v", a.AppConfig.Port), nil)
}

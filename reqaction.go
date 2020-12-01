package goBase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// ActionHandle 路由要处理的防反
type ActionHandle interface {
	Handle()
	Execute(*APP, http.ResponseWriter, *http.Request)
}

// Action 接口处理
type Action struct {
	W       http.ResponseWriter
	R       *http.Request
	App     *APP
	IsLogin bool
	UID     int64
}

// Execute 程序的入口
func (a *Action) Execute(app *APP, w http.ResponseWriter, req *http.Request) {
	a.W = w
	a.R = req
	a.App = app
	a.UID = a.GetUID()
	if a.IsLogin && a.UID < 1 {
		a.Fail(301, "账号没有登录")
	} else {
		a.Handle()
	}
}
func (a *Action) Handle() {}

// GetUID 获取用户的ID
func (a *Action) GetUID() int64 {
	data := a.GetAesCookie("guk")
	if len(data) < 1 {
		return 0
	}
	uid := a.GetAesCookie(data)
	if len(uid) < 1 {
		return 0
	}
	u, e := strconv.ParseInt(uid, 10, 64)
	if e != nil {
		return 0
	}
	return u

}

// SetUID 设置用户信息
func (a *Action) SetUID(uid int64) {
	code := a.App.RandString(10)
	a.SetAesCookie("guk", code)
	a.SetAesCookie(code, uid)
}

// SetCookie 设置Cookie
func (a *Action) SetCookie(val http.Cookie) {
	a.W.Header().Set("Set-Cookie", val.String())
}

// SetAesCookie AES Cookie
func (a *Action) SetAesCookie(key string, val interface{}) {
	data, e := a.App.Aes.Encode(fmt.Sprintf("%v", val))
	if e != nil {
		return
	}
	a.SetCookie(http.Cookie{Name: key, Value: data, Path: "/"})
}

// GetAesCookie 获取AES加密的KEY
func (a *Action) GetAesCookie(key string) string {
	data := a.GetCookie(key)
	if len(data) < 1 {
		return ""
	}
	ret, e := a.App.Aes.Decode(data)
	if e != nil {
		return ""
	}
	return ret
}

// GetCookie 获取Cookie信息
func (a *Action) GetCookie(key string) string {
	v, e := a.R.Cookie(key)
	if e != nil {
		return ""
	}
	return v.Value
}

// Fail 输出错误信息
func (a *Action) Fail(code int64, msg string) {
	a.W.Header().Add("Content-Type", "application/json")
	r := map[string]interface{}{"code": code, "msg": msg, "data": ""}
	rd, _ := json.Marshal(r)
	a.W.Write(rd)
}

//Success 成功时候的输出
func (a *Action) Success(data interface{}) {
	a.W.Header().Add("Content-Type", "application/json")
	r := map[string]interface{}{"code": 0, "msg": "", "data": data}
	rd, _ := json.Marshal(r)
	a.W.Write(rd)
}

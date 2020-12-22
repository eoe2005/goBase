package goBase

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ActionHandle 路由要处理的防反
type ActionHandle interface {
	Handle(req *GReq)
}

// GReq 接口处理
type GReq struct {
	W       http.ResponseWriter
	R       *http.Request
	App     *APP
	Context context.Context
}

// GetUID 获取用户的ID
func (a *GReq) GetUID() int64 {
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
func (a *GReq) SetUID(uid int64) {
	code := a.App.RandString(10)
	a.SetAesCookie("guk", code)
	a.SetAesCookie(code, uid)
}

// DeleteCookie 删除COOKIE
func (a *GReq) DeleteCookie(key ...string) {
	for i := range key {
		c := &http.Cookie{
			Name:   key[i],
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}
		http.SetCookie(a.W, c)
	}
}

// SetCookie 设置Cookie
func (a *GReq) SetCookie(key, val string) {
	c := &http.Cookie{
		Name:   key,
		Value:  val,
		Path:   "/",
		MaxAge: 86400 * 7,
		//Domain: "localhost",
		//Expires: time.Now().AddDate(0, 1, 0),
	}

	http.SetCookie(a.W, c)
	//a.W.Header().Add("set-cookie", c.String())
}

// SetAesCookie AES Cookie
func (a *GReq) SetAesCookie(key string, val interface{}) {
	data, e := a.App.Aes.Encode(fmt.Sprintf("%v", val))
	if e != nil {
		LogError("设置AESCookie 失败 : %v %v", key, e)
		return
	}
	a.SetCookie(key, data)
}

// GetIP 获取客户端IP
func (a *GReq) GetIP() string {
	xForwardedFor := a.R.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(a.R.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(a.R.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

// GetAesCookie 获取AES加密的KEY
func (a *GReq) GetAesCookie(key string) string {
	data := a.GetCookie(key)
	if len(data) < 1 {
		return ""
	}
	ret, e := a.App.Aes.Decode(data)
	if e != nil {
		LogError("AES解密错误 %v", e)
		return ""
	}
	return ret
}

// GetCookie 获取Cookie信息
func (a *GReq) GetCookie(key string) string {
	v, e := a.R.Cookie(key)
	if e != nil {
		return ""
	}
	return v.Value
}

// Fail 输出错误信息
func (a *GReq) Fail(code int64, msg string) {
	a.W.Header().Add("Content-Type", "application/json")
	r := map[string]interface{}{"code": code, "msg": msg, "data": ""}
	rd, _ := json.Marshal(r)
	a.W.Write(rd)
}

//Success 成功时候的输出
func (a *GReq) Success(data interface{}) {
	a.W.Header().Add("Content-Type", "application/json")
	r := map[string]interface{}{"code": 0, "msg": "", "data": data}
	rd, _ := json.Marshal(r)
	a.W.Write(rd)
}

// CheckPostParams 检查POST参数
func (a *GReq) CheckPostParams(d map[string]string) bool {
	for k, v := range d {
		t := a.R.FormValue(k)
		if !a.checkParams(t, v) {
			return false
		}
	}
	return true
}

// CheckGetParams 检查GET参数
func (a *GReq) CheckGetParams(d map[string]string) bool {
	data := a.R.URL.Query()
	for k, v := range d {
		t := data.Get(k)
		if !a.checkParams(t, v) {
			return false
		}
	}
	return true
}

// Post 获取POST数据
func (a *GReq) Post(key string) string {
	return a.R.PostFormValue(key)
}

func (a *GReq) checkParams(v, rules string) bool {
	vals := strings.Split(rules, "|")
	for i := range vals {
		switch vals[i] {
		case "required":
			if v == "" {
				a.Fail(201, "参数错误")
				return false
			}
		case "email":
			pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
			reg := regexp.MustCompile(pattern)
			if !reg.MatchString(v) {
				a.Fail(201, "邮箱格式错误")
				return false
			}
		case "code":
			if len(v) < 4 {
				a.Fail(201, "验证码错误")
				return false
			}

		case "mobile":
			pattern := `1(3|5|7|8)\d{9}` //匹配电子邮箱
			reg := regexp.MustCompile(pattern)
			if !reg.MatchString(v) {
				a.Fail(201, "手机号格式错误")
				return false
			}
		case "passwd":
			if len(v) < 6 || len(v) > 20 {
				a.Fail(201, "密码必须大于6位小于20位")
				return false
			}
		}
	}
	return true

}

// DbInsert 数据插入
func (a *GReq) DbInsert(format string, args ...interface{}) int64 {
	return DBInsert(a.App.GetMysqlCon("default"), format, args...)
}

//DbDelete 删除数据
func (a *GReq) DbDelete(format string, args ...interface{}) int64 {
	return DBDelete(a.App.GetMysqlCon("default"), format, args...)
}

// DbFetchRow 获取一行数据
func (a *GReq) DbFetchRow(format string, args ...interface{}) map[string]interface{} {
	return DBGetRow(a.App.GetMysqlCon("default"), format, args...)
}

// DbFetchAll 获取全部数据
func (a *GReq) DbFetchAll(format string, args ...interface{}) []map[string]interface{} {
	return DBGetAll(a.App.GetMysqlCon("default"), format, args...)
}

// DbUpdateData 更新数据
func (a *GReq) DbUpdateData(format string, args ...interface{}) int64 {
	return DBUpdate(a.App.GetMysqlCon("default"), format, args...)
}

// SetCacheNx 设置缓存
func (a *GReq) SetCacheNx(k string, val interface{}, timeout int64) {
	rdb := a.App.GetRedis("default")
	defer rdb.Close()
	rdb.SetNX(a.Context, k, val, time.Second*time.Duration(timeout))
}

// SetCache 没有超时时间
func (a *GReq) SetCache(k string, val interface{}) {
	rdb := a.App.GetRedis("default")
	defer rdb.Close()
	rdb.Set(a.Context, k, val, 0)
}

// CacheGet 获取Cache的配置
func (a *GReq) CacheGet(k string) string {
	rdb := a.App.GetRedis("default")
	defer rdb.Close()
	r, e := rdb.Get(a.Context, k).Result()
	if e != nil {
		return ""
	}
	return r
}

//CacheDel 删除Cache
func (a *GReq) CacheDel(k ...string) {
	rdb := a.App.GetRedis("default")
	defer rdb.Close()
	rdb.Del(a.Context, k...)
}

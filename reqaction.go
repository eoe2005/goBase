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
	"text/template"
	"time"

	"github.com/go-redis/redis/v8"
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

// GetAdminUID 获取登录的管理员的UID
func (a *GReq) GetAdminUID() int64 {
	data := a.GetAesCookie("aui")
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

// SetAdminUID 设置管理员的UID
func (a *GReq) SetAdminUID(adminUID int64) {
	code := a.App.RandString(10)
	a.SetAesCookie("aui", code, 0)
	a.SetAesCookie(code, adminUID, 0)
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
	a.SetAesCookie("guk", code, 0)
	a.SetAesCookie(code, uid, 0)
}

// GetTableDefault 获取数据库表的配置
func (a *GReq) GetTableDefault(tableName string) *DBTable {
	return a.GetTable(tableName, "default")
}

// GetTable 获取数据库表的配置
func (a *GReq) GetTable(tableName, conName string) *DBTable {
	return InitDBTable(tableName, conName, a.App.GetMysqlCon(conName))
}

// Display 直接渲染模板
func (a *GReq) Display(templatName string, data interface{}) {
	templateHtm := a.RedisDefaultGetFunc("tmp_name:"+templatName, func(args ...interface{}) interface{} {
		a := args[0].(*GReq)

		row := a.GetTableDefault("tb_template").FindByWhere("template_name=?", templatName)
		if row != nil {
			return ""
		}
		contents, ok := row["content"]
		if ok {
			return contents
		}
		return ""
	}, 5, a).(string)
	t, e := template.New(templatName).Parse(templateHtm)
	if e != nil {
		a.Fail(500, "解析数据失败")
		return
	}
	a.W.Header().Add("Content-Type", "text/html; charset=utf-8")
	t.Execute(a.W, data)

	return
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
func (a *GReq) SetCookie(key, val string, MaxAge int) {
	if MaxAge == 0 {
		MaxAge = 86400 * 365
	}
	c := &http.Cookie{
		Name:   key,
		Value:  val,
		Path:   "/",
		MaxAge: MaxAge,
	}

	http.SetCookie(a.W, c)
	//a.W.Header().Add("set-cookie", c.String())
}

// SetAesCookie AES Cookie
func (a *GReq) SetAesCookie(key string, val interface{}, maxAge int) {
	data, e := a.App.Aes.Encode(fmt.Sprintf("%v", val))
	if e != nil {
		LogError("设置AESCookie 失败 : %v %v", key, e)
		return
	}
	a.SetCookie(key, data, maxAge)
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
	a.W.Header().Add("Content-Type", "application/json; charset=utf-8")
	r := map[string]interface{}{"code": code, "msg": msg, "data": ""}
	rd, _ := json.Marshal(r)
	a.W.Write(rd)
}

//Success 成功时候的输出
func (a *GReq) Success(data interface{}) {
	a.W.Header().Add("Content-Type", "application/json; charset=utf-8")
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

// Get 获取Get参数
func (a *GReq) Get(key string) string {
	return a.R.URL.Query().Get(key)
}

// GetInt 获取GET参数
func (a *GReq) GetInt(key string) int64 {
	r, e := strconv.ParseInt(a.Get(key), 10, 64)
	if e != nil {
		return 0
	}
	return r
}

// PostInt 获取POSTInt 数据
func (a *GReq) PostInt(key string) int64 {
	r, e := strconv.ParseInt(a.R.PostFormValue(key), 10, 64)
	if e != nil {
		return 0
	}
	return r
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

// RedisDefaultGetFunc 设置缓存
func (a *GReq) RedisDefaultGetFunc(k string, callfunc func(args ...interface{}) interface{}, timeout int64, args ...interface{}) interface{} {
	rdb := a.App.GetRedis("default")
	defer rdb.Close()
	r, ger := rdb.Get(a.Context, k).Result()
	if ger != nil || r == "" {
		r2 := callfunc(args...)
		dd, e := json.Marshal(r2)
		if e == nil {
			r = string(dd)
			rdb.Set(a.Context, k, r, time.Second*time.Duration(timeout))

		}
	}
	var ret interface{}
	if nil != json.Unmarshal([]byte(r), &ret) {
		return nil
	}
	return ret
}

//RedisDefaultDel 删除Cache
func (a *GReq) RedisDefaultDel(k ...string) {
	rdb := a.App.GetRedis("default")
	defer rdb.Close()
	rdb.Del(a.Context, k...)
}

// RedisDefaultFunc redis 一般操作
func (a *GReq) RedisDefaultFunc(funcCall func(r *redis.Client, args ...interface{}) interface{}, args ...interface{}) interface{} {
	red := a.App.GetRedis("default")
	defer red.Close()
	return funcCall(red, args...)
}

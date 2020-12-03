# goBase
golang 基础类库

# 发送邮件
```go
package main
import(
	g "github.com/eoe2005/goBase"
)
func main(){
	sm := &g.Smtp{
		FromName : "测试账号",
		FromEmail : "notreplay@smtp.服务器域名",
		ToName : "大飞哥",
		ToEmail : "1234@qq.com",
		Subject : "这个邮件很重要",
		Text : "收到邮件必须要回的",
		Html : "<a href='http://测试域名'>点击这个链接回到站点</a>",
		Files : []string{"/home/cc-user/test.go","/home/cc-user/t.php"},
	}
	sm.Send()
}
```

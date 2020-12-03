package goBase

import (
	"bufio"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	"path"
	"strings"
	"time"
)

const STMP_LFRT = "\r\n"

type Smtp struct {
	FromEmail string
	FromName string
	ToEmail string
	ToName string
	Subject string
	Html string
	Text string
	Files []string
}
// Send 发送邮件
func (r *Smtp) Send()  {
	data := []string{}
	data = append(data,fmt.Sprintf("From: =?utf-8?B?%v?= <%v>%v",base64.StdEncoding.EncodeToString([]byte(r.FromName)),r.FromEmail,STMP_LFRT))
	data = append(data,fmt.Sprintf("To: =?utf-8?B?%v?= <%v>%v",base64.StdEncoding.EncodeToString([]byte(r.ToName)),r.ToEmail,STMP_LFRT))
	data = append(data,fmt.Sprintf("Subject: =?utf-8?B?%v?=%v",base64.StdEncoding.EncodeToString([]byte(r.Subject)),STMP_LFRT))

	data = append(data,fmt.Sprintf("%v%v","MIME-Version: 1.0",STMP_LFRT))
	h := md5.New()
	h.Write([]byte(time.Now().String()))
	boundary := hex.EncodeToString(h.Sum(nil))

	data = append(data,fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%v\"%v",boundary,STMP_LFRT))
	if r.Text != ""{
		data = append(data,STMP_LFRT)
		data = append(data,fmt.Sprintf("--%v%v",boundary,STMP_LFRT))
		data = append(data,fmt.Sprintf("Content-Type: text/plain; charset = \"utf-8\"%v",STMP_LFRT))
		data = append(data,fmt.Sprintf("Content-Transfer-Encoding: base64%v",STMP_LFRT))
		data = append(data,STMP_LFRT)
		data = append(data,fmt.Sprintf("%v%v",base64.StdEncoding.EncodeToString([]byte(r.Text)),STMP_LFRT))
	}
	if r.Html != ""{
		data = append(data,STMP_LFRT)
		data = append(data,fmt.Sprintf("--%v%v",boundary,STMP_LFRT))
		data = append(data,fmt.Sprintf("Content-Type: text/html; charset = \"utf-8\"%v",STMP_LFRT))
		data = append(data,fmt.Sprintf("Content-Transfer-Encoding: base64%v",STMP_LFRT))
		data = append(data,STMP_LFRT)
		data = append(data,fmt.Sprintf("%v%v",base64.StdEncoding.EncodeToString([]byte(r.Html)),STMP_LFRT))
	}
	if len(r.Files) > 0{
		for i := range r.Files{
			filePath := r.Files[i]
			fData,e := ioutil.ReadFile(filePath)
			if e != nil{
				continue
			}
			data = append(data,STMP_LFRT)
			data = append(data,fmt.Sprintf("--%v%v",boundary,STMP_LFRT))
			bName := path.Base(filePath)
			data = append(data,fmt.Sprintf("Content-Type: application/octet-stream; name=\"=?utf-8?B?%v?=%v",bName,STMP_LFRT))
			data = append(data,fmt.Sprintf("Content-Disposition: attachment; filename=\"=?utf-8?B?%v?=%v",bName,STMP_LFRT))
			data = append(data,fmt.Sprintf("Content-Transfer-Encoding: base64%v",STMP_LFRT))
			data = append(data,STMP_LFRT)
			data = append(data,fmt.Sprintf("%v%v",base64.StdEncoding.EncodeToString(fData),STMP_LFRT))
		}
	}
	data = append(data,STMP_LFRT)
	data = append(data,fmt.Sprintf("--%v--%v",boundary,STMP_LFRT))

	fmt.Printf("生成邮件内容：%v\n",data)
	sdomain := strings.Split(r.ToEmail,"@")
	if len(sdomain) < 2{
		return
	}
	domain := sdomain[1]
	lst,e := net.LookupMX(domain)
	if e!=nil{
		fmt.Printf("查询服务器失败：%v %v\n" ,domain,e)
		return
	}

	for i:= range lst{
		mx := lst[i]
		con,e := net.DialTimeout("tcp",fmt.Sprintf("%v:25",mx.Host),time.Second * 10)
		if e!= nil{
			fmt.Printf("链接失败：%v %v\n" ,mx.Host,e)
			continue
		}
		red := bufio.NewReader(con)
		red.ReadLine()
		con.Write([]byte(fmt.Sprintf("EHLO %v%v",domain,STMP_LFRT)))
		b,_,_:=red.ReadLine()
		fmt.Printf("EHLO 接收到内容：%v\n" ,string(b))
		con.Write([]byte(fmt.Sprintf("MAIL From:<%v>%v",r.FromEmail,STMP_LFRT)))
		b1,_,_:=red.ReadLine()
		fmt.Printf("MAIL From 接收到内容：%v\n" ,string(b1))
		con.Write([]byte(fmt.Sprintf("RCPT To:<%v>%v",r.ToEmail,STMP_LFRT)))
		b2,_,_:=red.ReadLine()
		fmt.Printf("RCPT To 接收到内容：%v\n" ,string(b2))
		con.Write([]byte(fmt.Sprintf("DATA%v",STMP_LFRT)))
		b3,_,_:=red.ReadLine()
		fmt.Printf("Data 接收到内容：%v\n" ,string(b3))
		con.Write([]byte(fmt.Sprintf("%v",STMP_LFRT)))
		con.Write([]byte(fmt.Sprintf("%v%v",strings.Join(data,""),STMP_LFRT)))
		b5,_,_:=red.ReadLine()
		fmt.Printf("Body 接收到内容：%v\n" ,string(b5))
		con.Write([]byte(fmt.Sprintf("%v",STMP_LFRT)))
		con.Write([]byte(fmt.Sprintf(".%v",STMP_LFRT)))
		b7,_,_:=red.ReadLine()
		fmt.Printf(".接收到内容：%v\n" ,string(b7))
		con.Write([]byte(fmt.Sprintf("QUIT%v",STMP_LFRT)))
		con.Write([]byte(fmt.Sprintf("%v",STMP_LFRT)))
		con.Close()
		fmt.Printf("接收到内容：发送完毕\n")
		return
	}
}
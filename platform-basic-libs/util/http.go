package util

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	. "github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

//获取真实的IP  1.1.1.1, 2.2.2.2, 3.3.3.3
func CtxClientIP(ctx *fasthttp.RequestCtx) string {
	clientIP := Bytes2str(ctx.Request.Header.Peek("X-Forwarded-For"))
	if index := strings.IndexByte(clientIP, ','); index >= 0 {
		clientIP = clientIP[0:index]
		//获取最开始的一个 即 1.1.1.1
	}
	clientIP = strings.TrimSpace(clientIP)
	if len(clientIP) > 0 {
		return clientIP
	}
	clientIP = strings.TrimSpace(Bytes2str(ctx.Request.Header.Peek("X-Real-Ip")))
	if len(clientIP) > 0 {
		return clientIP
	}
	return ctx.RemoteIP().String()
}

// PostJSON POST请求 BODY为JSON格式 ContentType=application/json
func PostJSON(URL string, v interface{}) ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(URL, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// WriteJSON 写入json字符串
func WriteJSON(ctx *fasthttp.RequestCtx, v interface{}) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	ctx.Response.Header.Add("Content-Type", "application/json")
	ctx.Response.SetBody(b)
	return nil
}

var regIPv4 = regexp.MustCompile(
	`^(((\d{1,2})|(1\d{1,2})|(2[0-4]\d)|(25[0-5]))\.){3}((\d{1,2})|(1\d{1,2})|(2[0-4]\d)|(25[0-5]))$`,
)

func GetToken(ctx *Ctx) (token string) {
	return ctx.Get("X-Token")
}

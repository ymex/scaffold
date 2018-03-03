package http

import (
	"testing"
	"fmt"
)

func TestHttpDo(t *testing.T) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	url := "https://sms-api.luosimao.com/v1/send.json"
	values := make(map[string]string)
	values["mobile"] = "15512340000"
	values["message"] = "测试短信发送接口。【铁壳网络】"
	result, err := HttpDo("POST", headers, url, values, "api", "key-see")
	fmt.Println(string(result), err)

}

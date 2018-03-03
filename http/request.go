package http

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func HttpGet(uri string) ([]byte, error) {
	resp, err := http.Get(uri)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func HttpPostForm(requestUrl string, values map[string]string) ([]byte, error) {

	var vs = url.Values{}
	for key, v := range values {
		vs[key] = []string{v}
	}

	resp, err := http.PostForm(requestUrl, vs)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil

}

// http 请求
// method 请求方法
// headers 请求头
// requestUrl 请求地址
// values 参数
// authName BasicAuth
// authPassword BasicAuth
func HttpDo(method string, headers map[string]string, requestUrl string,
	values map[string]string, authName string, authPassword string) ([]byte, error) {

	var vs = url.Values{}
	for key, v := range values {
		vs[key] = []string{v}
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, requestUrl, strings.NewReader(vs.Encode()))

	if err != nil {
		return nil, err
	}

	for key, v := range headers {
		req.Header.Set(key, v)
	}
	req.SetBasicAuth(authName, authPassword)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

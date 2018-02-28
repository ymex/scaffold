package http

import (
	"io/ioutil"
	"net/http"
	"net/url"
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

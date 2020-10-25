package helper

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// HTTPJsonRequest Json 请求体添加 json
func HTTPJsonRequest(url, method string, body []byte, header http.Header) ([]byte, error) {
	if !IsURL(url) {
		return nil, fmt.Errorf("不合法的 URL 地址 %s", url)
	}
	method = strings.ToUpper(method)
	if method != "POST" && method != "PUT" {
		return nil, fmt.Errorf("请输入正确的请求方法")
	}
	c := new(http.Client)
	r, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	header.Set("Content-Type", "application/json")
	r.Header = header
	resp, err := c.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// HTTPFormRequest 表单请求
func HTTPFormRequest(url, method string, v url.Values, header http.Header) ([]byte, error) {
	if !IsURL(url) {
		return nil, fmt.Errorf("不合法的 URL 地址 %s", url)
	}
	method = strings.ToUpper(method)
	if method != "GET" && method != "POST" && method != "DELETE" && method != "PUT" {
		return nil, fmt.Errorf("请输入正确的请求方法")
	}
	c := new(http.Client)
	var (
		r   *http.Request
		err error
	)
	switch method {
	case "GET", "DELETE":
		r, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
		r.URL.RawQuery = v.Encode()
	case "POST", "PUT":
		r, err = http.NewRequest(method, url, bytes.NewReader([]byte(v.Encode())))
		if err != nil {
			return nil, err
		}
	}
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header = header
	resp, err := c.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

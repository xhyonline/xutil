package helper

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
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
	// 证书信任
	c := new(http.Client)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c.Transport = tr
	r, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if header == nil {
		header = http.Header{}
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
	// 证书信任
	c := new(http.Client)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c.Transport = tr
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
	if header == nil {
		header = http.Header{}
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

// HTTPRequest HTTP 请求
func HTTPRequest(url, method string, header http.Header, body io.Reader, times uint8) ([]byte, error) {
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, nil
	}
	// 重试机制
	client := &http.Client{}
	var count uint8
	for times != count {
		resp, err := client.Do(r)
		if err != nil {
			count++
			continue
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return b, nil
	}
	return nil, fmt.Errorf("超过重试次数")
}

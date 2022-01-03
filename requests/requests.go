package requests

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/xhyonline/xutil/helper"
)

// HTTPJsonRequest Json 请求体添加 json
func HTTPJsonRequest(url, method string, body []byte, header http.Header) ([]byte, error) {
	if !helper.IsURL(url) {
		return nil, fmt.Errorf("不合法的 URL 地址 %s", url)
	}
	method = strings.ToUpper(method)
	if method != "POST" && method != "PUT" {
		return nil, fmt.Errorf("请输入正确的请求方法")
	}
	// 证书信任
	c := new(http.Client)
	tr := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
	}
	c.Transport = tr
	c.Timeout = time.Second * 10
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
	if !helper.IsURL(url) {
		return nil, fmt.Errorf("不合法的 URL 地址 %s", url)
	}
	method = strings.ToUpper(method)
	if method != "GET" && method != "POST" && method != "DELETE" && method != "PUT" {
		return nil, fmt.Errorf("请输入正确的请求方法")
	}
	// 证书信任
	c := new(http.Client)
	tr := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
	}
	c.Timeout = time.Second * 10
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
func HTTPRequest(url, method string, header http.Header, body io.Reader) ([]byte, error) {
	// validate
	if url == "" || method == "" || header == nil {
		return nil, fmt.Errorf("请检查请求参数")
	}
	method = strings.ToUpper(method)
	client := new(http.Client)
	tr := &http.Transport{
		DisableKeepAlives: true,
	}
	client.Transport = tr
	client.Timeout = time.Second * 10
	r, err := http.NewRequestWithContext(context.Background(), method, url, body)
	if err != nil {
		return nil, err
	}
	r.Header = header
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

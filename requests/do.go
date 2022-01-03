package requests

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

type DoContext interface {
	// 获取响应体
	GetBody() []byte
	// 获取响应体(字符串)
	GetBodyString() string
	// 获取响应头
	GetHeader() http.Header
	// 获取状态码
	GetCode() int
	// Location returns the URL of the response's "Location" header,
	// if present. Relative redirects are resolved relative to
	// the Response's Request. ErrNoLocation is returned if no
	// Location header is present.
	Location() *url.URL
	// Cookies parses and returns the cookies set in the Set-Cookie headers.
	Cookies() []*http.Cookie
}

type doResponse struct {
	body     []byte
	header   http.Header
	code     int
	length   int64
	location *url.URL
	cookie   []*http.Cookie
}

// GetBody 请求体
func (s *doResponse) GetBody() []byte {
	return s.body
}

// GetBodyString 请求内容
func (s *doResponse) GetBodyString() string {
	return string(s.body)
}

// GetCode 响应码
func (s *doResponse) GetCode() int {
	return s.code
}

// GetHeader 响应头
func (s *doResponse) GetHeader() http.Header {
	return s.header
}

// GetContentLength
func (s *doResponse) GetContentLength() int64 {
	return s.length
}

// Location returns the URL of the response's "Location" header,
// if present. Relative redirects are resolved relative to
// the Response's Request. ErrNoLocation is returned if no
// Location header is present.
func (s *doResponse) Location() *url.URL {
	return s.location
}

// Cookies parses and returns the cookies set in the Set-Cookie headers.
func (s *doResponse) Cookies() []*http.Cookie {
	return s.cookie
}

type Option func(client *http.Client)

// WithTimeout 设置超时时间
func WithTimeout(timout time.Duration) Option {
	return func(client *http.Client) {
		client.Timeout = timout
	}
}

// Do 执行 HTTP 请求,对响应结构的封装,不再需要 close body
func Do(req *http.Request, fn func(c DoContext) error, options ...Option) error {
	tr := &http.Transport{
		DisableKeepAlives: true,
	}
	client := new(http.Client)
	client.Transport = tr
	client.Timeout = time.Second * 5 // 默认超时时间
	for _, f := range options {
		f(client)
	}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	locationHeader, err := resp.Location()
	if err != nil && err != http.ErrNoLocation {
		return err
	}
	return fn(&doResponse{
		body:     body,
		header:   resp.Header,
		code:     resp.StatusCode,
		length:   resp.ContentLength,
		location: locationHeader,
		cookie:   resp.Cookies(),
	})
}

// NSQ 消息中间件工具包
package mq

import (
	"time"
)

// Provider 服务提供商
type Provider int

// all providers
const (
	ProviderNSQ Provider = iota
)

// Config 配置 兼容各种服务
type Config struct {
	PubHost string `default:"nsqd"`
	PubTCP  string `default:"4150"`
	PubHTTP string `default:"4151"`
	SubHost string `default:"nsqlookupd"`
	SubTCP  string `default:"4160"` // 订阅走的统一是 lookup
	SubHTTP string `default:"4161"`
}

// Client xobj client
type Client interface {
	// 发布消息
	Pub(topic string, payload interface{}) error
	// 延迟发布消息
	Delay(topic string, payload interface{}, delay time.Duration) error
	// 订阅
	Sub(topic, channel string, f HandlerFunc)
	// 创建主题
	CreateTopic(topic string) error
	// 清理网络资源
	Close()
}

// New 新建存储客户端，为了混用不同的基础施舍，供应商和bucket在调用时填写，不放在设置中。
func New(provider Provider, config Config) Client {
	switch provider {
	case ProviderNSQ:
		return newNSQClient(config)
	default:
		panic("invalid provider")
	}
}

// Context 继承了队列接收消息后的上下文，包括消息内容和对消息的一些控制。
type Context interface {
	// Bind binds the payload body into provided type `i`. The default binder
	// is based on json.
	Bind(i interface{}) error
	// Data show the origin payload data body in message
	Data() []byte
	// String convert the origin payload data body to string
	String() string
}

// HandlerFunc 订阅者处理消息的函数
type HandlerFunc func(c Context) error

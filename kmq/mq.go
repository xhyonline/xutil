package kmq

import "time"

type Address struct {
	// 主机地址
	Host string
	// 端口地址
	Port string
}

// Config 卡夫卡集群配置信息
type Config struct {
	Address []Address
}

// Client 客户端
type Client interface {
	// 发布消息
	Pub(topic, key string, payload interface{}) error
	// 订阅
	Sub(topic, group string, f HandlerFunc) error
	// 创建主题,
	CreateTopic(topic string, partition, replicas int) error
	// 查看所有主题
	GetTopics() ([]string, error)
	// 删除主题
	RemoveTopics(topics ...string) error
	// 创建消费者组
	CreateGroup(name string, topics ...string) error
	// 清理网络资源
	Close()
}

func NewClient(c Config) Client {
	return newKafka(c)
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
	// GetOffset 获取这条数据所在的偏移量
	GetOffset() int
	// GetPartition 获取改数据所在分区
	GetPartition() int
	// GetKey 获取 key
	GetKey() string
	// 获取时间
	GetTime() time.Time
}

// HandlerFunc 订阅者处理消息的函数
type HandlerFunc func(c Context) error

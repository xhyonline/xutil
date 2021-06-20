package nmq

import (
	"encoding/json"

	"time"

	"github.com/levigross/grequests"
	nsq "github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
	"github.com/xhyonline/xutil/xlog"
)

var log = xlog.Get().Debug()

// NSQ context
type nsqContext struct {
	msg *nsq.Message
}

// Bind extract payload data to v
func (c *nsqContext) Bind(v interface{}) error {
	return json.Unmarshal(c.msg.Body, v)
}

// Data origin payload data
func (c *nsqContext) Data() []byte {
	return c.msg.Body
}

// String convert origin payload data to string
func (c *nsqContext) String() string {
	return string(c.msg.Body)
}

// NSQ 客户端
type nsqClient struct {
	// 只有一个生产者
	producer *nsq.Producer
	// 记录所有的消费者，结束时stop
	consumers []*nsq.Consumer
	config    Config
}

// 新建客户端
func newNSQClient(config Config) Client {
	var err error

	var c = &nsqClient{
		config: config,
	}

	c.producer, err = nsq.NewProducer(config.PubHost+":"+config.PubTCP, nsq.NewConfig())
	if err != nil {
		logrus.WithError(err).Panic("init nsq producer failed")
	}
	c.producer.SetLogger(NewLogrusLoggerAtLevel(logrus.WarnLevel))
	log.Info("NSQ Producer 初始化完成。")
	c.consumers = make([]*nsq.Consumer, 0)
	return c
}

// 为了对外隐藏nsq包的对象，把handler转换成使用context接口
func decorate(f HandlerFunc) nsq.HandlerFunc {
	return func(msg *nsq.Message) error {
		c := &nsqContext{
			msg: msg,
		}
		return f(c)
	}
}

// Reg 注册一个消费者处理函数 func(msg *nsq.Message) error
func (c *nsqClient) Sub(topic, channel string, handler HandlerFunc) {
	q, err := nsq.NewConsumer(topic, channel, nsq.NewConfig())
	if err != nil {
		log.WithError(err).Panic("init nsq consumer failed")
	}
	c.consumers = append(c.consumers, q)
	q.AddHandler(decorate(handler))
	q.SetLogger(NewLogrusLoggerAtLevel(logrus.WarnLevel))
	err = q.ConnectToNSQLookupd(c.config.SubHost + ":" + c.config.SubHTTP)
	if err != nil {
		log.WithError(err).Panic("nsq consumer connect to lookupd failed")
	}
	log.Infof("订阅 nsq topic %s by %s", topic, channel)
}

// Pub 发布消息，用json编码
// 因为 json Marshal 一个字符串不会报错，所以也可传入字符串，接收时使用 String() 方法接收
func (c *nsqClient) Pub(topic string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return c.producer.Publish(topic, data)
}

// Delay 延迟发布消息，用json编码
func (c *nsqClient) Delay(topic string, payload interface{}, delay time.Duration) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return c.producer.DeferredPublish(topic, delay, data)
}

// CreateTopic create a topic on nsqd by http request
func (c *nsqClient) CreateTopic(topic string) error {
	_, err := grequests.Post("http://"+c.config.PubHost+":"+c.config.PubHTTP+"/topic/create",
		&grequests.RequestOptions{Params: map[string]string{"topic": topic}})
	if err != nil {
		log.WithError(err).Errorf("创建主题 %s 出错", topic)
		return err
	}
	log.Infof("创建主题 %s 成功", topic)
	return nil
}

// Close graceful shutdown
func (c *nsqClient) Close() {
	for _, consumer := range c.consumers {
		consumer.Stop()
	}
	c.producer.Stop()
}

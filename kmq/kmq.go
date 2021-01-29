// kmq 包是消息中间件 kafka 的工具包

package kmq

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"context"

	"github.com/segmentio/kafka-go"
	"github.com/xhyonline/xutil/xlog"
)

var log = xlog.Get(true)

// consumerAttr 消费者属性
type consumerAttr struct {
	// 最大允许消费者个数,这是根据一个主题下的 partition 分片与之对应的,我们要保证同一个消费组下的每一个消费者对应一个分片
	maxAllowConsumerCount int
	// 当前消费者个数
	currentConsumerCount int
	// 记录该主题下的所有消费者
	consumerArr []*kafka.Reader
}

type kmq struct {
	// 卡夫卡总控
	conn *kafka.Conn
	// 有多个生产者,对应不同的主题 key 为主题名称 value 为 *kafka.Writer
	producer sync.Map
	// 配置信息
	config Config
	// 消费者集合与主题集合   key 为主题名 value 该主题下的消费者信息
	consumers map[string]consumerAttr
}

// CreateTopic 创建一个主题
// 请注意,主题消费者个数与 partition 分片对应
func (k *kmq) CreateTopic(topic string, partition, replicas int) error {
	configs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     partition,
			ReplicationFactor: replicas,
		},
	}
	return k.conn.CreateTopics(configs...)
}

// GetTopics 获取主题
func (k *kmq) GetTopics() ([]string, error) {
	p, err := k.conn.ReadPartitions()
	if err != nil {
		return nil, err
	}

	m := map[string]struct{}{}
	var topics = make([]string, 0)

	for _, v := range p {
		if _, ok := m[v.Topic]; ok {
			continue
		}
		m[v.Topic] = struct{}{}
		topics = append(topics, v.Topic)
	}

	return topics, nil
}

// CreateGroup 创建消费者组
func (k *kmq) CreateGroup(name string, topics ...string) error {
	_, err := kafka.NewConsumerGroup(kafka.ConsumerGroupConfig{
		ID:      name,
		Brokers: nil,
		Dialer: &kafka.Dialer{
			DialFunc: func(ctx context.Context, network string, address string) (conn net.Conn, e error) {
				return k.conn, nil
			},
		},
		Topics: topics,
	})
	return err
}

// Pub 推送数据
func (k *kmq) Pub(topic string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	// 涉及到多个推送
	producer, _ := k.producer.LoadOrStore(topic, &kafka.Writer{
		Addr:  kafka.TCP(k.config.Host + ":" + k.config.Port),
		Topic: topic,
	})
	return producer.(*kafka.Writer).WriteMessages(context.Background(), kafka.Message{Value: body})
}

// Sub 订阅
func (k *kmq) Sub(topic, group string, handle HandlerFunc) error {
	// 在创建消费者前,我们先要判断该主题有多少 partition
	v, ok := k.consumers[topic]

	if ok && v.maxAllowConsumerCount < v.currentConsumerCount+1 {
		return fmt.Errorf("该主题下的消费者个数超出分片个数。当前消费者个数为: %d 个 分片个数为 %d",
			v.currentConsumerCount, v.maxAllowConsumerCount)
	}
	// 创建消费者
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{k.config.Host + ":" + k.config.Port},
		GroupID:  group,
		Topic:    topic,
		MinBytes: 1,    //
		MaxBytes: 10e6, // 10MB
	})

	if ok {
		v.currentConsumerCount++
		v.consumerArr = append(v.consumerArr, r)
		go handler(r, handle)
		return nil
	}

	ps, err := k.conn.ReadPartitions(topic)
	if err != nil {
		return fmt.Errorf("订阅时,检查该主题下可允许的最大消费者个数失败 %s", err)
	}

	k.consumers = make(map[string]consumerAttr)
	k.consumers[topic] = consumerAttr{
		maxAllowConsumerCount: len(ps),
		currentConsumerCount:  1,
		consumerArr:           make([]*kafka.Reader, 0),
	}

	go handler(r, handle)
	return nil

}

// RemoveTopics 删除主题
func (k *kmq) RemoveTopics(topic ...string) error {
	return k.conn.DeleteTopics(topic...)
}

// Close 清理资源
func (k *kmq) Close() {
	k.producer.Range(func(key, value interface{}) bool {
		_ = value.(*kafka.Writer).Close()
		return true
	})
	for _, v := range k.consumers {
		for _, item := range v.consumerArr {
			_ = item.Close()
		}
	}
}

func handler(r *kafka.Reader, handlerFunc HandlerFunc) {
	defer r.Close()
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("kafka 读取数据时发生错误 %s", err)
		}
		err = handlerFunc(m)
		if err != nil {
			log.Fatalf("%s", err)
		}
	}
}

// newKafka 新建卡夫卡实例
func newKafka(c Config) *kmq {
	var mq = new(kmq)
	for {
		conn, err := kafka.Dial("tcp", c.Host+":"+c.Port)
		if err != nil {
			log.Errorf("kafka 连接失败 %s", err)
			time.Sleep(time.Second)
			continue
		}
		mq.conn = conn
		mq.config = c
		return mq
	}
}

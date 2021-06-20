// 工具包 创建 Redis

package kv

import (
	"C"
	"time"

	"github.com/vmihailenco/msgpack/v4"
	"github.com/xhyonline/xutil/xtype"

	"github.com/go-redis/cache/v7"
	"github.com/go-redis/redis/v7"
	"github.com/xhyonline/xutil/xlog"
)

var log = xlog.Get().Debug()

// Config 数据库配置，可以被主配置直接引用
type Config struct {
	Host     string `default:"redis"`
	Port     string `default:"6379"`
	Password string
	DB       int `default:"0"`
	// 连接池相关配置如下,如果你不用连接池,默认 0 就好
	PoolSize     int `default:"10"` // 连接池的数量,官方默认设定是 10
	MinIdleConns int // 最小空闲数量,如果你不想使用连接池,这个值设定为 0 就好,假设你设置了某个值,它会提前建立对应的连接个数放在池子里
}

// New 用配置生成一个 redis 数据库 RClient,若目标数据库未启动会一直等待
func New(config Config) *redis.Client {
	var kv = redis.NewClient(&redis.Options{
		Addr:         config.Host + ":" + config.Port,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
	})

	for {
		pong, err := kv.Ping().Result()
		if err == redis.Nil {
			log.Info("Redis error")
		} else if err != nil {
			log.Info("Redis connect error:", err)
		} else {
			log.Info(pong)
			break
		}
		time.Sleep(time.Second * 3)
	}

	log.Info("Redis connect successful.")

	return kv
}

// Client is a cache RClient interface
type Client interface {
	// Set 写缓存
	Set(key string, object interface{}, exp time.Duration)
	// MustSet 写缓存,并检查错误
	MustSet(key string, object interface{}, exp time.Duration) error
	// Get 读缓存
	Get(key string, pointer interface{}) error
	// SetString 写 string 缓存
	SetString(key string, s string, exp time.Duration) error
	// GetString 读 string 缓存
	GetString(key string) (string, error)
	// Exists 是否存在
	Exists(key string) bool
	// Expire 刷新过期时间
	Expire(key string, ex time.Duration) error
	// Delete 清缓存
	Delete(key string)
	// Clean 批量清除一类缓存
	Clean(cate string)
	// LGet 获取列表
	LGet(key string) (xtype.Strings, error)
	// LPush 为列表右侧增加一个元素
	LPush(key, item string) error
	// LLen 获取列表长度
	LLen(list string) (int64, error)
	// LPop 左侧弹出一个元素
	LPop(key string) (string, error)
	// RPop 右侧弹出一个元素
	RPop(key string) (string, error)
	// LRemove 删除列表中的指定元素
	LRemove(key, item string) error
	// HGet 获取 Hash 类型中的值
	HGet(key, filed string) (string, error)
	// HSet 设置一个 Hash 类型的值
	HSet(key, filed string) (bool, error, interface{})
	// HMSet 设置多个 Hash 值
	HMSet(key string, fields map[string]interface{}) (string, error)
	// HGetAll 获取所有的 Hash 返回 Map
	HGetAll(key string) (map[string]string, error)
	// 判断 Hash 是否存在某个 Key
	HExists(key, filed string) (bool, error)
	// 删除 Hash 中的多个元素
	HDel(key string, filed ...string) error
	// SGet 获取集合
	SGet(key string) (xtype.Strings, error)
	// SAdd 为集合增加一个元素
	SAdd(key, item string) error
	// SAdd 为集合增加一个元素，并刷新过期时间
	SAddEx(key, item string, ex time.Duration) error
	// SRemove 删除集合中的指定元素
	SRemove(key, item string) error
}

// RClient Redis 客户端
type RClient struct {
	Kv    *redis.Client // TODO 可以通过 Kv.Conn() 方法获取连接池中的一个连接,通过 close 后自动放回连接池
	codec *cache.Codec
}

// Set 写缓存
func (c *RClient) Set(key string, object interface{}, exp time.Duration) error {
	err := c.codec.Set(&cache.Item{
		Key:        key,
		Object:     object,
		Expiration: exp,
	})
	if err != nil {
		return err
	}
	return nil
}

// MustSet 写缓存,检查并返回错误
func (c *RClient) MustSet(key string, object interface{}, exp time.Duration) error {
	return c.codec.Set(&cache.Item{
		Key:        key,
		Object:     object,
		Expiration: exp,
	})
}

// Get 读缓存
func (c *RClient) Get(key string, pointer interface{}) error {
	return c.codec.Get(key, pointer)
}

// Set 写 string 缓存,当过期时间为 0 时,就是永久
func (c *RClient) SetString(key string, s string, exp time.Duration) error {
	return c.Kv.Set(key, s, exp).Err()
}

// Get 读 string 缓存
func (c *RClient) GetString(key string) (string, error) {
	return c.Kv.Get(key).Result()
}

// Exists 是否存在
func (c *RClient) Exists(key string) bool {
	return c.Kv.Exists(key).Val() != 0
}

// Delete 清缓存
func (c *RClient) Delete(key string) {
	err := c.codec.Delete(key)
	if err == cache.ErrCacheMiss {
		return
	} else if err != nil {
		log.WithError(err).WithField("key", key).Error("delete cache failed")
	}
}

// Expire 刷新过期时间
func (c *RClient) Expire(key string, ex time.Duration) error {
	return c.Kv.Expire(key, ex).Err()
}

// Clean 批量清除一类缓存
func (c *RClient) Clean(cate string) {
	if cate == "" {
		log.Error("someone try to clean all cache keys")
		return
	}
	i := 0
	for _, key := range c.Kv.Keys(cate + "*").Val() {
		err := c.codec.Delete(key)
		if err != nil {
			log.WithError(err).WithField("key", key).Error("delete cache failed,stop batch delete")
			break
		}
		i++
	}
	log.Infof("delete %d %s cache", i, cate)
}

// LPush 为列表右侧增加一个元素
func (c *RClient) LPush(list, item string) error {
	cmd := c.Kv.LPush(list, item)
	return cmd.Err()
}

// LRemove 删除列表中所有的指定元素 从表头开始向表尾搜索
func (c *RClient) LRemove(list, item string) error {
	cmd := c.Kv.LRem(list, 0, item)
	return cmd.Err()
}

// LPop 左侧弹出一个元素
func (c *RClient) LPop(list string) (string, error) {
	cmd := c.Kv.LPop(list)
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	str := cmd.Val()
	return str, nil
}

// RPop 右侧弹出一个元素
func (c *RClient) RPop(list string) (string, error) {
	cmd := c.Kv.RPop(list)

	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	str := cmd.Val()
	return str, nil
}

// LLen 返回列表长度
func (c *RClient) LLen(list string) (int64, error) {
	cmd := c.Kv.LLen(list)

	if cmd.Err() != nil {
		return 0, cmd.Err()
	}
	return cmd.Val(), nil
}

// LGet 获取列表所有元素
func (c *RClient) LGet(list string) (xtype.Strings, error) {
	cmd := c.Kv.LRange(list, 0, -1)

	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return xtype.Strings(cmd.Val()), nil
}

// HGet 获取 Hash 类型的值
func (c *RClient) HGet(key, filed string) (string, error) {
	cmd := c.Kv.HGet(key, filed)

	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return cmd.Val(), nil
}

// HSet 设置 Hash 类型的值 注: interface 类型别传一个指针结构体....它是解析不了的
func (c *RClient) HSet(key, filed string, value interface{}) (bool, error) {
	cmd := c.Kv.HSet(key, filed, value)

	if cmd.Err() != nil {
		return false, cmd.Err()
	}
	return true, nil
}

// HMSet 设置多个 Hash ,如果成功返回字符串 OK
func (c *RClient) HMSet(key string, fields map[string]interface{}) (string, error) {
	cmd := c.Kv.HMSet(key, fields)

	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return "OK", nil
}

// HGetAll 获取所有的 Hash
func (c *RClient) HGetAll(key string) (map[string]string, error) {
	cmd := c.Kv.HGetAll(key)

	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return cmd.Val(), nil
}

// HExists 判断 Hash 中某个 key 是否存在
func (c *RClient) HExists(key, field string) (bool, error) {
	cmd := c.Kv.HExists(key, field)

	if cmd.Err() != nil {
		return false, cmd.Err()
	}
	return cmd.Val(), nil
}

// HDel Hash 删除
func (c *RClient) HDel(key, field string) (int64, error) {
	cmd := c.Kv.HDel(key, field)

	if cmd.Err() != nil {
		return 0, cmd.Err()
	}
	return cmd.Val(), nil
}

// ============= 连接池相关 ====================

// PoolStatus 查看连接池状态
func (c *RClient) PoolStatus() *redis.PoolStats {
	return c.Kv.PoolStats()
}

// 暂时就不多加取和放回的操作
// 请使用 conn:=client.Kv.Conn 从连接池中获取连接 conn.Close() 放回连接池去

// =========== 哨兵相关 ====================

type Sentinel struct {
	MasterName   string // 哨兵配置 master 节点的名字
	SentinelIP   string
	SentinelPort string
}

// SentinelGetClient 通过哨兵,直接主节点客户端
func SentinelGetClient(config *Sentinel) (*RClient, error) {

	sentinel := redis.NewSentinelClient(&redis.Options{
		Network: "tcp",
		Addr:    config.SentinelIP + ":" + config.SentinelPort,
	})
	cmd := sentinel.GetMasterAddrByName(config.MasterName)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	// 获取主节点的信息
	info, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	master := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    info[0] + ":" + info[1],
	})
	client := &RClient{
		Kv: master,
	}
	client.codec = &cache.Codec{
		Redis: client.Kv,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
	return client, nil
}

// NewRedisClient 新建一个 Redis 客户端 它使用的是 go-redis 包
func NewRedisClient(config Config) *RClient {
	c := &RClient{
		Kv: New(config),
	}
	// 必须等到 redis 建立完毕
	c.codec = &cache.Codec{
		Redis: c.Kv,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
	return c
}

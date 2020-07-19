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

// 工具包 创建 Redis
var log = xlog.Get()

// Config 数据库配置，可以被主配置直接引用
type Config struct {
	Host     string `default:"redis"`
	Port     string `default:"6379"`
	Password string
	DB       int `default:"0"`
}

// New 用配置生成一个 redis 数据库 client,若目标数据库未启动会一直等待
func New(config Config) *redis.Client {
	var kv = redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.DB,
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

// Client is a cache client interface
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

// client Redis 客户端
type client struct {
	kv    *redis.Client
	codec *cache.Codec
}

// Set 写缓存
func (c *client) Set(key string, object interface{}, exp time.Duration) {
	err := c.codec.Set(&cache.Item{
		Key:        key,
		Object:     object,
		Expiration: exp,
	})
	if err != nil {
		log.WithError(err).WithField("key", key).Error("set cache failed")
	}
}

// MustSet 写缓存,检查并返回错误
func (c *client) MustSet(key string, object interface{}, exp time.Duration) error {
	return c.codec.Set(&cache.Item{
		Key:        key,
		Object:     object,
		Expiration: exp,
	})
}

// Get 读缓存
func (c *client) Get(key string, pointer interface{}) error {
	return c.codec.Get(key, pointer)
}

// Set 写 string 缓存
func (c *client) SetString(key string, s string, exp time.Duration) error {
	return c.kv.Set(key, s, exp).Err()
}

// Get 读 string 缓存
func (c *client) GetString(key string) (string, error) {
	return c.kv.Get(key).Result()
}

// Exists 是否存在
func (c *client) Exists(key string) bool {
	return c.kv.Exists(key).Val() != 0
}

// Delete 清缓存
func (c *client) Delete(key string) {
	err := c.codec.Delete(key)
	if err == cache.ErrCacheMiss {
		return
	} else if err != nil {
		log.WithError(err).WithField("key", key).Error("delete cache failed")
	}
}

// Expire 刷新过期时间
func (c *client) Expire(key string, ex time.Duration) error {
	return c.kv.Expire(key, ex).Err()
}

// Clean 批量清除一类缓存
func (c *client) Clean(cate string) {
	if cate == "" {
		log.Error("someone try to clean all cache keys")
		return
	}
	i := 0
	for _, key := range c.kv.Keys(cate + "*").Val() {
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
func (c *client) LPush(list, item string) error {
	cmd := c.kv.LPush(list, item)
	return cmd.Err()
}

// LRemove 删除列表中所有的指定元素 从表头开始向表尾搜索
func (c *client) LRemove(list, item string) error {
	cmd := c.kv.LRem(list, 0, item)
	return cmd.Err()
}

// LPop 左侧弹出一个元素
func (c *client) LPop(list string) (string, error) {
	strCmd := c.kv.LPop(list)
	err := strCmd.Err()
	if err != nil {
		return "", err
	}
	str := strCmd.Val()
	return str, nil
}

// RPop 右侧弹出一个元素
func (c *client) RPop(list string) (string, error) {
	cmd := c.kv.RPop(list)
	err := cmd.Err()
	if err != nil {
		return "", err
	}
	str := cmd.Val()
	return str, nil
}

// LLen 返回列表长度
func (c *client) LLen(list string) (int64, error) {
	cmd := c.kv.LLen(list)
	err := cmd.Err()
	if err != nil {
		return 0, err
	}
	return cmd.Val(), nil
}

// LGet 获取列表所有元素
func (c *client) LGet(list string) (xtype.Strings, error) {
	cmd := c.kv.LRange(list, 0, -1)
	err := cmd.Err()
	if err != nil {
		return nil, err
	}
	return xtype.Strings(cmd.Val()), nil
}

// HGet 获取 Hash 类型的值
func (c *client) HGet(key, filed string) (string, error) {
	cmd := c.kv.HGet(key, filed)
	err := cmd.Err()
	if err != nil {
		return "", err
	}
	return cmd.Val(), nil
}

// HSet 设置 Hash 类型的值 注: interface 类型别传一个指针结构体....它是解析不了的
func (c *client) HSet(key, filed string, value interface{}) (bool, error) {
	cmd := c.kv.HSet(key, filed, value)
	err := cmd.Err()
	if err != nil {
		return false, err
	}
	return cmd.Val(), nil
}

// HMSet 设置多个 Hash ,如果成功返回字符串 OK
func (c *client) HMSet(key string, fields map[string]interface{}) (string, error) {
	cmd := c.kv.HMSet(key, fields)
	err := cmd.Err()
	if err != nil {
		return "", err
	}
	return cmd.Val(), nil
}

// HGetAll 获取所有的 Hash
func (c *client) HGetAll(key string) (map[string]string, error) {
	cmd := c.kv.HGetAll(key)
	err := cmd.Err()
	if err != nil {
		return nil, err
	}
	return cmd.Val(), nil
}

// HExists 判断 Hash 中某个 key 是否存在
func (c *client) HExists(key, field string) (bool, error) {
	cmd := c.kv.HExists(key, field)
	err := cmd.Err()
	if err != nil {
		return false, err
	}
	return cmd.Val(), nil
}

// HDel Hash 删除
func (c *client) HDel(key, field string) (int64, error) {
	cmd := c.kv.HDel(key, field)
	err := cmd.Err()
	if err != nil {
		return 0, err
	}
	return cmd.Val(), nil
}

// NewRedisClient 新建一个 Redis 客户端
func NewRedisClient(config Config) *client {
	c := &client{
		kv: New(config),
	}
	// 必须等到 redis 建立完毕
	c.codec = &cache.Codec{
		Redis: c.kv,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
	return c
}

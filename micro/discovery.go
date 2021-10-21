package micro

import (
	"encoding/json"
	"strings"

	"github.com/xhyonline/xutil/helper"

	"sync"

	"context"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/xhyonline/xutil/logger"
	"go.etcd.io/etcd/clientv3"
)

// MicroMicroServiceDiscovery 微服务发现实例
type MicroMicroServiceDiscovery struct {
	// etcd 服务实例
	client *clientv3.Client
	// etcd 前缀
	prefix string
	// 各服务节点, Key 为服务名 val 为节点
	nodes map[string][]*Node
	// 保证 nodes 并发安全
	lock sync.RWMutex
}

// NewMicroServiceDiscovery 实例化一个服务发现实例
func NewMicroServiceDiscovery(client *clientv3.Client, prefix string) *MicroMicroServiceDiscovery {
	s := &MicroMicroServiceDiscovery{
		client: client,
		prefix: "/" + strings.Trim(prefix, "/") + "/",
		nodes:  make(map[string][]*Node),
		lock:   sync.RWMutex{},
	}
	go s.watch()
	return s
}

// watch 监听事件
func (s *MicroMicroServiceDiscovery) watch() {
	// 先获取先前前缀下的所有服务
	kv := clientv3.KV(s.client)
	logger.Info("服务发现开始监听前缀为:" + s.prefix)
	getResp, err := kv.Get(context.Background(), s.prefix, clientv3.WithPrefix())
	if err != nil {
		logger.Errorf("etcd 监听服务发生错误,服务停止 %s", err)
		return
	}
	s.lock.Lock()
	// 将所有服务全部都扔进集合中
	for _, item := range getResp.Kvs {
		err = s.addService(item.Key, item.Value)
		if err != nil {
			logger.Warnf("发现不合规的节点key:%s value:", string(item.Key), string(item.Value))
			continue
		}
	}
	s.lock.Unlock()

	// 开始正式监听

	watcher := clientv3.NewWatcher(s.client)

	defer watcher.Close()
	// 从 getResp.Header.Revision+1 开始,监听后续所有的以 s.prefix 前缀开头的 key 事件
	c := watcher.Watch(context.Background(), s.prefix, clientv3.WithPrefix(), clientv3.WithRev(getResp.Header.Revision+1))
	// 从管道持续读取
	for watchResp := range c {
		for _, event := range watchResp.Events {
			s.lock.Lock()
			switch event.Type {
			case mvccpb.PUT:
				if err = s.addService(event.Kv.Key, event.Kv.Value); err != nil {
					s.lock.Unlock()
					// 跳过这个
					continue
				}
			case mvccpb.DELETE:
				s.removeService(event.Kv.Key)
			}
			s.lock.Unlock()
		}
	}
}

// addService 新增服务
func (s *MicroMicroServiceDiscovery) addService(key, value []byte) error {
	// 获取服务名
	name := s.getServerName(key)

	node := new(Node)
	if err := json.Unmarshal(value, node); err != nil {
		return err
	}
	if err := node.Validate(); err != nil {
		return err
	}
	nodes, ok := s.nodes[name]
	if !ok {
		nodes = make([]*Node, 0)
	}
	s.nodes[name] = append(nodes, node)
	return nil
}

// removeService 删除服务
func (s *MicroMicroServiceDiscovery) removeService(key []byte) {
	name := s.getServerName(key)
	nodes, ok := s.nodes[name]
	if !ok {
		return
	}
	address := s.getNodeByKey(key)
	for k, node := range nodes {
		if node.Host == address.Host && node.Port == address.Port {
			nodes = append(nodes[0:k], nodes[k+1:]...)
			break
		}
	}
	logger.Infof("服务发现删除节点 服务名:%s 节点地址:%s", name, address.Host+":"+address.Port)
	if len(nodes) == 0 {
		delete(s.nodes, name)
		return
	}
	s.nodes[name] = nodes
}

// GetService 获取服务
func (s *MicroMicroServiceDiscovery) GetService(name string) *Node {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if nodes, ok := s.nodes[name]; ok {
		return nodes[helper.GetRandom(len(nodes))]
	}
	return nil
}

// getServerName 根据规则获取服务名
func (s *MicroMicroServiceDiscovery) getServerName(key []byte) string {
	return strings.Split(strings.Replace(string(key), s.prefix, "", 1), "/")[0]
}

// getNodeByKey 通过 key 获取节点
func (s *MicroMicroServiceDiscovery) getNodeByKey(key []byte) Node {
	result := strings.Split(strings.Split(strings.Replace(string(key), s.prefix, "", 1), "/")[1], ":")
	return Node{
		Host: result[0],
		Port: result[1],
	}
}

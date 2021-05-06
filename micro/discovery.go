package micro

import (
	"encoding/json"

	"sync"

	"context"
	"github.com/coreos/etcd/mvcc/mvccpb"
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
	return &MicroMicroServiceDiscovery{
		client: client,
		prefix: prefix,
		nodes:  make(map[string][]*Node),
		lock:   sync.RWMutex{},
	}
}

// Watch 监听事件
func (s *MicroMicroServiceDiscovery) Watch() error {
	// 先获取先前前缀下的所有服务
	kv := clientv3.KV(s.client)
	getResp, err := kv.Get(context.Background(), s.prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	s.lock.Lock()
	// 将所有服务全部都扔进集合中
	for _, item := range getResp.Kvs {
		err = s.addService(item.Key, item.Value)
		if err != nil {
			return err
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
	return nil
}

// addService 新增服务
func (s *MicroMicroServiceDiscovery) addService(key, value []byte) error {
	node := new(Node)
	if err := json.Unmarshal(value, node); err != nil {
		return err
	}
	if err := node.Validate(); err != nil {
		return err
	}
	nodes, ok := s.nodes[string(key)]
	if !ok {
		nodes = make([]*Node, 0)
	}
	nodes = append(nodes, node)
	return nil
}

// removeService 删除服务
func (s *MicroMicroServiceDiscovery) removeService(key []byte) {
	nodes,ok:=s.nodes[string(key)]
	if !ok {
		return
	}
	//if len(nodes)==1 && nodes[0].Port== {
	//
	//}
	delete(s.nodes, string(key))
}

// GetService 获取服务
func (s *MicroMicroServiceDiscovery) GetService(name string) *Node {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if v, ok := s.nodes[name]; ok {
		return v
	}
	return nil
}

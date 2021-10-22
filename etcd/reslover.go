package etcd

import (
	"log"
	"sync"

	"github.com/xhyonline/xutil/micro"

	"github.com/xhyonline/xutil/logger"
	"go.etcd.io/etcd/clientv3"

	"google.golang.org/grpc/resolver"
)

var schema = "/micro/server/"

type Resolver struct {
	// etcd 客户端
	cli *clientv3.Client
	// 负载均衡器
	cc resolver.ClientConn
	// 服务列表
	serverList map[string]resolver.Address
	// 服务列表锁
	lock sync.Mutex
	// 服务发现组件实例
	discoverInstance *micro.MicroMicroServiceDiscovery
	// 服务名
	name string
}

// Build 实现了第三方方法 resolver.Register() 的入参接口 resolver.Builder
// 当客户端使用 grpc.Dial 时,将会自动触发该函数,有点像一个 hook 钩子
// 参数释义:
// target : 当客户端调用 grpc.Dial() 方法时,会将入参解析到 target 中,例如 grpc.Dial("dns://some_authority/foo.bar") 就会解析成  &Target{Scheme: "dns", Authority: "some_authority", Endpoint: "foo.bar"}
// cc 负载均衡器
func (s *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	logger.Infof("grpc 启动触发")
	s.cc = cc
	prefix := target.Scheme + target.Endpoint + "/"
	logger.Info("触发 " + prefix)
	s.discoverInstance = micro.NewMicroServiceDiscovery(s.cli, prefix)
	s.discoverInstance.SetAfterAddServiceHook(s.SetHook)
	s.discoverInstance.SetAfterRemoveServiceHook(s.RemoveHook)
	return s, nil
}

// SetServiceList 设置服务
func (s *Resolver) SetServiceList(key, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	// 写入服务
	s.serverList[key] = resolver.Address{Addr: val}
	// 写入完毕,将改地址加入负载均衡器
	s.cc.UpdateState(resolver.State{
		Addresses: s.GetServices(),
	})
	logger.Info("新增服务地址:" + val + " 并已经加入负载均衡器")
}

// DelServiceList 删除服务地址
func (s *Resolver) DelServiceList(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	addr := s.serverList[key]
	delete(s.serverList, key)
	// 删除完毕,更新负载均衡器
	s.cc.UpdateState(resolver.State{
		Addresses: s.GetServices(),
	})
	log.Println("删除服务地址:" + addr.Addr + " 并移除负载均衡器")
}

// GetServices 获取当前的服务列表
func (s *Resolver) GetServices() []resolver.Address {
	address := make([]resolver.Address, 0, len(s.serverList))
	for _, v := range s.serverList {
		address = append(address, v)
	}
	return address
}

// Close 实现 resolver.Resolver 的关闭接口
func (s *Resolver) Close() {
}

// Scheme 实现了第三方方法 resolver.Register() 的入参接口 resolver.Builder
func (s *Resolver) Scheme() string {
	return schema
}

// ResolveNow 实现第三方 resolver.Resolver 的接口,监视目标更新
func (s *Resolver) ResolveNow(rn resolver.ResolveNowOptions) {

}

// SetHook
func (s *Resolver) SetHook(key, _ string, node *micro.Node) {
	s.SetServiceList(key, node.Host+":"+node.Port)
}

// RemoveHook
func (s *Resolver) RemoveHook(key string) {
	s.DelServiceList(key)
}

func NewResolver(name string, client *clientv3.Client) resolver.Builder {
	return &Resolver{
		cli:        client,
		serverList: make(map[string]resolver.Address),
		lock:       sync.Mutex{},
		name:       name,
	}
}

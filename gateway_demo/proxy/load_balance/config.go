package load_balance

import (
	"fmt"

	"github.com/hwg1999/go_gateway/gateway_demo/proxy/zookeeper"
)

// 配置主体
type LoadBalanceConf interface {
	Attach(o Observer)
}

func (s *LoadBalanceZkConf) Attach(o Observer) {
	s.observers = append(s.observers, o)
}

func (s *LoadBalanceZkConf) NotifyAllObservers() {
	for _, obj := range s.observers {
		obj.Update()
	}
}

func (s *LoadBalanceZkConf) GetConf() []string {
	confList := []string{}
	for _, ip := range s.activeList {
		weight, ok := s.confIpWeight[ip]
		if !ok {
			weight = "50" // 默认权重
		}
		confList = append(confList, fmt.Sprintf(s.format, ip)+","+weight)
	}

	return confList
}

type Observer interface {
	Update()
}

type LoadBalanceObserver struct {
	ModuleConf *LoadBalanceZkConf
}

func (l *LoadBalanceObserver) Update() {
	fmt.Println("Update get conf:", l.ModuleConf.GetConf())
}

func NewLoadBalanceObserver(conf *LoadBalanceZkConf) *LoadBalanceObserver {
	return &LoadBalanceObserver{
		ModuleConf: conf,
	}
}

type LoadBalanceZkConf struct {
	observers    []Observer
	path         string
	zkHosts      []string
	confIpWeight map[string]string
	activeList   []string
	format       string
}

// 更新配置时，通过观察者也更新
func (s *LoadBalanceZkConf) WatchConf() {
	zkManager := zookeeper.NewZkManager(s.zkHosts)
	zkManager.GetConnect()
	fmt.Println("watchConf")
	chanList, chanErr := zkManager.WatchServerListByPath(s.path)
	go func() {
		defer zkManager.Close()
		for {
			select {
			case changeErr := <-chanErr:
				fmt.Println("changeErr", changeErr)
			case changedList := <-chanList:
				fmt.Println("watch node changed")
				s.UpdateConf(changedList)
			}
		}
	}()
}

// 更新配置时，通知观察者也更新
func (s *LoadBalanceZkConf) UpdateConf(conf []string) {
	s.activeList = conf
	for _, obs := range s.observers {
		obs.Update()
	}
}

func NewLoadBalanceZkConf(format, path string, zkHosts []string, conf map[string]string) (*LoadBalanceZkConf, error) {
	zkManager := zookeeper.NewZkManager(zkHosts)
	zkManager.GetConnect()
	defer zkManager.Close()
	zlist, err := zkManager.GetServerListByPath(path)
	if err != nil {
		return nil, err
	}
	mConf := &LoadBalanceZkConf{format: format, activeList: zlist, confIpWeight: conf, zkHosts: zkHosts, path: path}
	mConf.WatchConf()
	return mConf, nil
}

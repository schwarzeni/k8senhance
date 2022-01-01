package netproxy

import "sync"

// TODO: 这种是 cloud 为单机模式的，不利于多副本拓展
type ProxyChannelManager struct {
	data map[string]chan *ProxyRequestWrapper
	lock sync.RWMutex
}

func (sm *ProxyChannelManager) Set(id string, ch chan *ProxyRequestWrapper) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	sm.data[id] = ch
}

func (sm *ProxyChannelManager) Get(id string) (chan *ProxyRequestWrapper, bool) {
	sm.lock.RLock()
	defer sm.lock.RUnlock()
	v, ok := sm.data[id]
	return v, ok
}

func (sm *ProxyChannelManager) Del(id string) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	delete(sm.data, id)
}

func NewProxyChannelManager() *ProxyChannelManager {
	return &ProxyChannelManager{
		data: make(map[string]chan *ProxyRequestWrapper),
	}
}

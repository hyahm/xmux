package xmux

import "sync"

type cacheTable struct {
	cache map[string]*rt
	mu    *sync.RWMutex
}

var ctLocker *cacheTable

func init() {
	ctLocker = &cacheTable{
		cache: make(map[string]*rt),
		mu:    &sync.RWMutex{},
	}
}

func (ct cacheTable) Get(key string) (*rt, bool) {
	ct.mu.RLock()
	defer ct.mu.RUnlock()
	v, ok := ct.cache[key]
	return v, ok
}

func (ct cacheTable) Set(key string, value *rt) {
	ct.mu.Lock()
	ct.cache[key] = value
	defer ct.mu.Unlock()

}

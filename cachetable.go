package xmux

import "sync"

type cacheTable map[string]*rt

var ctLocker *sync.RWMutex

func init() {
	ctLocker = &sync.RWMutex{}
}

func (ct cacheTable) Get(key string) (*rt, bool) {
	ctLocker.RLock()
	defer ctLocker.RUnlock()
	value, ok := ct[key]
	if ok {
		return value, true
	} else {
		return nil, false
	}

}

func (ct cacheTable) Set(key string, value *rt) {
	ctLocker.Lock()
	ct[key] = value
	defer ctLocker.Unlock()

}

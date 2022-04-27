package cache

import "sync"

// 简单缓存， 先进后出

type SimpleCache struct {
	cache map[string][]byte
	order []string
	mu    sync.RWMutex
	size  int
}

func NewCache(max int) *SimpleCache {
	if max <= 0 {
		max = 10000
	}
	return &SimpleCache{
		cache: make(map[string][]byte),
		order: make([]string, max),
		mu:    sync.RWMutex{},
		size:  max,
	}
}

func (sc *SimpleCache) Add(key string, value []byte) (string, bool) {
	if _, ok := sc.cache[key]; ok {
		sc.cache[key] = value
		return key, false
	} else {
		sc.cache[key] = value
		sc.order = append(sc.order, key)
		if len(sc.order) >= sc.size {
			rmKey := sc.order[sc.size-1]
			sc.order = sc.order[1:]
			return rmKey, true
		}
		return key, false
	}
}

func (sc *SimpleCache) Get(key string) ([]byte, bool) {
	if v, ok := sc.cache[key]; ok {
		return v, ok
	}
	return nil, false
}

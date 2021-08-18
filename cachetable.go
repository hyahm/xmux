package xmux

import (
	"github.com/hyahm/lru"
)

var cache *lru.List

func InitCache(count uint64) {
	if count == 0 {
		count = 10000
	}
	cache = lru.Init(count)
}

func Get(key string) (*rt, bool) {
	if cache.Exsit(key) {
		return cache.Get(key).(*rt), true
	}
	return nil, false
}

func Set(key string, value *rt) {
	cache.Add(key, value)
}

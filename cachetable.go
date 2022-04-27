package xmux

import (
	"github.com/hyahm/cache"
)

var urlCache cache.Cacher

func initUrlCache(count int) {
	if count == 0 {
		count = 10000
	}
	urlCache = cache.NewCache(count, cache.LRU)
}

func getUrlCache(key string) (*rt, bool) {
	value := urlCache.Get(key)
	if value == nil {
		return nil, false
	}
	return value.(*rt), true
}

func setUrlCache(key string, value *rt) {
	urlCache.Add(key, value)
}

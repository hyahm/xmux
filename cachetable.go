package xmux

import (
	"github.com/hyahm/gocache"
)

var urlCache gocache.Cacher[string, *rt]

func initUrlCache(count int) {
	if count == 0 {
		count = 10000
	}
	urlCache = gocache.NewCache[string, *rt](count, gocache.LRU)
}

func getUrlCache(key string) (*rt, bool) {
	return urlCache.Get(key)
}

func setUrlCache(key string, value *rt) {
	urlCache.Add(key, value)
}

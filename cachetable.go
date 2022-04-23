package xmux

import "github.com/hyahm/lru"

var urlCache *lru.List

func initUrlCache(count uint64) {
	if count == 0 {
		count = 10000
	}
	urlCache = lru.Init(count)
}

func getUrlCache(key string) (*rt, bool) {
	if urlCache.Exsit(key) {
		return urlCache.Get(key).(*rt), true
	}
	return nil, false
}

func setUrlCache(key string, value *rt) {
	urlCache.Add(key, value)
}

package xmux

import "github.com/hyahm/xmux/cache"

const CacheKey = "CacheKey"

type Cacher[K string, V []byte] interface {
	Add(key K, value V) (K, bool) // 添加值， 如果返回 k, true 说明有删除值，并返回删除的key
	Remove(key K)                 // 移除k
	Len() int                     // 长度
	OrderPrint(int)               // 顺序打印
	Get(key K) (V, bool)          // 获取值
	LastKey() K                   // 获取最先要删除的key
}

func ResponseCache() *RouteGroup {
	cache.InitResponseCache()
	cache := NewRouteGroup()
	// cache.Post("/-/cache/size", list)
	// cache.Post("/-/cache/add", list)
	// cache.Post("/-/cache/del", list)
	// cache.Post("/-/cache/clear", list)
	return cache
}

package xmux

import "github.com/hyahm/xmux/cache"

const CacheKey = "CacheKey"

func ResponseCache() *RouteGroup {
	cache.InitResponseCache()
	cache := NewRouteGroup()
	// cache.Post("/-/cache/size", list)
	// cache.Post("/-/cache/add", list)
	// cache.Post("/-/cache/del", list)
	// cache.Post("/-/cache/clear", list)
	return cache
}

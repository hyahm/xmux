package xmux

import "github.com/hyahm/xmux/cache"

func ResponseCache() *GroupRoute {
	cache.InitResponseCache()
	cache := NewGroupRoute()
	// cache.Post("/-/cache/size", list)
	// cache.Post("/-/cache/add", list)
	// cache.Post("/-/cache/del", list)
	// cache.Post("/-/cache/clear", list)
	return cache
}

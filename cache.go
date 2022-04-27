package xmux

import (
	"fmt"
	"sync"
)

const CacheKey = "CacheKey"

var rc *responseCache

// 返回的缓存
type responseCache struct {
	store  Cacher
	status map[string]int // 0 说明是缓存   1  说明是正在更新   2： 说明需要更新
	mu     sync.RWMutex
}
type Cacher interface {
	Add(string, []byte) (string, bool)
	Get(string) ([]byte, bool)
}

// 设置缓存值
func SetCache(key string, value []byte) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	if _, ok := rc.status[key]; ok {
		rc.status[key] = 0
	}
	rk, ok := rc.store.Add(key, value)
	if ok {
		// 如果删除了值， 也要删除对应的update
		delete(rc.status, rk)
	}
}

// 如果key存在，就设置缓存
// func SetCacheIfExsits(key string, value []byte) {
// 	rc.mu.RLock()
// 	defer rc.mu.RUnlock()
// 	if _, ok := rc.status[key]; ok {
// 		rc.status[key] = false
// 		rc.store.Add(key, value)
// 	}
// }

// 获取缓存值, 如果不存在返回nil
func GetCache(key string) []byte {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	value, ok := rc.store.Get(key)
	if ok {
		return value
	}
	return nil
}

type CacheStatus string

const (
	NotFoundCache    CacheStatus = "Not found cache"
	CacheIsUpdateing CacheStatus = "Cache is Updating"
	CacheNeedUpdate  CacheStatus = "Cache need Updating"
	CacheHit         CacheStatus = "cache hit"
)

// 获取缓存，如果正在更新
// 如果返回 NotFoundCache    说明不存在这个缓存
// 如果返回 CacheIsUpdateing  说明当前还在更新中， 还不是最新的缓存
// 如果返回 CacheNeedUpdate  说明缓存需要更新
// 如果返回 CacheHit 说明是最新的，可以直接返回
func GetCacheIfUpdating(key string) ([]byte, CacheStatus) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	if status, ok := rc.status[key]; ok {
		// 判断是否存在，存在， 如果正在更新，并且值是nil
		value, gok := rc.store.Get(key)
		if !gok {
			fmt.Println("somthing wrong")
		}
		switch status {
		case 0:
			return value, CacheHit
		case 1:
			return nil, CacheIsUpdateing
		default:
			return nil, CacheNeedUpdate
		}

	} else {
		return nil, NotFoundCache
	}

}

// 是否存在缓存
func ExsitsCache(key string) bool {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	_, ok := rc.status[key]
	return ok
}

// 是否在更新缓存
func IsUpdate(key string) bool {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	if v, ok := rc.status[key]; ok {
		return v == 1
	}
	return false
}

// need update cache
func NeedUpdate(key string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.status[key] = 2
	rc.store.Add(key, nil)
}

func SetUpdate(key string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.status[key] = 1
	rc.store.Add(key, []byte(""))
}

func InitResponseCache(cache Cacher) {
	rc = &responseCache{
		store:  cache,
		status: make(map[string]int),
		mu:     sync.RWMutex{},
	}
}

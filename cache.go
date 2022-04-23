package xmux

import (
	"net/http"
	"sync"
	"time"
)

const CacheKey = "CacheKey"

func DefaultCacheTemplateCacheWithResponse(w http.ResponseWriter, r *http.Request) bool {
	// 获取唯一id
	// 建议 url + uid 或者 MD5(url + uid), 如果跟uid无关， 可以只用url
	// 先要判断一下是否存在缓存
	ck := GetInstance(r).Get(CacheKey)
	if ck == nil {
		return false
	}
	cacheKey := ck.(string)
	_, cacheStatus := GetCacheIfUpdating(cacheKey)
	switch cacheStatus {
	case CacheHit:
		return true
	case CacheIsUpdateing:
		for {
			select {
			case <-time.After(time.Second):
				return false
			default:
				time.Sleep(time.Millisecond)
				if !IsUpdate(cacheKey) {
					return true
				}
			}

		}
	case CacheNeedUpdate:
		SetNeedUpdateToUpdate(cacheKey)
		return false
	case NotFoundCache:
		SetUpdate(cacheKey)
		return false
	default:
		return false
	}
}

func DefaultCacheTemplateCacheWithoutResponse(w http.ResponseWriter, r *http.Request) bool {
	// 获取唯一id
	// 建议 url + uid 或者 MD5(url + uid), 如果跟uid无关， 可以只用url
	// 先要判断一下是否存在缓存
	ck := GetInstance(r).Get(CacheKey)
	if ck == nil {
		// 没有启用缓存
		return false
	}
	cacheKey := ck.(string)
	cb, cacheStatus := GetCacheIfUpdating(cacheKey)
	switch cacheStatus {
	case CacheHit:
		w.Write(cb)
		return true
	case CacheIsUpdateing:
		for {
			select {
			case <-time.After(time.Second):
				return false
			default:
				time.Sleep(time.Millisecond)
				if !IsUpdate(cacheKey) {
					w.Write(GetCache(cacheKey))
					return true
				}
			}

		}
	case CacheNeedUpdate:
		SetNeedUpdateToUpdate(cacheKey)
		return false
	case NotFoundCache:
		SetUpdate(cacheKey)
		return false
	default:
		return false
	}

}

// cache 是一个接口
type cacher interface {
	SetCache(key string, cache []byte) // 设置缓存保存的结构体
	ClearCache(key string)             // 清除某个key的缓存
	ClearAll()                         // 清除所有key
}

// 返回的缓存
type responseCache struct {
	store map[string]*cacheStruct
	mu    sync.RWMutex
}

// 设置缓存值
func SetCache(key string, cache []byte) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.store[key].response = cache
	rc.store[key].isUpdate = false
	rc.store[key].update = time.Now()
}

func SetNeedUpdateToUpdate(key string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	if _, ok := rc.store[key]; ok {
		rc.store[key].response = []byte("")
	}
}

// 如果key存在，就设置缓存
func SetCacheIfExsits(key string, cache []byte) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	if _, ok := rc.store[key]; ok {
		rc.store[key].response = cache
		rc.store[key].isUpdate = false
		rc.store[key].update = time.Now()
	}
}

// 获取缓存值
func GetCache(key string) []byte {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	if _, ok := rc.store[key]; ok {
		return rc.store[key].response
	}
	return nil
}

// 获取缓存，如果存在
func GetCacheIfExsits(key string) ([]byte, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	if _, ok := rc.store[key]; ok {
		return rc.store[key].response, true
	}
	return nil, false
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
// 如果返回 NotFoundCache 说明是最新的，可以直接返回
func GetCacheIfUpdating(key string) ([]byte, CacheStatus) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	if _, ok := rc.store[key]; ok {
		if rc.store[key].isUpdate {
			if rc.store[key].response == nil {
				return nil, CacheNeedUpdate
			} else {
				return rc.store[key].response, CacheIsUpdateing
			}
		}
		return rc.store[key].response, CacheHit
	}

	return nil, NotFoundCache
}

// 是否存在缓存
func ExsitsCache(key string) bool {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	_, ok := rc.store[key]
	return ok
}

// 是否在更新缓存
func IsUpdate(key string) bool {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	if _, ok := rc.store[key]; ok {
		return rc.store[key].isUpdate
	}
	return false
}

func NeedUpdate(key string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	if _, ok := rc.store[key]; ok {
		rc.store[key].response = nil
	}
}

func SetUpdate(key string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	if _, ok := rc.store[key]; ok {
		rc.store[key].isUpdate = true
	} else {
		rc.store[key] = &cacheStruct{
			isUpdate: true,
			update:   time.Now(),
		}
	}
}

type cacheStruct struct {
	response   []byte
	update     time.Time // 最后一次更新的时间， 用来判断最后更新的时间
	isUpdate   bool      // 判断是否在刷新缓存中
	needUpdate bool      // 设置需要更新
}

var rc *responseCache

func InitResponseCache() {
	rc = &responseCache{
		store: make(map[string]*cacheStruct),
		mu:    sync.RWMutex{},
	}
}

func ResponseCache() *GroupRoute {
	InitResponseCache()
	cache := NewGroupRoute()
	// cache.Post("/-/cache/size", list)
	// cache.Post("/-/cache/add", list)
	// cache.Post("/-/cache/del", list)
	// cache.Post("/-/cache/clear", list)
	return cache
}

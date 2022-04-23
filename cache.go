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
	cb, ok := GetCacheIfUpdating(cacheKey)
	if cb == nil && !ok {
		SetUpdate(cacheKey)
		return false
	}
	if ok {
		// 如果在更新， 那么等待更新完毕
		for {
			select {
			case <-time.After(time.Second * 10):
				SetUpdate(cacheKey)
				return false
			default:
				time.Sleep(time.Millisecond)
				if !IsUpdate(cacheKey) {
					return true
				}
			}

		}
	}
	return true
}

func DefaultCacheTemplateCacheWithoutResponse(w http.ResponseWriter, r *http.Request) bool {
	// 获取唯一id
	// 建议 url + uid 或者 MD5(url + uid), 如果跟uid无关， 可以只用url
	// 先要判断一下是否存在缓存
	ck := GetInstance(r).Get(CacheKey)
	if ck == nil {
		return false
	}
	cacheKey := ck.(string)
	cb, ok := GetCacheIfUpdating(cacheKey)
	if cb == nil && !ok {
		SetUpdate(cacheKey)
		return false
	}
	if ok {
		// 如果在更新， 那么等待更新完毕
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
	}
	w.Write(cb)
	return true
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

// 获取缓存，如果正在更新
// 如果返回 nil, false    说明不存在这个缓存
// 如果返回 []byte, true  说明当前还在更新中， 还不是最新的缓存
// 如果返回 []byte, false 说明是最新的，可以直接返回
func GetCacheIfUpdating(key string) ([]byte, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	if _, ok := rc.store[key]; ok {
		return rc.store[key].response, rc.store[key].isUpdate
	}
	return nil, false
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
	response []byte
	update   time.Time // 最后一次更新的时间， 用来判断最后更新的时间
	isUpdate bool      // 判断是否在刷新缓存中
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

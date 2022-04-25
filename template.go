package xmux

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/hyahm/xmux/cache"
)

func DefaultCacheTemplateCacheWithResponse(w http.ResponseWriter, r *http.Request) bool {
	// 获取唯一id
	// 建议 url + uid 或者 MD5(url + uid), 如果跟uid无关， 可以只用url
	// 先要判断一下是否存在缓存
	ck := GetInstance(r).Get(cache.CacheKey)
	if ck == nil {
		return false
	}
	cacheKey := ck.(string)
	_, cacheStatus := cache.GetCacheIfUpdating(cacheKey)
	switch cacheStatus {
	case cache.CacheHit:
		return true
	case cache.CacheIsUpdateing:
		for {
			select {
			case <-time.After(time.Second):
				return false
			default:
				time.Sleep(time.Millisecond)
				if !cache.IsUpdate(cacheKey) {
					return true
				}
			}

		}
	case cache.CacheNeedUpdate:
		cache.SetNeedUpdateToUpdate(cacheKey)
		return false
	case cache.NotFoundCache:
		cache.SetUpdate(cacheKey)
		return false
	default:
		return false
	}
}

func DefaultCacheTemplateCacheWithoutResponse(w http.ResponseWriter, r *http.Request) bool {
	// 获取唯一id
	// 建议 url + uid 或者 MD5(url + uid), 如果跟uid无关， 可以只用url
	// 先要判断一下是否存在缓存
	ck := GetInstance(r).Get(cache.CacheKey)
	if ck == nil {
		// 没有启用缓存
		return false
	}
	cacheKey := ck.(string)
	cb, cacheStatus := cache.GetCacheIfUpdating(cacheKey)
	switch cacheStatus {
	case cache.CacheHit:
		w.Write(cb)
		return true
	case cache.CacheIsUpdateing:
		for {
			select {
			case <-time.After(time.Second):
				return false
			default:
				time.Sleep(time.Millisecond)
				if !cache.IsUpdate(cacheKey) {
					w.Write(cache.GetCache(cacheKey))
					return true
				}
			}

		}
	case cache.CacheNeedUpdate:
		cache.SetNeedUpdateToUpdate(cacheKey)
		return false
	case cache.NotFoundCache:
		cache.SetUpdate(cacheKey)
		return false
	default:
		return false
	}

}

func exit(start time.Time, w http.ResponseWriter, r *http.Request) {
	r.Body.Close()
	var send []byte
	if GetInstance(r).Response != nil && GetInstance(r).Get(STATUSCODE).(int) == 200 {

		ck := GetInstance(r).Get(cache.CacheKey)

		if ck != nil {
			cacheKey := ck.(string)
			if cache.IsUpdate(cacheKey) {
				// 如果没有设置缓存，还是以前的处理方法
				send, err := json.Marshal(GetInstance(r).Response)
				if err != nil {
					log.Println(err)
				}

				// 如果之前是更新的状态，那么就修改
				cache.SetCache(cacheKey, send)
				w.Write(send)
			} else {
				// 如果不是更新的状态， 那么就不用更新，而是直接从缓存取值
				send = cache.GetCache(cacheKey)
				w.Write(send)
			}
		} else {
			send, err := json.Marshal(GetInstance(r).Response)
			if err != nil {
				log.Println(err)
			}
			w.Write(send)
		}

	}
	log.Printf("connect_id: %d,method: %s\turl: %s\ttime: %f\t status_code: %v, body: %v\n",
		GetInstance(r).GetConnectId(),
		r.Method,
		r.URL.Path, time.Since(start).Seconds(),
		GetInstance(r).Get(STATUSCODE),
		string(send))
}

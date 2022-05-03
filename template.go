package xmux

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

func DefaultCacheTemplateCacheWithResponse(w http.ResponseWriter, r *http.Request) bool {
	// 获取唯一id
	// 建议 url + uid 或者 MD5(url + uid), 如果跟uid无关， 可以只用url
	// 先要判断一下是否存在缓存
	cacheKey := GetInstance(r).CacheKey
	if cacheKey == "" {
		return false
	}
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
				time.Sleep(time.Millisecond * 10)
				if !IsUpdate(cacheKey) {
					SetUpdate(cacheKey)
					return true
				}
			}

		}
	case CacheNeedUpdate:
		SetUpdate(cacheKey)
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
	cacheKey := GetInstance(r).CacheKey
	if cacheKey == "" {
		// 没有启用缓存
		return false
	}

	cb, cacheStatus := GetCacheIfUpdating(cacheKey)
	switch cacheStatus {
	case CacheHit:
		w.Write(cb)
		return true
	case CacheIsUpdateing:
		// 如果在更新中，那么等待更新完毕再返回缓存， 如果等待1秒了还没返回就不等待缓存
		for {
			select {
			case <-time.After(time.Second):
				return false
			default:
				time.Sleep(time.Millisecond * 10)
				if !IsUpdate(cacheKey) {
					w.Write(GetCache(cacheKey))
					return true
				}
			}
		}
	case CacheNeedUpdate:
		SetUpdate(cacheKey)
		return false
	case NotFoundCache:
		SetUpdate(cacheKey)
		return false
	default:
		return false
	}

}

func exit(start time.Time, w http.ResponseWriter, r *http.Request) {
	// r.Body.Close()
	var send []byte
	var err error
	if GetInstance(r).Response != nil && GetInstance(r).StatusCode == 200 {
		cacheKey := GetInstance(r).CacheKey

		if cacheKey != "" {
			if IsUpdate(cacheKey) {
				// 如果没有设置缓存，还是以前的处理方法
				send, err = json.Marshal(GetInstance(r).Response)
				if err != nil {
					log.Println(err)
				}
				// 如果之前是更新的状态，那么就修改
			} else {
				// 如果不是更新的状态， 那么就不用更新，而是直接从缓存取值
				send = GetCache(cacheKey)
			}
			SetCache(cacheKey, send)
			w.Write(send)
		} else {
			send, err = json.Marshal(GetInstance(r).Response)
			if err != nil {
				log.Println(err)
			}
			w.Write(send)
		}

	}
	log.Printf("connect_id: %d,method: %s\turl: %s\ttime: %f\t status_code: %v, response: %v\n",
		GetInstance(r).GetConnectId(),
		r.Method,
		r.URL.Path, time.Since(start).Seconds(),
		GetInstance(r).StatusCode,
		string(send))
}

func DefaultPermissionTemplate(w http.ResponseWriter, r *http.Request) (post bool) {
	// 如果是管理员的，直接就过
	// if uid == <adminId> {
	// 	retrun false
	// }

	pages := GetInstance(r).GetPageKeys()
	// 如果长度为0的话，说明任何人都可以访问
	if len(pages) == 0 {
		return false
	}
	// todo: get user.role
	role := ""
	if _, ok := pages[role]; !ok {
		return true
	}

	// roles := []string{"env", "important"}
	// 内置的方法最大支持8种权限，如果想要更多可以通过 GetPerm 实现
	// 设置权限列表
	var pl = []string{"Read", "Create", "Update", "Delete"}
	// map 的key 对应页面的value  value 对应二进制位置(从右到左)
	permissionMap := make(map[string]int)
	for k, v := range pl {
		permissionMap[v] = k
	}
	// todo: perm  type uint8
	perm := 14
	result := GetPerm(pl, uint8(perm))
	for index, ok := range result {
		currFuncName := GetInstance(r).GetFuncName()
		if ok && strings.Contains(currFuncName, pl[index]) {
			return false
		}
	}
	// no permission
	return true
	// 假如权限拿到二进制对应的10进制数据是下面
	// perm = 14       // 00001110   {"Delete", "Create", "Update"}
	// perm = 10       // 00001010   {"Create", "Delete"}
	// perm = 4        // 00000100   {"Update"}

	//

	//  请求/project/read     map[admin:{} project:{}]
	// 判断 pages 是否存在 perm
	// 注意点： 这里的页面权限本应该只会匹配到一个， 这个是对于的页面权限的值

	// permMap := make(map[string]bool)

	// 先拿到pl 对应名称的 索引
	//         8        4        2          1
	//		 delete	 update	 create		read
	//  bit   0        0       0         0
	/*
		用户表
		id
		1
		权限表
		id      uid   roles                       perm
		1       1     "env"                       0-15
		2       1     "important"
	*/
}

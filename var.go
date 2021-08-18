package xmux

import (
	"net/http"
	"sync"
)

// 保存url里面的参数
type params map[string]string // url 参数对应的值

var allparams map[string]params // 保存的url 参数
var paramsLocker sync.RWMutex

func init() {
	allparams = make(map[string]params)
	paramsLocker = sync.RWMutex{}
}
func Var(r *http.Request) params {
	return GetParams(r.URL.Path)
}

func GetParams(key string) params {
	paramsLocker.Lock()
	defer paramsLocker.Unlock()
	return allparams[key]
}

func SetParams(key string, params params) {
	paramsLocker.Lock()
	allparams[key] = params
	paramsLocker.Unlock()
}

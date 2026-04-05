package xmux

import (
	"errors"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

const (
	zero  = byte('0')
	one   = byte('1')
	lsb   = byte('[') // left square brackets
	rsb   = byte(']') // right square brackets
	space = byte(' ')
)

var uint8arr [8]uint8

// ErrBadStringFormat represents a error of input string's format is illegal .
var ErrBadStringFormat = errors.New("bad string format")

// ErrEmptyString represents a error of empty input string.
var ErrEmptyString = errors.New("empty string")

var ErrTypeUnsupport = errors.New("data type is unsupported")

func init() {
	// for _, value := range os.Environ() {
	// 	if strings.Contains(value, "GOEXPERIMENT") {
	// 		fmt.Println(value)
	// 	}

	// }

	uint8arr[0] = 128
	uint8arr[1] = 64
	uint8arr[2] = 32
	uint8arr[3] = 16
	uint8arr[4] = 8
	uint8arr[5] = 4
	uint8arr[6] = 2
	uint8arr[7] = 1
}

// 保存url里面的参数
type params map[string]string // url 参数对应的值

var allparams map[string]params // 保存的url 参数
var paramsLocker sync.RWMutex

func init() {
	allparams = make(map[string]params)
	paramsLocker = sync.RWMutex{}
}
func Var(r *http.Request) params {
	return getParams(r.URL.Path)
}

func getParams(key string) params {
	paramsLocker.RLock()
	defer paramsLocker.RUnlock()
	return allparams[key]
}

func setParams(key string, params params) {
	paramsLocker.Lock()
	allparams[key] = params
	paramsLocker.Unlock()
}

// router.go 中添加
// PageKeyFuncMap 返回 pagekey -> []funcName 的映射
func (r *Router) PageKeyFuncMap() map[string][]string {
	result := make(map[string][]string)
	// 遍历所有精确匹配路由
	for _, route := range r.urlRoute {
		funcName := getFuncName(route.handle)
		for pk := range route.pagekeys {
			result[pk] = append(result[pk], funcName)
		}
	}
	// 遍历所有正则匹配路由
	for _, route := range r.urlTpl {
		funcName := getFuncName(route.handle)
		for pk := range route.pagekeys {
			result[pk] = append(result[pk], funcName)
		}
	}
	return result
}

// 辅助函数：从 handle 提取 funcName
func getFuncName(handle http.Handler) string {
	name := runtime.FuncForPC(reflect.ValueOf(handle).Pointer()).Name()
	n := strings.LastIndex(name, ".")
	if n >= 0 {
		return name[n+1:]
	}
	return name
}

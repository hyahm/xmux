package xmux

import (
	"errors"
	"net/http"
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
	paramsLocker.Lock()
	defer paramsLocker.Unlock()
	return allparams[key]
}

func setParams(key string, params params) {
	paramsLocker.Lock()
	allparams[key] = params
	paramsLocker.Unlock()
}

package xmux

import (
	"net/http"
	"sync"
)

// bridge  数据二次封装

type Data struct {
	Data interface{}            // 处理后的数据
	ctx  map[string]interface{} // 用来传递自定义值
	mu   *sync.RWMutex
}

var allconn map[*http.Request]*Data
var dataLock *sync.RWMutex

func init() {
	allconn = make(map[*http.Request]*Data)
	dataLock = &sync.RWMutex{}
}

func GetData(r *http.Request) *Data {
	if r == nil {
		return nil
	}
	dataLock.RLock()
	defer dataLock.RUnlock()
	if v, ok := allconn[r]; ok {
		return v
	}
	return nil
}

func (data *Data) Set(k string, v interface{}) {
	if data.ctx == nil {
		data.ctx = make(map[string]interface{})
		data.mu = &sync.RWMutex{}
	}

	data.mu.Lock()
	data.ctx[k] = v
	data.mu.Unlock()
}

func (data *Data) Get(k string) interface{} {
	if data.ctx == nil {
		return nil
	}
	data.mu.RLock()
	defer data.mu.RUnlock()
	return data.ctx[k]
}

func (data *Data) Del(k string) {
	if data.ctx == nil {
		return
	}
	data.mu.Lock()
	delete(data.ctx, k)
	data.mu.Unlock()
}

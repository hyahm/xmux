package xmux

import (
	"net/http"
	"sync"
)

// instance  数据二次封装, 用户各模块之间的数据传递

type FlowData struct {
	Data interface{}            // 处理后的数据
	ctx  map[string]interface{} // 用来传递自定义值
	mu   *sync.RWMutex
}

var allconn map[*http.Request]*FlowData
var dataLock *sync.RWMutex

func init() {
	allconn = make(map[*http.Request]*FlowData)
	dataLock = &sync.RWMutex{}
}

func GetInstance(r *http.Request) *FlowData {
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

func (data *FlowData) Set(k string, v interface{}) {
	if data.ctx == nil {
		data.ctx = make(map[string]interface{})
	}

	data.mu.Lock()
	data.ctx[k] = v
	data.mu.Unlock()
}

func (data *FlowData) Get(k string) interface{} {
	if data.ctx == nil {
		return nil
	}
	data.mu.RLock()
	defer data.mu.RUnlock()
	if v, ok := data.ctx[k]; ok {
		return v
	}
	return nil
}

func (data *FlowData) Del(k string) {
	data.mu.Lock()
	delete(data.ctx, k)
	data.mu.Unlock()
}

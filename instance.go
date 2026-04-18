package xmux

import (
	"net/http"
)

const xmux_context = "XMUX_CONTEXT"

type FlowData struct {
	Data any
	// 处理后的数据
	ctx map[string]interface{} // 用来传递自定义值
	// mu         *sync.RWMutex
	Response   interface{} // 返回的数据结构
	connectId  int64
	funcName   string
	pages      map[string]struct{}
	StatusCode int
	Body       []byte
	cacheKey   string
	module     []func(w http.ResponseWriter, r *http.Request) bool
}

func GetInstance(r *http.Request) *FlowData {

	return r.Context().Value(xmux_context).(*FlowData)
}

func (data *FlowData) Set(k string, v interface{}) {
	data.ctx[k] = v
}

func (data *FlowData) SetCacheKey(key string) {
	data.cacheKey = key
}

func (data *FlowData) GetCacheKey() string {
	return data.cacheKey
}

func (data *FlowData) GetConnectId() int64 {
	return data.connectId
}

func (data *FlowData) GetFuncName() string {
	return data.funcName
}

func (data *FlowData) GetModules() []func(http.ResponseWriter, *http.Request) bool {
	return data.module
}

func (data *FlowData) GetPageKeys() map[string]struct{} {
	return data.pages
}

func (data *FlowData) Get(k string) interface{} {
	if v, ok := data.ctx[k]; ok {
		return v
	}
	return nil
}

func (data *FlowData) Del(k string) {
	delete(data.ctx, k)
}

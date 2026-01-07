package xmux

import (
	"net/http"
	"sync"
)

// instance  数据二次封装, 用户各模块之间的数据传递

type FlowData struct {
	Data any
	// 处理后的数据
	ctx        map[string]interface{} // 用来传递自定义值
	mu         *sync.RWMutex
	Response   interface{} // 返回的数据结构
	connectId  int64
	funcName   string
	pages      map[string]struct{}
	StatusCode int
	Body       []byte
	CacheKey   string
}

type conns struct {
	conn map[*http.Request]*FlowData
	mu   *sync.RWMutex
}

var allconn *conns

// var dataLock *sync.RWMutex

func init() {
	allconn = &conns{
		conn: make(map[*http.Request]*FlowData),
		mu:   &sync.RWMutex{},
	}
}

func (conns *conns) Set(r *http.Request, fd *FlowData) {
	conns.mu.Lock()
	defer conns.mu.Unlock()
	conns.conn[r] = fd
}

func (conns *conns) Del(r *http.Request) {
	conns.mu.Lock()
	defer conns.mu.Unlock()
	delete(conns.conn, r)
}

func (conns *conns) Get(r *http.Request) *FlowData {
	conns.mu.RLock()
	defer conns.mu.RUnlock()
	if v, ok := conns.conn[r]; ok {
		return v
	}
	return nil
}

func GetInstance(r *http.Request) *FlowData {
	if r == nil {
		return nil
	}
	allconn.mu.RLock()
	defer allconn.mu.RUnlock()
	if v, ok := allconn.conn[r]; ok {
		return v
	}
	return nil
}

func (data *FlowData) Set(k string, v interface{}) {
	data.mu.Lock()
	data.ctx[k] = v
	data.mu.Unlock()
}

func (data *FlowData) GetConnectId() int64 {
	data.mu.RLock()
	defer data.mu.RUnlock()
	return data.connectId
}

func (data *FlowData) GetFuncName() string {
	data.mu.RLock()
	defer data.mu.RUnlock()
	return data.funcName
}

func (data *FlowData) GetPageKeys() map[string]struct{} {
	data.mu.RLock()
	defer data.mu.RUnlock()
	return data.pages
}

func (data *FlowData) Get(k string) interface{} {
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

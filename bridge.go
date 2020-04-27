package xmux

import (
	"net/http"
	"sync"
)

// bridge  数据二次封装

type Data struct {
	Var    map[string]string // 参数
	Header map[string]string
	Data   interface{} // 处理后的数据
	End    interface{}
	ctx    map[string]interface{} // 用来传递自定义值
	mu     *sync.RWMutex
}

var Bridge map[string]*Data

func init() {
	Bridge = make(map[string]*Data)
}

func GetData(r *http.Request) *Data {
	url := slash(r.URL.Path)
	return Bridge[url]
}

func (data *Data) Set(k string, v interface{}) {
	if data.mu == nil {
		data.mu = &sync.RWMutex{}
	}
	if data.ctx == nil {
		data.ctx = make(map[string]interface{})
	}
	data.mu.Lock()
	data.ctx[k] = v
	data.mu.Unlock()
}

func (data *Data) Get(k string) interface{} {
	if data.mu == nil {
		data.mu = &sync.RWMutex{}
	}
	if data.ctx == nil {
		data.ctx = make(map[string]interface{})
	}
	data.mu.Lock()
	defer data.mu.Unlock()
	return data.ctx[k]
}

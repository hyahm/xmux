package xmux

import (
	"net/http"
	"sync"
)

// bridge  数据二次封装

type Data struct {
	Data interface{}            // 处理后的数据
	ctx  map[string]interface{} // 用来传递自定义值
	// Header map[string]string
	mu  *sync.RWMutex
	End interface{}
}

type params map[string]string // url 参数对应的值

var allparams map[string]params // 保存的url 参数

var allconn map[*http.Request]*Data

func Var(r *http.Request) params {
	url := slash(r.URL.Path)
	return allparams[url]
}

func init() {
	allparams = make(map[string]params)
	allconn = make(map[*http.Request]*Data)
}

func GetData(r *http.Request) *Data {
	if r == nil {
		return nil
	}

	return allconn[r]
}

func (data *Data) Set(k string, v interface{}) {

	data.mu.Lock()
	data.ctx[k] = v
	data.mu.Unlock()
}

func (data *Data) Get( k string) interface{} {

	data.mu.Lock()
	defer data.mu.Unlock()
	return data.ctx[k]
}

func (data *Data) Del(k string) {

	data.mu.Lock()
	delete(data.ctx, k)
	data.mu.Unlock()
}

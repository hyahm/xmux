package xmux

import "net/http"

// bridge  数据二次封装

type Data struct {
	Var    map[string]string // 参数
	Header map[string]string
	Data   interface{} // 处理后的数据
	End    interface{}
	Ctx    map[string]interface{} // 用来传递自定义值
}

var Bridge map[string]*Data

func init() {
	Bridge = make(map[string]*Data)
}

func GetData(r *http.Request) *Data {
	url := slash(r.URL.Path)
	return Bridge[url]
}

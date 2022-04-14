package xmux

import (
	"net/http"
	"sync"
)

// 初始化临时使用， 最后会合并到 router
type Route struct {
	// 组里面也包括路由 后面的其实还是patter和handle, 还没到handle， 这里的key是个method
	new       bool
	handle    http.Handler // handle
	methods   map[string]struct{}
	module    *module             // 增加的 modules
	delmodule map[string]struct{} // 删除的modules

	pagekeys    map[string]struct{} // 页面权限
	delPageKeys map[string]struct{} // 删除的权限

	header    map[string]string   // 请求头
	delheader map[string]struct{} // 删除的请求头

	describe         string      // 接口描述
	responseData     interface{} // 接口返回实例
	bindResponseData bool

	bindType   bindType    // 数据绑定格式
	dataSource interface{} // 数据源
}

func (rt *Route) GetHeader() map[string]string {
	if !rt.new {
		panic("can not support init")
	}
	return rt.header
}

func (rt *Route) BindResponse(response interface{}) *Route {
	if !rt.new {
		panic("can not support init")
	}
	rt.responseData = response
	rt.bindResponseData = true
	return rt
}

func (rt *Route) AddPageKeys(pagekeys ...string) *Route {
	if !rt.new {
		panic("can not support init")
	}
	// 退出文档的组
	for _, v := range pagekeys {
		if rt.pagekeys == nil {
			rt.pagekeys = make(map[string]struct{})
		}
		rt.pagekeys[v] = struct{}{}
	}
	return rt
}

func (rt *Route) DelPageKeys(pagekeys ...string) *Route {
	if !rt.new {
		panic("can not support init")
	}
	if rt.delPageKeys == nil {
		rt.delPageKeys = make(map[string]struct{})
	}
	for _, key := range pagekeys {
		rt.delPageKeys[key] = struct{}{}
	}
	return rt
}

// 数据绑定
func (rt *Route) Bind(s interface{}) *Route {
	if !rt.new {
		panic("can not support init")
	}
	rt.dataSource = s
	return rt
}

// json数据绑定
func (rt *Route) BindJson(s interface{}) *Route {
	if !rt.new {
		panic("can not support init")
	}
	// 接口补充说明
	rt.dataSource = s
	rt.bindType = jsonT
	return rt
}

// yaml数据绑定
func (rt *Route) BindYaml(s interface{}) *Route {
	if !rt.new {
		panic("can not support init")
	}
	// 接口补充说明
	rt.dataSource = s
	rt.bindType = yamlT
	return rt
}

// xml数据绑定
func (rt *Route) BindXml(s interface{}) *Route {
	if !rt.new {
		panic("can not support init")
	}
	// 接口补充说明

	rt.dataSource = s
	rt.bindType = xmlT
	return rt
}

// 组里面也包括路由 后面的其实还是patter和handle

func (rt *Route) SetHeader(k, v string) *Route {
	if !rt.new {
		panic("can not support init")
	}
	if rt.header == nil {
		rt.header = map[string]string{}
	}
	rt.header[k] = v
	return rt
}

func (rt *Route) DelHeader(dh ...string) *Route {
	if !rt.new {
		panic("can not support init")
	}
	if rt.delheader == nil {
		rt.delheader = make(map[string]struct{})
	}
	for _, v := range dh {
		rt.delheader[v] = struct{}{}
	}
	return rt
}

func (rt *Route) AddModule(handles ...func(http.ResponseWriter, *http.Request) bool) *Route {
	if !rt.new {
		panic("can not support init")
	}
	if rt.module == nil {
		rt.module = &module{
			filter:    make(map[string]struct{}),
			funcOrder: make([]func(w http.ResponseWriter, r *http.Request) bool, 0),
			mu:        sync.RWMutex{},
		}
	}
	rt.module.add(handles...)
	return rt
}

func (rt *Route) DelModule(handles ...func(http.ResponseWriter, *http.Request) bool) *Route {
	if !rt.new {
		panic("can not support init")
	}
	if rt.delmodule == nil {
		rt.delmodule = make(map[string]struct{})
	}
	for _, handle := range handles {
		rt.delmodule[GetFuncName(handle)] = struct{}{}
	}
	return rt
}

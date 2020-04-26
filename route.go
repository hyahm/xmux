package xmux

import (
	"log"
	"net/http"
	"strings"
)

// 初始化临时使用， 最后会合并到 router
type Route struct {
	// 组里面也包括路由 后面的其实还是patter和handle, 还没到handle， 这里的key是个method
	method     map[string]http.Handler
	header     http.Header
	args       []string // 保存正则的变量名
	midware    []func(http.ResponseWriter, *http.Request) bool
	delmidware []func(http.ResponseWriter, *http.Request) bool
	describe   string // 接口描述

	request        string      // 请求的请求示例
	dataSource     interface{} // 数据源
	response       string      // 接口返回示例
	st_request     interface{}
	params_request map[string]string
	st_response    interface{}
	reqHeader      map[string]string
	supplement     string
	delheader      []string
}

func (rt *Route) Bind(s interface{}) *Route {
	// 接口补充说明
	rt.dataSource = s
	return rt
}

func (rt *Route) ApiSupplement(s string) *Route {
	// 接口补充说明
	rt.supplement = s
	return rt
}

func (rt *Route) ApiReqStruct(s interface{}) *Route {
	// 接口返回数据的结构
	rt.st_request = s
	return rt
}

func (rt *Route) ApiReqParams(s map[string]string) *Route {
	// 接口返回数据的结构
	rt.params_request = s
	return rt
}

func (rt *Route) ApiResStruct(s interface{}) *Route {
	// 接口接收数据的结构
	rt.st_response = s
	return rt
}

func (rt *Route) makeDoc() Document {
	doc := Document{
		Opt:     make([]Opt, 0),
		Callbak: make([]Opt, 0),
	}
	doc.Describe = rt.describe
	doc.Header = rt.reqHeader

	if rt.st_response != nil {
		doc.Callbak = PostOpt(rt.st_response)
	}
	doc.Request = rt.request
	doc.Response = rt.response
	return doc
}

func (rt *Route) ApiDescribe(s string) *Route {
	// 接口的简单描述
	rt.describe = s
	return rt
}

func (rt *Route) ApiReqHeader(head map[string]string) *Route {
	// 接口的请求头

	rt.reqHeader = head
	return rt
}

func (rt *Route) ApiRequestTemplate(s string) *Route {
	// 接口的请求实例， 一般是json的字符串
	rt.request = s
	return rt
}

func (rt *Route) ApiResponseTemplate(s string) *Route {
	// 接口的返回实例， 一般是json的字符串
	rt.response = s
	return rt
}

// 组里面也包括路由 后面的其实还是patter和handle
func (r *Router) Pattern(pattern string) *Route {
	// 格式化路径
	if r.route == nil {
		r.route = make(map[string]*Route)
	}
	if r.pattern == nil {
		r.pattern = make(map[string]int)
	}
	if r.tplpattern == nil {
		r.tplpattern = make(map[string]int)
	}
	pattern = slash(pattern)
	if _, ok := r.pattern[pattern]; ok {
		log.Fatalf("Pattern Duplicate for %s", pattern)
	}

	if pattern == "" || pattern[0:1] != "/" || strings.ContainsAny(pattern, " \t\n") {
		log.Fatalf("Pattern Error for %s", pattern)
	}
	route := &Route{
		method: make(map[string]http.Handler),
		header: http.Header{},
		args:   make([]string, 0),
	}
	// 增加pattern 判断

	if v, listvar := match(pattern); len(listvar) > 0 {
		if _, ok := r.pattern[v]; ok {
			log.Fatalf("Pattern Duplicate for %s", v)
		}
		r.tpl[v] = route
		r.tplpattern[v] = 1
		r.tpl[v].args = append(r.tpl[v].args, listvar...)
		return r.tpl[v]
	}
	r.pattern[pattern] = 0
	r.route[pattern] = route
	return r.route[pattern]
}

func (rt *Route) Post(handler func(http.ResponseWriter, *http.Request)) *Route {
	if _, ok := rt.method[http.MethodPost]; ok {
		log.Fatal("method post duplicate")
	}
	rt.method[http.MethodPost] = http.HandlerFunc(handler)
	return rt
}

func (rt *Route) Get(handler func(http.ResponseWriter, *http.Request)) *Route {
	if _, ok := rt.method[http.MethodGet]; ok {
		log.Fatal("method get duplicate")
	}
	rt.method[http.MethodGet] = http.HandlerFunc(handler)
	return rt
}

func (rt *Route) Delete(handler func(http.ResponseWriter, *http.Request)) *Route {
	if _, ok := rt.method[http.MethodDelete]; ok {
		log.Fatal("method Delete duplicate")
	}
	rt.method[http.MethodDelete] = http.HandlerFunc(handler)
	return rt
}

func (rt *Route) Head(handler func(http.ResponseWriter, *http.Request)) *Route {
	if _, ok := rt.method[http.MethodHead]; ok {
		log.Fatal("method Head duplicate")
	}
	rt.method[http.MethodHead] = http.HandlerFunc(handler)
	return rt
}

func (rt *Route) WebSocket(ws WsHandler) *Route {
	if _, ok := rt.method[http.MethodGet]; ok {
		log.Fatal("method Get duplicate")
	}
	rt.method[http.MethodGet] = http.HandlerFunc(ws.Websocket)
	return rt
}

func (rt *Route) Options(handler func(http.ResponseWriter, *http.Request)) *Route {
	if _, ok := rt.method[http.MethodOptions]; ok {
		log.Fatal("method Options duplicate")
	}
	rt.method[http.MethodOptions] = http.HandlerFunc(handler)
	return rt
}

func (rt *Route) Connect(handler func(http.ResponseWriter, *http.Request)) *Route {
	if _, ok := rt.method[http.MethodConnect]; ok {
		log.Fatal("method Connect duplicate")
	}
	rt.method[http.MethodConnect] = http.HandlerFunc(handler)
	return rt
}

func (rt *Route) Patch(handler func(http.ResponseWriter, *http.Request)) *Route {
	if _, ok := rt.method[http.MethodPatch]; ok {
		log.Fatal("method Patch duplicate")
	}
	rt.method[http.MethodPatch] = http.HandlerFunc(handler)
	return rt
}

func (rt *Route) Trace(handler func(http.ResponseWriter, *http.Request)) *Route {
	if _, ok := rt.method[http.MethodTrace]; ok {
		log.Fatal("method Trace duplicate")
	}
	rt.method[http.MethodTrace] = http.HandlerFunc(handler)
	return rt
}

func (rt *Route) Put(handler func(http.ResponseWriter, *http.Request)) *Route {
	if _, ok := rt.method[http.MethodPut]; ok {
		log.Fatal("method put duplicate")
	}
	rt.method[http.MethodPut] = http.HandlerFunc(handler)
	return rt
}

func (rt *Route) AddHeader(k, v string) *Route {
	if rt.header == nil {
		rt.header = http.Header{}
	}
	rt.header.Add(k, v)
	return rt
}

func (rt *Route) DelHeader(k string) *Route {
	if rt.delheader == nil {
		rt.delheader = make([]string, 0)
	}
	rt.delheader = append(rt.delheader, k)
	return rt
}

func (rt *Route) AddMidware(handle func(http.ResponseWriter, *http.Request) bool) *Route {
	if rt.midware == nil {
		rt.midware = make([]func(http.ResponseWriter, *http.Request) bool, 0)
	}
	rt.midware = append(rt.midware, handle)
	return rt
}

func (rt *Route) DelMidware(handle func(http.ResponseWriter, *http.Request) bool) *Route {
	if rt.midware == nil {
		rt.delmidware = make([]func(http.ResponseWriter, *http.Request) bool, 0)
	}
	rt.delmidware = append(rt.delmidware, handle)
	return rt
}

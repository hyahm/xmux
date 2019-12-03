package xmux

import (
	"log"
	"net/http"
	"strings"
)

type Route struct {
	// 组里面也包括路由 后面的其实还是patter和handle, 还没到handle， 这里的key是个method
	method map[string]http.Handler
	header map[string]string
	args   []string // 保存正则的变量名
}

// 组里面也包括路由 后面的其实还是patter和handle
func (r *Router) Pattern(pattern string) *Route {
	// 格式路径
	pattern = slash(pattern)
	if _, ok := r.route[pattern]; ok {
		log.Fatalf("pattern duplicate for %s", pattern)
	}
	pattern = strings.Trim(pattern, " ")
	if pattern == "" || pattern[0:1] != "/" {
		log.Fatalf("pattern error for %s", pattern)
	}
	route := &Route{
		method: make(map[string]http.Handler),
		header: make(map[string]string),
		args:   make([]string, 0),
	}
	lv := make([]string, 0)
	if v, listvar, ok := match(pattern, "^", lv); ok {
		r.tpl[v] = route
		r.tpl[v].args = append(r.tpl[v].args, listvar...)
		return r.tpl[v]
	}

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

func (rt *Route) SetHeader(k, v string) *Route {
	rt.header[k] = v
	return rt
}

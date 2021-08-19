package xmux

import (
	"log"
	"net/http"
)

// string 对应的是method
type MethodsRoute map[string]*Route

func (mr MethodsRoute) getRoute(pattern, method string, handler func(http.ResponseWriter, *http.Request),
	gr *GroupRoute) *Route {
	if _, ok := mr[method]; ok {
		log.Fatalf("url: %s, method: %s is duplicate", pattern, method)
	}

	mr[method] = &Route{
		handle:    http.HandlerFunc(handler),
		midware:   gr.midware,
		pagekeys:  gr.pagekeys,
		module:    gr.module,
		delmodule: gr.delmodule,
	}
	return mr[method]
}

func (mr MethodsRoute) getRouterRoute(pattern, method string,
	handler func(http.ResponseWriter, *http.Request), router *Router) *Route {
	if _, ok := mr[method]; ok {
		log.Fatalf("url: %s, method: %s is duplicate", pattern, method)
	}

	mr[method] = &Route{
		handle:   http.HandlerFunc(handler),
		midware:  router.midware,
		pagekeys: router.pagekeys,
		module:   router.module,
	}
	return mr[method]
}

func (mr MethodsRoute) SetHeader(key, value string) MethodsRoute {
	for _, route := range mr {
		route.SetHeader(key, value)
	}
	return mr
}

func (mr MethodsRoute) MiddleWare(midware func(handle func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request)) MethodsRoute {
	for _, route := range mr {
		route.MiddleWare(midware)
	}
	return mr
}

func (mr MethodsRoute) AddModule(handles ...func(http.ResponseWriter, *http.Request) bool) MethodsRoute {
	for _, route := range mr {
		route.AddModule(handles...)
	}
	return mr
}

func (mr MethodsRoute) AddPageKeys(pagekeys ...string) MethodsRoute {
	for _, route := range mr {
		route.AddPageKeys(pagekeys...)
	}
	return mr
}

func (mr MethodsRoute) DelHeader(key string) MethodsRoute {
	for _, route := range mr {
		route.DelHeader(key)
	}
	return mr
}

func (mr MethodsRoute) DelMiddleWare(midware func(handle func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request)) MethodsRoute {
	for _, route := range mr {
		route.DelMiddleWare(midware)
	}
	return mr
}

func (mr MethodsRoute) DelModule(handles ...func(http.ResponseWriter, *http.Request) bool) MethodsRoute {
	for _, route := range mr {
		route.DelModule(handles...)
	}
	return mr
}

func (mr MethodsRoute) DelPageKeys(pagekeys ...string) MethodsRoute {
	for _, route := range mr {
		route.DelPageKeys(pagekeys...)
	}
	return mr
}

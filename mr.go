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

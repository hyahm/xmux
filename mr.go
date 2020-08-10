package xmux

import (
	"log"
	"net/http"
)

type MethodsRoute map[string]*Route

func (mr MethodsRoute) getRoute(method string, handler func(http.ResponseWriter, *http.Request)) *Route {

	if _, ok := mr[method]; ok {
		log.Fatal("method post duplicate")
	}
	mr[method] = &Route{
		handle: http.HandlerFunc(handler),
	}
	return mr[method]
}

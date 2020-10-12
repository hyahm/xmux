package xmux

import (
	"net/http"
)

// string 对应的是method
type MethodsRoute map[string]*Route

func (mr MethodsRoute) getRoute(method string, handler func(http.ResponseWriter, *http.Request),
	midware func(handle func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request)) *Route {
	// if _, ok := mr[method]; !ok {
	// 	mr = make(map[string]*Route)
	// }
	mr[method] = &Route{
		handle:  http.HandlerFunc(handler),
		midware: midware,
	}
	return mr[method]
}

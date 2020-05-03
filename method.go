package xmux

import (
	"log"
	"net/http"
)

type MethodsRoute map[string]*Route

func (mr MethodsRoute) Post(handler func(http.ResponseWriter, *http.Request)) *Route {
	if _, ok := mr[http.MethodPost]; ok {
		log.Fatal("method post duplicate")
	}
	route := &Route{
		handle: http.HandlerFunc(handler),
	}
	mr[http.MethodPost] = route
	return route
}

func (mr MethodsRoute) Get(handler func(http.ResponseWriter, *http.Request)) *Route {
	if _, ok := mr[http.MethodGet]; ok {
		log.Fatal("method get duplicate")
	}

	route := &Route{
		handle: http.HandlerFunc(handler),
	}
	mr[http.MethodGet] = route
	return route
}

func (mr MethodsRoute) Delete(handler func(http.ResponseWriter, *http.Request)) *Route {
	if _, ok := mr[http.MethodDelete]; ok {
		log.Fatal("method Delete duplicate")
	}
	route := &Route{
		handle: http.HandlerFunc(handler),
	}
	mr[http.MethodDelete] = route
	return route
}

func (mr MethodsRoute) Head(handler func(http.ResponseWriter, *http.Request)) *Route {
	if _, ok := mr[http.MethodHead]; ok {
		log.Fatal("method Head duplicate")
	}
	route := &Route{
		handle: http.HandlerFunc(handler),
	}
	mr[http.MethodHead] = route
	return route
}

// func (mr methodsRoute) WebSocket(ws WsHandler) *Route {
// 	if _, ok := mr[http.MethodGet]; ok {
// 		log.Fatal("method Get duplicate")
// 	}
// 	return &Route{
// 		handle: http.HandlerFunc(handler),
// 	}
// }

func (mr MethodsRoute) Options(handler func(http.ResponseWriter, *http.Request)) *Route {
	if _, ok := mr[http.MethodOptions]; ok {
		log.Fatal("method Options duplicate")
	}
	route := &Route{
		handle: http.HandlerFunc(handler),
	}
	mr[http.MethodOptions] = route
	return route
}

func (mr MethodsRoute) Connect(handler func(http.ResponseWriter, *http.Request)) *Route {
	if _, ok := mr[http.MethodConnect]; ok {
		log.Fatal("method Connect duplicate")
	}
	route := &Route{
		handle: http.HandlerFunc(handler),
	}
	mr[http.MethodConnect] = route
	return route
}

func (mr MethodsRoute) Patch(handler func(http.ResponseWriter, *http.Request)) *Route {
	if _, ok := mr[http.MethodPatch]; ok {
		log.Fatal("method Patch duplicate")
	}
	route := &Route{
		handle: http.HandlerFunc(handler),
	}
	mr[http.MethodPatch] = route
	return route
}

func (mr MethodsRoute) Trace(handler func(http.ResponseWriter, *http.Request)) *Route {
	if _, ok := mr[http.MethodTrace]; ok {
		log.Fatal("method Trace duplicate")
	}
	route := &Route{
		handle: http.HandlerFunc(handler),
	}
	mr[http.MethodTrace] = route
	return route
}

func (mr MethodsRoute) Put(handler func(http.ResponseWriter, *http.Request)) *Route {
	if _, ok := mr[http.MethodPut]; ok {
		log.Fatal("method put duplicate")
	}
	route := &Route{
		handle: http.HandlerFunc(handler),
	}
	mr[http.MethodPut] = route
	return route
}

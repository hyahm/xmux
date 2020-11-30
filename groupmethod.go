package xmux

import (
	"net/http"
)

func (gr *GroupRoute) method(pattern string, handler func(http.ResponseWriter, *http.Request), method string) *Route {
	if pt, ok := gr.makeRoute(pattern); ok {
		return gr.tpl[pt].getRoute(method, handler, gr.midware)
	} else {
		return gr.route[pt].getRoute(method, handler, gr.midware)
	}
}

func (gr *GroupRoute) Post(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return gr.method(pattern, handler, http.MethodPost)
}

func (gr *GroupRoute) All(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return gr.method(pattern, handler, "*")
}

func (gr *GroupRoute) Get(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return gr.method(pattern, handler, http.MethodGet)
}

func (gr *GroupRoute) Delete(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return gr.method(pattern, handler, http.MethodDelete)
}

func (gr *GroupRoute) Head(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return gr.method(pattern, handler, http.MethodHead)
}

func (gr *GroupRoute) Options(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return gr.method(pattern, handler, http.MethodOptions)
}

func (gr *GroupRoute) Connect(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return gr.method(pattern, handler, http.MethodConnect)
}

func (gr *GroupRoute) Patch(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return gr.method(pattern, handler, http.MethodPatch)
}

func (gr *GroupRoute) Trace(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return gr.method(pattern, handler, http.MethodTrace)
}

func (gr *GroupRoute) Put(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return gr.method(pattern, handler, http.MethodPut)
}

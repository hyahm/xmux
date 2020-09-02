package xmux

import (
	"net/http"
)

func (r *Router) method(pattern string, handler func(http.ResponseWriter, *http.Request), method string) *Route {
	if pt, ok := r.makeRoute(pattern); ok {
		return r.tpl[pt].getRoute(method, handler, r.midware)
	} else {
		return r.route[pt].getRoute(method, handler, r.midware)
	}
}

func (r *Router) Post(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return r.method(pattern, handler, http.MethodPost)
}

func (r *Router) Get(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return r.method(pattern, handler, http.MethodGet)
}

func (r *Router) Delete(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return r.method(pattern, handler, http.MethodDelete)
}

func (r *Router) Head(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return r.method(pattern, handler, http.MethodHead)
}

func (r *Router) Options(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return r.method(pattern, handler, http.MethodOptions)
}

func (r *Router) Connect(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return r.method(pattern, handler, http.MethodConnect)
}

func (r *Router) Patch(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return r.method(pattern, handler, http.MethodPatch)
}

func (r *Router) Trace(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return r.method(pattern, handler, http.MethodTrace)
}

func (r *Router) Put(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return r.method(pattern, handler, http.MethodPut)
}

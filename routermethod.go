package xmux

import (
	"net/http"
)

func (r *Router) Post(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {

	if pt, ok := r.makeRoute(pattern); ok {
		return r.tpl[pt].getRoute(http.MethodPost, handler)
	} else {
		return r.route[pt].getRoute(http.MethodPost, handler)
	}
}

func (r *Router) Get(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if pt, ok := r.makeRoute(pattern); ok {
		return r.tpl[pt].getRoute(http.MethodGet, handler)
	} else {
		return r.route[pt].getRoute(http.MethodGet, handler)
	}
}

func (r *Router) Delete(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if pt, ok := r.makeRoute(pattern); ok {
		return r.tpl[pt].getRoute(http.MethodDelete, handler)
	} else {
		return r.route[pt].getRoute(http.MethodDelete, handler)
	}
}

func (r *Router) Head(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if pt, ok := r.makeRoute(pattern); ok {
		return r.tpl[pt].getRoute(http.MethodHead, handler)
	} else {
		return r.route[pt].getRoute(http.MethodHead, handler)
	}
}

func (r *Router) Option(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if pt, ok := r.makeRoute(pattern); ok {
		return r.tpl[pt].getRoute(http.MethodOptions, handler)
	} else {
		return r.route[pt].getRoute(http.MethodOptions, handler)
	}
}

func (r *Router) Connect(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if pt, ok := r.makeRoute(pattern); ok {
		return r.tpl[pt].getRoute(http.MethodConnect, handler)
	} else {
		return r.route[pt].getRoute(http.MethodConnect, handler)
	}
}

func (r *Router) Patch(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if pt, ok := r.makeRoute(pattern); ok {
		return r.tpl[pt].getRoute(http.MethodPatch, handler)
	} else {
		return r.route[pt].getRoute(http.MethodPatch, handler)
	}
}

func (r *Router) Trace(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if pt, ok := r.makeRoute(pattern); ok {
		return r.tpl[pt].getRoute(http.MethodTrace, handler)
	} else {
		return r.route[pt].getRoute(http.MethodTrace, handler)
	}
}

func (r *Router) Put(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if pt, ok := r.makeRoute(pattern); ok {
		return r.tpl[pt].getRoute(http.MethodPut, handler)
	} else {
		return r.route[pt].getRoute(http.MethodPut, handler)
	}
}

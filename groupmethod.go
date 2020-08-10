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

func (gr *GroupRoute) Post(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if pt, ok := gr.makeRoute(pattern); ok {
		return gr.tpl[pt].getRoute(http.MethodPost, handler)

	} else {
		return gr.route[pt].getRoute(http.MethodPost, handler)
	}
}

func (gr *GroupRoute) Get(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if pt, ok := gr.makeRoute(pattern); ok {
		return gr.tpl[pt].getRoute(http.MethodGet, handler)
	} else {
		return gr.route[pt].getRoute(http.MethodGet, handler)
	}
}

func (gr *GroupRoute) Delete(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if pt, ok := gr.makeRoute(pattern); ok {
		return gr.tpl[pt].getRoute(http.MethodDelete, handler)
	} else {
		return gr.route[pt].getRoute(http.MethodDelete, handler)
	}
}

func (gr *GroupRoute) Head(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if pt, ok := gr.makeRoute(pattern); ok {
		return gr.tpl[pt].getRoute(http.MethodHead, handler)
	} else {
		return gr.route[pt].getRoute(http.MethodHead, handler)
	}

}

func (gr *GroupRoute) Options(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if pt, ok := gr.makeRoute(pattern); ok {
		return gr.tpl[pt].getRoute(http.MethodOptions, handler)
	} else {
		return gr.route[pt].getRoute(http.MethodOptions, handler)
	}

}

func (gr *GroupRoute) Connect(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if pt, ok := gr.makeRoute(pattern); ok {
		return gr.tpl[pt].getRoute(http.MethodConnect, handler)
	} else {
		return gr.route[pt].getRoute(http.MethodConnect, handler)
	}
}

func (gr *GroupRoute) Patch(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if pt, ok := gr.makeRoute(pattern); ok {
		return gr.tpl[pt].getRoute(http.MethodPatch, handler)
	} else {
		return gr.route[pt].getRoute(http.MethodPatch, handler)
	}
}

func (gr *GroupRoute) Trace(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if pt, ok := gr.makeRoute(pattern); ok {
		return gr.tpl[pt].getRoute(http.MethodTrace, handler)
	} else {
		return gr.route[pt].getRoute(http.MethodTrace, handler)
	}
}

func (gr *GroupRoute) Put(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if pt, ok := gr.makeRoute(pattern); ok {
		return gr.tpl[pt].getRoute(http.MethodPut, handler)
	} else {
		return gr.route[pt].getRoute(http.MethodPut, handler)
	}
}

package xmux

import (
	"log"
	"net/http"
	"sync"
)

// get this route
func (gr *GroupRoute) defindMethod(pattern string, handler func(http.ResponseWriter, *http.Request), methods ...string) *Route {

	temphead := make(map[string]string)
	for k, v := range gr.header {
		temphead[k] = v
	}

	tempPages := make(map[string]struct{})
	for k := range gr.pagekeys {
		tempPages[k] = struct{}{}
	}
	newRoute := &Route{
		handle:   http.HandlerFunc(handler),
		pagekeys: make(map[string]struct{}),
		module: &module{
			filter:    make(map[string]struct{}),
			funcOrder: make([]func(w http.ResponseWriter, r *http.Request) bool, 0),
			mu:        sync.RWMutex{},
		},
		new:         true,
		methods:     make(map[string]struct{}),
		header:      make(map[string]string),
		delmodule:   make(map[string]struct{}),
		delPageKeys: make(map[string]struct{}),
		delheader:   make(map[string]struct{}),
	}
	url, vars, ok := gr.makeRoute(pattern)
	gr.params[url] = vars
	if ok {
		if gr.tpl == nil {
			gr.tpl = make(map[string]*Route)
		}
		for _, method := range methods {
			if _, methodOk := newRoute.methods[method]; methodOk {
				// 如果也存在， 那么method重复了
				log.Fatal("method : " + method + "  duplicate, url: " + url)
			}
			newRoute.methods[method] = struct{}{}
		}
		gr.tpl[url] = newRoute
	} else {
		if gr.route == nil {
			gr.route = make(map[string]*Route)
		}
		// 如果存在就判断是否存在method
		for _, method := range methods {
			if _, methodOk := newRoute.methods[method]; methodOk {
				// 如果也存在， 那么method重复了
				log.Fatal("method : " + method + "  duplicate, url: " + url)
			}
			newRoute.methods[method] = struct{}{}
		}
		// 如果不存在就创建一个 route
		gr.route[url] = newRoute
	}
	return newRoute
}

func (gr *GroupRoute) Post(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !gr.new {
		panic("must be init by NewGroupRoute")
	}
	return gr.defindMethod(pattern, handler, http.MethodPost)
}

func (gr *GroupRoute) Any(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !gr.new {
		panic("must be init by NewGroupRoute")
	}
	return gr.defindMethod(pattern, handler, http.MethodConnect,
		http.MethodDelete, http.MethodGet, http.MethodHead,
		http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodTrace,
	)
}

func (gr *GroupRoute) Requests(pattern string, handler func(http.ResponseWriter, *http.Request), methods ...string) *Route {
	return gr.defindMethod(pattern, handler, methods...)
}

func (gr *GroupRoute) Get(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !gr.new {
		panic("must be init by NewGroupRoute")
	}
	return gr.defindMethod(pattern, handler, http.MethodGet)
}

func (gr *GroupRoute) Request(pattern string, handler func(http.ResponseWriter, *http.Request), methods ...string) *Route {
	if !gr.new {
		panic("must be init by NewGroupRoute")
	}
	return gr.defindMethod(pattern, handler, methods...)
}

func (gr *GroupRoute) Delete(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !gr.new {
		panic("must be init by NewGroupRoute")
	}
	return gr.defindMethod(pattern, handler, http.MethodDelete)
}

func (gr *GroupRoute) Head(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !gr.new {
		panic("must be init by NewGroupRoute")
	}
	return gr.defindMethod(pattern, handler, http.MethodHead)
}

func (gr *GroupRoute) Options(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !gr.new {
		panic("must be init by NewGroupRoute")
	}
	return gr.defindMethod(pattern, handler, http.MethodOptions)
}

func (gr *GroupRoute) Connect(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !gr.new {
		panic("must be init by NewGroupRoute")
	}
	return gr.defindMethod(pattern, handler, http.MethodConnect)
}

func (gr *GroupRoute) Patch(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !gr.new {
		panic("must be init by NewGroupRoute")
	}
	return gr.defindMethod(pattern, handler, http.MethodPatch)
}

func (gr *GroupRoute) Trace(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !gr.new {
		panic("must be init by NewGroupRoute")
	}
	return gr.defindMethod(pattern, handler, http.MethodTrace)
}

func (gr *GroupRoute) Put(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !gr.new {
		panic("must be init by NewGroupRoute")
	}
	return gr.defindMethod(pattern, handler, http.MethodPut)
}

package xmux

import (
	"log"
	"net/http"
	"sync"
)

func (gr *GroupRoute) any(pattern string, handler func(http.ResponseWriter, *http.Request), methods ...string) MethodsRoute {

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
		header:      make(map[string]string),
		delmodule:   make(map[string]struct{}),
		delPageKeys: make(map[string]struct{}),
		delheader:   make(map[string]struct{}),
	}
	url, vars, ok := gr.makeRoute(pattern)
	gr.params[url] = vars
	if ok {
		if gr.tpl[url] == nil {
			gr.tpl[url] = make(MethodsRoute)
		}
		for _, method := range methods {
			if _, methodOk := gr.tpl[url][method]; methodOk {
				// 如果也存在， 那么method重复了
				log.Fatal("method : " + method + "  duplicate, url: " + url)
			}
			if gr.tpl[url] == nil {
				gr.tpl[url] = make(MethodsRoute)
			}
			// newRoute.methods[method] = struct{}{}
			newRoute.url = url
			newRoute.params = vars
			gr.tpl[url][method] = newRoute
		}
		return gr.tpl[url]
	} else {
		if gr.route[url] == nil {
			gr.route[url] = make(MethodsRoute)
		}

		// 如果存在就判断是否存在method
		for _, method := range methods {
			if _, methodOk := gr.route[url][method]; methodOk {
				// 如果也存在， 那么method重复了
				log.Fatal("method : " + method + "  duplicate, url: " + url)
			}
			if gr.route[url] == nil {
				gr.route[url] = make(MethodsRoute)
			}
			newRoute.url = url
			gr.route[url][method] = newRoute
		}
		// 如果不存在就创建一个 route
		return gr.route[url]
	}

}

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
		header:      make(map[string]string),
		delmodule:   make(map[string]struct{}),
		delPageKeys: make(map[string]struct{}),
		delheader:   make(map[string]struct{}),
	}
	url, vars, ok := gr.makeRoute(pattern)
	gr.params[url] = vars
	if ok {
		if gr.tpl[url] == nil {
			gr.tpl[url] = make(map[string]*Route)
		}
		for _, method := range methods {
			if _, methodOk := gr.tpl[url][method]; methodOk {
				// 如果也存在， 那么method重复了
				log.Fatal("method : " + method + "  duplicate, url: " + url)
			}
			if gr.tpl[url] == nil {
				gr.tpl[url] = make(MethodsRoute)
			}
			// newRoute.methods[method] = struct{}{}
			newRoute.url = url
			newRoute.params = vars
			gr.tpl[url][method] = newRoute
		}

	} else {
		if gr.route[url] == nil {
			gr.route[url] = make(map[string]*Route)
		}
		// 如果存在就判断是否存在method
		for _, method := range methods {
			if _, methodOk := gr.route[url][method]; methodOk {
				// 如果也存在， 那么method重复了
				log.Fatal("method : " + method + "  duplicate, url: " + url)
			}
			newRoute.url = url
			gr.route[url][method] = newRoute
		}
		// 如果不存在就创建一个 route

	}
	return newRoute
}

func (gr *GroupRoute) Post(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !gr.new {
		panic("must be init by NewGroupRoute")
	}
	return gr.defindMethod(pattern, handler, http.MethodPost)
}

func (gr *GroupRoute) Any(pattern string, handler func(http.ResponseWriter, *http.Request)) MethodsRoute {
	if !gr.new {
		panic("must be init by NewGroupRoute")
	}
	return gr.any(pattern, handler, http.MethodConnect,
		http.MethodDelete, http.MethodGet, http.MethodHead,
		http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodTrace,
	)
}

func (gr *GroupRoute) Requests(pattern string, handler func(http.ResponseWriter, *http.Request), methods ...string) MethodsRoute {
	return gr.any(pattern, handler, methods...)
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

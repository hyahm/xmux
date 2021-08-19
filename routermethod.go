package xmux

import (
	"log"
	"net/http"
)

// get this route
func (r *Router) defindMethod(pattern string, handler func(http.ResponseWriter, *http.Request), methods ...string) MethodsRoute {
	if len(methods) == 0 {
		panic(pattern + " have not any methods")
	}
	mr := make(MethodsRoute)
	for _, method := range methods {
		mr[method] = r.method(pattern, handler, method)
	}
	return mr
}

func (r *Router) method(pattern string, handler func(http.ResponseWriter, *http.Request), method string) *Route {
	// 判断是否是正则
	temphead := make(map[string]string)
	for k, v := range r.header {
		temphead[k] = v
	}

	tempPages := make(map[string]struct{})
	for k := range r.pagekeys {
		tempPages[k] = struct{}{}
	}
	newRoute := &Route{
		handle:      http.HandlerFunc(handler),
		midware:     r.midware,
		pagekeys:    tempPages,
		module:      r.module,
		isRoot:      true,
		header:      temphead,
		delheader:   make([]string, 0),
		delmodule:   delModule{},
		delPageKeys: make([]string, 0),
	}
	url, ok := r.makeRoute(pattern)
	if ok {
		if _, urlok := r.tpl[url]; urlok {
			// 如果存在就判断是否存在method
			if _, methodok := r.tpl[url][method]; methodok {
				// 如果也存在， 那么method重复了
				log.Fatal(ErrMethodDuplicate)
				return nil
			} else {
				// 如果不存在就创建一个 route
				r.tpl[url][method] = newRoute
				return newRoute
			}
		} else {
			// 如果不存在就创建一个PatternMR
			r.tpl = make(PatternMR)
			mr := make(MethodsRoute)
			mr[method] = newRoute
			r.tpl[url] = mr
			return newRoute
		}
		// 判断是不是正则
		// return gr.tpl[url].getRoute(url, method, handler, gr)
	} else {
		if _, urlok := r.route[url]; urlok {
			// 如果存在就判断是否存在method
			if _, methodok := r.route[url][method]; methodok {
				// 如果也存在， 那么method重复了
				log.Fatal(ErrMethodDuplicate)
				return nil
			} else {
				// 如果不存在就创建一个 route
				r.route[url][method] = newRoute
				return newRoute
			}
		} else {
			// 如果不存在就创建一个PatternMR
			r.route = make(PatternMR)
			mr := make(MethodsRoute)
			mr[method] = newRoute
			r.route[url] = mr
			return newRoute
		}
	}
}

func (r *Router) Post(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.method(pattern, handler, http.MethodPost)
}

func (r *Router) Any(pattern string, handler func(http.ResponseWriter, *http.Request)) MethodsRoute {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.defindMethod(pattern, handler, http.MethodConnect,
		http.MethodDelete, http.MethodGet, http.MethodHead,
		http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodTrace,
	)
}

func (r *Router) Get(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.method(pattern, handler, http.MethodGet)
}

func (r *Router) Request(pattern string, handler func(http.ResponseWriter, *http.Request), methods ...string) MethodsRoute {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.defindMethod(pattern, handler, methods...)
}

func (r *Router) Delete(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.method(pattern, handler, http.MethodDelete)
}

func (r *Router) Head(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.method(pattern, handler, http.MethodHead)
}

func (r *Router) Options(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.method(pattern, handler, http.MethodOptions)
}

func (r *Router) Connect(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.method(pattern, handler, http.MethodConnect)
}

func (r *Router) Patch(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.method(pattern, handler, http.MethodPatch)
}

func (r *Router) Trace(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.method(pattern, handler, http.MethodTrace)
}

func (r *Router) Put(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.method(pattern, handler, http.MethodPut)
}

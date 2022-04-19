package xmux

import (
	"log"
	"net/http"
)

// get this route
func (r *Router) defindMethod(pattern string, handler func(http.ResponseWriter, *http.Request), methods ...string) *Route {
	if len(methods) == 0 {
		panic(pattern + " have not any methods")
	}
	temphead := make(map[string]string)
	for k, v := range r.header {
		temphead[k] = v
	}

	tempPages := make(map[string]struct{})
	for k := range r.pagekeys {
		tempPages[k] = struct{}{}
	}
	newRoute := &Route{
		handle:       http.HandlerFunc(handler),
		pagekeys:     tempPages,
		module:       r.module.cloneMudule(),
		new:          true,
		responseData: r.responseData,
		header:       temphead,
		delheader:    make(map[string]struct{}),
		delmodule:    make(map[string]struct{}),
		delPageKeys:  make(map[string]struct{}),
	}
	// 判断是否是正则
	url, vars, ok := r.makeRoute(pattern)
	r.params[url] = vars
	if ok {
		// 正则匹配的
		for _, method := range methods {

			if _, ok := r.tpl[url]; ok {
				if _, ok := r.tpl[url][method]; ok {
					log.Fatal("method : " + method + "  duplicate, url: " + url)
				}
			}
			if r.tpl[url] == nil {
				r.tpl[url] = make(MethodsRoute)
			}
			r.tpl[url][method] = newRoute
		}

		// 如果不存在就创建一个 route

	} else {
		// 直接匹配
		for _, method := range methods {
			// 如果存在就判断是否存在method
			if _, ok := r.route[url]; ok {
				if _, ok := r.route[url][method]; ok {
					log.Fatal("method : " + method + "  duplicate, url: " + url)
				}
			}
			if r.route[url] == nil {
				r.route[url] = make(MethodsRoute)
			}
			r.route[url][method] = newRoute
		}

	}
	return newRoute
}

func (r *Router) any(pattern string, handler func(http.ResponseWriter, *http.Request), methods ...string) MethodsRoute {
	if len(methods) == 0 {
		panic(pattern + " have not any methods")
	}
	temphead := make(map[string]string)
	for k, v := range r.header {
		temphead[k] = v
	}

	tempPages := make(map[string]struct{})
	for k := range r.pagekeys {
		tempPages[k] = struct{}{}
	}
	newRoute := &Route{
		handle:       http.HandlerFunc(handler),
		pagekeys:     tempPages,
		module:       r.module.cloneMudule(),
		new:          true,
		responseData: r.responseData,
		header:       temphead,
		delheader:    make(map[string]struct{}),
		delmodule:    make(map[string]struct{}),
		delPageKeys:  make(map[string]struct{}),
	}
	// 判断是否是正则
	url, vars, ok := r.makeRoute(pattern)

	r.params[url] = vars
	if ok {
		// 正则匹配的
		if r.tpl[url] == nil {
			r.tpl[url] = make(MethodsRoute)
		}
		for _, method := range methods {
			if _, ok := r.tpl[url]; ok {
				if _, ok := r.tpl[url][method]; ok {
					log.Fatal("method : " + method + "  duplicate, url: " + url)
				}
			}
			if r.tpl[url] == nil {
				r.tpl[url] = make(MethodsRoute)
			}
			r.tpl[url][method] = newRoute
		}

		// 如果不存在就创建一个 route
		return r.tpl[url]
	} else {
		// 直接匹配
		if r.route[url] == nil {
			r.route[url] = make(MethodsRoute)
		}
		r.route[url] = make(MethodsRoute)
		for _, method := range methods {
			// 如果存在就判断是否存在method
			if _, ok := r.route[url]; ok {
				if _, ok := r.route[url][method]; ok {
					log.Fatal("method : " + method + "  duplicate, url: " + url)
				}
			}
			if r.route[url] == nil {
				r.route[url] = make(MethodsRoute)
			}
			r.route[url][method] = newRoute
		}
		return r.route[url]
	}
}

func (r *Router) Post(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.defindMethod(pattern, handler, http.MethodPost)
}

func (r *Router) Any(pattern string, handler func(http.ResponseWriter, *http.Request)) MethodsRoute {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.any(pattern, handler, http.MethodConnect,
		http.MethodDelete, http.MethodGet, http.MethodHead,
		http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodTrace,
	)
}

// func (mr MethodsRoute) Bind(data interface{}) MethodsRoute {
// 	for _, route := range mr {
// 		route.dataSource = data
// 	}
// 	return mr
// }

func (r *Router) Get(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.defindMethod(pattern, handler, http.MethodGet)
}

func (r *Router) Request(pattern string, handler func(http.ResponseWriter, *http.Request), methods ...string) MethodsRoute {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.any(pattern, handler, methods...)
}

func (r *Router) Delete(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.defindMethod(pattern, handler, http.MethodDelete)
}

func (r *Router) Head(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.defindMethod(pattern, handler, http.MethodHead)
}

func (r *Router) Options(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.defindMethod(pattern, handler, http.MethodOptions)
}

func (r *Router) Connect(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.defindMethod(pattern, handler, http.MethodConnect)
}

func (r *Router) Patch(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.defindMethod(pattern, handler, http.MethodPatch)
}

func (r *Router) Trace(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.defindMethod(pattern, handler, http.MethodTrace)
}

func (r *Router) Put(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.defindMethod(pattern, handler, http.MethodPut)
}

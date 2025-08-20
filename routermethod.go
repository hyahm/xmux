package xmux

import (
	"log"
	"net/http"
	"path"
)

// get this route
func (r *Router) defindMethod(pattern string, handler func(http.ResponseWriter, *http.Request), method ...string) *Route {
	temphead := make(map[string]string)
	for k, v := range r.header {
		temphead[k] = v
	}
	tempPages := make(map[string]struct{})
	for k := range r.pagekeys {
		tempPages[k] = struct{}{}
	}
	// route.methods = append(route.methods, http.MethodGet)
	// subprefixs := SubtractSliceMap(prefix, route.delprefix)
	// subprefix := append(subprefixs, route.prefixs...)
	// allurl := path.Join(subprefix...)
	// allurl = PrettySlash(allurl + route.url)
	// url, vars, ok := makeRoute(allurl)
	newRoute := &Route{
		handle:       http.HandlerFunc(handler),
		pagekeys:     tempPages,
		module:       r.module.cloneMudule(),
		new:          true,
		responseData: r.responseData,
		methods:      append(make([]string, 0), method...),
		header:       temphead,
		delheader:    make(map[string]struct{}),
		delmodule:    make(map[string]struct{}),
		delPageKeys:  make(map[string]struct{}),
		prefixs:      make([]string, 0),
		delprefix:    map[string]struct{}{},
	}
	// prefix := path.Join(r.prefix...)
	// prefix = path.Join(prefix, pattern)
	// 判断是否是正则
	url, vars, ok := makeRoute(pattern)
	newRoute.params = vars
	if ok {
		// 正则匹配的
		if _, ok := r.urlTpl[url]; ok {
			if _, ok := r.urlTpl[url]; ok {
				m, exsit := SliceExsit(r.urlTpl[url].methods, method)
				if exsit {
					log.Fatal("method : " + m + "  duplicate, url: " + url)
				}
			}
		}
		if len(r.prefix) > 0 {
			url = r.mergePrefix(newRoute, url)

		}
		r.urlTpl[url] = newRoute
	} else {
		// 直接匹配
		// 如果存在就判断是否存在method
		if _, ok := r.urlRoute[url]; ok {
			if _, ok := r.urlRoute[url]; ok {
				m, exsit := SliceExsit(r.urlRoute[url].methods, method)
				if exsit {
					log.Fatal("method : " + m + "  duplicate, url: " + url)
				}
			}
		}
		if len(r.prefix) > 0 {
			allPrefix := append(r.prefix, url)
			url = path.Join(allPrefix...)
		}
		r.urlRoute[url] = newRoute

	}
	return newRoute
}

// func (r *Router) any(pattern string, handler func(http.ResponseWriter, *http.Request), methods ...string) MethodsRoute {
// 	if len(methods) == 0 {
// 		panic(pattern + " have not any methods")
// 	}
// 	temphead := make(map[string]string)
// 	for k, v := range r.header {
// 		temphead[k] = v
// 	}

// 	tempPages := make(map[string]struct{})
// 	for k := range r.pagekeys {
// 		tempPages[k] = struct{}{}
// 	}
// 	newRoute := &Route{
// 		handle:       http.HandlerFunc(handler),
// 		pagekeys:     tempPages,
// 		module:       r.module.cloneMudule(),
// 		new:          true,
// 		responseData: r.responseData,
// 		header:       temphead,
// 		delheader:    make(map[string]struct{}),
// 		delmodule:    make(map[string]struct{}),
// 		delPageKeys:  make(map[string]struct{}),
// 		prefixs:      make([]string, 0),
// 		delprefix:    map[string]struct{}{},
// 	}
// 	// 判断是否是正则
// 	prefix := path.Join(r.prefix...)
// 	prefix = path.Join(prefix, pattern)
// 	url, vars, ok := makeRoute(prefix)
// 	newRoute.params = vars
// 	newRoute.url = url
// 	if ok {
// 		// 正则匹配的
// 		for _, method := range methods {
// 			if _, ok := r.tpl[url]; ok {
// 				if _, ok := r.tpl[url][method]; ok {
// 					log.Fatal("method : " + method + "  duplicate, url: " + url)
// 				}
// 			} else {
// 				r.tpl[url] = make(MethodsRoute)
// 			}

// 			r.tpl[url][method] = newRoute
// 		}

// 		// 如果不存在就创建一个 route
// 		return r.tpl[url]
// 	} else {
// 		// 直接匹配
// 		for _, method := range methods {
// 			// 如果存在就判断是否存在method
// 			if _, ok := r.route[url]; ok {
// 				if _, ok := r.route[url][method]; ok {
// 					log.Fatal("method : " + method + "  duplicate, url: " + url)
// 				}
// 			} else {
// 				r.route[url] = make(MethodsRoute)
// 			}
// 			r.route[url][method] = newRoute
// 		}
// 		return r.route[url]
// 	}
// }

func (r *Router) Post(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.defindMethod(pattern, handler, http.MethodPost)
}

func (r *Router) Any(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.defindMethod(pattern, handler,
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

func (r *Router) Request(pattern string, handler func(http.ResponseWriter, *http.Request), methods ...string) *Route {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	return r.defindMethod(pattern, handler, methods...)
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

// func (r *Router) Connect(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
// 	if !r.new {
// 		panic("must be use get router by NewRouter()")
// 	}
// 	return r.defindMethod(pattern, handler, http.MethodConnect)
// }

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

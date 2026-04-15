package xmux

import (
	"log"
	"net/http"
	"path"
)

// get this route
func (r *router) defindMethod(pattern string, handler func(http.ResponseWriter, *http.Request), method ...string) *Route {
	if !r.DisableOption {
		var exsitOption bool
		for _, v := range method {
			if v == http.MethodOptions {
				exsitOption = true
				break
			}
		}
		if !exsitOption {
			method = append(method, http.MethodOptions)
		}
	}
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
		handle:        http.HandlerFunc(handler),
		pagekeys:      tempPages,
		module:        r.module.cloneMudule(),
		postModule:    r.postModule.cloneMudule(),
		new:           true,
		responseData:  r.responseData,
		methods:       append(make([]string, 0), method...),
		header:        temphead,
		delheader:     make(map[string]struct{}),
		delmodule:     make(map[string]struct{}),
		delPostModule: make(map[string]struct{}),
		delPageKeys:   make(map[string]struct{}),
		prefixs:       make([]string, 0),
		delprefix:     map[string]struct{}{},
		middleware: onion{
			mws: make([]Middleware, 0),
		},
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
				m, exsit := SliceExist(r.urlTpl[url].methods, method)
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
				m, exsit := SliceExist(r.urlRoute[url].methods, method)
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

func (r *router) Post(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return r.defindMethod(pattern, handler, http.MethodPost)
}

func (r *router) Any(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
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

func (r *router) Get(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {

	return r.defindMethod(pattern, handler, http.MethodGet)
}

func (r *router) Connect(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {

	return r.defindMethod(pattern, handler, http.MethodConnect)
}

func (r *router) Request(pattern string, handler func(http.ResponseWriter, *http.Request), methods ...string) *Route {
	return r.defindMethod(pattern, handler, methods...)
}

func (r *router) Delete(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return r.defindMethod(pattern, handler, http.MethodDelete)
}

func (r *router) Head(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return r.defindMethod(pattern, handler, http.MethodHead)
}

func (r *router) Options(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return r.defindMethod(pattern, handler, http.MethodOptions)
}

// func (r *Router) Connect(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
// 	if !r.new {
// 		panic("must be use get router by NewRouter()")
// 	}
// 	return r.defindMethod(pattern, handler, http.MethodConnect)
// }

func (r *router) Patch(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return r.defindMethod(pattern, handler, http.MethodPatch)
}

func (r *router) Trace(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return r.defindMethod(pattern, handler, http.MethodTrace)
}

func (r *router) Put(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return r.defindMethod(pattern, handler, http.MethodPut)
}

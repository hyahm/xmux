package xmux

import (
	"log"
	"net/http"
)

// func (gr *RouteGroup) any(pattern string, handler func(http.ResponseWriter, *http.Request), methods ...string) MethodsRoute {

// 	temphead := make(map[string]string)
// 	for k, v := range gr.header {
// 		temphead[k] = v
// 	}

// 	tempPages := make(map[string]struct{})
// 	for k := range gr.pagekeys {
// 		tempPages[k] = struct{}{}
// 	}

// 	mr := make(map[string]*Route)
// 	for _, method := range methods {
// 		newRoute := &Route{
// 			handle:   http.HandlerFunc(handler),
// 			pagekeys: make(map[string]struct{}),
// 			module: &module{
// 				filter:    make(map[string]struct{}),
// 				funcOrder: make([]func(w http.ResponseWriter, r *http.Request) bool, 0),
// 				mu:        sync.RWMutex{},
// 			},
// 			new:         true,
// 			methods:     methods,
// 			url:         pattern,
// 			header:      make(map[string]string),
// 			delmodule:   make(map[string]struct{}),
// 			delPageKeys: make(map[string]struct{}),
// 			delheader:   make(map[string]struct{}),
// 			prefixs:     make([]string, 0),
// 			delprefix:   map[string]struct{}{},
// 		}
// 		if gr.module != nil {
// 			newRoute.module = gr.module.cloneMudule()
// 		} else {
// 			newRoute.module = &module{
// 				filter:    make(map[string]struct{}),
// 				funcOrder: make([]func(w http.ResponseWriter, r *http.Request) bool, 0),
// 				mu:        sync.RWMutex{},
// 			}
// 		}
// 		mr[method] = newRoute
// 		gr.routes = append(gr.routes, newRoute)
// 	}
// 	return mr
// 	// url, vars, ok := gr.makeRoute(pattern)
// 	// gr.params[url] = vars
// 	// if ok {
// 	// 	if gr.tpl[url] == nil {
// 	// 		gr.tpl[url] = make(MethodsRoute)
// 	// 	}
// 	// 	for _, method := range methods {
// 	// 		if _, methodOk := gr.tpl[url][method]; methodOk {
// 	// 			// 如果也存在， 那么method重复了
// 	// 			log.Fatal("method : " + method + "  duplicate, url: " + url)
// 	// 		}
// 	// 		if gr.tpl[url] == nil {
// 	// 			gr.tpl[url] = make(MethodsRoute)
// 	// 		}
// 	// 		// newRoute.methods[method] = struct{}{}
// 	// 		newRoute.url = url
// 	// 		newRoute.params = vars

// 	// 		gr.tpl[url][method] = newRoute
// 	// 	}
// 	// 	return gr.tpl[url]
// 	// } else {
// 	// 	if gr.route[url] == nil {
// 	// 		gr.route[url] = make(MethodsRoute)
// 	// 	}

// 	// 	// 如果存在就判断是否存在method
// 	// 	for _, method := range methods {
// 	// 		if _, methodOk := gr.route[url][method]; methodOk {
// 	// 			// 如果也存在， 那么method重复了
// 	// 			log.Fatal("method : " + method + "  duplicate, url: " + url)
// 	// 		}
// 	// 		if gr.route[url] == nil {
// 	// 			gr.route[url] = make(MethodsRoute)
// 	// 		}
// 	// 		newRoute.url = url
// 	// 		gr.route[url][method] = newRoute
// 	// 	}
// 	// 	// 如果不存在就创建一个 route
// 	// 	return gr.route[url]
// 	// }

// }

// get this route
func (gr *RouteGroup) defindMethod(pattern string, handler func(http.ResponseWriter, *http.Request), methods ...string) *Route {

	temphead := make(map[string]string)
	for k, v := range gr.header {
		temphead[k] = v
	}

	tempPages := make(map[string]struct{})
	for k := range gr.pagekeys {
		tempPages[k] = struct{}{}
	}
	newRoute := &Route{
		handle:      http.HandlerFunc(handler),
		pagekeys:    make(map[string]struct{}),
		new:         true,
		methods:     methods,
		header:      make(map[string]string),
		delmodule:   make(map[string]struct{}),
		delPageKeys: make(map[string]struct{}),
		delheader:   make(map[string]struct{}),
		prefixs:     make([]string, 0),
		delprefix:   map[string]struct{}{},
	}
	if gr.module != nil {
		newRoute.module = gr.module.cloneMudule()
	} else {
		newRoute.module = &module{
			filter:    make(map[string]struct{}),
			funcOrder: make([]func(w http.ResponseWriter, r *http.Request) bool, 0),
		}
	}
	url, vars, ok := makeRoute(pattern)
	if ok {
		if _, ok := gr.urlTpl[url]; ok {
			m, exsit := SliceExsit(gr.urlTpl[url].methods, newRoute.methods)
			if exsit {
				log.Fatal("method : " + m + "  duplicate, url: " + url)
			}
			// url 存在， 但是method不存在， 提醒下url 重复了，
			log.Fatalf("Found that the URL(%s) has multiple request methods. Please use Request method to merge the processing\n", url)
		}
		// for _, method := range methods {
		// 	if _, methodOk := gr.tpl[url][method]; methodOk {
		// 		// 如果也存在， 那么method重复了
		// 		log.Fatal("method : " + method + "  duplicate, url: " + url)
		// 	}
		// 	if gr.tpl[url] == nil {
		// 		gr.tpl[url] = make(MethodsRoute)
		// 	}
		// newRoute.methods[method] = struct{}{}
		newRoute.params = vars

		gr.urlTpl[url] = newRoute
		// }

	} else {
		if _, ok := gr.urlRoute[url]; ok {
			m, exsit := SliceExsit(gr.urlRoute[url].methods, newRoute.methods)
			if exsit {
				log.Fatal("method : " + m + "  duplicate, url: " + url)
			}
			// url 存在， 但是method不存在， 提醒下url 重复了，
			log.Fatalf("Found that the URL(%s) has multiple request methods. Please use Request method to merge the processing\n", url)
		}
		// 如果存在就判断是否存在method
		// for _, method := range methods {
		// 	if _, methodOk := gr.route[url][method]; methodOk {
		// 		// 如果也存在， 那么method重复了
		// 		log.Fatal("method : " + method + "  duplicate, url: " + url)
		// 	}
		gr.urlRoute[url] = newRoute
		// }
		// 	// 如果不存在就创建一个 route

	}
	return newRoute
}

func (gr *RouteGroup) Post(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !gr.new {
		panic("must be init by NewRouteGroup")
	}
	return gr.defindMethod(pattern, handler, http.MethodPost)
}

func (gr *RouteGroup) Any(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !gr.new {
		panic("must be init by NewRouteGroup")
	}
	return gr.defindMethod(pattern, handler, http.MethodConnect,
		http.MethodDelete, http.MethodGet, http.MethodHead,
		http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodTrace,
	)
}

func (gr *RouteGroup) Get(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !gr.new {
		panic("must be init by NewRouteGroup")
	}
	return gr.defindMethod(pattern, handler, http.MethodGet)
}

func (gr *RouteGroup) Request(pattern string, handler func(http.ResponseWriter, *http.Request), methods ...string) *Route {
	if !gr.new {
		panic("must be init by NewRouteGroup")
	}
	return gr.defindMethod(pattern, handler, methods...)
}

func (gr *RouteGroup) Delete(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !gr.new {
		panic("must be init by NewRouteGroup")
	}
	return gr.defindMethod(pattern, handler, http.MethodDelete)
}

func (gr *RouteGroup) Head(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !gr.new {
		panic("must be init by NewRouteGroup")
	}
	return gr.defindMethod(pattern, handler, http.MethodHead)
}

func (gr *RouteGroup) Options(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !gr.new {
		panic("must be init by NewRouteGroup")
	}
	return gr.defindMethod(pattern, handler, http.MethodOptions)
}

func (gr *RouteGroup) Connect(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !gr.new {
		panic("must be init by NewRouteGroup")
	}
	return gr.defindMethod(pattern, handler, http.MethodConnect)
}

func (gr *RouteGroup) Patch(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !gr.new {
		panic("must be init by NewRouteGroup")
	}
	return gr.defindMethod(pattern, handler, http.MethodPatch)
}

func (gr *RouteGroup) Trace(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !gr.new {
		panic("must be init by NewRouteGroup")
	}
	return gr.defindMethod(pattern, handler, http.MethodTrace)
}

func (gr *RouteGroup) Put(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	if !gr.new {
		panic("must be init by NewRouteGroup")
	}
	return gr.defindMethod(pattern, handler, http.MethodPut)
}

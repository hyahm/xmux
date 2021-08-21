package xmux

import (
	"log"
	"net/http"
)

// get this route
func (gr *GroupRoute) method(pattern string, handler func(http.ResponseWriter, *http.Request), method string) *Route {

	temphead := make(map[string]string)
	for k, v := range gr.header {
		temphead[k] = v
	}

	tempDelHead := make([]string, 0)

	tempDelHead = append(tempDelHead, gr.delheader...)

	tempPages := make(map[string]struct{})
	for k := range gr.pagekeys {
		tempPages[k] = struct{}{}
	}
	tempDelPageKeys := make([]string, 0)
	tempDelPageKeys = append(tempDelPageKeys, gr.delPageKeys...)
	newRoute := &Route{
		handle:      http.HandlerFunc(handler),
		midware:     gr.midware,
		pagekeys:    tempPages,
		delmidware:  gr.delmidware,
		module:      gr.module,
		delmodule:   gr.delmodule,
		header:      temphead,
		delheader:   tempDelHead,
		delPageKeys: tempDelPageKeys,
	}
	url, ok := gr.makeRoute(pattern)
	if ok {
		if _, urlOk := gr.tpl[url]; urlOk {
			// 如果存在就判断是否存在method
			if _, methodOk := gr.tpl[url][method]; methodOk {
				// 如果也存在， 那么method重复了
				log.Fatal("method : " + method + "  duplicate, url: " + url)
			} else {
				// 如果不存在就创建一个 route

				gr.tpl[url][method] = newRoute
				return newRoute
			}
		} else {
			// 如果不存在就创建一个PatternMR
			gr.tpl = make(PatternMR)
			mr := make(MethodsRoute)
			mr[method] = newRoute
			gr.tpl[url] = mr
			return newRoute
		}
		// 判断是不是正则
		// return gr.tpl[url].getRoute(url, method, handler, gr)
	} else {
		if _, urlOk := gr.route[url]; urlOk {
			// 如果存在就判断是否存在method
			if _, methodOk := gr.route[url][method]; methodOk {
				// 如果也存在， 那么method重复了
				log.Fatal("method : " + method + "  duplicate, url: " + url)
			} else {
				// 如果不存在就创建一个 route
				gr.route[url][method] = newRoute
				return newRoute
			}
		} else {
			// 如果不存在就创建一个PatternMR
			gr.route = make(PatternMR)
			mr := make(MethodsRoute)
			mr[method] = newRoute
			gr.route[url] = mr
			return newRoute
		}
	}
	return nil
}

// get this route
func (gr *GroupRoute) defindMethod(pattern string, handler func(http.ResponseWriter, *http.Request), methods ...string) MethodsRoute {
	if len(methods) == 0 {
		panic(pattern + " have not any methods")
	}
	mr := make(MethodsRoute)
	for _, method := range methods {
		mr[method] = gr.method(pattern, handler, method)
	}
	return mr
}

func (gr *GroupRoute) Post(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return gr.method(pattern, handler, http.MethodPost)
}

func (gr *GroupRoute) Any(pattern string, handler func(http.ResponseWriter, *http.Request)) MethodsRoute {
	return gr.defindMethod(pattern, handler, http.MethodConnect,
		http.MethodDelete, http.MethodGet, http.MethodHead,
		http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodTrace,
	)
}

// func (gr *GroupRoute) Requests(pattern string, handler func(http.ResponseWriter, *http.Request), methods ...string) *Route {
// 	return gr.method(pattern, handler, methods...)
// }

func (gr *GroupRoute) Get(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return gr.method(pattern, handler, http.MethodGet)
}

func (gr *GroupRoute) Request(pattern string, handler func(http.ResponseWriter, *http.Request), methods ...string) MethodsRoute {
	return gr.defindMethod(pattern, handler, methods...)
}

func (gr *GroupRoute) Delete(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return gr.method(pattern, handler, http.MethodDelete)
}

func (gr *GroupRoute) Head(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return gr.method(pattern, handler, http.MethodHead)
}

func (gr *GroupRoute) Options(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return gr.method(pattern, handler, http.MethodOptions)
}

func (gr *GroupRoute) Connect(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return gr.method(pattern, handler, http.MethodConnect)
}

func (gr *GroupRoute) Patch(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return gr.method(pattern, handler, http.MethodPatch)
}

func (gr *GroupRoute) Trace(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return gr.method(pattern, handler, http.MethodTrace)
}

func (gr *GroupRoute) Put(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return gr.method(pattern, handler, http.MethodPut)
}

package xmux

import (
	"errors"
	"log"
	"net/http"
)

var ErrMethodDuplicate = errors.New("method duplicate")

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
		module:      gr.module,
		delmodule:   gr.delmodule,
		header:      temphead,
		delheader:   tempDelHead,
		delPageKeys: tempDelPageKeys,
	}
	url, ok := gr.makeRoute(pattern)
	if ok {
		if _, urlok := gr.tpl[url]; urlok {
			// 如果存在就判断是否存在method
			if _, methodok := gr.tpl[url][method]; methodok {
				// 如果也存在， 那么method重复了
				log.Fatal(ErrMethodDuplicate)
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
		if _, urlok := gr.route[url]; urlok {
			// 如果存在就判断是否存在method
			if _, methodok := gr.route[url][method]; methodok {
				// 如果也存在， 那么method重复了
				log.Fatal(ErrMethodDuplicate)
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

func (gr *GroupRoute) Post(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return gr.method(pattern, handler, http.MethodPost)
}

func (gr *GroupRoute) Any(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return gr.method(pattern, handler, "*")
}

// func (gr *GroupRoute) Requests(pattern string, handler func(http.ResponseWriter, *http.Request), methods ...string) *Route {
// 	return gr.method(pattern, handler, methods...)
// }

func (gr *GroupRoute) Get(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return gr.method(pattern, handler, http.MethodGet)
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

package xmux

import (
	"log"
	"net/http"
	"regexp"
)

// get this route
func (gr *routeGroup) defindMethod(pattern string, handler func(http.ResponseWriter, *http.Request), methods ...string) *route {

	temphead := make(map[string]string)
	for k, v := range gr.header {
		temphead[k] = v
	}

	tempPages := make(map[string]struct{})
	for k := range gr.pagekeys {
		tempPages[k] = struct{}{}
	}
	newRoute := &route{
		handle:        http.HandlerFunc(handler),
		pagekeys:      make(map[string]struct{}),
		methods:       methods,
		module:        gr.module.cloneMudule(),
		postModule:    gr.postModule.cloneMudule(),
		header:        make(map[string]string),
		delmodule:     make(map[string]struct{}),
		delPostModule: make(map[string]struct{}),
		delPageKeys:   make(map[string]struct{}),
		delheader:     make(map[string]struct{}),
		prefixs:       gr.prefix,
		delprefix:     gr.delprefix,
		denyPrefix:    gr.denyPrefix,
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
		newRoute.regex = regexp.MustCompile(url)
		if _, ok := gr.urlTpl[url]; ok {
			m, exsit := SliceExist(gr.urlTpl[url].methods, newRoute.methods)
			if exsit {
				log.Fatal("method : " + m + "  duplicate, url: " + url)
			}
			// url 存在， 但是method不存在， 提醒下url 重复了，
			log.Fatalf("Found that the URL(%s) has multiple request methods. Please use Request method to merge the processing\n", url)
		}
		newRoute.params = vars
		gr.urlTpl[url] = newRoute
		// }

	} else {
		if _, ok := gr.urlRoute[url]; ok {
			m, exsit := SliceExist(gr.urlRoute[url].methods, newRoute.methods)
			if exsit {
				log.Fatal("method : " + m + "  duplicate, url: " + url)
			}
			// url 存在， 但是method不存在， 提醒下url 重复了，
			log.Fatalf("Found that the URL(%s) has multiple request methods. Please use Request method to merge the processing\n", url)
		}

		gr.urlRoute[url] = newRoute
		// 	// 如果不存在就创建一个 route

	}

	if enableRouterTree {
		// 启用了路由树， 生成路由树
		gr.routerTrees.Metas = append(gr.routerTrees.Metas, Meta{
			Url:     url,
			Methods: newRoute.methods,
		})
	}
	return newRoute
}

func (gr *routeGroup) Post(pattern string, handler func(http.ResponseWriter, *http.Request)) *route {
	return gr.defindMethod(pattern, handler, http.MethodPost)
}

func (gr *routeGroup) Any(pattern string, handler func(http.ResponseWriter, *http.Request)) *route {
	return gr.defindMethod(pattern, handler, http.MethodConnect,
		http.MethodDelete, http.MethodGet, http.MethodHead,
		http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodTrace,
	)
}

func (gr *routeGroup) Get(pattern string, handler func(http.ResponseWriter, *http.Request)) *route {
	return gr.defindMethod(pattern, handler, http.MethodGet)
}

func (gr *routeGroup) Connect(pattern string, handler func(http.ResponseWriter, *http.Request)) *route {
	return gr.defindMethod(pattern, handler, http.MethodConnect)
}

func (gr *routeGroup) Request(pattern string, handler func(http.ResponseWriter, *http.Request), methods ...string) *route {
	return gr.defindMethod(pattern, handler, methods...)
}

func (gr *routeGroup) Delete(pattern string, handler func(http.ResponseWriter, *http.Request)) *route {
	return gr.defindMethod(pattern, handler, http.MethodDelete)
}

func (gr *routeGroup) Head(pattern string, handler func(http.ResponseWriter, *http.Request)) *route {
	return gr.defindMethod(pattern, handler, http.MethodHead)
}

func (gr *routeGroup) Options(pattern string, handler func(http.ResponseWriter, *http.Request)) *route {
	return gr.defindMethod(pattern, handler, http.MethodOptions)
}

func (gr *routeGroup) Patch(pattern string, handler func(http.ResponseWriter, *http.Request)) *route {
	return gr.defindMethod(pattern, handler, http.MethodPatch)
}

func (gr *routeGroup) Trace(pattern string, handler func(http.ResponseWriter, *http.Request)) *route {
	return gr.defindMethod(pattern, handler, http.MethodTrace)
}

func (gr *routeGroup) Put(pattern string, handler func(http.ResponseWriter, *http.Request)) *route {
	return gr.defindMethod(pattern, handler, http.MethodPut)
}

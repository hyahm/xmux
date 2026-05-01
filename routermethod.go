package xmux

import (
	"log"
	"net/http"
	"path"
	"regexp"

	"github.com/google/uuid"
)

// get this route
func (r *router) defindMethod(pattern string, handler func(http.ResponseWriter, *http.Request), method ...string) *route {
	// 			exsitOption = true
	// 			break
	// 		}
	// 	}
	// 	if !exsitOption {
	// 		method = append(method, http.MethodOptions)
	// 	}
	// }
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
	newRoute := &route{
		handle:        http.HandlerFunc(handler),
		pagekeys:      tempPages,
		module:        r.module.cloneMudule(),
		postModule:    r.postModule.cloneMudule(),
		responseData:  r.responseData,
		methods:       append(make([]string, 0), method...),
		header:        temphead,
		uuid:          uuid.New().String(),
		delheader:     make(map[string]struct{}),
		delmodule:     make(map[string]struct{}),
		delPostModule: make(map[string]struct{}),
		delPageKeys:   make(map[string]struct{}),
		prefixs:       make([]string, 0),
		delprefix:     map[string]struct{}{},
		parentUuid:    "root",
		middleware: onion{
			mws: make([]Middleware, 0),
		},
	}
	// prefix := path.Join(r.prefix...)
	// prefix = path.Join(prefix, pattern)
	// 判断是否是正则
	url, vars, ok := parsePath(pattern)
	newRoute.params = vars
	if ok {
		// 预编译正则
		newRoute.regex = regexp.MustCompile(url)
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

func (r *router) Post(pattern string, handler func(http.ResponseWriter, *http.Request)) *route {
	return r.defindMethod(pattern, handler, http.MethodPost)
}

func (r *router) Any(pattern string, handler func(http.ResponseWriter, *http.Request)) *route {
	return r.defindMethod(pattern, handler,
		http.MethodDelete, http.MethodGet, http.MethodHead,
		http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodTrace,
	)
}

func (r *router) Get(pattern string, handler func(http.ResponseWriter, *http.Request)) *route {

	return r.defindMethod(pattern, handler, http.MethodGet)
}

func (r *router) Connect(pattern string, handler func(http.ResponseWriter, *http.Request)) *route {

	return r.defindMethod(pattern, handler, http.MethodConnect)
}

func (r *router) Request(pattern string, handler func(http.ResponseWriter, *http.Request), methods ...string) *route {
	return r.defindMethod(pattern, handler, methods...)
}

func (r *router) Delete(pattern string, handler func(http.ResponseWriter, *http.Request)) *route {
	return r.defindMethod(pattern, handler, http.MethodDelete)
}

func (r *router) Head(pattern string, handler func(http.ResponseWriter, *http.Request)) *route {
	return r.defindMethod(pattern, handler, http.MethodHead)
}

func (r *router) Options(pattern string, handler func(http.ResponseWriter, *http.Request)) *route {
	return r.defindMethod(pattern, handler, http.MethodOptions)
}

func (r *router) Patch(pattern string, handler func(http.ResponseWriter, *http.Request)) *route {
	return r.defindMethod(pattern, handler, http.MethodPatch)
}

func (r *router) Trace(pattern string, handler func(http.ResponseWriter, *http.Request)) *route {
	return r.defindMethod(pattern, handler, http.MethodTrace)
}

func (r *router) Put(pattern string, handler func(http.ResponseWriter, *http.Request)) *route {
	return r.defindMethod(pattern, handler, http.MethodPut)
}

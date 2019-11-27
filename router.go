package xmux

import (
	"net/http"
	"regexp"
	"strings"
)

var Var map[string]map[string]string

func init() {
	// 存储变量
	Var = make(map[string]map[string]string)
	// 存储正则的地址
	reUrl = make(map[string]*reroute)
}

type reroute struct {
	R      *Route
	name   []string
	header map[string]string
}

type rt struct {
	Handle http.Handler
	Header map[string]string
}

type Router struct {
	g              map[string]*GroupRoute // 组路由
	s              map[string]*Route      // 单实例路由
	IgnoreIco      bool                   // 是否忽略 /favicon.ico 请求。 默认忽略
	Options        http.Handler           // 预请求 处理函数， 如果存在， 优先处理, 前后端分离后， 前段可能会先发送一个预请求
	NotFound       http.Handler
	MethodNotAllow http.Handler
	HandleNotFound http.Handler
	groupKey       map[string]bool // 组路由
	routeTable     map[string]*rt  // 路由表
	header         map[string]string
}

func (r *Router) Group(patter string) *GroupRoute {
	//   /article if /article/ to /article;  if article to /article
	if patter[0:1] != "/" {
		patter = "/" + patter
	}
	if patter[len(patter)-1:len(patter)] == "/" {
		patter = patter[:len(patter)-1]
	}
	g := &GroupRoute{
		prefix: patter,
		suffix: make(map[string]*Route),
		header: make(map[string]string),
	}

	r.g[patter] = g
	r.groupKey[patter] = true
	return g
}

func (r *Router) AddGroup(groute *GroupRoute) *Router {
	r.g[groute.prefix] = groute
	r.groupKey[groute.prefix] = true
	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	if r.Options != nil && req.Method == http.MethodOptions {
		r.Options.ServeHTTP(w, req)
		return
	}

	key := req.URL.Path
	if r.IgnoreIco && key == "/favicon.ico" {
		return
	}

	// 先进行路由表缓存寻找
	if route, ok := r.routeTable[key+req.Method]; ok {
		for k, v := range route.Header {
			w.Header().Set(k, v)
		}
		route.Handle.ServeHTTP(w, req)
		return
	}
	// 先判断有几段
	i := strings.Count(key, "/")
	// 分出一级路径

	if i > 1 {
		end := strings.Index(key[1:], "/")
		first_key := key[:end+1]
		// 判断是不是组成员
		if _, ok := r.groupKey[first_key]; ok {
			//如果存在就是组成员， 继续判断二段路径是否存在
			if route, subok := r.g[first_key].suffix[key[end+2:]]; subok {
				if handle, metok := route.method[req.Method]; metok {

					r.routeTable[key+req.Method] = &rt{
						Handle: handle,
						Header: route.header,
					}
					for k, v := range route.header {
						w.Header().Set(k, v)
					}
					handle.ServeHTTP(w, req)
					return
				} else {
					r.routeTable[key+req.Method] = &rt{
						Handle: r.HandleNotFound,
					}
					r.HandleNotFound.ServeHTTP(w, req)
					return
				}
			}
		}

	}
	// 单一的路径，  不是组成员
	if route, ok := r.s[key]; ok {
		if handle, metok := route.method[req.Method]; metok {
			r.routeTable[key+req.Method] = &rt{
				Handle: handle,
				Header: route.header,
			}
			for k, v := range route.header {
				w.Header().Add(k, v)
			}
			handle.ServeHTTP(w, req)
			return
		} else {
			r.routeTable[key+req.Method] = &rt{
				Handle: r.HandleNotFound,
			}
			r.HandleNotFound.ServeHTTP(w, req)
			return
		}
	}
	// 最后正则里面寻找路由
	for k, route := range reUrl {
		re := regexp.MustCompile(k)
		if re.MatchString(key) {
			// 获取var
			x := re.FindStringSubmatch(key)
			myvar := make(map[string]string)
			for i, v := range route.name {
				myvar[v] = x[i+1]
			}
			Var[key] = myvar
			if handle, metok := route.R.method[req.Method]; metok {
				r.routeTable[key+req.Method] = &rt{
					Handle: handle,
					Header: route.header,
				}
				for k, v := range route.header {
					w.Header().Set(k, v)
				}
				handle.ServeHTTP(w, req)
				return
			} else {
				r.routeTable[key+req.Method] = &rt{
					Handle: r.HandleNotFound,
				}
				r.HandleNotFound.ServeHTTP(w, req)
				return
			}
		}

	}

	r.routeTable[key+req.Method] = &rt{
		Handle: r.NotFound,
	}
	r.NotFound.ServeHTTP(w, req)
	return

}

// 组里面也包括路由 后面的其实还是patter和handle
func (r *Router) HandleFunc(pattern string) *Route {
	if pattern[0:1] != "/" {
		pattern = "/" + pattern
	}
	if pattern[len(pattern)-1:len(pattern)] == "/" {
		pattern = pattern[:len(pattern)-1]
	}

	route := &Route{
		method: make(map[string]http.Handler),
		header: make(map[string]string),
	}
	lv := make([]string, 0)
	if v, listvar, ok := match(pattern, "^", lv); ok {
		reUrl[v] = &reroute{
			R:      route,
			name:   listvar,
			header: route.header,
		}
		return route
	}
	r.s[pattern] = route
	return route
}

func NewRouter() *Router {
	return &Router{
		g:              make(map[string]*GroupRoute),
		s:              make(map[string]*Route),
		IgnoreIco:      true,
		Options:        options(),
		NotFound:       notFound(),
		MethodNotAllow: methodNotAllowed(),
		groupKey:       make(map[string]bool),
		HandleNotFound: notHandle(),
		routeTable:     make(map[string]*rt),
		header:         make(map[string]string),
	}
}

func methodNotAllowed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	})
}

func notFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		return
	})
}

func options() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})
}

func notHandle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not found handle"))
		return
	})
}

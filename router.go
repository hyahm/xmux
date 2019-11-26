package xmux

import (
	"net/http"
	"strings"
)

type Router struct {
	G              map[string]*GroupRoute // 组路由
	S              map[string]*Route      // 单实例路由
	IgnoreIco      bool                   // 是否忽略 /favicon.ico 请求。 默认忽略
	Options        http.Handler           // 预请求 处理函数， 如果存在， 优先处理, 前后端分离后， 前段可能会先发送一个预请求
	NotFound       http.Handler
	MethodNotAllow http.Handler
	GroupKey       map[string]bool
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
	}

	r.G[patter] = g
	r.GroupKey[patter] = true
	return g
}

func (r *Router) AddGroup(groute *GroupRoute) *Router {
	r.G[groute.prefix] = groute
	r.GroupKey[groute.prefix] = true
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
	// 先判断有几段
	i := strings.Count(key, "/")
	// 分出一级路径

	if i > 1 {
		end := strings.Index(key[1:], "/")
		first_key := key[:end+1]
		// 判断是不是组成员
		if _, ok := r.GroupKey[first_key]; ok {
			//如果存在就是组成员， 继续判断二段路径是否存在
			if route, subok := r.G[first_key].suffix[key[end+2:]]; subok {
				if route.allHandle == nil {
					// 如果不存在不限制的路由， 那么一定限制了method
					if handle, metok := route.method[req.Method]; metok {
						handle.ServeHTTP(w, req)
						return
					}
				} else {
					route.allHandle.ServeHTTP(w, req)
					return
				}
			} else {
				r.NotFound.ServeHTTP(w, req)
				return
			}
		}

	}
	// 单一的路径，  不是组成员
	if route, ok := r.S[key]; ok {
		if route.allHandle == nil {
			if handle, metok := route.method[req.Method]; metok {
				handle.ServeHTTP(w, req)
				return
			} else {
				r.MethodNotAllow.ServeHTTP(w, req)
				return
			}
		} else {
			route.allHandle.ServeHTTP(w, req)
			return
		}
	} else {
		// 没找到路由
		r.NotFound.ServeHTTP(w, req)
		return
	}
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
	}
	r.S[pattern] = route
	return route
}

func NewRouter() *Router {
	return &Router{
		G:              make(map[string]*GroupRoute),
		S:              make(map[string]*Route),
		IgnoreIco:      true,
		Options:        options(),
		NotFound:       notFound(),
		MethodNotAllow: methodNotAllowed(),
		GroupKey:       make(map[string]bool),
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

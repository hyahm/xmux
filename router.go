package xmux

import (
	"net/http"
	"regexp"
	"sync"
)

var Var map[string]map[string]string

func init() {
	// 存储变量
	Var = make(map[string]map[string]string)
}

type reroute struct {
	R      *Route
	name   []string // 保存的变量名
	header map[string]string
}

type rt struct {
	Handle http.Handler
	Header map[string]string
}

type Router struct {
	IgnoreIco        bool         // 是否忽略 /favicon.ico 请求。 默认忽略
	DisableOption    bool         // 禁止全局option
	Slash            bool         // 是否检测请求的url
	Options          http.Handler // 预请求 处理函数， 如果存在， 优先处理, 前后端分离后， 前段可能会先发送一个预请求
	NotFound         http.Handler
	HandleNotFound   http.Handler
	MethodNotAllowed http.Handler
	route            map[string]*Route            // 单实例路由
	groupKey         map[string]map[string]string // 组路由, 存的组路由的请求头
	routeTable       *sync.Map                    // 路由表
	header           map[string]string            // 全局路由头
	tpl              map[string]*Route            // 正则路由
}

func (r *Router) SetHeader(k, v string) *Router {
	if r.header == nil {
		panic("please use xmux.NewRouter()")
	}
	r.header[k] = v
	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.URL.Path
	// 格式路径
	if r.Slash {
		key = slash(key)
	}

	if r.IgnoreIco && key == "/favicon.ico" {
		return
	}

	// 先进行路由表缓存寻找
	if route, ok := r.routeTable.Load(key + req.Method); ok {
		for k, v := range route.(*rt).Header {
			w.Header().Set(k, v)
		}
		route.(*rt).Handle.ServeHTTP(w, req)
		return
	}
	tmpheader := make(map[string]string)

	for k, v := range r.header {
		tmpheader[k] = v
		w.Header().Set(k, v)
	}
	var thisHandle http.Handler
	// option 请求处理
	if !r.DisableOption && req.Method == http.MethodOptions {
		if r.Options == nil {
			r.Options = options()
		}
		r.Options.ServeHTTP(w, req)
		return

	}
	// 先匹配
	if _, ok := r.route[key]; ok {
		// 判断是不是组成员
		if header, gok := r.groupKey[key]; gok {
			//是组成员的话， 3头合一
			for k, v := range header {
				tmpheader[k] = v
				w.Header().Set(k, v)
			}
		}
		//然后就是本身的头
		for k, v := range r.route[key].header {
			tmpheader[k] = v
			w.Header().Set(k, v)
		}
		// 是否能找到方法
		if handle, metok := r.route[key].method[req.Method]; metok {
			thisHandle = handle
		} else {
			if r.route[key] != nil {
				thisHandle = r.MethodNotAllowed
			} else {
				if r.HandleNotFound == nil {
					thisHandle = handleNotFound()
				} else {
					thisHandle = r.HandleNotFound
				}
			}
		}
	} else {
		// 最后正则里面寻找路由
		for reurl, route := range r.tpl {
			re := regexp.MustCompile(reurl)
			if re.MatchString(key) {
				vl := re.FindStringSubmatch(key)
				vm := make(map[string]string)
				for i, v := range route.args {
					vm[v] = vl[i+1]
				}
				Var[key] = vm
				// 获取var
				//判断是不是组路由

				if header, ok := r.groupKey[reurl]; ok {
					for k, v := range header {
						tmpheader[k] = v
						w.Header().Set(k, v)
					}
				}
				//然后就是本身的头
				for k, v := range route.header {
					tmpheader[k] = v
					w.Header().Set(k, v)
				}

				// 是否能找到方法
				if handle, metok := route.method[req.Method]; metok {
					//保存到路由表
					thisHandle = handle
				} else {
					if r.route[key] != nil {
						thisHandle = r.MethodNotAllowed
					} else {
						if r.HandleNotFound == nil {
							thisHandle = handleNotFound()
						} else {
							thisHandle = r.HandleNotFound
						}
					}
				}
			}

		}
		if thisHandle == nil {
			if r.NotFound == nil {
				thisHandle = notFound()
			} else {
				thisHandle = r.NotFound
			}
		}

	}

	r.routeTable.Store(key+req.Method, &rt{
		Handle: thisHandle,
		Header: tmpheader,
	})
	thisHandle.ServeHTTP(w, req)
	return

}

func NewRouter() *Router {
	return &Router{
		IgnoreIco:        true,
		Slash:            false,
		Options:          options(),
		NotFound:         notFound(),
		HandleNotFound:   handleNotFound(),
		MethodNotAllowed: methodNotAllowed(),
		groupKey:         make(map[string]map[string]string),
		routeTable:       &sync.Map{},
		header:           make(map[string]string),
		route:            make(map[string]*Route),
		tpl:              make(map[string]*Route),
	}
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

func handleNotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1>when you see this page, it means you forget set handle in " + r.URL.Path + "<h1>"))
		return
	})
}

func methodNotAllowed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	})
}

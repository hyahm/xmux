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

// type Midware

type reroute struct {
	R      *Route
	name   []string // 保存的变量名
	header map[string]string
}

type rt struct {
	Handle  http.Handler
	Header  map[string]string
	Midware []func(http.ResponseWriter, *http.Request) bool
}

type Router struct {
	IgnoreIco        bool // 是否忽略 /favicon.ico 请求。 默认忽略
	HanleFavicon     http.Handler
	DisableOption    bool         // 禁止全局option
	Options          http.Handler // 预请求 处理函数， 如果存在， 优先处理, 前后端分离后， 前段可能会先发送一个预请求
	UrlNotFound      http.Handler
	HandleNotFound   http.Handler
	MethodNotAllowed http.Handler
	Slash            bool
	route            map[string]*Route // 单实例路由
	tpl              map[string]*Route // 正则路由

	group map[string]*GroupRoute // 组路由

	//  标记所有的pattern， 防止有重复的pattern， 0: route 1, tpl, 2, groupRouter 3, groupTpl

	pattern    map[string]int // 完全匹配
	tplpattern map[string]int // 正则匹配

	groupname map[string]string // 根据 pattern 寻找 组
	// cacheMidware     map[string][]http.Handler    // 组路由, 存的组路由的请求头
	header  map[string]string                               // 全局路由头
	midware []func(http.ResponseWriter, *http.Request) bool // 全局中间件

	routeTable *sync.Map // 路由表

	once *sync.Once
}

func (r *Router) SetHeader(k, v string) *Router {
	if r.header == nil {
		r.header = make(map[string]string)
	}
	r.header[k] = v
	return r
}

func (r *Router) AddMidware(handle func(http.ResponseWriter, *http.Request) bool) *Router {
	if r.midware == nil {
		r.midware = make([]func(http.ResponseWriter, *http.Request) bool, 0)
	}
	r.midware = append(r.midware, handle)
	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.once == nil {
		r.once = &sync.Once{}
	}
	if r.routeTable == nil {
		r.routeTable = &sync.Map{}
	}
	r.once.Do(func() {
		// 初始化默认路由
		r.initHandler()
	})

	// 去掉路径多余的斜杠

	url := req.URL.Path
	if r.Slash {
		url = slash(req.URL.Path)
	}

	// /favicon.ico  和 Option 请求， 不支持自定义请求头和中间件
	if r.IgnoreIco && url == "/favicon.ico" {
		r.HanleFavicon.ServeHTTP(w, req)
		return
	}

	// option 请求处理
	if !r.DisableOption && req.Method == http.MethodOptions {
		r.Options.ServeHTTP(w, req)
		return
	}

	// 先进行路由表缓存寻找
	if route, ok := r.routeTable.Load(url + req.Method); ok {
		// 设置请求头
		for k, v := range route.(*rt).Header {
			w.Header().Set(k, v)
		}
		// 请求中间件
		for _, v := range route.(*rt).Midware {
			v(w, req)
		}
		route.(*rt).Handle.ServeHTTP(w, req)

	} else {
		// 获取handler
		// fmt.Println("no cached")
		r.serveHTTP(url, w, req)
	}
}

func (r *Router) serveHTTP(url string, w http.ResponseWriter, req *http.Request) {
	// 应该弄成中间件形式
	var thisHandle http.Handler
	var tp int = -1
	var matchurl string
	///  寻找路由   ///
	// 先寻找完全匹配的
	if _, ok := r.pattern[url]; ok {
		matchurl = url
		if r.pattern[url] == 0 {
			// 匹配的单路由
			// 是否能找到方法
			if handle, mok := r.route[url].method[req.Method]; mok {
				tp = 0
				thisHandle = handle
			} else {
				if r.route[url] != nil {
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
			// r.pattern[url] 肯定等于 2
			group := r.group[r.groupname[matchurl]].route[url]
			if handle, mok := group.method[req.Method]; mok {
				tp = 2
				thisHandle = handle
			} else {
				if group != nil {
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

	} else {
		// 最后正则里面寻找路由

		for reUrl, n := range r.tplpattern {
			re := regexp.MustCompile(reUrl)
			if re.MatchString(url) {
				matchurl = reUrl
				vl := re.FindStringSubmatch(url)
				vm := make(map[string]string)
				if n == 1 {
					// 单路由
					if handle, mok := r.tpl[matchurl].method[req.Method]; mok {
						for i, v := range r.tpl[matchurl].args {
							vm[v] = vl[i+1]
						}
						// 获取var
						if r.Slash {
							Var[url] = vm
						} else {
							Var[req.URL.Path] = vm
						}
						tp = 1
						thisHandle = handle
						goto endloop
					} else {
						if r.route[url] != nil {
							thisHandle = r.MethodNotAllowed
						} else {
							if r.HandleNotFound == nil {
								thisHandle = handleNotFound()
							} else {
								thisHandle = r.HandleNotFound
							}
						}
						goto endloop
					}
				} else {
					group := r.group[r.groupname[matchurl]].tpl[matchurl]
					if handle, mok := group.method[req.Method]; mok {
						for i, v := range group.args {
							vm[v] = vl[i+1]
						}
						// 获取var
						if r.Slash {
							Var[url] = vm
						} else {
							Var[req.URL.Path] = vm
						}
						tp = 3
						thisHandle = handle
						goto endloop
					} else {
						if group != nil {
							thisHandle = r.MethodNotAllowed
						} else {
							if r.HandleNotFound == nil {
								thisHandle = handleNotFound()
							} else {
								thisHandle = r.HandleNotFound
							}
						}
						goto endloop
					}
				}

			}

		}
		// 没有匹配到
		thisHandle = r.HandleNotFound
	}
endloop:
	tmpHeader := make(map[string]string)
	for k, v := range r.header {
		tmpHeader[k] = v
	}

	tmpMidware := make([]func(http.ResponseWriter, *http.Request) bool, 0)

	for _, v := range r.midware {
		tmpMidware = append(tmpMidware, v)
	}
	///  结束寻找路由     ///
	// 合并请求头和中间件
	switch tp {
	case 0, 2:
		if tp == 2 {
			group := r.group[r.groupname[url]].route[url]
			for _, v := range group.midware {
				tmpMidware = append(tmpMidware, v)
			}
			for k, v := range group.header {
				tmpHeader[k] = v
				w.Header().Set(k, v)
			}
		}
		for _, v := range r.route[url].midware {
			tmpMidware = append(tmpMidware, v)
		}
		for k, v := range r.route[url].header {
			tmpHeader[k] = v
			w.Header().Set(k, v)
		}
	case 1, 3:
		if tp == 3 {
			group := r.group[r.groupname[url]].tpl[matchurl]
			for _, v := range group.midware {
				tmpMidware = append(tmpMidware, v)
			}
			for k, v := range group.header {
				tmpHeader[k] = v
				w.Header().Set(k, v)
			}
		}
		for _, v := range r.tpl[matchurl].midware {
			tmpMidware = append(tmpMidware, v)
		}
		for k, v := range r.tpl[matchurl].header {
			tmpHeader[k] = v
			w.Header().Set(k, v)
		}

	default:
	}
	// 执行 中间件
	for _, v := range tmpMidware {
		var ok bool
		ok = v(w, req)
		if ok {
			return
		}
	}

	// 缓存handler

	if r.Slash && req.URL.Path != url {
		r.routeTable.Store(req.URL.Path+req.Method, &rt{
			Handle:  thisHandle,
			Header:  tmpHeader,
			Midware: tmpMidware,
		})

	}

	r.routeTable.Store(url+req.Method, &rt{
		Handle:  thisHandle,
		Header:  tmpHeader,
		Midware: tmpMidware,
	})

	thisHandle.ServeHTTP(w, req)
}

func (r *Router) initHandler() {
	// 匹配完成后，最先执行这个， 初始化当前方法
	if r.MethodNotAllowed == nil {
		r.MethodNotAllowed = methodNotAllowed()
	}

	if r.HandleNotFound == nil {
		r.HandleNotFound = handleNotFound()
	}

	if r.HanleFavicon == nil {
		r.HanleFavicon = favicon()
	}

	if r.Options == nil {
		r.Options = options()
	}

}

func NewRouter() *Router {
	return &Router{
		IgnoreIco: true,
		// group:      make(map[string]map[string]string),
		Slash:      true,
		routeTable: &sync.Map{},
		header:     make(map[string]string),
		route:      make(map[string]*Route),
		tpl:        make(map[string]*Route),
		once:       &sync.Once{},
	}
}

func urlNotFound() http.Handler {
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

func favicon() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})
}

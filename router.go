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
	IgnoreIco        bool // 是否忽略 /favicon.ico 请求。 默认忽略
	HanleFavicon     http.Handler
	DisableOption    bool         // 禁止全局option
	Slash            bool         // 是否检测请求的url， 如果是， 匹配前自动去除多余的斜杠
	Options          http.Handler // 预请求 处理函数， 如果存在， 优先处理, 前后端分离后， 前段可能会先发送一个预请求
	UrlNotFound      http.Handler
	HandleNotFound   http.Handler
	MethodNotAllowed http.Handler
	route            map[string]*Route            // 单实例路由
	groupKey         map[string]map[string]string // 组路由, 存的组路由的请求头
	routeTable       *sync.Map                    // 路由表
	header           map[string]string            // 全局路由头
	tpl              map[string]*Route            // 正则路由
	midware          map[string][]http.Handler    // 全局中间件
	once             sync.Once
}

func (r *Router) SetHeader(k, v string) *Router {
	if r.header == nil {
		panic("please use xmux.NewRouter()")
	}
	r.header[k] = v
	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	r.once.Do(func() {
		r.initHandler()
	})
	url := req.URL.Path
	if r.Slash {
		url = slash(url)
	}
	// 更新请求头
	tmpHeader := make(map[string]string)

	for k, v := range r.header {
		tmpHeader[k] = v
		w.Header().Set(k, v)
	}
	// tmpHeader := r.addHeader(url, w)

	if r.assetHandler(url, w, req) {
		return
	}
	// 中间件预留

	// 获取handler
	// 有没有过来
	r.serveHTTP(url, tmpHeader, w, req)
	return

}

// 这个作为中间件测试
func (r *Router) assetHandler(url string, w http.ResponseWriter, req *http.Request) bool {
	if r.IgnoreIco && url == "/favicon.ico" {
		r.HanleFavicon.ServeHTTP(w, req)
		return true
	}

	// 先进行路由表缓存寻找
	if route, ok := r.routeTable.Load(url + req.Method); ok {
		route.(*rt).Handle.ServeHTTP(w, req)
		return true
	}

	// option 请求处理
	if !r.DisableOption && req.Method == http.MethodOptions {
		r.Options.ServeHTTP(w, req)
		return true
	}
	return false
}

func (r *Router) serveHTTP(url string, tmpHeader map[string]string, w http.ResponseWriter, req *http.Request) {
	// 应该弄成中间件形式

	var thisHandle http.Handler
	// 先寻找完全匹配的
	if _, ok := r.route[url]; ok {
		// 是否能找到方法
		if handle, mok := r.route[url].method[req.Method]; mok {
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
		// 最后正则里面寻找路由
		for reUrl, route := range r.tpl {
			re := regexp.MustCompile(reUrl)
			if re.MatchString(url) {
				vl := re.FindStringSubmatch(url)
				vm := make(map[string]string)
				for i, v := range route.args {
					vm[v] = vl[i+1]
				}
				// 获取var
				Var[url] = vm

				// 是否能找到方法
				if handle, mok := route.method[req.Method]; mok {
					//保存到路由表
					thisHandle = handle
					goto endloop
				} else {
					if r.route[url] != nil {
						thisHandle = r.MethodNotAllowed
						goto endloop
					}
				}
			}

		}
		// 没有匹配到
		thisHandle = r.HandleNotFound

	}
endloop:
	if header, gok := r.groupKey[url]; gok {
		//是组成员的话， 3头合一
		for k, v := range header {
			tmpHeader[k] = v
			w.Header().Set(k, v)
		}
	}
	//然后就是本身的头
	if r.route[url] != nil {
		for k, v := range r.route[url].header {
			tmpHeader[k] = v
			w.Header().Set(k, v)
		}
	}

	// 缓存handler

	r.routeTable.Store(url+req.Method, &rt{
		Handle: thisHandle,
		Header: tmpHeader,
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
		IgnoreIco:  true,
		Slash:      false,
		groupKey:   make(map[string]map[string]string),
		routeTable: &sync.Map{},
		header:     make(map[string]string),
		route:      make(map[string]*Route),
		tpl:        make(map[string]*Route),
		once:       sync.Once{},
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

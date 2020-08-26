package xmux

import (
	"log"
	"net/http"
	"regexp"
	"sync"
)

var connections int

func GetConnents() int {
	return connections
}

type reroute struct {
	R      *Route
	name   []string // 保存的变量名
	header map[string]string
}

type rt struct {
	Handle     http.Handler
	Header     map[string]string
	Module     []func(http.ResponseWriter, *http.Request) bool
	dataSource interface{}
}

type Router struct {
	IgnoreIco      bool // 是否忽略 /favicon.ico 请求。 默认忽略
	HanleFavicon   http.Handler
	DisableOption  bool         // 禁止全局option
	HandleOptions  http.Handler // 预请求 处理函数， 如果存在， 优先处理, 前后端分离后， 前段可能会先发送一个预请求
	HandleNotFound http.Handler
	Slash          bool
	route          PatternMR                                       // 单实例路由， 组路由最后也会合并过来
	tpl            PatternMR                                       // 正则路由， 组路由最后也会合并过来
	pattern        map[string][]string                             // 记录所有路由， []string 是正则匹配的参数
	header         map[string]string                               // 全局路由头
	module         []func(http.ResponseWriter, *http.Request) bool // 全局模块
	routeTable     map[string]*rt                                  // 路由表
	mu             *sync.RWMutex                                   // 传递参数的锁
	rtLock         *sync.RWMutex                                   // 缓存表的锁
	end            func(http.ResponseWriter, *http.Request)
}

func (r *Router) makeRoute(pattern string) (string, bool) {
	// 格式化路径
	// 创建 methodsRoute
	// 格式路径
	if r.mu == nil {
		r.mu = &sync.RWMutex{}
	}
	if r.rtLock == nil {
		r.rtLock = &sync.RWMutex{}
	}
	if r.routeTable == nil {
		r.routeTable = make(map[string]*rt)
	}
	if r.Slash {
		pattern = slash(pattern)
	}
	if r.route == nil {
		r.route = make(map[string]MethodsRoute)
	}
	if r.pattern == nil {
		r.pattern = make(map[string][]string)
	}

	if v, listvar := match(pattern); len(listvar) > 0 {
		if _, ok := r.pattern[v]; ok {
			log.Fatalf("pattern %s is Duplication", pattern)
		}
		r.tpl[v] = make(map[string]*Route)
		r.pattern[v] = listvar
		return v, true
		// 判断是否重复
	} else {
		if _, ok := r.pattern[pattern]; ok {
			log.Fatalf("pattern %s is Duplication", pattern)
		}
		r.route[pattern] = make(map[string]*Route)
		r.pattern[pattern] = make([]string, 0)
		return pattern, false
	}
}

func (r *Router) SetHeader(k, v string) *Router {
	if r.header == nil {
		r.header = map[string]string{}
	}
	r.header[k] = v
	return r
}

func (r *Router) AddModule(handle func(http.ResponseWriter, *http.Request) bool) *Router {
	if r.module == nil {
		r.module = make([]func(http.ResponseWriter, *http.Request) bool, 0)
	}
	r.module = append(r.module, handle)
	return r
}

func (r *Router) EndModule(handle func(http.ResponseWriter, *http.Request)) *Router {
	r.end = handle
	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	connections++
	r.mu.Lock()
	allconn[req] = &Data{}
	r.mu.Unlock()
	defer func() {
		if r.end != nil && req.URL.Path != "/favicon.ico" && !r.IgnoreIco || req.Method != http.MethodOptions {
			r.end(w, req)
		}

		r.mu.Lock()
		delete(allconn, req)
		r.mu.Unlock()
		connections--
	}()
	if r.Slash {
		req.URL.Path = slash(req.URL.Path)
	}

	// /favicon.ico  和 Option 请求， 不支持自定义请求头和模块
	if req.URL.Path == "/favicon.ico" {
		if r.IgnoreIco {
			return
		} else {
			for k, v := range r.header {
				w.Header().Set(k, v)
			}
			r.HanleFavicon.ServeHTTP(w, req)
			return
		}
	}

	// option 请求处理
	if !r.DisableOption && req.Method == http.MethodOptions {
		for k, v := range r.header {
			w.Header().Set(k, v)
		}
		r.HandleOptions.ServeHTTP(w, req)
		return
	}

	// 先进行路由表缓存寻找
	if route, ok := r.routeTable[req.URL.Path+req.Method]; ok {
		// 设置请求头
		r.readFromCache(route, w, req)

	} else {
		// 获取handler
		r.serveHTTP(w, req)
	}
}

func (r *Router) readFromCache(route *rt, w http.ResponseWriter, req *http.Request) {

	if route.dataSource != nil {
		allconn[req].Data = route.dataSource
	}

	for k, v := range route.Header {
		w.Header().Set(k, v)
	}
	// 请求模块
	var ok bool
	for _, v := range route.Module {
		ok = v(w, req)
		if ok {
			return
		}
	}
	route.Handle.ServeHTTP(w, req)
}

// url 是匹配的路径， 可能不是规则的路径
func (r *Router) serveHTTP(w http.ResponseWriter, req *http.Request) {
	var thisRoute *Route

	if _, ok := r.route[req.URL.Path]; ok {
		thisRoute = r.route[req.URL.Path][req.Method]
	} else {
		for reUrl := range r.tpl {
			re := regexp.MustCompile(reUrl)
			if re.MatchString(req.URL.Path) {
				thisRoute = r.tpl[reUrl][req.Method]
				ap := make(map[string]string, 0)
				vl := re.FindStringSubmatch(req.URL.Path)
				for i, v := range r.pattern[reUrl] {
					ap[v] = vl[i+1]
				}
				allparams[req.URL.Path] = ap
				goto endloop

			}

		}
		r.HandleNotFound.ServeHTTP(w, req)
		return
	}
endloop:
	if thisRoute.dataSource != nil {
		allconn[req].Data = thisRoute.dataSource
	}

	// 全局的请求头
	tmpHeader := make(map[string]string)
	for k, v := range r.header {
		tmpHeader[k] = v
		w.Header().Set(k, v)
	}

	// 全局的模块
	tmpModule := make([]func(http.ResponseWriter, *http.Request) bool, 0)

	for _, v := range r.module {
		tmpModule = append(tmpModule, v)
	}

	// 增加单路由的请求头和模块
	for _, v := range thisRoute.module {
		tmpModule = append(tmpModule, v)
	}
	for k, v := range thisRoute.header {
		tmpHeader[k] = v
		w.Header().Set(k, v)
	}
	for _, v := range thisRoute.delheader {
		delete(tmpHeader, v)
		w.Header().Del(v)
	}
	// 删除多余的模块
	for _, v := range thisRoute.delmodule {
		tmp := make([]func(http.ResponseWriter, *http.Request) bool, len(tmpModule)-1)
		for i, tmd := range tmpModule {
			if CompareFunc(v, tmd) {
				copy(tmp[i:], tmpModule[i+1:])
				break
			}
			tmp[i] = tmd
		}
		tmpModule = tmp
	}

	// 缓存handler
	thisRouter := &rt{
		Handle:     thisRoute.handle,
		Header:     tmpHeader,
		Module:     tmpModule,
		dataSource: thisRoute.dataSource,
	}
	r.rtLock.Lock()
	r.routeTable[req.URL.Path+req.Method] = thisRouter
	r.rtLock.Unlock()
	for _, v := range tmpModule {
		ok := v(w, req)
		if ok {
			return
		}
	}
	thisRoute.handle.ServeHTTP(w, req)

}

func (r *Router) Run(opt ...string) error {
	addr := ":8080"
	if len(opt) > 0 {
		addr = opt[0]
	}
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	return srv.ListenAndServe()
}

func NewRouter() *Router {
	return &Router{
		IgnoreIco:      true,
		routeTable:     make(map[string]*rt),
		header:         map[string]string{},
		route:          make(map[string]MethodsRoute),
		tpl:            make(map[string]MethodsRoute),
		pattern:        make(map[string][]string),
		HanleFavicon:   handleFavicon(),
		HandleOptions:  handleOptions(),
		HandleNotFound: handleNotFound(),
		mu:             &sync.RWMutex{},
		rtLock:         &sync.RWMutex{},
	}
}

func handleNotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		return
	})
}

func handleOptions() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})
}

func methodNotAllowed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	})
}

func handleFavicon() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})
}

// 组路由添加到router
func (r *Router) AddGroup(group *GroupRoute) *Router {

	// 将路由的所有变量全部移交到route
	if group.pattern == nil && group.route == nil {
		return nil
	}
	if r.header == nil {
		r.header = make(map[string]string)
	}
	if r.pattern == nil {
		r.pattern = make(map[string][]string)
	}

	if r.route == nil {
		r.route = make(map[string]MethodsRoute)
	}
	if r.tpl == nil {
		r.tpl = make(map[string]MethodsRoute)
	}
	for url, args := range group.pattern {
		r.pattern[url] = args

		if len(args) == 0 {
			for m := range group.route[url] {
				if _, ok := r.route[url][m]; ok {
					log.Fatalf("%s %s is Duplication", url, m)
				}

				merge(group, group.route[url][m])

			}
			if r.route[url] == nil {
				r.route[url] = make(map[string]*Route)
			}

			r.route[url] = group.route[url]

		} else {
			for m := range group.tpl[url] {
				if _, ok := r.tpl[url][m]; ok {
					log.Fatalf("%s %s is Duplication", url, m)
				}

				merge(group, group.tpl[url][m])

			}
			if r.tpl[url] == nil {
				r.tpl[url] = make(map[string]*Route)
			}
			r.tpl[url] = group.tpl[url]

		}

	}
	return r
}

func merge(group *GroupRoute, route *Route) {
	// 合并 delheader
	if group.delheader != nil {
		//
		if route.delheader == nil {
			route.delheader = group.delheader
		} else {
			tmpdelheader := make([]string, 0)
			tmpdelheader = append(tmpdelheader, group.delheader...)
			tmpdelheader = append(tmpdelheader, route.delheader...)
			route.delheader = tmpdelheader
		}

	}

	// 合并 delmodule
	if group.delmodule != nil {
		//
		if route.delmodule == nil {
			route.delmodule = group.delmodule
		} else {
			tmpdelmodule := make([]func(http.ResponseWriter, *http.Request) bool, 0)
			tmpdelmodule = append(tmpdelmodule, group.delmodule...)
			tmpdelmodule = append(tmpdelmodule, route.delmodule...)
			route.module = tmpdelmodule
		}
	}
	// 合并 groupKey
	if route.groupKey == "" {
		route.groupKey = group.groupKey
		route.groupLable = group.groupLable
		route.groupTitle = group.groupTitle
	}
	// 合并 module
	if group.module != nil {
		//
		if route.module == nil {
			route.module = group.module
		} else {
			tmpmodule := make([]func(http.ResponseWriter, *http.Request) bool, 0)
			tmpmodule = append(tmpmodule, group.module...)
			tmpmodule = append(tmpmodule, route.module...)
			route.module = tmpmodule
		}

	}
	// 合并 reqHeader
	if group.reqHeader != nil {
		//
		tmpReqHeader := make(map[string]string)
		for k, v := range group.reqHeader {
			tmpReqHeader[k] = v
		}

		if route.reqHeader != nil {
			for k, v := range route.reqHeader {
				tmpReqHeader[k] = v
			}

		}

		route.reqHeader = tmpReqHeader

	}
	// 合并 header
	if group.header != nil {
		//
		tmpHeader := make(map[string]string)
		for k, v := range group.header {
			tmpHeader[k] = v
		}
		if route.header != nil {
			for k, v := range route.header {
				tmpHeader[k] = v
			}

		}

		route.header = tmpHeader

	}
	// 合并 codeMsg
	if group.codeMsg != nil {
		//
		tmpCodeMsg := make(map[string]string)
		for k, v := range group.codeMsg {
			tmpCodeMsg[k] = v
		}
		if route.codeMsg != nil {
			for k, v := range route.codeMsg {
				if v == "" {
					delete(tmpCodeMsg, k)
				} else {
					tmpCodeMsg[k] = v
				}

			}

		}

		route.codeMsg = tmpCodeMsg
	}
	// 合并 codeField
	if route.codeField == "" {
		route.codeField = group.codeField
	}

}

func (r *Router) DebugRoute() {
	for url, mr := range r.route {
		for k, v := range mr {
			log.Printf("url: %s, method: %s, header: %+v\n", url, k, v.header)
		}

	}
}

func (r *Router) DebugTpl() {
	for url, mr := range r.tpl {
		for k, v := range mr {
			log.Printf("url: %s, method: %s, header: %+v\n", url, k, v.header)
		}

	}
}

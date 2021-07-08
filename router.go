package xmux

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var connections int32

var stop bool

func GetConnents() int32 {
	return connections
}

type rt struct {
	Handle     http.Handler
	Header     map[string]string
	Module     []func(http.ResponseWriter, *http.Request) bool
	dataSource interface{}
	midware    func(handle func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request)
}

type Router struct {
	new            bool // 判断是否是通过newRouter 来初始化的
	ReadTimeout    time.Duration
	IgnoreIco      bool // 是否忽略 /favicon.ico 请求。 默认忽略
	HanleFavicon   http.Handler
	DisableOption  bool         // 禁止全局option
	HandleOptions  http.Handler // 预请求 处理函数， 如果存在， 优先处理, 前后端分离后， 前段可能会先发送一个预请求
	HandleNotFound http.Handler
	Slash          bool                                            // 忽略地址多个斜杠， 默认不忽略
	route          PatternMR                                       // 单实例路由， 组路由最后也会合并过来
	tpl            PatternMR                                       // 正则路由， 组路由最后也会合并过来
	pattern        map[string][]string                             // 记录所有路由， []string 是正则匹配的参数
	header         map[string]string                               // 全局路由头
	module         []func(http.ResponseWriter, *http.Request) bool // 全局模块
	// routeTable     *cacheTable                                     // 路由表
	midware func(handle func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request)
}

func (r *Router) MiddleWare(midware func(handle func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request)) *Router {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	r.midware = midware
	return r
}

func (r *Router) makeRoute(pattern string) (string, bool) {
	// 格式化路径
	// 创建 methodsRoute

	if r.Slash {
		pattern = slash(pattern)
	}

	if v, listvar := match(pattern); len(listvar) > 0 {
		if r.tpl == nil {
			r.tpl = make(map[string]MethodsRoute)
			r.tpl[v] = make(map[string]*Route)
		}
		if _, ok := r.tpl[v]; !ok {
			r.tpl[v] = make(map[string]*Route)
		}
		r.pattern[v] = listvar
		return v, true
		// 判断是否重复
	} else {

		if r.route == nil {
			r.route = make(map[string]MethodsRoute)
			r.route[pattern] = make(map[string]*Route)

		}
		if _, ok := r.route[pattern]; !ok {
			r.route[pattern] = make(map[string]*Route)

		}
		r.pattern[pattern] = make([]string, 0)
		return pattern, false
	}
}

func (r *Router) SetHeader(k, v string) *Router {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	r.header[k] = v
	return r
}

func (r *Router) AddModule(handle func(http.ResponseWriter, *http.Request) bool) *Router {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	if r.module == nil {
		r.module = make([]func(http.ResponseWriter, *http.Request) bool, 0)
	}
	r.module = append(r.module, handle)
	return r
}

func (r *Router) readFromCache(route *rt, w http.ResponseWriter, req *http.Request) {
	if route.dataSource != nil {
		allconn.Get(req)
		fd := allconn.Get(req)
		if fd != nil {
			fd.Data = route.dataSource
		}
	}

	for k, v := range route.Header {
		w.Header().Set(k, v)
	}
	// 请求模块
	for _, v := range route.Module {
		ok := v(w, req)
		if ok {
			return
		}
	}
	if route.midware != nil {
		route.midware(route.Handle.ServeHTTP, w, req)
	} else {
		route.Handle.ServeHTTP(w, req)
	}

}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	if stop {
		return
	}
	atomic.AddInt32(&connections, 1)

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
	fd := &FlowData{
		mu: &sync.RWMutex{},
	}
	allconn.Set(req, fd)
	defer func() {
		allconn.Del(req)
		atomic.AddInt32(&connections, -1)
	}()
	// 先进行路由表缓存寻找
	// route, ok := r.routeTable.Get(req.URL.Path + req.Method)
	// if ok {
	// 	r.readFromCache(route, w, req)
	// } else {
	// 获取handler
	r.serveHTTP(w, req)
	// }
}

// url 是匹配的路径， 可能不是规则的路径
func (r *Router) serveHTTP(w http.ResponseWriter, req *http.Request) {
	var thisRoute *Route
	if _, ok := r.route[req.URL.Path]; ok {
		thisRoute, ok = r.route[req.URL.Path][req.Method]
		if !ok {
			thisRoute, ok = r.route[req.URL.Path]["*"]
			if !ok {
				r.HandleNotFound.ServeHTTP(w, req)
				return
			}
		}
	} else {
		for reUrl := range r.tpl {

			re := regexp.MustCompile(reUrl)
			req.URL.Path = strings.Trim(req.URL.Path, " ")
			if re.MatchString(req.URL.Path) {
				thisRoute, ok = r.tpl[reUrl][req.Method]
				if !ok {
					thisRoute, ok = r.tpl[reUrl]["*"]
					if !ok {
						r.HandleNotFound.ServeHTTP(w, req)
						return
					}

				}
				ap := make(map[string]string)
				vl := re.FindStringSubmatch(req.URL.Path)
				for i, v := range r.pattern[reUrl] {
					ap[v] = vl[i+1]
				}
				SetParams(req.URL.Path, ap)
				goto endloop
			}
		}
		r.HandleNotFound.ServeHTTP(w, req)
		return
	}
endloop:
	if thisRoute.dataSource != nil {
		allconn.Get(req)
		fd := allconn.Get(req)
		if fd != nil {
			fd.Data = thisRoute.dataSource
		}

	}

	// 全局的请求头
	tmpHeader := make(map[string]string)
	for k, v := range r.header {
		tmpHeader[k] = v
		w.Header().Set(k, v)
	}

	// 全局的模块
	tmpModule := make([]func(http.ResponseWriter, *http.Request) bool, 0)

	tmpModule = append(tmpModule, r.module...)

	// 增加单路由的请求头和模块
	tmpModule = append(tmpModule, thisRoute.module...)
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
			if compareFunc(v, tmd) {
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
		midware:    thisRoute.midware,
	}
	// r.routeTable.Set(req.URL.Path+req.Method, thisRouter)
	r.readFromCache(thisRouter, w, req)
}

func (r *Router) Run(opt ...string) error {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	addr := ":8080"
	if len(opt) > 0 {
		addr = opt[0]
	}
	srv := &http.Server{
		Addr:        addr,
		ReadTimeout: r.ReadTimeout,
		Handler:     r,
	}
	fmt.Printf("listen on %s\n", addr)
	return srv.ListenAndServe()
}

func (r *Router) RunTLS(crt, key string, opt ...string) error {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	addr := ":443"
	if len(opt) > 0 {
		addr = opt[0]
	}
	srv := &http.Server{
		Addr:        addr,
		ReadTimeout: r.ReadTimeout,
		Handler:     r,
	}
	return srv.ListenAndServeTLS(crt, key)
}

func NewRouter() *Router {
	return &Router{
		new:       true,
		IgnoreIco: true,
		// routeTable: &cacheTable{
		// 	cache: make(map[string]*rt),
		// 	mu:    &sync.RWMutex{},
		// },
		header:         map[string]string{},
		pattern:        make(map[string][]string),
		HanleFavicon:   handleFavicon(),
		HandleOptions:  handleOptions(),
		HandleNotFound: handleNotFound(),
	}
}

func handleNotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusNotFound)
	})
}

func handleOptions() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func handleFavicon() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
	})
}

// 组路由添加到router
func (r *Router) AddGroup(group *GroupRoute) *Router {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	// 将路由的所有变量全部移交到route
	if group.pattern == nil && group.route == nil {
		return nil
	}

	if r.route == nil {
		r.route = make(map[string]MethodsRoute)
	}
	if r.tpl == nil {
		r.tpl = make(map[string]MethodsRoute)
	}
	if r.midware != nil && !reflect.DeepEqual(group.delmidware, r.midware) && group.midware == nil {
		group.midware = r.midware
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
			r.route[url] = group.route[url]

		} else {
			for m := range group.tpl[url] {
				if _, ok := r.tpl[url][m]; ok {
					log.Fatalf("%s %s is Duplication", url, m)
				}

				merge(group, group.tpl[url][m])

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
	if group.midware != nil && route.midware == nil {
		route.midware = group.midware
	}
	if group.midware != nil && !reflect.DeepEqual(route.delmidware, group.midware) && route.midware == nil {
		route.midware = group.midware
	}
}

func (r *Router) DebugRoute() {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	for url, mr := range r.route {
		for k, v := range mr {
			log.Printf("url: %s, method: %s, header: %+v\n", url, k, v.header)
		}

	}
}

func (r *Router) DebugTpl() {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	for url, mr := range r.tpl {
		for k, v := range mr {
			log.Printf("url: %s, method: %s, header: %+v\n", url, k, v.header)
		}

	}
}

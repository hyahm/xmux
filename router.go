package xmux

import (
	"log"
	"net/http"
	"regexp"
	"sync"

	"github.com/hyahm/golog"
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
	Midware    []func(http.ResponseWriter, *http.Request) bool
	dataSource interface{}
}

type Router struct {
	IgnoreIco      bool // 是否忽略 /favicon.ico 请求。 默认忽略
	HanleFavicon   http.Handler
	DisableOption  bool         // 禁止全局option
	Options        http.Handler // 预请求 处理函数， 如果存在， 优先处理, 前后端分离后， 前段可能会先发送一个预请求
	HandleNotFound http.Handler
	Slash          bool
	route          PatternMR                                       // 单实例路由， 组路由最后也会合并过来
	tpl            PatternMR                                       // 正则路由， 组路由最后也会合并过来
	pattern        map[string][]string                             // 记录所有路由， value 是正则匹配的参数
	header         map[string]string                               // 全局路由头
	midware        []func(http.ResponseWriter, *http.Request) bool // 全局中间件
	routeTable     map[string]*rt                                  // 路由表
	mu             *sync.RWMutex
}

func (r *Router) makeRoute(pattern string) (string, bool) {
	// 格式化路径
	// 创建 methodsRoute
	// 格式路径
	if r.Slash {
		pattern = slash(pattern)
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

func (r *Router) AddMidware(handle func(http.ResponseWriter, *http.Request) bool) *Router {
	if r.midware == nil {
		r.midware = make([]func(http.ResponseWriter, *http.Request) bool, 0)
	}
	r.midware = append(r.midware, handle)
	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	connections++
	if r.mu == nil {
		r.mu = &sync.RWMutex{}
	}
	if r.routeTable == nil {
		r.routeTable = make(map[string]*rt)
	}
	defer func() {
		connections--
		delete(allconn, req)
	}()
	url := req.URL.Path
	if r.Slash {
		url = slash(req.URL.Path)
	}

	// /favicon.ico  和 Option 请求， 不支持自定义请求头和中间件
	if r.IgnoreIco && url == "/favicon.ico" {
		for k, v := range r.header {
			w.Header().Set(k, v)
		}
		r.HanleFavicon.ServeHTTP(w, req)
		return
	}
	// option 请求处理
	if !r.DisableOption && req.Method == http.MethodOptions {
		for k, v := range r.header {
			w.Header().Set(k, v)
		}
		r.Options.ServeHTTP(w, req)
		return
	}

	// 先进行路由表缓存寻找
	if route, ok := r.routeTable[url+req.Method]; ok {
		// 设置请求头
		if route.dataSource != nil {
			allconn[req] = &Data{
				ctx: make(map[string]interface{}),
				mu:  &sync.RWMutex{},
			}
			allconn[req].Data = route.dataSource
		}

		for k, v := range route.Header {
			w.Header().Set(k, v)
		}
		// 请求中间件
		var ok bool
		for _, v := range route.Midware {
			ok = v(w, req)
			if ok {
				return
			}
		}
		route.Handle.ServeHTTP(w, req)

	} else {
		// 获取handler
		r.serveHTTP(url, w, req)
	}
}

// url 是匹配的路径， 可能不是规则的路径
func (r *Router) serveHTTP(url string, w http.ResponseWriter, req *http.Request) {

	var this_route *Route
	if _, ok := r.route[url]; ok {
		this_route = r.route[url][req.Method]

	} else {
		for reUrl, mr := range r.tpl {
			re := regexp.MustCompile(reUrl)
			if re.MatchString(url) {
				this_route = mr[req.Method]
				ap := make(map[string]string, 0)
				vl := re.FindStringSubmatch(url)
				for i, v := range r.pattern[reUrl] {
					ap[v] = vl[i+1]
				}
				allparams[url] = ap
				goto endloop

			}

		}
		r.HandleNotFound.ServeHTTP(w, req)
		return
	}
endloop:
	if this_route.dataSource != nil {
		allconn[req] = &Data{
			ctx: make(map[string]interface{}),
			mu:  &sync.RWMutex{},
		}
		allconn[req].Data = this_route.dataSource
	}

	// 全局的请求头
	tmpHeader := make(map[string]string)
	for k, v := range r.header {
		tmpHeader[k] = v
		w.Header().Set(k, v)
	}

	// 全局的中间件
	tmpMidware := make([]func(http.ResponseWriter, *http.Request) bool, 0)

	for _, v := range r.midware {
		tmpMidware = append(tmpMidware, v)
	}

	// 增加单路由的请求头和中间件
	for _, v := range this_route.midware {
		tmpMidware = append(tmpMidware, v)
	}
	for k, v := range this_route.header {
		tmpHeader[k] = v
		w.Header().Set(k, v)
	}
	for _, v := range this_route.delheader {
		delete(tmpHeader, v)
		w.Header().Del(v)
	}
	// 删除多余的中间件
	for _, v := range this_route.delmidware {
		for i, tmd := range tmpMidware {
			if CompareFunc(v, tmd) {
				tmp := make([]func(http.ResponseWriter, *http.Request) bool, 0)
				tmp = append(tmp, tmpMidware[0:i]...)
				tmp = append(tmp, tmpMidware[i+1:]...)
				tmpMidware = tmp
				break
			}
		}

	}

	// 缓存handler
	thisRouter := &rt{
		Handle:     this_route.handle,
		Header:     tmpHeader,
		Midware:    tmpMidware,
		dataSource: this_route.dataSource,
	}
	r.mu.Lock()
	r.routeTable[url+req.Method] = thisRouter
	r.mu.Unlock()
	for _, v := range tmpMidware {
		ok := v(w, req)
		if ok {
			return
		}
	}
	this_route.handle.ServeHTTP(w, req)

}

func (r *Router) Debug() {
	golog.Infof("%+v", r)
	golog.Infof("%+v", r.pattern)
	golog.Infof("%+v", r.route)
	golog.Infof("%+v", r.tpl)
}

func NewRouter() *Router {
	return &Router{
		IgnoreIco:      true,
		Slash:          true,
		routeTable:     make(map[string]*rt),
		header:         map[string]string{},
		route:          make(map[string]MethodsRoute),
		tpl:            make(map[string]MethodsRoute),
		pattern:        make(map[string][]string),
		HanleFavicon:   favicon(),
		Options:        options(),
		HandleNotFound: handleNotFound(),
	}
}

func handleNotFound() http.Handler {
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
	golog.Infof("%+v", *r)
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

	// 合并 delmidware
	if group.delmidware != nil {
		//
		if route.delmidware == nil {
			route.delmidware = group.delmidware
		} else {
			tmpdelmidware := make([]func(http.ResponseWriter, *http.Request) bool, 0)
			tmpdelmidware = append(tmpdelmidware, group.delmidware...)
			tmpdelmidware = append(tmpdelmidware, route.delmidware...)
			route.midware = tmpdelmidware
		}
	}
	// 合并 groupKey
	if route.groupKey == "" {
		route.groupKey = group.groupKey
		route.groupLable = group.groupLable
		route.groupTitle = group.groupTitle
	}
	// 合并 midware
	if group.midware != nil {
		//
		if route.midware == nil {
			route.midware = group.midware
		} else {
			tmpmidware := make([]func(http.ResponseWriter, *http.Request) bool, 0)
			tmpmidware = append(tmpmidware, group.midware...)
			tmpmidware = append(tmpmidware, route.midware...)
			route.midware = tmpmidware
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

package xmux

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/yaml.v3"
)

var connections int32

const CURRFUNCNAME = "CURRFUNCNAME"
const PAGES = "PAGES"

var stop bool

func GetConnents() int32 {
	return connections
}

func StopService() {
	stop = true
}

func StartService() {
	stop = false
}

type rt struct {
	Handle     http.Handler
	Header     map[string]string
	pagekeys   map[string]struct{}
	module     []func(http.ResponseWriter, *http.Request) bool
	dataSource interface{} // 绑定数据结构，
	bindType   bindType
	midware    func(handle func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request)
	// instance   map[*http.Request]interface{} // 解析到这里
}

func DefaultExitTemplate(now time.Time, w http.ResponseWriter, r *http.Request) {
}

type Router struct {
	Exit           func(time.Time, http.ResponseWriter, *http.Request)
	new            bool                                     // 判断是否是通过newRouter 来初始化的
	Enter          func(http.ResponseWriter, *http.Request) // 当有请求进入时候的执行
	ReadTimeout    time.Duration
	IgnoreIco      bool // 是否忽略 /favicon.ico 请求。 默认忽略
	HanleFavicon   http.Handler
	DisableOption  bool         // 禁止全局option
	HandleOptions  http.Handler // 预请求 处理函数， 如果存在， 优先处理, 前后端分离后， 前段可能会先发送一个预请求
	HandleNotFound http.Handler
	Slash          bool                // 忽略地址多个斜杠， 默认不忽略
	route          PatternMR           // 单实例路由， 组路由最后也会合并过来
	tpl            PatternMR           // 正则路由， 组路由最后也会合并过来
	params         map[string][]string // 记录所有路由， []string 是正则匹配的参数
	header         map[string]string   // 全局路由头
	module         module              // 全局模块
	// routeTable     *rt                                             // 路由表
	pagekeys map[string]struct{}
	midware  func(handle func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request)
}

func (r *Router) MiddleWare(midware func(handle func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request)) *Router {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	r.midware = midware
	return r
}

// 判断是否是正则路径， 返回一个路径 string 和是否是正则的 bool
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
		r.params[v] = listvar
		return v, true
	} else {
		if r.route == nil {
			r.route = make(map[string]MethodsRoute)
			r.route[pattern] = make(map[string]*Route)

		}
		if _, ok := r.route[pattern]; !ok {
			r.route[pattern] = make(map[string]*Route)

		}
		r.params[pattern] = make([]string, 0)
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

func (r *Router) AddPageKeys(pagekeys ...string) *Router {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	if r.pagekeys == nil {
		r.pagekeys = make(map[string]struct{})
	}
	for _, v := range pagekeys {
		r.pagekeys[v] = struct{}{}
	}
	return r
}

func (r *Router) AddModule(handles ...func(http.ResponseWriter, *http.Request) bool) *Router {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}

	r.module = r.module.add(handles...)
	return r
}

func (r *Router) readFromCache(start time.Time, route *rt, w http.ResponseWriter, req *http.Request) {
	fd := &FlowData{
		ctx: make(map[string]interface{}),
		mu:  &sync.RWMutex{},
	}
	atomic.AddInt32(&connections, 1)
	allconn.Set(req, fd)
	defer func() {
		atomic.AddInt32(&connections, -1)
		allconn.Del(req)
		if r.Exit != nil {
			r.Exit(start, w, req)
		}
		if err := recover(); err != nil {
			log.Println(req.URL.Path, "---------", err)
		}
	}()
	if route.dataSource != nil {
		base := reflect.TypeOf(route.dataSource)
		// 支持bind 指针和结构体
		if base.Kind() == reflect.Ptr {
			fd.Data = reflect.New(reflect.TypeOf(route.dataSource).Elem()).Interface()
		} else {
			fd.Data = reflect.New(reflect.TypeOf(route.dataSource)).Interface()
		}

		// 数据绑定
		switch route.bindType {
		case jsonT:
			err := json.NewDecoder(req.Body).Decode(&fd.Data)
			if err != nil {
				if err != io.EOF {
					log.Println(err)
				}
				break
			}
			defer req.Body.Close()
		case yamlT:
			err := yaml.NewDecoder(req.Body).Decode(&fd.Data)
			if err != nil {
				if err != io.EOF {
					log.Println(err)
				}
				break
			}
			defer req.Body.Close()
		case xmlT:
			err := xml.NewDecoder(req.Body).Decode(&fd.Data)
			if err != nil {
				if err != io.EOF {
					log.Println(err)
				}
				break
			}
			defer req.Body.Close()
		}
	}
	for k, v := range route.Header {
		w.Header().Set(k, v)
	}

	// 权限导入
	// pages
	fd.Set(PAGES, route.pagekeys)
	// 当前函数名去掉目录层级后的
	name := runtime.FuncForPC(reflect.ValueOf(route.Handle).Pointer()).Name()
	n := strings.LastIndex(name, ".")

	fd.Set(CURRFUNCNAME, name[n+1:])
	// 请求模块
	for _, module := range route.module {
		ok := module(w, req)
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
	start := time.Now()
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	if r.Enter != nil {
		r.Enter(w, req)
	}

	if stop {
		w.WriteHeader(http.StatusLocked)
		return
	}
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

	// 正在寻址中的话， 就等待寻址完成
	// 先进行路由表缓存寻找
	route, ok := Get(req.URL.Path + req.Method)
	if ok {
		r.readFromCache(start, route, w, req)
	} else {
		// 寻址
		r.serveHTTP(start, w, req)
	}
}

// url 是匹配的路径， 可能不是规则的路径, 寻址的时候还是要加锁
func (r *Router) serveHTTP(start time.Time, w http.ResponseWriter, req *http.Request) {
	var thisRoute *Route
	if _, ok := r.route[req.URL.Path]; ok {
		thisRoute, ok = r.route[req.URL.Path][req.Method]
		if !ok {
			r.HandleNotFound.ServeHTTP(w, req)
			atomic.AddInt32(&connections, -1)
			return
		}
	} else {
		for reUrl := range r.tpl {
			re := regexp.MustCompile(reUrl)
			req.URL.Path = strings.Trim(req.URL.Path, " ")
			if re.MatchString(req.URL.Path) {
				thisRoute, ok = r.tpl[reUrl][req.Method]
				if !ok {
					r.HandleNotFound.ServeHTTP(w, req)
					atomic.AddInt32(&connections, -1)
					return
				}
				ap := make(map[string]string)
				vl := re.FindStringSubmatch(req.URL.Path)
				for i, v := range r.params[reUrl] {
					ap[v] = vl[i+1]
				}
				setParams(req.URL.Path, ap)
				goto endloop
			}
		}
		r.HandleNotFound.ServeHTTP(w, req)
		atomic.AddInt32(&connections, -1)
		return
	}
endloop:

	// 缓存handler
	thisRouter := &rt{
		Handle:     thisRoute.handle,
		Header:     thisRoute.header,
		module:     thisRoute.module.getMuduleList(),
		dataSource: thisRoute.dataSource,
		midware:    thisRoute.midware,
		pagekeys:   thisRoute.pagekeys,
		bindType:   thisRoute.bindType,
	}
	// 设置缓存
	Set(req.URL.Path+req.Method, thisRouter)
	r.readFromCache(start, thisRouter, w, req)
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
	fmt.Printf("listen on %s\n", addr)
	return srv.ListenAndServeTLS(crt, key)
}

func NewRouter(cache ...uint64) *Router {
	var c uint64
	if len(cache) > 0 {
		c = cache[0]
	}
	InitCache(c)
	return &Router{
		new:            true,
		IgnoreIco:      true,
		route:          make(map[string]MethodsRoute),
		tpl:            make(map[string]MethodsRoute),
		header:         map[string]string{},
		params:         make(map[string][]string),
		HanleFavicon:   handleFavicon(),
		HandleOptions:  handleOptions(),
		HandleNotFound: handleNotFound(),
		module:         module{},
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

// 组路由添加到router里面,
// 挂载到group之前， 全局的变量已经挂载到route 里面了， 所以不用再管组变量了
func (r *Router) AddGroup(group *GroupRoute) *Router {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	// 将路由的所有变量全部移交到route
	if group.params == nil && group.route == nil {
		return nil
	}

	for url, args := range group.params {
		r.params[url] = args

		if len(args) == 0 {
			for method := range group.route[url] {
				if _, ok := r.route[url][method]; ok {
					log.Fatalf("%s %s is Duplication", url, method)
				}
				mergeDoc(group, group.route[url][method])
				if group.route[url][method].midware == nil {
					log.Println("2222222222222222222222")
					group.route[url][method].midware = r.midware
				}
				group.route[url][method].module = r.module.addModule(group.route[url][method].module)
				for key := range group.route[url][method].delmodule.modules {
					group.route[url][method].module = group.route[url][method].module.deleteKey(key)
				}
				// 合并head
				tempHeader := make(map[string]string)
				for k, v := range r.header {
					tempHeader[k] = v
				}
				for k, v := range group.route[url][method].header {
					tempHeader[k] = v
				}
				group.route[url][method].header = tempHeader
				// 合并 delheader
				for _, k := range group.route[url][method].delheader {
					delete(group.route[url][method].header, k)
				}

				// 合并 pagekeys
				tempPages := make(map[string]struct{})
				for k := range r.pagekeys {
					tempPages[k] = struct{}{}
				}

				for k := range group.route[url][method].pagekeys {
					tempPages[k] = struct{}{}
				}
				group.route[url][method].pagekeys = tempPages
				// 删除 pagekeys

				for _, k := range group.route[url][method].delPageKeys {
					delete(group.route[url][method].pagekeys, k)
				}
				// delete midware
				if group.route[url][method].delmidware != nil && GetFuncName(group.route[url][method].delmidware) == GetFuncName(r.midware) {
					group.route[url][method].midware = nil
				}
			}
			r.route[url] = group.route[url]

		} else {
			for method := range group.tpl[url] {
				if _, ok := r.tpl[url][method]; ok {
					log.Fatalf("%s %s is Duplication", url, method)
				}

				mergeDoc(group, group.tpl[url][method])
				if group.tpl[url][method].midware == nil {
					group.tpl[url][method].midware = r.midware
				}

				group.tpl[url][method].module = r.module.addModule(group.tpl[url][method].module)
				for key := range group.tpl[url][method].delmodule.modules {
					group.tpl[url][method].module = group.tpl[url][method].module.deleteKey(key)
				}
				// 合并head
				tempHeader := make(map[string]string)
				for k, v := range r.header {
					tempHeader[k] = v
				}
				for k, v := range group.tpl[url][method].header {
					tempHeader[k] = v
				}
				group.tpl[url][method].header = tempHeader

				// 合并 delheader
				for _, k := range group.tpl[url][method].delheader {
					delete(group.tpl[url][method].header, k)
				}

				// 合并 pagekeys
				tempPages := make(map[string]struct{})
				for k := range r.pagekeys {
					tempPages[k] = struct{}{}
				}

				for k := range group.tpl[url][method].pagekeys {
					tempPages[k] = struct{}{}
				}
				group.tpl[url][method].pagekeys = tempPages

				// 删除pagekey
				for _, k := range group.tpl[url][method].delPageKeys {
					delete(group.tpl[url][method].pagekeys, k)
				}
				// delete midware
				if group.tpl[url][method].delmidware != nil && GetFuncName(group.tpl[url][method].delmidware) == GetFuncName(r.midware) {
					group.tpl[url][method].midware = nil
				}
			}

			r.tpl[url] = group.tpl[url]

		}

	}

	return r
}

// 将路由组的信息合并到路由
func mergeDoc(group *GroupRoute, route *Route) {

	// 合并 groupKey
	if route.groupKey == "" {
		route.groupKey = group.groupKey
		route.groupLabel = group.groupLabel
		route.groupTitle = group.groupTitle
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

	// 合并请求头

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

func debugPrint(url string, mr MethodsRoute) {
	for k, v := range mr {
		log.Printf("url: %s, method: %s, header: %+v, module: %#v, midware: %#v \n",
			url, k, v.header, v.module.funcOrder, GetFuncName(v.midware))
	}
}

func (r *Router) DebugRoute() {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	for url, mr := range r.route {
		debugPrint(url, mr)
	}
}

func (r *Router) DebugAssignRoute(thisurl string) {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	for url, mr := range r.route {
		if thisurl == url {
			debugPrint(url, mr)
			return
		}
	}
}

func (r *Router) DebugTpl(pattern string) {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	for url, mr := range r.tpl {
		debugPrint(url, mr)
	}
}

func (r *Router) DebugIncludeTpl(pattern string) {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	for url, mr := range r.tpl {
		if strings.Contains(url, pattern) {
			debugPrint(url, mr)
		}
	}
}

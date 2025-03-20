package xmux

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hyahm/xmux/helper"
)

var connections int32

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
	Handle       http.Handler
	Header       map[string]string
	pagekeys     map[string]struct{}
	module       []func(http.ResponseWriter, *http.Request) bool
	dataSource   interface{} // 绑定数据结构，
	bindType     bindType
	responseData interface{}
	// instance   map[*http.Request]interface{} // 解析到这里
}

func requestBytes(reqbody []byte, r *http.Request) {
	log.Printf("connect_id: %d\tdata: %s", GetInstance(r).GetConnectId(), string(reqbody))
}

type Router struct {
	addr                 string
	prefix               []string
	MaxPrintLength       int
	Exit                 func(time.Time, http.ResponseWriter, *http.Request)
	new                  bool                                          // 判断是否是通过newRouter 来初始化的
	EnableConnect        bool                                          // 判断是否是通过newRouter 来初始化的
	Enter                func(http.ResponseWriter, *http.Request) bool // 当有请求进入时候的执行
	ReadTimeout          time.Duration
	HanleFavicon         func(http.ResponseWriter, *http.Request)
	DisableOption        bool                                     // 禁止全局option
	HandleOptions        func(http.ResponseWriter, *http.Request) // 预请求 处理函数， 如果存在， 优先处理, 前后端分离后， 前段可能会先发送一个预请求
	HandleNotFound       func(http.ResponseWriter, *http.Request)
	HandleConnect        func(http.ResponseWriter, *http.Request)
	NotFoundRequireField func(string, http.ResponseWriter, *http.Request) bool
	UnmarshalError       func(error, http.ResponseWriter, *http.Request) bool
	IgnoreSlash          bool                // 忽略地址多个斜杠， 默认不忽略
	route                UMR                 // 单实例路由， 组路由最后也会合并过来
	tpl                  UMR                 // 正则路由， 组路由最后也会合并过来
	params               map[string][]string // 记录所有路由， []string 是正则匹配的参数
	header               map[string]string   // 全局路由头
	module               *module             // 全局模块
	responseData         interface{}
	pagekeys             mstringstruct

	SwaggerTitle       string
	SwaggerDescription string
	SwaggerVersion     string
}

func (r *Router) BindResponse(response interface{}) *Router {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	r.responseData = response
	return r
}

// 判断是否是正则路径， 返回一个路径 string 和是否是正则的 bool
func (r *Router) makeRoute(pattern string) (string, []string, bool) {
	// 格式化路径
	// 创建 methodsRoute

	if r.IgnoreSlash {
		pattern = PrettySlash(pattern)
	}

	if v, listvar := match(pattern); len(listvar) > 0 {
		r.params[v] = listvar
		return v, listvar, true
	} else {

		r.params[pattern] = make([]string, 0)
		return pattern, nil, false
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
	r.module.add(handles...)
	return r
}

func (r *Router) readFromCache(start time.Time, route *rt, w http.ResponseWriter, req *http.Request, fd *FlowData) {

	if route.responseData != nil {
		fd.Response = Clone(route.responseData)
	}
	if route.dataSource != nil {
		base := reflect.TypeOf(route.dataSource)
		// 支持bind 指针和结构体
		if base.Kind() == reflect.Ptr {
			fd.Data = reflect.New(reflect.TypeOf(route.dataSource).Elem()).Interface()
		} else {
			fd.Data = reflect.New(reflect.TypeOf(route.dataSource)).Interface()
		}

		if route.bindType != 0 {
			if r.bind(route, w, req, fd) {
				return
			}
		} else {
			GetInstance(req).Body = []byte("")
		}
	}

	for k, v := range route.Header {
		w.Header().Set(k, v)
	}

	// 权限导入
	// pages
	fd.pages = route.pagekeys
	// 当前函数名去掉目录层级后的
	name := runtime.FuncForPC(reflect.ValueOf(route.Handle).Pointer()).Name()
	n := strings.LastIndex(name, ".")
	fd.funcName = name[n+1:]
	// 请求模块
	for _, module := range route.module {
		ok := module(w, req)
		if ok {
			return
		}
	}

	if route.Handle.(http.HandlerFunc) != nil {
		route.Handle.ServeHTTP(w, req)
	}

}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}

	atomic.AddInt32(&connections, 1)
	defer atomic.AddInt32(&connections, -1)

	if req.Method == http.MethodConnect {
		if r.EnableConnect {
			if r.HandleConnect == nil {
				r.HandleConnect = handleConnect
			}
		} else {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		r.HandleConnect(w, req)
		return
	}

	ci := time.Now().UnixNano()
	fd := &FlowData{
		ctx:        make(map[string]interface{}),
		mu:         &sync.RWMutex{},
		connectId:  ci,
		StatusCode: 200,
	}

	allconn.Set(req, fd)
	defer allconn.Del(req)
	start := time.Now()
	if r.Enter != nil {
		if r.Enter(w, req) {
			return
		}
	}
	if r.Exit != nil {
		defer r.Exit(start, w, req)
	}
	if stop {
		fd.StatusCode = http.StatusLocked
		w.WriteHeader(http.StatusLocked)
		return
	}

	if r.IgnoreSlash {
		req.URL.Path = PrettySlash(req.URL.Path)
	}
	// /favicon.ico 请求
	if req.URL.Path == "/favicon.ico" {
		for k, v := range r.header {
			w.Header().Set(k, v)
		}
		r.HanleFavicon(w, req)
		return
	}
	// option 请求处理
	if !r.DisableOption && req.Method == http.MethodOptions {
		for k, v := range r.header {
			w.Header().Set(k, v)
		}
		r.HandleOptions(w, req)
		return
	}

	// 正在寻址中的话， 就等待寻址完成
	// 先进行路由表缓存寻找
	route, ok := getUrlCache(req.URL.Path + req.Method)
	if ok {
		r.readFromCache(start, route, w, req, fd)
	} else {
		// 寻址
		r.serveHTTP(start, w, req, fd)
	}
}

// url 是匹配的路径， 可能不是规则的路径, 寻址的时候还是要加锁
func (r *Router) serveHTTP(start time.Time, w http.ResponseWriter, req *http.Request, fd *FlowData) {
	var thisRoute *Route
	if _, ok := r.route[req.URL.Path]; ok {
		route, ok := r.route[req.URL.Path][req.Method]
		if !ok {
			r.HandleNotFound(w, req)
			atomic.AddInt32(&connections, -1)
			return
		}
		thisRoute = route
	} else {
		for reUrl := range r.tpl {
			re := regexp.MustCompile(reUrl)
			// req.URL.Path = strings.Trim(req.URL.Path, " ")
			if re.MatchString(req.URL.Path) {
				route, ok := r.tpl[reUrl][req.Method]
				if ok {
					ap := make(map[string]string)
					vl := re.FindStringSubmatch(req.URL.Path)
					for i, v := range r.params[reUrl] {
						ap[v] = vl[i+1]
					}
					thisRoute = route
					setParams(req.URL.Path, ap)
					goto endloop
				}

			}
		}
		r.HandleNotFound(w, req)
		atomic.AddInt32(&connections, -1)
		return
	}
endloop:
	// 缓存handler
	thisRouter := &rt{
		Handle:       thisRoute.handle,
		Header:       thisRoute.header,
		module:       thisRoute.module.GetModules(),
		dataSource:   thisRoute.dataSource,
		pagekeys:     thisRoute.pagekeys,
		bindType:     thisRoute.bindType,
		responseData: thisRoute.responseData,
	}
	// 设置缓存
	setUrlCache(req.URL.Path+req.Method, thisRouter)
	r.readFromCache(start, thisRouter, w, req, fd)
}

func (r *Router) SetAddr(addr string) {
	r.addr = addr
}

func (r *Router) Run() error {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	svc := &http.Server{
		Addr:        r.addr,
		ReadTimeout: r.ReadTimeout,
		Handler:     r,
	}
	fmt.Printf("listen on %s\n", r.addr)
	return svc.ListenAndServe()
}

func (r *Router) Debug(ctx context.Context) {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	svc := &http.Server{
		Addr:        r.addr,
		ReadTimeout: r.ReadTimeout,
		Handler:     r,
	}
	fmt.Printf("listen on %s\n", r.addr)
	go svc.ListenAndServe()
	select {
	case <-ctx.Done():
		svc.Close()
		return
	}

}

func SetPem(name string) string {
	return name
}

func SetKey(name string) string {
	return name
}

type Opt interface {
	SetKey() string
	SetPem() string
	SetAddr() string
}

func (r *Router) RunUnsafeTLS(opt ...string) error {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	addr := ":443"
	if len(opt) > 0 {
		addr = opt[0]
	}

	svc := &http.Server{
		Addr:        addr,
		ReadTimeout: r.ReadTimeout,
		Handler:     r,
	}
	keyfile := "keys/server.key"
	pemfile := "keys/server.pem"
	// 如果key文件不存在那么就自动生成
	_, err1 := os.Stat(keyfile)
	_, err2 := os.Stat(pemfile)
	if os.IsNotExist(err1) || os.IsNotExist(err2) {
		createTLS()
		if err := svc.ListenAndServeTLS(pemfile, keyfile); err != nil {
			log.Fatal(err)
		}
	}
	if err := svc.ListenAndServeTLS(pemfile, keyfile); err != nil {
		log.Fatal(err)
	}
	fmt.Println("listen on " + addr + " over https")
	return nil
}

func (r *Router) RunTLS(certFile, keyFile string) error {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}

	if strings.Trim(keyFile, "") == "" {
		keyFile = "server.key"
	}
	if strings.Trim(certFile, "") == "" {
		certFile = "server.crt" // 证书文件
	}

	svc := &http.Server{
		Addr:        r.addr,
		ReadTimeout: r.ReadTimeout,
		Handler:     r,
	}

	// 如果key文件不存在那么就自动生成
	_, err1 := os.Stat(keyFile)
	_, err2 := os.Stat(certFile)
	if os.IsNotExist(err1) || os.IsNotExist(err2) {
		GenRsa(keyFile, certFile)
	}
	fmt.Println("listen on ", r.addr, "over https")
	return svc.ListenAndServeTLS(certFile, keyFile)
}

func NewRouter(cacheSize ...int) *Router {
	var c int
	if len(cacheSize) > 0 {
		c = cacheSize[0]
	}
	initUrlCache(c)
	return &Router{
		addr:           ":8080",
		MaxPrintLength: 2 << 10, // 默认的form最大2k
		new:            true,
		prefix:         []string{"/"},
		route:          make(UMR),
		tpl:            make(UMR),
		header:         map[string]string{},
		params:         make(map[string][]string),
		Exit:           exit,
		module: &module{
			filter:    make(map[string]struct{}),
			funcOrder: make([]func(w http.ResponseWriter, r *http.Request) bool, 0),
			mu:        sync.RWMutex{},
		},
		HanleFavicon:         handleFavicon,
		HandleOptions:        handleOptions,
		HandleNotFound:       handleNotFound,
		NotFoundRequireField: notFoundRequireField,

		UnmarshalError: unmarshalError,
	}
}

func unmarshalError(err error, w http.ResponseWriter, r *http.Request) bool {
	fmt.Println(err)
	return false
}

func notFoundRequireField(key string, w http.ResponseWriter, r *http.Request) bool {
	fmt.Println("required field not found", key)
	return false
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	GetInstance(r).StatusCode = http.StatusNotFound
	w.WriteHeader(http.StatusNotFound)
}

func handleConnect(w http.ResponseWriter, r *http.Request) {
	destConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer destConn.Close()

	// 向客户端返回成功响应
	w.WriteHeader(http.StatusOK)

	// 使用 Hijacker 获取客户端的 TCP 连接
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer clientConn.Close()

	// 在客户端和目标服务器之间建立双向隧道
	go transfer(destConn, clientConn)
	transfer(clientConn, destConn)
}

// 数据传输函数
func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func handleOptions(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func handleFavicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
}

func (r *Router) cloneHeader() mstringstring {
	tempHeader := make(map[string]string)
	for k, v := range r.header {
		tempHeader[k] = v
	}
	return tempHeader
}

func (r *Router) Prefix(prefixs ...string) *Router {
	if len(prefixs) == 0 {
		return r
	}
	r.prefix = append(r.prefix, prefixs...)

	return r
}

// 组路由合并到每个单路由里面
func (r *Router) merge(group *RouteGroup, route *Route) {
	// 合并head
	tempHeader := r.cloneHeader()
	// 组的删除是为了删全局
	tempHeader.delete(group.delheader)
	// 添加组路由的
	tempHeader.add(group.header)

	// 私有路由删除组和全局的
	tempHeader.delete(route.delheader)
	// 添加个人路由
	tempHeader.add(route.header)
	// 最终请求头
	route.header = tempHeader

	// 合并返回
	// 本身要是绑定了数据，就不需要找上级了
	if !route.bindResponseData {
		if group.bindResponseData {
			route.responseData = Clone(group.responseData)
		} else {
			route.responseData = Clone(r.responseData)
		}
	}

	// 合并 pagekeys
	// 全局key
	tempPages := r.pagekeys.clone()

	// 组的删除为了删全局
	tempPages.delete(group.delPageKeys)
	// 添加组
	tempPages.add(group.pagekeys)
	// 个人的删除组
	// 删除单路由
	tempPages.delete(route.delPageKeys)
	// 添加个人
	tempPages.add(route.pagekeys)
	// 最终页面权限
	route.pagekeys = tempPages

	// 模块合并
	tempModules := r.module.cloneMudule()
	// 组删除模块为了删全局
	tempModules.delete(group.delmodule)
	// 添加组模块
	tempModules.add(group.module.funcOrder...)
	// 私有删除模块
	tempModules.delete(route.delmodule)
	// 添加私有模块
	tempModules.add(route.module.funcOrder...)
	route.module = tempModules
	// 与组的区别， 组里面这里是合并， 这里是删除
}

// 组路由添加到router里面,
// 挂载到group之前， 全局的变量已经挂载到route 里面了， 所以不用再管组变量了
func (r *Router) AddGroup(group *RouteGroup) *Router {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	// 将路由的所有变量全部移交到route
	if group.params == nil && group.route == nil {
		return nil
	}
	prefixs := SubtractSliceMap(r.prefix, group.delprefix)
	prefix := append(prefixs, group.prefix...)
	for _, route := range group.routes {
		subprefixs := SubtractSliceMap(prefix, route.delprefix)
		subprefix := append(subprefixs, route.prefixs...)
		allurl := path.Join(subprefix...)
		allurl = path.Join(allurl, route.url)
		url, vars, ok := makeRoute(allurl)
		route.params = vars
		if ok {
			if r.tpl[url] == nil {
				r.tpl[url] = make(map[string]*Route)
			}
			for _, method := range route.methods {
				if _, methodOk := r.tpl[url][method]; methodOk {
					// 如果也存在， 那么method重复了
					log.Fatal("method : " + method + "  duplicate, url: " + url)
				}
				if r.tpl[url] == nil {
					r.tpl[url] = make(MethodsRoute)
				}
				// newRoute.methods[method] = struct{}{}
				route.url = url
				route.params = vars
				r.merge(group, route)
			}

		} else {
			if r.route[url] == nil {
				r.route[url] = make(map[string]*Route)
			}
			// 如果存在就判断是否存在method
			for _, method := range route.methods {
				if _, methodOk := r.route[url][method]; methodOk {
					// 如果也存在， 那么method重复了
					log.Fatal("method : " + method + "  duplicate, url: " + url)
				}
				route.url = url
				r.route[url][method] = route
				r.merge(group, route)
			}
			// 如果不存在就创建一个 route

		}
	}
	// for url, args := range group.params {
	// 	r.params[url] = args
	// 	if len(args) == 0 {
	// 		for method := range group.route[url] {
	// 			if _, ok := r.route[url]; ok {
	// 				if _, gok := r.route[url][method]; gok {
	// 					log.Fatal("method : " + method + "  duplicate, url: " + url)
	// 				}
	// 			}
	// 			r.merge(group, group.route[url][method])
	// 		}

	// 		r.route[url] = group.route[url]

	// 	} else {
	// 		for method := range group.tpl[url] {
	// 			if _, ok := r.tpl[url]; ok {
	// 				if _, gok := r.tpl[url][method]; gok {
	// 					log.Fatal("method : " + method + "  duplicate, url: " + url)
	// 				}
	// 			}
	// 			r.merge(group, group.tpl[url][method])
	// 		}

	// 		r.tpl[url] = group.tpl[url]
	// 	}

	// }

	return r
}

// 将路由组的信息合并到路由

func debugPrint(url string, mr MethodsRoute) {
	for method, route := range mr {
		names := make([]string, 0)
		for _, v := range route.module.funcOrder {
			names = append(names, helper.GetFuncName(v))
		}
		log.Printf("url: %s, method: %s, header: %+v, module: %#v,  pages: %#v\n",
			url, method, route.header, names, route.pagekeys)
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
		}
	}
}

func (r *Router) DebugTpl() {
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

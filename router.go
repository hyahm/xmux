package xmux

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
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
	"github.com/quic-go/quic-go/http3"
)

var connections int32 = -1

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

type Router struct {
	addr                 string
	prefix               []string
	MaxPrintLength       int
	Exit                 func(time.Time, http.ResponseWriter, *http.Request)
	new                  bool                                          // 判断是否是通过newRouter 来初始化的
	Enter                func(http.ResponseWriter, *http.Request) bool // 当有请求进入时候的执行
	ReadTimeout          time.Duration
	HanleFavicon         func(http.ResponseWriter, *http.Request)
	DisableOption        bool                                     // 禁止全局option
	HandleOptions        func(http.ResponseWriter, *http.Request) // 预请求 处理函数， 如果存在， 优先处理, 前后端分离后， 前段可能会先发送一个预请求
	HandleNotFound       func(http.ResponseWriter, *http.Request)
	HandleAll            func(http.ResponseWriter, *http.Request) bool
	NotFoundRequireField func(string, http.ResponseWriter, *http.Request) bool
	UnmarshalError       func(error, http.ResponseWriter, *http.Request) bool
	IgnoreSlash          bool     // 忽略地址多个斜杠， 默认不忽略
	urlRoute             UrlRoute // 单实例路由， 组路由最后也会合并过来
	urlTpl               UrlRoute // 正则路由， 组路由最后也会合并过来
	// params               map[string][]string // 记录所有路由， map[string]string 是正则匹配的参数
	header             map[string]string // 全局路由头
	module             *module           // 全局模块
	responseData       interface{}
	pagekeys           mstringstruct
	SwaggerTitle       string
	SwaggerDescription string
	SwaggerVersion     string
}

func (r *Router) SetTimeout(t time.Duration) *Router {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	r.ReadTimeout = t
	return r
}
func (r *Router) BindResponse(response interface{}) *Router {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	r.responseData = response
	return r
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

func (r *Router) readFromCache(route *rt, w http.ResponseWriter, req *http.Request) {
	for k, v := range route.Header {
		w.Header().Set(k, v)
	}

	// 进入前的钩子函数
	if r.Enter != nil {
		if r.Enter(w, req) {
			return
		}
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
	// option 请求处理

	// 退出前的钩子函数
	if r.Exit != nil {
		defer r.Exit(start, w, req)
	}

	if !r.DisableOption && req.Method == http.MethodOptions {
		r.HandleOptions(w, req)
		return
	}
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

	if stop {
		// fd.StatusCode = http.StatusLocked
		w.WriteHeader(http.StatusLocked)
		return
	}
	if r.HandleAll != nil {
		if r.HandleAll(w, req) {
			return
		}
	}
	atomic.AddInt32(&connections, 1)
	defer atomic.AddInt32(&connections, -1)

	if r.IgnoreSlash {
		req.URL.Path = PrettySlash(req.URL.Path)
	}

	// 正在寻址中的话， 就等待寻址完成
	// 先进行路由表缓存寻找
	route, ok := getUrlCache(req.URL.Path + req.Method)

	if ok {
		r.readFromCache(route, w, req)
	} else {
		// 寻址
		r.serveHTTP(w, req)
	}
}

func (r *Router) setHeader(route *Route) map[string]string {
	// 设置请求头
	headers := make(map[string]string)
	for k, v := range r.header {
		headers[k] = v
	}
	for k, v := range route.header {
		headers[k] = v
	}
	for k := range route.delheader {
		delete(headers, k)
	}
	return headers
}

// url 是匹配的路径， 可能不是规则的路径, 寻址的时候还是要加锁
func (r *Router) serveHTTP(w http.ResponseWriter, req *http.Request) {
	var thisRoute *Route

	matchMethod := false
	if route, ok := r.urlRoute[req.URL.Path]; ok {

		for _, v := range route.methods {
			if v == req.Method {
				matchMethod = true
				break
			}
		}
		// route, ok := r.route[req.URL.Path][req.Method]
		if !ok {
			r.HandleNotFound(w, req)
			atomic.AddInt32(&connections, -1)
			return
		}
		if !matchMethod {
			r.HandleNotFound(w, req)
			return
		}
		thisRoute = route
	} else {
		for reUrl := range r.urlTpl {
			re := regexp.MustCompile(reUrl)
			// req.URL.Path = strings.Trim(req.URL.Path, " ")
			if re.MatchString(req.URL.Path) {
				route, ok := r.urlTpl[reUrl]
				if ok {
					// 匹配请求

					for _, v := range route.methods {
						if v == req.Method {
							matchMethod = true
							break
						}
					}
					if !matchMethod {
						r.HandleNotFound(w, req)
						return
					}
					ap := make(map[string]string)
					vl := re.FindStringSubmatch(req.URL.Path)
					for i, v := range route.params {
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

	// todo: 后面补充prefix
	// 缓存handler
	thisRouteCache := &rt{
		Handle:       thisRoute.handle,
		Header:       r.setHeader(thisRoute),
		module:       thisRoute.module.GetModules(),
		dataSource:   thisRoute.dataSource,
		pagekeys:     thisRoute.pagekeys,
		bindType:     thisRoute.bindType,
		responseData: thisRoute.responseData,
	}

	// 设置缓存
	setUrlCache(req.URL.Path+req.Method, thisRouteCache)
	r.readFromCache(thisRouteCache, w, req)
}

func (r *Router) SetAddr(addr string) *Router {
	r.addr = addr
	return r
}

func (r *Router) Run(addr ...string) error {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	if len(addr) > 0 {
		r.addr = addr[0]
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
	<-ctx.Done()
	svc.Close()

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

func (r *Router) RunQuic(certPemFile, keyPemFile string, addr ...string) error {
	certPem, err := os.ReadFile(certPemFile)
	if err != nil {
		return err
	}
	keyPem, err := os.ReadFile(keyPemFile)
	if err != nil {
		return err
	}
	// 1. 准备一个符合 TLS1.3 + ALPN=h3 的证书
	cert, err := tls.X509KeyPair(certPem, keyPem) // 也可使用自签或 ACME
	if err != nil {
		return err
	}

	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h3"}, // 必须声明
	}
	if len(addr) > 0 {
		r.addr = addr[0]
	}
	// 2. 启动 HTTP/3 服务器
	s := http3.Server{
		Addr:      r.addr,
		TLSConfig: tlsConf,
		Handler:   r,
	}

	log.Println("⇨  HTTP/3 server over https on " + r.addr)
	return s.ListenAndServe()
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
		urlRoute:       make(UrlRoute),
		urlTpl:         make(UrlRoute),
		header:         map[string]string{},
		// params:         make(map[string][]string),
		Exit: exit,
		module: &module{
			filter:    make(map[string]struct{}),
			funcOrder: make([]func(w http.ResponseWriter, r *http.Request) bool, 0),
		},
		HanleFavicon:         handleFavicon,
		HandleOptions:        handleOptions,
		HandleAll:            handleAll,
		HandleNotFound:       handleNotFound,
		NotFoundRequireField: notFoundRequireField,

		UnmarshalError: unmarshalError,
	}
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
func (r *Router) merge(group *RouteGroup, route *Route) *Route {
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
	return route
	// 与组的区别， 组里面这里是合并， 这里是删除
}

// 组路由添加到router里面,
// 挂载到group之前， 全局的变量已经挂载到route 里面了， 所以不用再管组变量了
func (r *Router) AddGroup(group *RouteGroup) *Router {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	// 将路由的所有变量全部移交到route
	if group.urlTpl == nil && group.urlRoute == nil {
		return nil
	}
	// 将前缀删除
	// prefixs := SubtractSliceMap(r.prefix, group.delprefix)
	// prefix := append(prefixs, group.prefix...)
	// 匹配理由的合并
	for url, route := range group.urlRoute {
		if _, ok := r.urlRoute[url]; ok {
			v, ok := SliceExsit(r.urlRoute[url].methods, route.methods)
			if ok {
				log.Fatal("method : " + v + "  duplicate, url: " + url)
			}
		}
		newRoute := r.merge(group, route)
		if !route.denyPrefix && len(r.prefix) > 0 {

			url = r.mergePrefix(route, url)
		}
		if !r.DisableOption {
			var exsitOption bool
			for _, v := range newRoute.methods {
				if v == http.MethodOptions {
					exsitOption = true
				}
			}
			if !exsitOption {
				newRoute.methods = append(newRoute.methods, http.MethodOptions)
			}
		}
		// 最终的prefix合并
		r.urlRoute[url] = newRoute
	}
	// 正则匹配的合并

	for url, tpl := range group.urlTpl {
		if _, ok := r.urlTpl[url]; ok {
			v, ok := SliceExsit(r.urlTpl[url].methods, tpl.methods)
			if ok {
				log.Fatal("method : " + v + "  duplicate, url: " + url)
			}
		}
		newRoute := r.merge(group, tpl)
		if len(r.prefix) > 0 {
			url = r.mergePrefix(tpl, url)
		}
		r.urlTpl[url] = newRoute

	}

	return r
}

func (r *Router) mergePrefix(newRoute *Route, url string) string {
	newAddPrefix := append(r.prefix, newRoute.prefixs...)
	prefixs := make([]string, 0)
	if len(newRoute.delprefix) > 0 {
		for _, v := range newAddPrefix {
			if _, ok := newRoute.delprefix[v]; !ok {
				prefixs = append(prefixs, v)
			}
		}
	} else {
		prefixs = newAddPrefix
	}

	if url[0:1] == "^" {
		url = url[1:]
		prefixs = append(prefixs, url)
		url = "^" + path.Join(prefixs...)
	} else {
		prefixs = append(prefixs, url)
		url = path.Join(prefixs...)
	}

	return url
}

// 将路由组的信息合并到路由

func debugPrint(url string, route *Route) {
	names := make([]string, 0)
	for _, v := range route.module.funcOrder {
		names = append(names, helper.GetFuncName(v))
	}
	log.Printf("url: %s, method: %s, header: %+v, module: %#v,  pages: %#v  responsedata: %v\n",
		url, route.methods, route.header, names, route.pagekeys, route.responseData)

}

func (r *Router) DebugRoute() {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	for url, mr := range r.urlRoute {
		debugPrint(url, mr)
	}
}

func (r *Router) DebugAssignRoute(thisurl string) {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	for url, mr := range r.urlRoute {
		if thisurl == url {
			debugPrint(url, mr)
		}
	}
}

func (r *Router) DebugTpl() {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	for url, mr := range r.urlTpl {
		debugPrint(url, mr)
	}
}

func (r *Router) DebugIncludeTpl(pattern string) {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	for url, mr := range r.urlTpl {
		if strings.Contains(url, pattern) {
			debugPrint(url, mr)
		}
	}
}

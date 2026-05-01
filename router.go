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
	"runtime"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"

	"github.com/hyahm/xmux/helper"
	"github.com/quic-go/quic-go/http3"
)

var connections int32 = -1

// const CTX = "XMUX_CTX"
const PAGES = "XMUX_PAGES"

// var stop bool

func GetConnents() int32 {
	return connections
}

type rt struct {
	Handle       http.Handler
	Header       map[string]string
	pagekeys     map[string]struct{}
	module       []func(http.ResponseWriter, *http.Request) bool
	postModule   []func(http.ResponseWriter, *http.Request) bool
	dataSource   interface{} // 绑定数据结构，
	bindType     bindType
	responseData interface{}
	middleware   onion
	funcName     string
	url          string
	methods      []string
	// instance   map[*http.Request]interface{} // 解析到这里
}

type router struct {
	addr           string
	prefix         []string
	MaxPrintLength int
	Exit           func(time.Time, http.ResponseWriter, *http.Request)
	Enter          func(http.ResponseWriter, *http.Request) bool // 当有请求进入时候的执行
	ReadTimeout    time.Duration
	HandleFavicon  func(http.ResponseWriter, *http.Request)
	DisableOption  bool                                     // 禁止全局option
	HandleOptions  func(http.ResponseWriter, *http.Request) // 预请求 处理函数， 如果存在， 优先处理, 前后端分离后， 前段可能会先发送一个预请求
	HandleNotFound func(http.ResponseWriter, *http.Request)
	HandleRecover  func(http.ResponseWriter, *http.Request)
	HandleAll      func(http.ResponseWriter, *http.Request) bool
	// NotFoundRequireField func(string, http.ResponseWriter, *http.Request) bool
	UnmarshalError func(error, http.ResponseWriter, *http.Request) bool
	IgnoreSlash    bool     // 忽略地址多个斜杠， 默认不忽略
	urlRoute       UrlRoute // 单实例路由， 组路由最后也会合并过来
	urlTpl         UrlRoute // 正则路由， 组路由最后也会合并过来
	// params               map[string][]string // 记录所有路由， map[string]string 是正则匹配的参数
	header             map[string]string // 全局路由头
	module             *module           // 全局模块
	postModule         *module           // 全局后置模块
	responseData       interface{}
	ModuleContinue     bool // 模块是否继续执行， 默认false， 只要有一个模块返回true就继续执行了， 取反之意
	pagekeys           mstringstruct
	SwaggerTitle       string
	SwaggerDescription string
	SwaggerVersion     string
	menuTree           []*MenuTree // 记录添加路由的顺序， 方便组件路由树
}

// 拿到所有路由
func (r *router) Routes() []MenuTree {
	routes := make([]MenuTree, 0)
	for url, v := range r.urlRoute {
		routes = append(routes, MenuTree{
			URL:        url,
			Uuid:       v.menuTree.Uuid,
			ParentUUID: v.menuTree.ParentUUID,
			Method:     strings.Join(v.methods, ","),
		})
	}
	for url, v := range r.urlTpl {
		routes = append(routes, MenuTree{
			URL:        url,
			Uuid:       v.menuTree.Uuid,
			ParentUUID: v.menuTree.ParentUUID,
			Method:     strings.Join(v.methods, ","),
		})
	}
	return routes
}

func (r *router) Menus() []MenuTree {
	routes := make([]MenuTree, 0)
	for url, v := range r.urlRoute {
		roles := make([]string, 0, len(v.pagekeys))
		for k := range v.pagekeys {
			roles = append(roles, k)
		}
		ri := MenuTree{
			URL:        url,
			Uuid:       v.menuTree.Uuid,
			ParentUUID: v.menuTree.ParentUUID,
			Method:     strings.Join(v.methods, ","),
			Roles:      roles,
			Meta:       v.menuTree.Meta,
		}
		ri.makeMenuId()
		routes = append(routes, ri)
	}
	for url, v := range r.urlTpl {
		roles := make([]string, 0, len(v.pagekeys))
		for k := range v.pagekeys {
			roles = append(roles, k)
		}
		ri := MenuTree{
			URL:        url,
			Uuid:       v.menuTree.Uuid,
			ParentUUID: v.menuTree.ParentUUID,
			Method:     strings.Join(v.methods, ","),
			Roles:      roles,
			Meta:       v.menuTree.Meta,
		}
		ri.makeMenuId()
		routes = append(routes, ri)
	}
	for _, v := range FlattenMenuTree(r.menuTree) {
		v.makeMenuId()
		routes = append(routes, *v)
	}
	return routes
}

// 设置超时的时候注意   只有 module和 postmodule 的函数块才会中断， enter, exit , handle 不受影响， 如果有耗时代码， 请放到 handle，
// 因为 postmodule 也会中断， 所以尽量不使用 postmodule， 可以写入到 exit
func (r *router) SetTimeout(t time.Duration) *router {
	r.ReadTimeout = t
	return r
}
func (r *router) BindResponse(response interface{}) *router {
	r.responseData = response
	return r
}

func (r *router) SetHeader(k, v string) *router {
	r.header[k] = v
	return r
}

func (r *router) AddPageKeys(pagekeys ...string) *router {
	if r.pagekeys == nil {
		r.pagekeys = make(map[string]struct{})
	}
	for _, v := range pagekeys {
		r.pagekeys[v] = struct{}{}
	}
	return r
}

func (r *router) AddModule(handles ...func(http.ResponseWriter, *http.Request) bool) *router {
	r.module.add(handles...)
	return r
}

func (r *router) AddPostModule(handles ...func(http.ResponseWriter, *http.Request) bool) *router {
	r.postModule.add(handles...)
	return r
}

func (r *router) readFromCache(route *rt, w http.ResponseWriter, req *http.Request) {

	for k, v := range route.Header {
		w.Header().Set(k, v)
	}

	// 进入前的钩子函数
	if r.Enter != nil {
		if r.ModuleContinue {
			if !r.Enter(w, req) {
				return
			}
		} else {
			if r.Enter(w, req) {
				return
			}
		}

	}
	ci := time.Now().UnixNano()
	fd := &FlowData{
		ctx: make(map[string]interface{}),
		// mu:         &sync.RWMutex{},
		connectId:  ci,
		StatusCode: 200,
		module:     route.module,
		funcName:   route.funcName,
		url:        route.url,
	}
	// add ctx data
	ctx := context.WithValue(req.Context(), xmux_context, fd)
	req = req.WithContext(ctx)
	// add module
	start := time.Now()
	// option 请求处理

	// 退出前的钩子函数
	if r.Exit != nil {
		defer r.Exit(start, w, req)
	}

	if route.responseData != nil {
		fd.Response = DeepCopy(route.responseData)
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
			if r.ModuleContinue {
				if !r.bind(route, w, req, fd) {
					return
				}
			} else {
				if r.bind(route, w, req, fd) {
					return
				}
			}

		} else {
			GetInstance(req).Body = []byte("")
		}
	}

	// 权限导入
	// pages
	fd.pages = route.pagekeys
	// 当前函数名去掉目录层级后的

	if r.ReadTimeout > 0 {
		ctx, cancel := context.WithTimeout(req.Context(), r.ReadTimeout)
		defer cancel()
		for _, module := range route.module {
			select {
			case <-ctx.Done():
				w.WriteHeader(http.StatusGatewayTimeout)
				return
			default:
				if r.ModuleContinue {
					if !module(w, req) {
						return
					}
				} else {
					if module(w, req) {
						return
					}
				}
			}

		}
		// 中间件
		if len(route.middleware.mws) > 0 && route.Handle != nil {
			route.Handle = route.middleware.ThenFunc(route.Handle)

		}
		if route.Handle.(http.HandlerFunc) != nil {
			route.Handle.ServeHTTP(w, req)
		}
		// 处理后置模块
		for _, module := range route.postModule {
			select {
			case <-ctx.Done():
				w.WriteHeader(http.StatusGatewayTimeout)
				return
			default:
				if r.ModuleContinue {
					if !module(w, req) {
						return
					}
				} else {
					if module(w, req) {
						return
					}
				}
			}
		}
		return
	}

	// 请求模块
	for _, module := range route.module {
		if r.ModuleContinue {
			if !module(w, req) {
				return
			}
		} else {
			if module(w, req) {
				return
			}
		}

	}
	// 中间件
	if len(route.middleware.mws) > 0 && route.Handle != nil {
		route.Handle = route.middleware.ThenFunc(route.Handle)

	}
	if route.Handle.(http.HandlerFunc) != nil {
		route.Handle.ServeHTTP(w, req)
	}
	// 处理后置模块
	for _, module := range route.postModule {
		if r.ModuleContinue {
			if !module(w, req) {
				return
			}
		} else {
			if module(w, req) {
				return
			}
		}
	}
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	defer func() {
		if err := recover(); err != nil && r.HandleRecover != nil {
			errStack := fmt.Errorf("panic: %v\n%s", err, debug.Stack())
			fmt.Println(errStack)
			r.HandleRecover(w, req)
		}
	}()

	r.shhandle(w, req)
}

func (r *router) shhandle(w http.ResponseWriter, req *http.Request) {

	if r.HandleAll != nil {
		if r.ModuleContinue {
			if !r.HandleAll(w, req) {
				return
			}
		} else {
			if r.HandleAll(w, req) {
				return
			}
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
		r.findRoute(w, req)
	}

}

func (r *router) setHeader(route *route) map[string]string {
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

// var re = regexp.MustCompile(`https?://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)

// url 是匹配的路径， 可能不是规则的路径, 寻址的时候还是要加锁
func (r *router) findRoute(w http.ResponseWriter, req *http.Request) {
	var thisRoute *route
	matchMethod := false
	for k, v := range r.header {
		w.Header().Set(k, v)
	}
	url := req.URL.Path
	if route, ok := r.urlRoute[url]; ok {

		for _, v := range route.methods {
			if v == req.Method {
				matchMethod = true
				break
			}
		}
		// route, ok := r.route[req.URL.Path][req.Method]
		if (!ok || !matchMethod) && r.HandleNotFound != nil {
			r.HandleNotFound(w, req)
			return
		}
		thisRoute = route
	} else {
		for subUrl, route := range r.urlTpl {

			if route.regex != nil && route.regex.MatchString(url) {

				// 匹配请求
				for _, v := range route.methods {
					if v == req.Method {
						matchMethod = true
						break
					}
				}
				if !matchMethod && r.HandleNotFound != nil {
					r.HandleNotFound(w, req)
					return
				}
				ap := make(map[string]string)
				vl := route.regex.FindStringSubmatch(url)
				for i, v := range route.params {
					ap[v] = vl[i+1]
				}
				thisRoute = route
				url = subUrl
				setParams(url, ap)
				goto endloop
			}
		}
		if r.HandleNotFound != nil {
			r.HandleNotFound(w, req)
			return
		}

	}
endloop:

	// todo: 后面补充prefix
	// 缓存handler
	name := runtime.FuncForPC(reflect.ValueOf(thisRoute.handle).Pointer()).Name()
	n := strings.LastIndex(name, ".")
	thisRouteCache := &rt{
		Handle:       thisRoute.handle,
		Header:       r.setHeader(thisRoute),
		module:       thisRoute.module.GetModules(),
		postModule:   thisRoute.postModule.GetModules(),
		dataSource:   thisRoute.dataSource,
		pagekeys:     thisRoute.pagekeys,
		bindType:     thisRoute.bindType,
		responseData: thisRoute.responseData,
		middleware:   thisRoute.middleware,
		funcName:     name[n+1:],
		url:          url,
		methods:      thisRoute.methods,
	}
	// 设置缓存
	setUrlCache(req.URL.Path+req.Method, thisRouteCache)
	r.readFromCache(thisRouteCache, w, req)
}

func (r *router) SetAddr(addr string) *router {
	r.addr = addr
	return r
}

func (r *router) Run(addr ...string) error {
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

func (r *router) Debug(ctx context.Context) {
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

func (r *router) RunUnsafeTLS(opt ...string) error {
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

func (r *router) RunQuic(certPemFile, keyPemFile string, addr ...string) error {
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

func (r *router) RunTLS(certFile, keyFile string) error {

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

func NewRouter(cacheSize ...int) *router {
	var c int
	if len(cacheSize) > 0 {
		c = cacheSize[0]
	}
	initUrlCache(c)
	r := &router{
		addr:           ":8080",
		MaxPrintLength: 2 << 10, // 默认的form最大2k
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
		postModule: &module{
			filter:    make(map[string]struct{}),
			funcOrder: make([]func(w http.ResponseWriter, r *http.Request) bool, 0),
		},
		HandleFavicon:  handleFavicon,
		HandleOptions:  handleOptions,
		HandleNotFound: handleNotFound,
		HandleRecover:  func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(500); w.Write([]byte("server panic")) },
		// NotFoundRequireField: notFoundRequireField,
		UnmarshalError: unmarshalError,
		menuTree:       make([]*MenuTree, 0),
	}
	return r
}

func (r *router) cloneHeader() mstringstring {
	tempHeader := make(map[string]string)
	for k, v := range r.header {
		tempHeader[k] = v
	}
	return tempHeader
}

func (r *router) Prefix(prefixs ...string) *router {
	if len(prefixs) == 0 {
		return r
	}
	r.prefix = append(r.prefix, prefixs...)

	return r
}

// 组路由合并到每个单路由里面
func (r *router) merge(group *RouteGroup, route *route) *route {
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
			route.responseData = DeepCopy(group.responseData)
		} else {
			route.responseData = DeepCopy(r.responseData)
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

	// 前置 模块合并
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

	// 后置模块合并
	tempPostModules := r.postModule.cloneMudule()
	// 组删除模块为了删全局
	tempPostModules.delete(group.delPostModule)
	// 添加组模块
	tempPostModules.add(group.postModule.funcOrder...)
	// 私有删除模块
	tempPostModules.delete(route.delPostModule)
	// 添加私有模块
	tempPostModules.add(route.postModule.funcOrder...)
	route.postModule = tempPostModules
	return route
	// 与组的区别， 组里面这里是合并， 这里是删除
}

// 组路由添加到router里面,
// 挂载到group之前， 全局的变量已经挂载到route 里面了， 所以不用再管组变量了
func (r *router) AddGroup(group *RouteGroup) *router {
	// 将路由的所有变量全部移交到route
	if group.urlTpl == nil && group.urlRoute == nil {
		return nil
	}
	group.menuTree.ParentUUID = "root"
	r.menuTree = append(r.menuTree, group.menuTree)
	// 将前缀删除
	// prefixs := SubtractSliceMap(r.prefix, group.delprefix)
	// prefix := append(prefixs, group.prefix...)
	// 匹配理由的合并
	for url, route := range group.urlRoute {
		if _, ok := r.urlRoute[url]; ok {
			v, ok := SliceExist(r.urlRoute[url].methods, route.methods)
			if ok {
				log.Fatal("method : " + v + "  duplicate, url: " + url)
			}
		}
		newRoute := r.merge(group, route)
		if !route.denyPrefix && len(r.prefix) > 0 {
			url = r.mergePrefix(route, url)
		}

		// 设置父级uuid为组的uuid， 方便组件路由树的生成
		r.urlRoute[url] = newRoute

	}
	// 正则匹配的合并

	for url, tpl := range group.urlTpl {
		if _, ok := r.urlTpl[url]; ok {
			v, ok := SliceExist(r.urlTpl[url].methods, tpl.methods)
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

func (r *router) mergePrefix(newRoute *route, url string) string {
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

func debugPrint(url string, route *route) {
	names := make([]string, 0, len(route.module.funcOrder))
	for _, v := range route.module.funcOrder {
		names = append(names, helper.GetFuncName(v))
	}
	postNames := make([]string, 0, len(route.postModule.funcOrder))
	for _, v := range route.postModule.funcOrder {
		postNames = append(postNames, helper.GetFuncName(v))
	}
	log.Printf("url: %s, method: %s, header: %+v, module: %#v, postModule: %#v,  pages: %#v  responsedata: %v\n",
		url, route.methods, route.header, names, postNames, route.pagekeys, route.responseData)

}

func (r *router) DebugRoute() {
	for url, mr := range r.urlRoute {
		debugPrint(url, mr)
	}
}

func (r *router) DebugAssignRoute(thisurl string) {
	for url, mr := range r.urlRoute {
		if thisurl == url {
			debugPrint(url, mr)
		}
	}
}

func (r *router) DebugTpl() {
	for url, mr := range r.urlTpl {
		debugPrint(url, mr)
	}
}

func (r *router) DebugIncludeTpl(pattern string) {
	for url, mr := range r.urlTpl {
		if strings.Contains(url, pattern) {
			debugPrint(url, mr)
		}
	}
}

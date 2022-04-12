package xmux

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	Handle       http.Handler
	Header       map[string]string
	pagekeys     map[string]struct{}
	module       []func(http.ResponseWriter, *http.Request) bool
	dataSource   interface{} // 绑定数据结构，
	bindType     bindType
	midware      http.Handler
	responseData interface{}
	// instance   map[*http.Request]interface{} // 解析到这里
}

func DefaultExitTemplate(now time.Time, w http.ResponseWriter, r *http.Request) {
}

type Router struct {
	Exit           func(time.Time, http.ResponseWriter, *http.Request)
	new            bool                                          // 判断是否是通过newRouter 来初始化的
	Enter          func(http.ResponseWriter, *http.Request) bool // 当有请求进入时候的执行
	ReadTimeout    time.Duration
	IgnoreIco      bool // 是否忽略 /favicon.ico 请求。 默认忽略
	HanleFavicon   func(http.ResponseWriter, *http.Request)
	DisableOption  bool                                     // 禁止全局option
	HandleOptions  func(http.ResponseWriter, *http.Request) // 预请求 处理函数， 如果存在， 优先处理, 前后端分离后， 前段可能会先发送一个预请求
	HandleNotFound func(http.ResponseWriter, *http.Request)
	IgnoreSlash    bool                // 忽略地址多个斜杠， 默认不忽略
	route          PatternMR           // 单实例路由， 组路由最后也会合并过来
	tpl            PatternMR           // 正则路由， 组路由最后也会合并过来
	params         map[string][]string // 记录所有路由， []string 是正则匹配的参数
	header         map[string]string   // 全局路由头
	module         module              // 全局模块
	responseData   interface{}
	// routeTable     *rt                                             // 路由表
	pagekeys map[string]struct{}
	midware  http.Handler
}

func (r *Router) MiddleWare(midware http.Handler) *Router {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	r.midware = midware
	return r
}

func (r *Router) BindResponse(response interface{}) *Router {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	r.responseData = response
	return r
}

// 判断是否是正则路径， 返回一个路径 string 和是否是正则的 bool
func (r *Router) makeRoute(pattern string) (string, bool) {
	// 格式化路径
	// 创建 methodsRoute

	if r.IgnoreSlash {
		pattern = prettySlash(pattern)
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

func bind(route *rt, req *http.Request, fd *FlowData) {
	// 数据绑定
	var err error
	defer req.Body.Close()
	switch route.bindType {
	case jsonT:
		err = json.NewDecoder(req.Body).Decode(&fd.Data)

	case yamlT:
		err = yaml.NewDecoder(req.Body).Decode(&fd.Data)

	case xmlT:
		err = xml.NewDecoder(req.Body).Decode(&fd.Data)
	}
	if err != nil {
		if err != io.EOF {
			log.Println(err)
		}
	}
}

func (r *Router) readFromCache(start time.Time, route *rt, w http.ResponseWriter, req *http.Request, fd *FlowData) {

	defer func() {
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

		if route.bindType != 0 {
			bind(route, req, fd)
		}

	}
	fd.Response = route.responseData
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
		route.midware.ServeHTTP(w, req)
	} else {
		route.Handle.ServeHTTP(w, req)
	}

}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	fd := &FlowData{
		ctx: make(map[string]interface{}),
		mu:  &sync.RWMutex{},
	}
	allconn.Set(req, fd)
	defer allconn.Del(req)
	fd.Set(STATUSCODE, 200)
	body := make([]byte, 0, 2048)
	fd.Set(RESPONSEBODY, body)
	fd.Set(CONNECTID, time.Now().UnixNano())
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
		fd.Set(STATUSCODE, http.StatusLocked)
		w.WriteHeader(http.StatusLocked)
		return
	}
	atomic.AddInt32(&connections, 1)
	defer atomic.AddInt32(&connections, -1)
	if r.IgnoreSlash {
		req.URL.Path = prettySlash(req.URL.Path)
	}
	// /favicon.ico  和 Option 请求， 不支持自定义请求头和模块
	if req.URL.Path == "/favicon.ico" {
		if r.IgnoreIco {
			return
		} else {
			for k, v := range r.header {
				w.Header().Set(k, v)
			}
			r.HanleFavicon(w, req)
			return
		}
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
	route, ok := Get(req.URL.Path + req.Method)
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
		thisRoute, ok = r.route[req.URL.Path][req.Method]
		if !ok {
			r.HandleNotFound(w, req)
			atomic.AddInt32(&connections, -1)
			return
		}
	} else {
		for reUrl := range r.tpl {
			re := regexp.MustCompile(reUrl)
			req.URL.Path = strings.Trim(req.URL.Path, " ")
			if re.MatchString(req.URL.Path) {
				if thisRoute, ok = r.tpl[reUrl][req.Method]; !ok {
					// 有匹配到，但是没找到方法就继续向下匹配
					continue
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
		r.HandleNotFound(w, req)
		atomic.AddInt32(&connections, -1)
		return
	}
endloop:
	// 缓存handler
	thisRouter := &rt{
		Handle:       thisRoute.handle,
		Header:       thisRoute.header,
		module:       thisRoute.module.getMuduleList(),
		dataSource:   thisRoute.dataSource,
		midware:      thisRoute.midware,
		pagekeys:     thisRoute.pagekeys,
		bindType:     thisRoute.bindType,
		responseData: thisRoute.responseData,
	}
	// 设置缓存
	Set(req.URL.Path+req.Method, thisRouter)
	r.readFromCache(start, thisRouter, w, req, fd)
}

func (r *Router) Run(opt ...string) error {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	addr := ":8080"
	if len(opt) > 0 {
		addr = opt[0]
	}
	svc := &http.Server{
		Addr:        addr,
		ReadTimeout: r.ReadTimeout,
		Handler:     r,
	}
	fmt.Printf("listen on %s\n", addr)
	return svc.ListenAndServe()
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
		CreateTLS()
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

func (r *Router) RunTLS(keyfile, pemfile string, opt ...string) error {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	if strings.Trim(keyfile, "") == "" {
		panic("keyfile is empty")
	}
	if strings.Trim(pemfile, "") == "" {
		panic("pemfile is empty")
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

	// 如果key文件不存在那么就自动生成
	_, err1 := os.Stat(keyfile)
	_, err2 := os.Stat(pemfile)
	if os.IsNotExist(err1) || os.IsNotExist(err2) {
		CreateTLS()
		if err := svc.ListenAndServeTLS(filepath.Join("keys", "server.pem"), filepath.Join("keys", "server.key")); err != nil {
			log.Fatal(err)
		}
	}
	if err := svc.ListenAndServeTLS(pemfile, keyfile); err != nil {
		log.Fatal(err)
	}
	fmt.Println("listen on " + addr + " over https")
	return nil
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
		Exit:           exit,
		params:         make(map[string][]string),
		HanleFavicon:   handleFavicon,
		HandleOptions:  handleOptions,
		HandleNotFound: handleNotFound,
		module:         module{},
	}
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	GetInstance(r).Set(STATUSCODE, http.StatusNotFound)
	w.WriteHeader(http.StatusNotFound)
}

func handleOptions(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func handleFavicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
}

// 组路由合并到每个单路由里面
func (r *Router) merge(group *GroupRoute, route *Route) {
	// 合并head
	tempHeader := make(map[string]string)
	for k, v := range r.header {
		tempHeader[k] = v
	}
	// 组的删除大于全局, 后于组
	for _, v := range group.delheader {
		delete(tempHeader, v)
	}
	// 添加个人路由
	for k, v := range route.header {
		tempHeader[k] = v
	}
	// 删除单路由
	for _, v := range route.delheader {
		delete(tempHeader, v)
	}
	route.header = tempHeader
	// 合并中间件
	if route.midware == nil {
		route.midware = group.midware
		if group.midware == nil {
			route.midware = r.midware
		}
	}
	// 合并返回
	if route.responseData == nil {
		route.responseData = group.responseData

		if group.responseData == nil {
			route.responseData = r.responseData
		}
	}

	// 合并 pagekeys
	tempPages := make(map[string]struct{})
	for k := range r.pagekeys {
		tempPages[k] = struct{}{}
	}
	// 组的删除大于全局, 后于组
	for _, v := range group.delPageKeys {
		delete(tempPages, v)
	}
	for k := range route.pagekeys {
		tempPages[k] = struct{}{}
	}
	// 删除单路由
	for _, v := range route.delPageKeys {
		delete(tempPages, v)
	}
	route.pagekeys = tempPages
	// delete midware, 如果router存在组路由，并且和delmidware相等，那么就删除
	if route.midware != nil && GetFuncName(route.delmidware) == GetFuncName(route.midware) {
		route.midware = nil
	}
	// 模块合并
	route.module = r.module.addModule(route.module)
	// 与组的区别， 组里面这里是合并， 这里是删除
	// 删除模块
	for key := range route.delmodule.modules {
		route.module = route.module.deleteKey(key)
	}

	merge(group, route)
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
				r.merge(group, group.route[url][method])
			}
			r.route[url] = group.route[url]

		} else {
			for method := range group.tpl[url] {
				if _, ok := r.tpl[url][method]; ok {
					log.Fatalf("%s %s is Duplication", url, method)
				}
				r.merge(group, group.tpl[url][method])
			}
			r.tpl[url] = group.tpl[url]
		}

	}

	return r
}

// 将路由组的信息合并到路由
func merge(group *GroupRoute, route *Route) {

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

}

func debugPrint(url string, mr MethodsRoute) {
	for k, v := range mr {
		log.Printf("url: %s, method: %s, header: %+v, module: %#v, midware: %#v , pages: %#v\n",
			url, k, v.header, v.module.funcOrder, GetFuncName(v.midware), v.pagekeys)
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

func (r *Router) GetAssignRoute(thisurl string) MethodsRoute {
	if !r.new {
		panic("must be use get router by NewRouter()")
	}
	for url, mr := range r.route {
		if thisurl == url {
			return mr
		}
	}
	return nil
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

package xmux

import (
	"fmt"
	"net/http"
	"regexp"
	"sync"

	"golang.org/x/net/context"
)

// type Midware

type reroute struct {
	R      *Route
	name   []string // 保存的变量名
	header map[string]string
}

type rt struct {
	ctx     context.Context
	Handle  http.Handler
	Header  map[string]string
	Midware []func(http.ResponseWriter, *http.Request) bool
	end     func(interface{})
}

type Router struct {
	IgnoreIco        bool // 是否忽略 /favicon.ico 请求。 默认忽略
	HanleFavicon     http.Handler
	DisableOption    bool         // 禁止全局option
	Options          http.Handler // 预请求 处理函数， 如果存在， 优先处理, 前后端分离后， 前段可能会先发送一个预请求
	UrlNotFound      http.Handler
	HandleNotFound   http.Handler
	MethodNotAllowed http.Handler
	Doc              http.Handler
	Slash            bool
	route            mr // 单实例路由
	tpl              mr // 正则路由

	group map[string]*GroupRoute // 组路由

	//  标记所有的pattern， 防止有重复的pattern， 0: route 1, tpl, 2, groupRouter 3, groupTpl

	pattern    map[string]int // 完全匹配
	tplpattern map[string]int // 正则匹配

	groupname map[string]string // 根据 pattern 寻找 组名

	header  map[string]string                               // 全局路由头
	midware []func(http.ResponseWriter, *http.Request) bool // 全局中间件

	routeTable *sync.Map // 路由表
	once       *sync.Once
}

func (r *Router) ShowApi(pattern string) *Route {

	return r.Pattern(pattern).Get(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		doc := &Doc{
			Api:   make([]Document, 0),
			Title: "xmux docs",
		}

		t := NewTemplate()
		// 单路由
		for url, v := range r.route {
			if url == pattern {
				continue
			}
			document := v.makeDoc()
			document.Url = url
			document.Supplement = v.supplement
			for mt, _ := range v.method {
				document.Method = mt
				if mt == http.MethodGet {
					if v.params_request != nil {
						document.Url += GetOpt(v.params_request)
					}
				} else {
					if v.st_request != nil {
						document.Opt = PostOpt(v.st_request)
					}
				}
				doc.Add(document)
			}
		}

		for url, v := range r.tpl {
			document := v.makeDoc()
			document.Url = url
			document.Supplement = v.supplement
			for mt, _ := range v.method {
				document.Method = mt
				if mt == http.MethodGet {
					if v.params_request != nil {
						document.Url += GetOpt(v.params_request)
					}
				} else {
					if v.st_request != nil {
						document.Opt = PostOpt(v.st_request)
					}
				}
				doc.Add(document)
			}
		}
		// 组路由

		for _, g := range r.group {
			for url, v := range g.route {
				document := v.makeDoc()
				document.Url = url
				document.Supplement = v.supplement
				for mt, _ := range v.method {
					document.Method = mt
					if mt == http.MethodGet {
						if v.params_request != nil {
							document.Url += GetOpt(v.params_request)
						}
					} else {
						if v.st_request != nil {
							document.Opt = PostOpt(v.st_request)
						}
					}
					doc.Add(document)
				}
			}

			for url, v := range g.tpl {
				document := v.makeDoc()
				document.Url = url
				document.Supplement = v.supplement
				for mt, _ := range v.method {
					document.Method = mt
					if mt == http.MethodGet {
						if v.params_request != nil {
							document.Url += GetOpt(v.params_request)
						}
					} else {
						if v.st_request != nil {
							document.Opt = PostOpt(v.st_request)
						}
					}
					doc.Add(document)
				}
			}
		}
		err := t.Execute(w, *doc)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		return
	}))

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
	if route, ok := r.routeTable.Load(url + req.Method); ok {
		// 设置请求头
		for k, v := range route.(*rt).Header {
			w.Header().Set(k, v)
		}
		defer delete(Ctx, url)
		// 请求中间件
		var ok bool
		for _, v := range route.(*rt).Midware {
			ok = v(w, req)
			if ok {
				return
			}
		}
		route.(*rt).Handle.ServeHTTP(w, req)
		if route.(*rt).end != nil {
			go route.(*rt).end(GetData(req).End)
		}

	} else {
		// 获取handler
		r.serveHTTP(url, w, req)
	}
}

// url 是匹配的路径， 可能不是规则的路径
func (r *Router) serveHTTP(url string, w http.ResponseWriter, req *http.Request) {
	// 应该弄成中间件形式
	var thisHandle http.Handler
	var tp int = -1
	var vl []string
	var matchurl string
	var this_route *Route
	data := &Data{}
	///  寻找路由   ///
	// 先寻找完全匹配的,  优化地方， 先找到路由， 然后找处理接口
	if this_tp, ok := r.pattern[url]; ok {
		matchurl = url
		tp = this_tp
		if r.pattern[url] == 0 {
			// 匹配的单路由
			// 是否能找到方法
			this_route = r.route[url]
		} else {
			// r.pattern[url] 肯定等于 2
			this_route = r.group[r.groupname[matchurl]].route[matchurl]
		}

	} else {
		// 最后正则里面寻找路由
		// reUrl 是一个正则的表达式路径， 是匹配路由的key
		for reUrl, n := range r.tplpattern {
			tp = n
			re := regexp.MustCompile(reUrl)
			if re.MatchString(url) {
				matchurl = reUrl
				vl = re.FindStringSubmatch(url)
				if n == 1 {
					// 单路由
					this_route = r.tpl[matchurl]
					goto endloop
				} else {
					// n == 3
					this_route = r.group[r.groupname[matchurl]].tpl[matchurl]
					goto endloop
				}

			}

		}
		this_route = nil

	}
endloop:

	// 根据路由找处理接口
	if this_route == nil {
		if r.HandleNotFound == nil {
			thisHandle = handleNotFound()
		} else {
			thisHandle = r.HandleNotFound
		}
	}

	if handle, ok := this_route.method[req.Method]; ok {
		// 判断是否有这个方法
		thisHandle = handle
		data.Data = this_route.dataSource
	} else {
		// 如果不存在, 返回MethodNotAllowed
		if r.MethodNotAllowed == nil {
			thisHandle = methodNotAllowed()
		} else {
			thisHandle = r.MethodNotAllowed
		}

	}

	// 如果是正则路由， 添加路由参数到全局里面
	if tp == 1 || tp == 3 {
		vm := make(map[string]string)
		for i, v := range this_route.args {
			vm[v] = vl[i+1]
		}
		data.Var = vm
	}
	slashurl := slash(url)
	Bridge[slashurl] = data

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
	///  结束寻找路由     ///

	// 这里是要是添加组路由的， 也就是tp 是 2或3的
	if tp == 2 || tp == 3 {
		vm := make(map[string]string)
		for i, v := range this_route.args {
			vm[v] = vl[i+1]
		}
		data.Var = vm
		group := r.group[r.groupname[matchurl]]

		// 添加中间件
		for _, v := range group.midware {
			tmpMidware = append(tmpMidware, v)
		}
		// 添加多余的请求头
		for k, v := range group.header {
			tmpHeader[k] = v
			w.Header().Set(k, v)
		}
		// 删除多余的header
		for _, v := range group.delheader {
			delete(tmpHeader, v)
			w.Header().Del(v)
		}
		// 删除多余的中间件
		for _, v := range group.delmidware {
			for i, tmd := range tmpMidware {

				if fmt.Sprintf("%v", v) == fmt.Sprintf("%v", tmd) {
					tmp := make([]func(http.ResponseWriter, *http.Request) bool, 0)
					tmp = append(tmp, tmpMidware[0:i]...)
					tmp = append(tmp, tmpMidware[i+1:]...)
					tmpMidware = tmp
					break
				}
			}

		}
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
			if fmt.Sprintf("%v", v) == fmt.Sprintf("%v", tmd) {
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
		ctx:     context.Background(),
		Handle:  thisHandle,
		Header:  tmpHeader,
		Midware: tmpMidware,
		end:     this_route.end,
	}

	cacheurl := url
	if r.Slash && req.URL.Path != url {
		cacheurl = req.URL.Path

	}
	r.routeTable.Store(cacheurl+req.Method, thisRouter)

	for _, v := range tmpMidware {
		ok := v(w, req)
		if ok {
			return
		}
	}

	thisHandle.ServeHTTP(w, req)
	if this_route.end != nil {
		go this_route.end(GetData(req).End)
	}
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
		header:     map[string]string{},
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

func apiDoc() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})
}

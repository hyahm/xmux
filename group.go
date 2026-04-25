package xmux

import (
	"log"
	"net/http"
	"strings"

	"github.com/hyahm/xmux/helper"
)

// 服务启动前的操作， 所以里面的map 都是单线程不需要加锁的
type routeGroup struct {
	// 感觉还没到method， 应该先uri后缀的
	urlRoute         UrlRoute // 完全匹配的路由对应的methodsroute
	header           mstringstring
	urlTpl           UrlRoute // 正则匹配的路由对应的methodsroute
	module           *module
	postModule       *module
	delmodule        map[string]struct{}
	delPostModule    map[string]struct{}
	responseData     interface{}
	bindResponseData bool
	delheader        map[string]struct{}
	pagekeys         mstringstruct // 页面权限
	prefix           []string
	delprefix        map[string]struct{}
	delPageKeys      map[string]struct{}
	denyPrefix       bool
	routerTrees      RouterTree
}

func NewRouteGroup() *routeGroup {
	group := &routeGroup{
		header: make(map[string]string),
		module: &module{
			filter:    make(map[string]struct{}),
			funcOrder: make([]func(w http.ResponseWriter, r *http.Request) bool, 0),
		},
		postModule: &module{
			filter:    make(map[string]struct{}),
			funcOrder: make([]func(w http.ResponseWriter, r *http.Request) bool, 0),
		},
		prefix:        []string{"/"},
		delprefix:     make(map[string]struct{}),
		delmodule:     make(map[string]struct{}),
		delPostModule: make(map[string]struct{}),
		// params:    make(map[string][]string),
		urlRoute: make(UrlRoute),
		urlTpl:   make(UrlRoute),

		// routes:    make([]*Route, 0),
	}
	if enableRouterTree {
		group.routerTrees = RouterTree{
			Metas: make([]Meta, 0),
		}
	}
	return group
}

func (g *routeGroup) BindResponse(response interface{}) *routeGroup {
	g.responseData = response
	g.bindResponseData = true
	return g
}

func (g *routeGroup) DebugAssignRoute(thisurl string) {
	for url, mr := range g.urlRoute {
		if thisurl == url {
			debugPrint(url, mr)
			return
		}
	}
}

func (g *routeGroup) DebugIncludeTpl(pattern string) {
	for url, mr := range g.urlTpl {
		if strings.Contains(url, pattern) {
			debugPrint(url, mr)
		}
	}
}

func (g *routeGroup) DenyPrefix() {
	g.denyPrefix = true
}

func (g *routeGroup) AddPageKeys(pagekeys ...string) *routeGroup {
	// 接口的请求头
	for _, v := range pagekeys {
		if g.pagekeys == nil {
			g.pagekeys = make(map[string]struct{})
		}
		g.pagekeys[v] = struct{}{}
	}
	return g
}

func (g *routeGroup) SetHeader(k, v string) *routeGroup {
	g.header[k] = v
	return g
}

func (g *routeGroup) DelHeader(headers ...string) *routeGroup {
	if g.delheader == nil {
		g.delheader = make(map[string]struct{})
	}
	for _, header := range headers {
		g.delheader[header] = struct{}{}
	}
	return g
}

func (g *routeGroup) DelPageKeys(pagekeys ...string) *routeGroup {
	if g.delPageKeys == nil {
		g.delPageKeys = make(map[string]struct{})
	}
	for _, key := range pagekeys {
		g.delPageKeys[key] = struct{}{}
	}
	return g
}

func (g *routeGroup) AddModule(handles ...func(http.ResponseWriter, *http.Request) bool) *routeGroup {
	g.module.add(handles...)
	return g
}

func (g *routeGroup) AddPostModule(handles ...func(http.ResponseWriter, *http.Request) bool) *routeGroup {
	g.postModule.add(handles...)
	return g
}

func (g *routeGroup) DelModule(handles ...func(http.ResponseWriter, *http.Request) bool) *routeGroup {
	for _, handle := range handles {
		g.delmodule[helper.GetFuncName(handle)] = struct{}{}
	}
	return g
}

func (g *routeGroup) DelPostModule(handles ...func(http.ResponseWriter, *http.Request) bool) *routeGroup {
	for _, handle := range handles {
		g.delPostModule[helper.GetFuncName(handle)] = struct{}{}
	}
	return g
}

func (g *routeGroup) Prefix(prefixs ...string) *routeGroup {
	g.prefix = append(g.prefix, prefixs...)
	return g
}

func (g *routeGroup) DelPrefix(prefixs ...string) *routeGroup {
	for _, prefix := range prefixs {
		g.delprefix[prefix] = struct{}{}
	}
	return g
}

// 第三个参数返回的true 就是正则
func makeRoute(pattern string) (string, []string, bool) {
	// 格式路径
	if v, listvar := match(pattern); len(listvar) > 0 {
		return v, listvar, true
		// 判断是否重复
	} else {
		return pattern, nil, false
	}
}

func (g *routeGroup) merge(group *routeGroup, route *route) *route {

	// 合并head
	tempHeader := g.header.clone()

	for k := range g.delPageKeys {
		route.delheader[k] = struct{}{}
	}
	// 子组的删除是为了删父辈
	for k := range group.delheader {
		delete(tempHeader, k)
		// 要删除的可能还在上上级， 所以要添加到子路由里面
		route.delheader[k] = struct{}{}
	}
	// 添加组路由的
	for k, v := range group.header {
		tempHeader[k] = v
	}
	// 删除私有路由的
	for k := range route.delheader {
		delete(tempHeader, k)
	}
	// 添加个人路由
	for k, v := range route.header {
		tempHeader[k] = v
	}

	// 最终请求头
	route.header = tempHeader

	// 合并返回
	if !route.bindResponseData {
		if group.bindResponseData {
			route.responseData = DeepCopy(group.responseData)
		} else {
			route.responseData = DeepCopy(g.responseData)
		}
	}

	// 合并 pagekeys
	// 全局key
	tempPages := g.pagekeys.clone()

	for k := range g.delPageKeys {
		route.delPageKeys[k] = struct{}{}
	}

	// 组的删除为了删全局
	for k := range group.delPageKeys {
		delete(tempPages, k)
		route.delPageKeys[k] = struct{}{}
	}
	// 添加组
	for k := range group.pagekeys {
		tempPages[k] = struct{}{}
	}
	// 删除单路由
	for k := range route.delPageKeys {
		delete(tempPages, k)
	}
	// 添加个人
	for k := range route.pagekeys {
		tempPages[k] = struct{}{}
	}
	// 最终页面权限
	route.pagekeys = tempPages

	// 模块合并
	tempModules := g.module.cloneMudule()

	for k := range g.delmodule {
		route.delmodule[k] = struct{}{}
	}
	// 组删除模块为了删全局
	tempModules.delete(group.delmodule)
	for k := range group.delmodule {
		route.delmodule[k] = struct{}{}
	}
	// 添加组模块
	tempModules.add(group.module.funcOrder...)
	// 私有删除模块
	tempModules.delete(route.delmodule)
	// 添加私有模块
	tempModules.add(route.module.funcOrder...)
	route.module = tempModules
	return route
}

// 组路由添加到组路由
func (g *routeGroup) AddGroup(group *routeGroup) *routeGroup {
	// 将路由的所有变量全部移交到route
	if group == nil || (group.urlTpl == nil && group.urlRoute == nil) {
		return g
	}

	// 缺少 请求头， 前缀， 模块， 响应数据 的合并
	for url, route := range group.urlRoute {
		if _, ok := g.urlRoute[url]; ok {
			log.Fatal("url : " + url + "  duplicate")
		}
		// 合并prefix, 主要是合并到 group 里面的路由里面
		newRoute := g.merge(group, route)
		newRoute.prefixs = append(g.prefix, newRoute.prefixs...)
		for key := range g.delprefix {
			newRoute.delprefix[key] = struct{}{}
		}
		g.urlRoute[url] = newRoute
	}

	for url, route := range group.urlTpl {
		if _, ok := g.urlTpl[url]; ok {
			log.Fatal("url : " + url + "  duplicate")
		}
		// 合并prefix, 主要是合并到 group 里面的路由里面
		newRoute := g.merge(group, route)
		newRoute.prefixs = append(g.prefix, newRoute.prefixs...)
		for key := range g.delprefix {
			newRoute.delprefix[key] = struct{}{}
		}
		g.urlTpl[url] = newRoute
	}
	if enableRouterTree {
		g.routerTrees.AddChild(&group.routerTrees)
	}
	return g
}

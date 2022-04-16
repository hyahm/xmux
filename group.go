package xmux

import (
	"log"
	"net/http"
	"strings"
	"sync"
)

// 服务启动前的操作， 所以里面的map 都是单线程不需要加锁的
type GroupRoute struct {
	new bool
	// 感觉还没到method， 应该先uri后缀的
	route            map[string]*Route // 完全匹配的路由对应的methodsroute
	header           map[string]string
	tpl              map[string]*Route // 正则匹配的路由对应的methodsroute
	module           *module
	delmodule        map[string]struct{}
	responseData     interface{}
	bindResponseData bool
	params           map[string][]string // value 是 args， 如果长度是0， 那就是完全匹配， 大于0就是正则匹配
	delheader        map[string]struct{}
	pagekeys         map[string]struct{} // 页面权限

	delPageKeys map[string]struct{}
}

func NewGroupRoute() *GroupRoute {
	return &GroupRoute{
		header: make(map[string]string),
		module: &module{
			filter:    make(map[string]struct{}),
			funcOrder: make([]func(w http.ResponseWriter, r *http.Request) bool, 0),
			mu:        sync.RWMutex{},
		},
		new:       true,
		delmodule: make(map[string]struct{}),
		params:    make(map[string][]string),
		route:     make(map[string]*Route),
		tpl:       make(map[string]*Route),
	}
}

func (g *GroupRoute) BindResponse(response interface{}) *GroupRoute {
	if !g.new {
		panic("must be init by NewGroupRoute()")
	}
	g.responseData = response
	g.bindResponseData = true
	return g
}

func (g *GroupRoute) DebugAssignRoute(thisurl string) {
	if !g.new {
		panic("must be use get router by NewRouter()")
	}
	for url, mr := range g.route {
		if thisurl == url {
			debugPrint(url, mr)
			return
		}
	}
}

func (g *GroupRoute) DebugIncludeTpl(pattern string) {
	if !g.new {
		panic("must be use get router by NewRouter()")
	}
	for url, mr := range g.tpl {
		if strings.Contains(url, pattern) {
			debugPrint(url, mr)
		}
	}
}

func (g *GroupRoute) AddPageKeys(pagekeys ...string) *GroupRoute {
	if !g.new {
		panic("must be init by NewGroupRoute()")
	}
	// 接口的请求头
	for _, v := range pagekeys {
		if g.pagekeys == nil {
			g.pagekeys = make(map[string]struct{})
		}
		g.pagekeys[v] = struct{}{}
	}
	return g
}

func (g *GroupRoute) AddHeader(k, v string) *GroupRoute {
	if !g.new {
		panic("must be init by NewGroupRoute()")
	}
	g.header[k] = v
	return g
}

func (g *GroupRoute) DelHeader(headers ...string) *GroupRoute {
	if !g.new {
		panic("must be init by NewGroupRoute()")
	}
	if g.delheader == nil {
		g.delheader = make(map[string]struct{})
	}
	for _, header := range headers {
		g.delheader[header] = struct{}{}
	}
	return g
}

func (g *GroupRoute) DelPageKeys(pagekeys ...string) *GroupRoute {
	if !g.new {
		panic("must be init by NewGroupRoute()")
	}
	if g.delPageKeys == nil {
		g.delPageKeys = make(map[string]struct{})
	}
	for _, key := range pagekeys {
		g.delPageKeys[key] = struct{}{}
	}
	return g
}

func (g *GroupRoute) AddModule(handles ...func(http.ResponseWriter, *http.Request) bool) *GroupRoute {
	if !g.new {
		panic("must be init by NewGroupRoute()")
	}
	g.module.add(handles...)
	return g
}

func (g *GroupRoute) DelModule(handles ...func(http.ResponseWriter, *http.Request) bool) *GroupRoute {
	if !g.new {
		panic("must be init by NewGroupRoute()")
	}
	for _, handle := range handles {
		g.delmodule[GetFuncName(handle)] = struct{}{}
	}
	return g
}

// 组里面也包括路由 后面的其实还是patter和handle,
// 根据路径来判断是不是正则表达式， 分别挂载到组路由的tpl 和 route 中
// 路径对应的 params 全部都在 pattern 中
// 返回url 和 是否是正则表达式
func (g *GroupRoute) makeRoute(pattern string) (string, []string, bool) {
	// 格式路径
	if v, listvar := match(pattern); len(listvar) > 0 {
		return v, listvar, true
		// 判断是否重复
	} else {
		return pattern, nil, false
	}
}

func (g *GroupRoute) merge(group *GroupRoute, route *Route) {
	// 合并head
	tempHeader := make(map[string]string)
	for k, v := range g.header {
		tempHeader[k] = v
	}

	for k := range g.delheader {
		route.delheader[k] = struct{}{}
	}
	// 组的删除是为了删全局
	for k := range group.delheader {
		delete(tempHeader, k)
		route.delheader[k] = struct{}{}
	}
	// 添加组路由的
	for k, v := range group.header {
		tempHeader[k] = v
	}
	// 私有路由删除组合全局的
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
			route.responseData = Clone(group.responseData)
		} else {
			route.responseData = Clone(g.responseData)
		}
	}

	// 合并 pagekeys
	tempPages := make(map[string]struct{})
	// 全局key
	for k := range g.pagekeys {
		tempPages[k] = struct{}{}
	}

	for k := range g.delPageKeys {
		route.delPageKeys[k] = struct{}{}
	}
	// 组的删除为了删全局
	for k := range group.delPageKeys {
		delete(tempPages, k)
		route.delheader[k] = struct{}{}
	}
	// 添加组
	for k := range group.pagekeys {
		tempPages[k] = struct{}{}
	}
	// 个人的删除组
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

	// 与组的区别， 组里面这里是合并， 这里是删除
}

// 组路由添加到组路由
func (g *GroupRoute) AddGroup(group *GroupRoute) *GroupRoute {
	if !g.new {
		panic("must be init by NewGroupRoute()")
	}
	// 将路由的所有变量全部移交到route
	if group == nil || (group.params == nil && group.route == nil) {
		return g
	}
	for url, args := range group.params {
		g.params[url] = args
		if len(args) == 0 {
			g.merge(group, group.route[url])
			if _, ok := g.route[url]; ok {
				if method, ok := exsitMethod(g.route[url].methods, group.route[url].methods); ok {
					log.Fatal("method : " + method + "  duplicate, url: " + url)
				}
			}
			g.route[url] = group.route[url]

		} else {
			g.merge(group, group.tpl[url])
			if _, ok := g.tpl[url]; ok {
				if method, ok := exsitMethod(g.tpl[url].methods, group.route[url].methods); ok {
					log.Fatal("method : " + method + "  duplicate, url: " + url)
				}
			}
			g.tpl[url] = group.tpl[url]
		}
	}
	return g
}

func exsitMethod(m1, m2 map[string]struct{}) (string, bool) {
	for k := range m1 {
		if _, ok := m2[k]; ok {
			return k, true
		}
	}
	return "", false
}

package xmux

import (
	"log"
	"net/http"
)

// 服务启动前的操作， 所以里面的map 都是单线程不需要加锁的
type GroupRoute struct {
	// 感觉还没到method， 应该先uri后缀的
	route        PatternMR // 完全匹配的路由对应的methodsroute
	ignoreSlash  bool
	header       map[string]string
	tpl          PatternMR // 正则匹配的路由对应的methodsroute
	module       module
	delmodule    delModule
	responseData interface{}
	params       map[string][]string // value 是 args， 如果长度是0， 那就是完全匹配， 大于0就是正则匹配
	delheader    []string
	pagekeys     map[string]struct{} // 页面权限
	groupKey     string
	delPageKeys  []string
	groupTitle   string
	groupLabel   string
	response     interface{}
	reqHeader    map[string]string
	codeMsg      map[string]string
	codeField    string
	midware      func(w http.ResponseWriter, r *http.Request)
	delmidware   func(w http.ResponseWriter, r *http.Request)
}

func NewGroupRoute() *GroupRoute {
	return &GroupRoute{
		header: make(map[string]string),
		module: module{},
	}
}

func (g *GroupRoute) BindResponse(response interface{}) *GroupRoute {
	g.responseData = response
	return g
}

func (g *GroupRoute) AddPageKeys(pagekeys ...string) *GroupRoute {
	// 接口的请求头
	for _, v := range pagekeys {
		if g.pagekeys == nil {
			g.pagekeys = make(map[string]struct{})
		}
		g.pagekeys[v] = struct{}{}
	}
	return g
}

func (g *GroupRoute) ApiReqHeader(k, v string) *GroupRoute {
	// 接口的请求头
	if g.reqHeader == nil {
		g.reqHeader = make(map[string]string)
	}
	g.reqHeader[k] = v
	return g
}

func (g *GroupRoute) MiddleWare(midware func(w http.ResponseWriter, r *http.Request)) *GroupRoute {
	// 接口的请求头
	g.midware = midware
	return g
}

func (g *GroupRoute) DelMiddleWare(midware func(w http.ResponseWriter, r *http.Request)) *GroupRoute {
	// 接口的请求头
	g.delmidware = midware
	return g
}

func (g *GroupRoute) AddHeader(k, v string) *GroupRoute {

	if g.header == nil {
		g.header = make(map[string]string)
	}
	g.header[k] = v
	return g
}

func (g *GroupRoute) ApiCodeMsg(k, v string) *GroupRoute {

	if g.codeMsg == nil {
		g.codeMsg = make(map[string]string)
	}
	g.codeMsg[k] = v
	return g
}

func (g *GroupRoute) ApiCodeField(name string) *GroupRoute {

	g.codeField = name
	return g
}

func (g *GroupRoute) DelHeader(k string) *GroupRoute {

	if g.delheader == nil {
		g.delheader = make([]string, 0)
	}
	g.delheader = append(g.delheader, k)
	return g
}

func (g *GroupRoute) DelPageKeys(pagekeys ...string) *GroupRoute {
	if g.delPageKeys == nil {
		if len(pagekeys) > 0 {
			g.delPageKeys = pagekeys
		}
		return g
	}
	g.delPageKeys = append(g.delPageKeys, pagekeys...)
	return g
}

func (g *GroupRoute) ApiCreateGroup(key string, title string, lable string) *GroupRoute {
	// 组api文档的key，组路由下面的全部会绑定到这个key下面, 如果key 为空， 则无效

	g.groupKey = key
	g.groupLabel = lable
	g.groupTitle = title
	return g
}

func (g *GroupRoute) AddModule(handles ...func(http.ResponseWriter, *http.Request) bool) *GroupRoute {
	g.module = g.module.add(handles...)
	return g
}

func (g *GroupRoute) DelModule(handles ...func(http.ResponseWriter, *http.Request) bool) *GroupRoute {
	g.delmodule = g.delmodule.addDeleteKey(handles...)
	return g
}

// 组里面也包括路由 后面的其实还是patter和handle,
// 根据路径来判断是不是正则表达式， 分别挂载到组路由的tpl 和 route 中
// 路径对应的 params 全部都在 pattern 中
// 返回url 和 是否是正则表达式
func (g *GroupRoute) makeRoute(pattern string) (string, bool) {
	// 格式路径
	if g.ignoreSlash {
		pattern = prettySlash(pattern)
	}

	if g.params == nil {
		g.params = make(map[string][]string)
	}

	if g.route == nil {
		g.route = make(map[string]MethodsRoute)
	}

	if g.tpl == nil {
		g.tpl = make(map[string]MethodsRoute)
	}

	if v, listvar := match(pattern); len(listvar) > 0 {
		if _, ok := g.tpl[v]; !ok {
			g.tpl[v] = make(map[string]*Route)
		}
		g.params[v] = listvar
		return v, true
		// 判断是否重复
	} else {
		if _, ok := g.route[pattern]; !ok {
			g.route[pattern] = make(map[string]*Route)
		}
		g.params[pattern] = make([]string, 0)
		return pattern, false
	}
}

func (g *GroupRoute) merge(group *GroupRoute, route *Route) {
	tempHeader := make(map[string]string)
	for k, v := range g.header {
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
			route.midware = g.midware
		}
	}
	// 合并返回
	if route.responseData == nil {
		route.responseData = group.responseData
		if group.responseData == nil {
			route.responseData = g.responseData
		}
	}

	// 合并 pagekeys
	tempPages := make(map[string]struct{})
	for k := range g.pagekeys {
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
	route.module = g.module.addModule(route.module)
	// 与组的区别， 组里面这里是合并， 这里是删除
	// 删除模块
	for key := range route.delmodule.modules {
		route.module = route.module.deleteKey(key)
	}

	merge(group, route)
}

// 组路由添加到组路由
func (g *GroupRoute) AddGroup(group *GroupRoute) *GroupRoute {
	// 将路由的所有变量全部移交到route
	if group == nil || (group.params == nil && group.route == nil) {
		return g
	}
	if g.header == nil {
		g.header = make(map[string]string)
	}
	if g.params == nil {
		g.params = make(map[string][]string)
	}
	if g.route == nil {
		g.route = make(map[string]MethodsRoute)
	}
	if g.tpl == nil {
		g.tpl = make(map[string]MethodsRoute)
	}

	for url, args := range group.params {
		g.params[url] = args
		if len(args) == 0 {
			for method := range group.route[url] {
				if _, ok := g.route[url][method]; ok {
					log.Fatalf("%s %s is Duplication", url, method)
				}
				g.merge(group, group.route[url][method])
			}
			g.route[url] = group.route[url]

		} else {
			for method := range group.tpl[url] {
				if _, ok := g.tpl[url][method]; ok {
					log.Fatalf("%s %s is Duplication", url, method)
				}
				g.merge(group, group.tpl[url][method])
			}
			g.tpl[url] = group.tpl[url]
		}
	}
	return g
}

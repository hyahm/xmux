package xmux

import (
	"log"
	"net/http"
	"reflect"
)

// 服务启动前的操作， 所以里面的map 都是单线程不需要加锁的
type GroupRoute struct {
	// 感觉还没到method， 应该先uri后缀的
	route      PatternMR // 路由对应的methodsroute
	slash      bool
	header     map[string]string
	tpl        PatternMR // 路由对应的methodsroute
	module     []func(http.ResponseWriter, *http.Request) bool
	delmodule  []func(http.ResponseWriter, *http.Request) bool
	pattern    map[string][]string // value 是 args， 如果长度是0， 那就是完全匹配， 大于0就是正则匹配
	delheader  []string
	pagekeys   map[int]struct{} // 页面权限
	groupKey   string
	groupTitle string
	groupLable string
	reqHeader  map[string]string
	codeMsg    map[string]string
	codeField  string
	midware    func(handle func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request)
	delmidware func(handle func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request)
}

func NewGroupRoute() *GroupRoute {
	return &GroupRoute{
		header: make(map[string]string),
		module: make([]func(http.ResponseWriter, *http.Request) bool, 0),
	}
}

func (g *GroupRoute) SetPageKeys(pagekeys []int) *GroupRoute {
	// 接口的请求头
	if g.pagekeys == nil {
		g.pagekeys = make(map[int]struct{})
	}
	for _, v := range pagekeys {
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

func (g *GroupRoute) MiddleWare(midware func(handle func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request)) *GroupRoute {
	// 接口的请求头
	g.midware = midware
	return g
}

func (g *GroupRoute) DelMiddleWare(midware func(handle func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request)) *GroupRoute {
	// 接口的请求头
	g.midware = midware
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

func (g *GroupRoute) ApiCreateGroup(key string, title string, lable string) *GroupRoute {
	// 组api文档的key，组路由下面的全部会绑定到这个key下面, 如果key 为空， 则无效

	g.groupKey = key
	g.groupLable = lable
	g.groupTitle = title
	return g
}

func (g *GroupRoute) AddModule(handle func(http.ResponseWriter, *http.Request) bool) *GroupRoute {

	if g.module == nil {
		g.module = make([]func(http.ResponseWriter, *http.Request) bool, 0)
	}
	g.module = append(g.module, handle)
	return g
}

func (g *GroupRoute) DelModule(handle func(http.ResponseWriter, *http.Request) bool) *GroupRoute {

	if g.delmodule == nil {
		g.delmodule = make([]func(http.ResponseWriter, *http.Request) bool, 0)
	}
	g.delmodule = append(g.delmodule, handle)
	return g
}

// 组里面也包括路由 后面的其实还是patter和handle, 是否是正则表达式
func (g *GroupRoute) makeRoute(pattern string) (string, bool) {

	// 格式路径
	if g.slash {
		pattern = slash(pattern)
	}

	if g.pattern == nil {
		g.pattern = make(map[string][]string)
	}

	if v, listvar := match(pattern); len(listvar) > 0 {
		if g.tpl == nil {
			g.tpl = make(map[string]MethodsRoute)

		}
		if _, ok := g.tpl[v]; !ok {
			g.tpl[v] = make(map[string]*Route)
		}
		g.pattern[v] = listvar
		return v, true
		// 判断是否重复
	} else {
		if g.route == nil {
			g.route = make(map[string]MethodsRoute)
		}
		if _, ok := g.route[pattern]; !ok {
			g.route[pattern] = make(map[string]*Route)
		}
		g.pattern[pattern] = make([]string, 0)
		return pattern, false
	}
}

func (g *GroupRoute) appendVarToRoute() {
	for _, mr := range g.route {
		// mr 是methodsRoute

		// mt 是 method
		for _, rt := range mr {
			if g.groupKey != "" && rt.groupKey == "" {
				rt.groupKey = g.groupKey
			}

			if g.reqHeader != nil {
				if rt.reqHeader == nil {
					rt.reqHeader = make(map[string]string)
				}
				for k, v := range g.reqHeader {
					rt.reqHeader[k] = v
				}
			}
		}

	}

	for _, mr := range g.tpl {
		for _, rt := range mr {
			if g.groupKey != "" && rt.groupKey == "" {
				rt.groupKey = g.groupKey
			}

			if g.reqHeader != nil {
				for k, v := range g.reqHeader {
					if rt.reqHeader == nil {
						rt.reqHeader = make(map[string]string)
					}
					rt.reqHeader[k] = v
				}
			}
		}
	}

}

// 组路由添加到组路由
func (g *GroupRoute) AddGroup(group *GroupRoute) *GroupRoute {
	// 将路由的所有变量全部移交到route
	for k := range group.pagekeys {
		g.pagekeys[k] = struct{}{}
	}
	if group.pattern == nil && group.route == nil {
		return nil
	}
	if g.header == nil {
		g.header = make(map[string]string)
	}
	if g.pattern == nil {
		g.pattern = make(map[string][]string)
	}
	if g.route == nil {
		g.route = make(map[string]MethodsRoute)
	}
	if g.tpl == nil {
		g.tpl = make(map[string]MethodsRoute)
	}
	if g.midware != nil && !reflect.DeepEqual(group.delmidware, g.midware) && group.midware == nil {
		group.midware = g.midware
	}
	for url, args := range group.pattern {
		g.pattern[url] = args

		if len(args) == 0 {
			for m := range group.route[url] {
				if _, ok := g.route[url][m]; ok {
					log.Fatalf("%s %s is Duplication", url, m)
				}
				merge(group, group.route[url][m])
			}
			g.route[url] = group.route[url]

		} else {
			for m := range group.tpl[url] {
				if _, ok := g.tpl[url][m]; ok {
					log.Fatalf("%s %s is Duplication", url, m)
				}

				merge(group, group.tpl[url][m])

			}

			g.tpl[url] = group.tpl[url]

		}

	}
	return g
}

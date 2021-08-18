package xmux

import (
	"log"
	"net/http"
)

// 服务启动前的操作， 所以里面的map 都是单线程不需要加锁的
type GroupRoute struct {
	// 感觉还没到method， 应该先uri后缀的
	route       PatternMR // 完全匹配的路由对应的methodsroute
	slash       bool
	header      map[string]string
	tpl         PatternMR // 正则匹配的路由对应的methodsroute
	module      module
	delmodule   delModule
	params      map[string][]string // value 是 args， 如果长度是0， 那就是完全匹配， 大于0就是正则匹配
	delheader   []string
	pagekeys    map[string]struct{} // 页面权限
	groupKey    string
	delPageKeys []string
	groupTitle  string
	groupLabel  string
	reqHeader   map[string]string
	codeMsg     map[string]string
	codeField   string
	midware     func(handle func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request)
	delmidware  func(handle func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request)
}

func NewGroupRoute() *GroupRoute {
	return &GroupRoute{
		header: make(map[string]string),
		module: module{},
	}
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
	if g.slash {
		pattern = slash(pattern)
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

func (g *GroupRoute) appendVarToRoute() {
	// for _, mr := range g.route {
	// mr 是methodsRoute

	// mt 是 method
	// for _, rt := range mr {
	// 	if g.groupKey != "" && rt.groupKey == "" {
	// 		rt.groupKey = g.groupKey
	// 	}

	// 	if g.reqHeader != nil {
	// 		if rt.reqHeader == nil {
	// 			rt.reqHeader = make(map[string]string)
	// 		}
	// 		for k, v := range g.reqHeader {
	// 			rt.reqHeader[k] = v
	// 		}
	// 	}
	// }

	// }

	// for _, mr := range g.tpl {
	// 	for _, rt := range mr {
	// 		if g.groupKey != "" && rt.groupKey == "" {
	// 			rt.groupKey = g.groupKey
	// 		}

	// 		if g.reqHeader != nil {
	// 			for k, v := range g.reqHeader {
	// 				if rt.reqHeader == nil {
	// 					rt.reqHeader = make(map[string]string)
	// 				}
	// 				rt.reqHeader[k] = v
	// 			}
	// 		}
	// 	}
	// }

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
				mergeDoc(group, group.route[url][method])
				if group.midware == nil {
					group.route[url][method].midware = g.midware
				}
				// 合并head
				tempHeader := make(map[string]string)
				for k, v := range g.header {
					tempHeader[k] = v
				}
				for k, v := range group.route[url][method].header {
					tempHeader[k] = v
				}
				group.route[url][method].header = tempHeader

				// 合并 delheader
				group.route[url][method].delheader = append(g.delheader, group.route[url][method].delheader...)

				// 合并 pagekeys
				tempPages := make(map[string]struct{})
				for k := range g.pagekeys {
					tempPages[k] = struct{}{}
				}

				for k := range group.route[url][method].pagekeys {
					tempPages[k] = struct{}{}
				}
				group.route[url][method].pagekeys = tempPages
				// 合并delpagekeys

				group.route[url][method].delPageKeys = append(g.delPageKeys, group.route[url][method].delPageKeys...)

				group.route[url][method].module = g.module.addModule(group.route[url][method].module)
			}
			g.route[url] = group.route[url]

		} else {
			for method := range group.tpl[url] {
				if _, ok := g.tpl[url][method]; ok {
					log.Fatalf("%s %s is Duplication", url, method)
				}
				if group.midware == nil {
					group.tpl[url][method].midware = g.midware
				}
				mergeDoc(group, group.tpl[url][method])
				// 合并header
				tempHeader := make(map[string]string)
				for k, v := range g.header {
					tempHeader[k] = v
				}
				for k, v := range group.tpl[url][method].header {
					tempHeader[k] = v
				}
				group.tpl[url][method].header = tempHeader

				// 合并delheader
				group.tpl[url][method].delheader = append(g.delheader, group.tpl[url][method].delheader...)

				// 合并 pagekeys
				tempPages := make(map[string]struct{})
				for k := range g.pagekeys {
					tempPages[k] = struct{}{}
				}

				for k := range group.tpl[url][method].pagekeys {
					tempPages[k] = struct{}{}
				}
				group.tpl[url][method].pagekeys = tempPages

				group.tpl[url][method].delPageKeys = append(g.delPageKeys, group.tpl[url][method].delPageKeys...)

				group.tpl[url][method].module = g.module.addModule(group.tpl[url][method].module)
			}

			g.tpl[url] = group.tpl[url]

		}

	}
	return g
}

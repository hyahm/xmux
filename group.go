package xmux

import (
	"log"
	"net/http"
)

//
type GroupRoute struct {
	// 感觉还没到method， 应该先uri后缀的
	route      PatternMR // 路由对应的methodsroute
	header     map[string]string
	tpl        PatternMR // 路由对应的methodsroute
	midware    []func(http.ResponseWriter, *http.Request) bool
	delmidware []func(http.ResponseWriter, *http.Request) bool
	pattern    map[string][]string // value 是 args， 如果长度是0， 那就是完全匹配， 大于0就是正则匹配
	delheader  []string
	groupKey   string
	groupTitle string
	groupLable string
	reqHeader  map[string]string
}

var reUrl map[string]*reroute

func NewGroupRoute() *GroupRoute {
	return &GroupRoute{
		route:   make(map[string]MethodsRoute),
		header:  make(map[string]string),
		tpl:     make(map[string]MethodsRoute),
		midware: make([]func(http.ResponseWriter, *http.Request) bool, 0),
	}
}

func (g *GroupRoute) ApiReqHeader(k, v string) *GroupRoute {
	// 接口的请求头
	if g.reqHeader == nil {
		g.reqHeader = make(map[string]string)
	}
	g.reqHeader[k] = v
	return g
}

func (g *GroupRoute) AddHeader(k, v string) *GroupRoute {

	if g.header == nil {
		g.header = make(map[string]string)
	}
	g.header[k] = v
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

func (g *GroupRoute) AddMidware(handle func(http.ResponseWriter, *http.Request) bool) *GroupRoute {

	if g.midware == nil {
		g.midware = make([]func(http.ResponseWriter, *http.Request) bool, 0)
	}
	g.midware = append(g.midware, handle)
	return g
}

func (g *GroupRoute) DelMidware(handle func(http.ResponseWriter, *http.Request) bool) *GroupRoute {

	if g.delmidware == nil {
		g.delmidware = make([]func(http.ResponseWriter, *http.Request) bool, 0)
	}
	g.delmidware = append(g.delmidware, handle)
	return g
}

func (g *GroupRoute) FirstMidware(handle func(http.ResponseWriter, *http.Request) bool) *GroupRoute {

	if g.midware == nil {
		g.midware = make([]func(http.ResponseWriter, *http.Request) bool, 0)
	}
	g.midware = append(g.midware, handle)
	return g
}

// 组里面也包括路由 后面的其实还是patter和handle
func (g *GroupRoute) Pattern(pattern string) MethodsRoute {

	// 格式路径
	pattern = slash(pattern)

	mr := make(map[string]*Route)

	if g.pattern == nil {
		g.pattern = make(map[string][]string)
	}

	if _, ok := g.pattern[pattern]; ok {
		log.Fatalf("pattern %s is Duplication", pattern)
	}

	if v, listvar := match(pattern); len(listvar) > 0 {
		if _, ok := g.pattern[v]; ok {
			log.Fatalf("Pattern Duplicate for %s", v)
		}

		if _, ok := g.pattern[v]; ok {
			log.Fatalf("pattern %s is Duplication", v)
		}
		g.tpl[v] = mr
		// 判断是否重复
		g.pattern[v] = listvar
		return mr
	}
	g.pattern[pattern] = make([]string, 0)
	g.route[pattern] = mr
	return mr
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

package xmux

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// 临时的， 最后会合并到route
type GroupRoute struct {
	// 感觉还没到method， 应该先uri后缀的
	route      mr
	name       string
	header     map[string]string
	tpl        mr
	midware    []func(http.ResponseWriter, *http.Request) bool
	delmidware []func(http.ResponseWriter, *http.Request) bool
	pattern    map[string]int
	delheader  []string
	groupKey   string
	groupTitle string
	groupLable string
	reqHeader  map[string]string
}

var reUrl map[string]*reroute

func NewGroupRoute(name string) *GroupRoute {
	return &GroupRoute{
		name:    name,
		route:   make(map[string]*Route),
		header:  make(map[string]string),
		tpl:     make(map[string]*Route),
		midware: make([]func(http.ResponseWriter, *http.Request) bool, 0),
		pattern: make(map[string]int),
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

func (g *GroupRoute) SetName(name string) *GroupRoute {
	if g.name == "" {
		g.name = time.Now().String()
	}
	if name != "" {
		g.name = name
	}

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
func (g *GroupRoute) Pattern(pattern string) *Route {
	if g.name == "" {
		g.name = time.Now().String()
	}
	if g.route == nil {
		g.route = make(map[string]*Route)
	}
	if g.tpl == nil {
		g.tpl = make(map[string]*Route)
	}
	if g.pattern == nil {
		g.pattern = make(map[string]int)
	}
	if g.pattern == nil {
		g.pattern = make(map[string]int)
	}
	// 格式路径
	pattern = slash(pattern)

	if _, ok := g.pattern[pattern]; ok {
		log.Fatalf("Pattern Duplicate for %s", pattern)
	}

	if pattern == "" || pattern[0:1] != "/" || strings.ContainsAny(pattern, " \t\n") {
		log.Fatalf("Pattern Error for %s", pattern)
	}
	route := &Route{
		method:     make(map[string]http.Handler),
		header:     make(map[string]string),
		args:       make([]string, 0),
		groupKey:   g.groupKey,
		groupLable: g.groupLable,
		groupTitle: g.groupTitle,
	}
	if v, listvar := match(pattern); len(listvar) > 0 {
		if _, ok := g.pattern[v]; ok {
			log.Fatalf("Pattern Duplicate for %s", v)
		}
		g.tpl[v] = route
		g.tpl[v].args = append(g.tpl[v].args, listvar...)
		g.pattern[v] = 3
		return g.tpl[v]
	}
	g.pattern[pattern] = 2
	g.route[pattern] = route
	return g.route[pattern]
}

// 组路由起的名字
func (r *Router) AddGroup(group *GroupRoute) *Router {
	if r.group == nil {
		r.group = make(map[string]*GroupRoute)
	}
	if group.name == "" {
		group.name = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	if r.header == nil {
		r.header = make(map[string]string)
	}
	if r.pattern == nil {
		r.pattern = make(map[string]int)
	}
	if r.tplpattern == nil {
		r.tplpattern = make(map[string]int)
	}
	if r.groupname == nil {
		r.groupname = make(map[string]string)
	}

	// 计算pattern 是否有重复的
	for k, v := range group.pattern {
		if _, ok := r.pattern[k]; ok {
			log.Fatalf("Pattern Duplicate for %s", k)
		}
		if v == 2 {
			r.pattern[k] = v
		} else {
			r.tplpattern[k] = v
		}
		r.groupname[k] = group.name
	}
	r.group[group.name] = group
	// group 的 组文档继承到 route
	for _, rt := range group.route {

		if group.groupKey != "" && rt.groupKey == "" {
			rt.groupKey = group.groupKey
		}

		if group.reqHeader != nil {
			if rt.reqHeader == nil {
				rt.reqHeader = make(map[string]string)
			}
			for k, v := range group.reqHeader {
				rt.reqHeader[k] = v
			}
		}
	}

	for _, rt := range group.tpl {
		if group.groupKey != "" && rt.groupKey == "" {
			rt.groupKey = group.groupKey
		}

		if group.reqHeader != nil {
			for k, v := range group.reqHeader {
				if rt.reqHeader == nil {
					rt.reqHeader = make(map[string]string)
				}
				rt.reqHeader[k] = v
			}
		}

	}
	return r
}

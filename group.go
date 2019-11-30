package xmux

import (
	"log"
	"net/http"
	"strings"
)

type GroupRoute struct {
	// 感觉还没到method， 应该先uri后缀的
	route  map[string]*Route
	name   string // 组名
	header map[string]string
	tpl    map[string]*Route
}

var reUrl map[string]*reroute

func NewGroupRoute() *GroupRoute {
	return &GroupRoute{
		route:  make(map[string]*Route),
		header: make(map[string]string),
		tpl:    make(map[string]*Route),
	}
}

func (g *GroupRoute) SetHeader(k, v string) *GroupRoute {
	g.header[k] = v
	return g
}

// 组里面也包括路由 后面的其实还是patter和handle
func (g *GroupRoute) Pattern(pattern string) *Route {

	// 格式路径
	pattern = slash(pattern)

	lv := make([]string, 0)
	route := &Route{
		method: make(map[string]http.Handler),
		header: make(map[string]string),
		args:   make([]string, 0),
	}
	if v, listvar, ok := match(pattern, "^", lv); ok {
		g.tpl[v] = route
		g.tpl[v].args = append(g.tpl[v].args, listvar...)
		return g.tpl[v]
	}
	g.route[pattern] = route
	return g.route[pattern]
}

func (r *Router) Group(name string) *GroupRoute {
	//   /article if /article/ to /article;  if article to /article
	name = strings.Trim(name, " ")
	g := &GroupRoute{
		name:   name,
		route:  make(map[string]*Route),
		header: make(map[string]string),
	}

	r.groupKey[name] = g.header
	return g
}

func (r *Router) AddGroup(groute *GroupRoute) *Router {
	for k, v := range groute.route {

		if _, ok := r.tpl[k]; ok {
			//路径检测
			log.Fatalf("pattern duplicate for %s", k)
		}
		r.groupKey[k] = groute.header
		r.route[k] = v
	}
	for k, v := range groute.tpl {
		if _, ok := r.tpl[k]; ok {
			//路径检测
			log.Fatalf("pattern duplicate for %s", k)
		}
		r.tpl[k] = v
		r.groupKey[k] = groute.header
	}

	return r
}

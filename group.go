package xmux

import (
	"log"
	"net/http"
)

type GroupRoute struct {
	// 感觉还没到method， 应该先uri后缀的
	suffix map[string]*Route
	prefix string // 组才有

}

var reUrl map[string]*reroute

func NewGroupRoute(pattern string) *GroupRoute {
	return &GroupRoute{
		suffix: make(map[string]*Route),
		prefix: pattern,
	}
}

// 组里面也包括路由 后面的其实还是patter和handle
func (g *GroupRoute) HandleFunc(pattern string) *Route {
	// name   if /name to name ; if name/shdk/ to name/shdk
	if pattern == "" {
		log.Fatal("pattern is error")
	}

	if pattern[0:1] == "/" {
		pattern = pattern[1:]
		if pattern == "" {
			log.Fatal("pattern is error")
		}
	}
	if pattern[len(pattern)-1:len(pattern)] == "/" {
		pattern = pattern[:len(pattern)-1]
	}
	route := &Route{
		method: make(map[string]http.Handler),
	}
	lv := make([]string, 0)
	pattern = g.prefix + "/" + pattern
	if v, listvar, ok := match(pattern, "^", lv); ok {
		reUrl[v] = &reroute{
			R:   route,
			Var: listvar,
		}
		return route
	}
	g.suffix[pattern] = route
	return route
}

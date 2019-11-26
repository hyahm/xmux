package xmux

import "net/http"

type GroupRoute struct {
	// 感觉还没到method， 应该先uri后缀的
	suffix map[string]*Route
	prefix string // 组才有
}

func NewGroupRoute(pattern string) *GroupRoute {
	return &GroupRoute{
		suffix: make(map[string]*Route),
		prefix: pattern,
	}
}

// 组里面也包括路由 后面的其实还是patter和handle
func (g *GroupRoute) HandleFunc(pattern string) *Route {
	// name   if /name to name ; if name/shdk/ to name/shdk
	if pattern[0:1] == "/" {
		pattern = pattern[1:]
	}
	if pattern[len(pattern)-1:len(pattern)] == "/" {
		pattern = pattern[:len(pattern)-1]
	}
	route := &Route{
		method: make(map[string]http.Handler),
	}
	g.suffix[pattern] = route
	return route
}

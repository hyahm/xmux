package xmux

import "net/http/pprof"

func (r *Router) Pprof() *GroupRoute {
	debug := NewGroupRoute()
	debug.Get("/debug/pprof/", pprof.Index)
	debug.Get("/debug/pprof/cmdline", pprof.Cmdline)
	debug.Get("/debug/pprof/profile", pprof.Profile)
	debug.Get("/debug/pprof/symbol", pprof.Symbol)
	debug.Get("/debug/pprof/trace", pprof.Trace)
	return debug
}

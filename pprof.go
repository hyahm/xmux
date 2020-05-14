package xmux

import "net/http/pprof"

func Pprof() *GroupRoute {
	debug := NewGroupRoute()
	debug.Pattern("/debug/pprof/").Get(pprof.Index)
	debug.Pattern("/debug/pprof/cmdline").Get(pprof.Cmdline)
	debug.Pattern("/debug/pprof/profile").Get(pprof.Profile)
	debug.Pattern("/debug/pprof/symbol").Get(pprof.Symbol)
	debug.Pattern("/debug/pprof/trace").Get(pprof.Trace)
	return debug
}

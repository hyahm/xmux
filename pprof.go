package xmux

import (
	"net/http/pprof"
)

func Pprof() *RouteGroup {
	pp := NewRouteGroup().BindResponse(nil)

	pp.Get("/debug/pprof/{all:name}", pprof.Index)
	pp.Get("/debug/pprof/cmdline", pprof.Cmdline)
	pp.Get("/debug/pprof/profile", pprof.Profile)
	pp.Get("/debug/pprof/symbol", pprof.Symbol)
	pp.Get("/debug/pprof/trace", pprof.Trace)
	return pp
}

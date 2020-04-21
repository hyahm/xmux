package xmux

import "golang.org/x/net/context"

var Ctx map[string]context.Context

func init() {
	Ctx = make(map[string]context.Context)
}

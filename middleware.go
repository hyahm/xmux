package xmux

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

// 洋葱容器
type onion struct {
	mws []Middleware
}

// 新建一颗洋葱
func New(mws ...Middleware) *onion {
	return &onion{mws: mws}
}

// 追加一层（最外层）
func (o *onion) Use(mw ...Middleware) *onion {
	o.mws = append(o.mws, mw...)
	return o
}

// 把最后一棒 Handler 包成 http.Handler
func (o *onion) Then(final http.Handler) http.Handler {
	// 从后往前包，保证执行顺序是洋葱圈
	for i := len(o.mws) - 1; i >= 0; i-- {
		final = o.mws[i](final)
	}
	return final
}

// 方便直接传函数
func (o *onion) ThenFunc(final http.Handler) http.Handler {
	return o.Then(final)
}

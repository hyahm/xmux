package xmux

import (
	"context"
	"net/http"
)

// UniversalFlow 通用流程基类：只做编排、不做业务
type BaseFlow struct {
	W   http.ResponseWriter
	R   *http.Request
	Ctx context.Context
	Ins *FlowData // xmux 内置实例（Get/Set 数据）
	err error
}

func (b *BaseFlow) Init(w http.ResponseWriter, r *http.Request) {
	b.W = w
	b.R = r
	b.Ctx = r.Context()
	b.Ins = GetInstance(r)
}

// 链式执行，出错 = 中断流程
func (b *BaseFlow) Chain(fns ...func() error) {
	if b.err != nil {
		return
	}
	for _, fn := range fns {
		if err := fn(); err != nil {
			b.err = err
			return
		}
	}
}

func (b *BaseFlow) Err() error { return b.err }

type Flow interface {
	Init(w http.ResponseWriter, r *http.Request)
	Run()
}

func Adapt(newFlow func() Flow) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// 每个请求：新建结构体实例 → 最强隔离
		f := newFlow()

		// 初始化
		f.Init(w, r)

		// 执行业务链式流程
		f.Run()

	}
}

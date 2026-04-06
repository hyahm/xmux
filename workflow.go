package xmux

import (
	"context"
	"net/http"
)

// ====================== 增强版流程基类：支持 顺序 + 分支 ======================
type BaseFlow struct {
	W   http.ResponseWriter
	R   *http.Request
	Ctx context.Context
	Ins *FlowData
	err error
}

func (b *BaseFlow) Init(w http.ResponseWriter, r *http.Request) {
	b.W = w
	b.R = r
	b.Ctx = r.Context()
	b.Ins = GetInstance(r)
	b.err = nil
}

// ---------------------- 核心能力 ----------------------
// 1. 顺序执行（原来的）
func (b *BaseFlow) Then(fns ...func() error) *BaseFlow {
	if b.err != nil {
		return b
	}
	for _, fn := range fns {
		if err := fn(); err != nil {
			b.err = err
			return b
		}
	}
	return b
}

// 2. 条件分支：满足条件才执行（解决你的判断/分支需求）
func (b *BaseFlow) If(cond bool, fns ...func() error) *BaseFlow {
	if b.err != nil || !cond {
		return b
	}
	return b.Then(fns...)
}

// 3. 二选一分支：if / else
func (b *BaseFlow) IfElse(cond bool, do func() error, elseDo func() error) *BaseFlow {
	if b.err != nil {
		return b
	}
	if cond {
		return b.Then(do)
	}
	return b.Then(elseDo)
}

// 4. 任意条件跳过（满足则跳过本段）
func (b *BaseFlow) SkipIf(cond bool, fns ...func() error) *BaseFlow {
	if b.err != nil || cond {
		return b
	}
	return b.Then(fns...)
}

// ---------------------- 错误 ----------------------
func (b *BaseFlow) Err() error { return b.err }

// ====================== 接口与适配器（不变） ======================
type Flow interface {
	Init(w http.ResponseWriter, r *http.Request)
	Run()
}

func Adapt(newFlow func() Flow) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f := newFlow()
		f.Init(w, r)
		f.Run()
	}
}

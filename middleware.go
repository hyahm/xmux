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

// Logger 记录方法、路径、状态、耗时
// func Logger() Middleware {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			start := time.Now()
// 			next.ServeHTTP(w, r)
// 			log.Printf("[%d] %s %s ----------------------------------- (%v)", 200, r.Method, r.URL.Path, time.Since(start))
// 		})
// 	}
// }

// Recovery 捕获 panic
// func Recovery() Middleware {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			defer func() {
// 				if err := recover(); err != nil {
// 					log.Printf("panic: %v", err)
// 					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 				}
// 			}()
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

// // 1. 一个请求组，负责合并等待
// type requestGroup struct {
// 	done chan []byte // 关闭即代表业务完成
// 	// val      []byte        // 业务结果
// 	connects int
// 	callback chan struct{}
// }

// // 2. 全局中心：相同 key 复用同一个 group

// type combineHandlers struct {
// 	mu   sync.Mutex
// 	inFL map[string]*requestGroup
// }

// var requestCoalescing = &combineHandlers{
// 	inFL: make(map[string]*requestGroup),
// 	mu:   sync.Mutex{},
// }

// func getCombineHandlers(key string) (*requestGroup, bool) {
// 	requestCoalescing.mu.Lock()
// 	defer requestCoalescing.mu.Unlock()
// 	g, ok := requestCoalescing.inFL[key]
// 	return g, ok
// }

// func setCombineHandlers(key string) {
// 	requestCoalescing.mu.Lock()
// 	if g, ok := requestCoalescing.inFL[key]; ok {
// 		g.connects++
// 	}
// 	requestCoalescing.mu.Unlock()

// }

// func delCombineHandlers(key string) {
// 	// 外面已经加锁了， 再加就死锁了
// 	delete(requestCoalescing.inFL, key)
// }

// func initCombineHandlers(key string) *requestGroup {
// 	requestCoalescing.mu.Lock()
// 	g := &requestGroup{
// 		done:     make(chan []byte, 100), // 足够大即可
// 		connects: 0,
// 		callback: make(chan struct{}, 100),
// 	}
// 	requestCoalescing.inFL[key] = g
// 	requestCoalescing.mu.Unlock()
// 	return g
// }

// func OptimizerModule(w http.ResponseWriter, r *http.Request) bool {
// 	key := r.URL.String()
// 	g, ok := getCombineHandlers(key)
// 	if ok {
// 		// 已有进行中的请求，直接等待
// 		setCombineHandlers(key)
// 		rb := <-g.done
// 		g.callback <- struct{}{}
// 		w.Write(rb)
// 		return true
// 	} else {
// 		// 第一个请求，建组并启动
// 		g := initCombineHandlers(key)
// 		GetInstance(r).Set("xmux_comb", g)
// 		return false
// 	}

// }

// const RESPONSEBYTES = "xmux_reponsebytes"

// func CombineHandlers() Middleware {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			key := r.URL.String()
// 			if g, ok := getCombineHandlers(key); ok {
// 				// 已有进行中的请求，直接等待
// 				setCombineHandlers(key)
// 				rb := <-g.done
// 				g.callback <- struct{}{}
// 				w.Write(rb)
// 				return
// 			}

// 			// 第一个请求，建组并启动
// 			g := initCombineHandlers(key)
// 			next.ServeHTTP(w, r)
// 			requestCoalescing.mu.Lock()
// 			defer requestCoalescing.mu.Unlock()
// 			n := g.connects
// 			if n > 0 {
// 				requestCoalescing.mu.Lock()
// 				responseBody := GetInstance(r).Get(RESPONSEBYTES).([]byte)
// 				wg := &sync.WaitGroup{}

// 				// 1. 锁内：把当前等待者全部“预置”到通道里
// 				go func() {
// 					// 2. 锁外：等他们全部拿走（可选，如果不需要回执可删掉）
// 					for i := 0; i < n; i++ {
// 						wg.Go(func() {
// 							<-g.callback
// 						})
// 					}
// 				}()

// 				for i := 0; i < n; i++ {
// 					g.done <- responseBody
// 				}
// 				wg.Wait()
// 				delCombineHandlers(key)

// 				w.Write(responseBody)
// 			} else {
// 				delCombineHandlers(key)
// 			}
// 		})

// 	}
// }

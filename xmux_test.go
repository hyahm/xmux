package xmux

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"testing"
)

func home(w http.ResponseWriter, r *http.Request) {
	// name := Var(r)["name"]
	// time.Sleep(time.Millisecond * 30)
	// fmt.Println(1111)
	// GetInstance(r).Set(RESPONSEBYTES, []byte("asjdlfjalsdf"))

}

func grouphome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("mmmmmmmm group")
	w.Write([]byte("grouphome" + Var(r)["name"] + "-" + Var(r)["age"]))
}

func adminhandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("mmmmmmmm admin")
}

func adminGroup() *RouteGroup {
	admin := NewRouteGroup().Prefix("test")
	admin.Get("/admin/{bbb}", home)
	admin.Get("/admin", adminhandle)
	admin.Get("/aaa/adf{re:([a-z]{1,4})sf([0-9]{0,10})sd: name, age}", grouphome)
	return admin
}

func userGroup() *RouteGroup {
	user := NewRouteGroup()
	// user.Get("/group", home).Use(CombineHandlers())
	user.Get("/group", home)
	user.AddGroup(adminGroup()).DelPostModule(postModule)
	return user
}

func TestMain(t *testing.T) {
	// pool := NewPool()
	router := NewRouter()
	// router.HandleAll = LimitFixedWindowCounterTemplate
	router.HandleRecover = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("服务器错误"))
	}
	router.SetHeader("Access-Control-Allow-Origin", "*")
	router.SetHeader("Content-Type", "application/x-www-form-urlencoded,application/json; charset=UTF-8")
	router.SetHeader("Access-Control-Allow-Headers", "Content-Type")
	router.SetHeader("Access-Control-Max-Age", "1728000")
	// router.SetHeader("Access-Control-Allow-Origin", "*").
	// 	SetHeader("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	router.AddGroup(Pprof())
	router.Enter = enter
	// router.Prefix("/api")
	// router.EnableConnect = true
	router.Get("/test", pp)
	// router.Get("/post", pp).Use(pool.Middleware(heavyHandler))
	router.HandleAll = nil
	// router.SetAddr(":8080")
	// router.AddGroup(userGroup())
	log.Fatal(router.SetAddr(":9999").Run())
}

type Binding struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func Recovery(key string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("panic: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			key := r.URL.Path
			e := getEntry(key)

			mu.Lock()
			//---------- 生产者路径 ----------
			if !e.done {
				e.count++   // 正在计算
				mu.Unlock() // 放锁去 IO

				// 模拟耗时计算
				// resp := []byte("这是 " + key + " 的处理结果")
				next.ServeHTTP(w, r)
				mu.Lock()
				e.response = GetInstance(r).Get(key).([]byte)
				e.done = true
				e.cond.Broadcast()
				mu.Unlock()
				return
			}
			//---------- 消费者路径 ----------
			e.count++
			for !e.done {
				e.cond.Wait() // 内部会临时放锁
			}
			data := append([]byte(nil), e.response...) // 锁内深拷贝
			e.count--
			if e.count == 0 { // 最后一个离开
				delete(entries, key) // 安全删除
			}
			mu.Unlock()

			w.Write(data)

		})
	}
}

// Middleware 返回 http.Handler 中间件
// opts: 可传入 KeyFunc，默认用 r.URL.Path
func (p *Pool) Middleware(next http.HandlerFunc, opts ...KeyFunc) http.HandlerFunc {
	keyFn := func(r *http.Request) string { return r.URL.Path }
	if len(opts) > 0 && opts[0] != nil {
		keyFn = opts[0]
	}

	return func(w http.ResponseWriter, r *http.Request) {
		key := keyFn(r)
		e := p.getEntry(key)

		p.mu.Lock()
		// 生产者路径
		if !e.done {
			e.count++
			p.mu.Unlock()

			// 执行业务 handler 拿结果
			rec := &responseRecorder{ResponseWriter: w, status: 200}
			next(rec, r)

			p.mu.Lock()
			e.response = rec.body
			e.done = true
			e.cond.Broadcast()
			p.mu.Unlock()
			return
		}

		// 消费者路径
		e.count++
		for !e.done {
			e.cond.Wait()
		}
		data := append([]byte(nil), e.response...)
		e.count--
		if e.count == 0 {
			delete(p.data, key)
		}
		p.mu.Unlock()

		// 把缓存内容写回客户端
		w.Write(data)
	}
}

// 内部拿 entry（锁内）
func (p *Pool) getEntry(key string) *entry {
	p.mu.Lock()
	defer p.mu.Unlock()
	e := p.data[key]
	if e == nil {
		e = &entry{cond: sync.NewCond(&p.mu)}
		p.data[key] = e
	}
	return e
}

// KeyFunc 允许调用方自定义分组 key
type KeyFunc func(r *http.Request) string

// responseRecorder 把响应内容截下来复用
type responseRecorder struct {
	http.ResponseWriter
	body   []byte
	status int
}

func (r *responseRecorder) Write(p []byte) (int, error) {
	r.body = append(r.body, p...)
	return len(p), nil
}

func postModule(w http.ResponseWriter, r *http.Request) bool {

	fmt.Println("这是一个后置模块")
	return false
}

// Pool 按 key 缓存请求结果
type Pool struct {
	mu   sync.Mutex
	data map[string]*entry
}

type entry struct {
	done     bool
	response []byte
	count    int32
	cond     *sync.Cond
}

func NewPool() *Pool {
	return &Pool{data: make(map[string]*entry)}
}

var (
	mu      sync.Mutex // 只用一把 Lock，简化生命周期
	entries = make(map[string]*entry)
)

// 拿到或新建 entry（锁内）
func getEntry(key string) *entry {
	mu.Lock()
	defer mu.Unlock()
	e := entries[key]
	if e == nil {
		e = &entry{cond: sync.NewCond(&mu)}
		entries[key] = e
	}
	return e
}

func pp(w http.ResponseWriter, r *http.Request) {

	a := make([]string, 0)
	fmt.Println(a[8])
	key := r.URL.Path
	e := getEntry(key)

	mu.Lock()
	//---------- 生产者路径 ----------
	if !e.done {
		e.count++   // 正在计算
		mu.Unlock() // 放锁去 IO

		// 模拟耗时计算
		resp := []byte("这是 " + key + " 的处理结果")

		mu.Lock()
		e.response = resp
		e.done = true
		e.cond.Broadcast()
		mu.Unlock()
		return
	}
	//---------- 消费者路径 ----------
	e.count++
	for !e.done {
		e.cond.Wait() // 内部会临时放锁
	}
	data := append([]byte(nil), e.response...) // 锁内深拷贝
	e.count--
	if e.count == 0 { // 最后一个离开
		delete(entries, key) // 安全删除
	}
	mu.Unlock()

	w.Write(data)
}

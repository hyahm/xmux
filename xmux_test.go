package xmux

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var count int32

func home(w http.ResponseWriter, r *http.Request) {
	// name := Var(r)["name"]
	time.Sleep(time.Millisecond * 10)
	g := GetInstance(r).Get("xmux_comb").(*requestGroup)
	responseBody := []byte("home admin")
	fmt.Println("home admin")
	key := r.URL.String()
	atomic.AddInt32(&count, 1)
	fmt.Println(atomic.LoadInt32(&count))
	// 1. 锁内：把当前等待者全部“预置”到通道里
	requestCoalescing.mu.Lock()
	n := g.connects
	for i := 0; i < n; i++ {
		g.done <- responseBody
	}
	delCombineHandlers(key)
	requestCoalescing.mu.Unlock()

	// 2. 锁外：等他们全部拿走（可选，如果不需要回执可删掉）
	for i := 0; i < n; i++ {
		<-g.callback
	}

	w.Write(responseBody)
	// go func(g *requestGroup, reponseBody []byte) {
	// 真正执行业务
	// requestCoalescing.mu.Lock()

	// defer requestCoalescing.mu.Unlock()

	// if g.connects > 0 {
	// 	all := make(chan struct{})
	// 	go func() {
	// 		for range g.connects {
	// 			<-g.callback
	// 		}
	// 		fmt.Println("处理完一轮")
	// 		all <- struct{}{}
	// 	}()
	// 	// 直接锁定， 新的请求等待解锁再进行
	// 	for range g.connects {
	// 		g.done <- reponseBody
	// 	}
	// 	<-all
	// 	// close(g.done) // 广播完成
	// 	fmt.Println("g.connects:", g.connects)

	// }
	// delCombineHandlers(key)
	// // }(g, reponseBody)

	// w.Write(responseBody)
	// 把 key 从 map 移除，后续请求可再次触发新业务
	// w.Write([]byte("home admin"))
	// time.Sleep(1 * time.Second)
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
	admin.Get("/admin", adminhandle).DelPostModule()
	admin.Get("/aaa/adf{re:([a-z]{1,4})sf([0-9]{0,10})sd: name, age}", grouphome)
	return admin
}

func userGroup() *RouteGroup {
	user := NewRouteGroup().Prefix("test")
	user.Get("/group", grouphome)
	user.AddGroup(adminGroup()).DelPostModule(postModule)
	return user
}

func TestMain(t *testing.T) {

	router := NewRouter()
	// router.HandleAll = LimitFixedWindowCounterTemplate
	router.SetHeader("Access-Control-Allow-Origin", "*")
	router.SetHeader("Content-Type", "application/x-www-form-urlencoded,application/json; charset=UTF-8")
	router.SetHeader("Access-Control-Allow-Headers", "Content-Type")
	router.SetHeader("Access-Control-Max-Age", "1728000")
	// router.SetHeader("Access-Control-Allow-Origin", "*").
	// 	SetHeader("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	router.AddGroup(Pprof()).AddPostModule(postModule)

	// router.Prefix("/api")
	// router.EnableConnect = true
	router.Get("/test", nil)
	router.Get("/post", pp)
	router.HandleAll = nil
	// router.SetAddr(":8080")
	router.AddGroup(userGroup())
	log.Fatal(router.SetAddr(":9999").Run())
}

func postModule(w http.ResponseWriter, r *http.Request) bool {
	fmt.Println("这是一个后置模块")
	return false
}

func pp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("这是一个处理函数")
}

// 1. 一个请求组，负责合并等待
type requestGroup struct {
	done chan []byte // 关闭即代表业务完成
	// val      []byte        // 业务结果
	connects int
	callback chan struct{}
}

// 2. 全局中心：相同 key 复用同一个 group

type combineHandlers struct {
	mu   sync.Mutex
	inFL map[string]*requestGroup
}

var requestCoalescing = &combineHandlers{
	inFL: make(map[string]*requestGroup),
	mu:   sync.Mutex{},
}

func getCombineHandlers(key string) (*requestGroup, bool) {
	requestCoalescing.mu.Lock()
	defer requestCoalescing.mu.Unlock()
	g, ok := requestCoalescing.inFL[key]
	return g, ok
}

func setCombineHandlers(key string) {
	requestCoalescing.mu.Lock()
	if g, ok := requestCoalescing.inFL[key]; ok {
		g.connects++
	}
	requestCoalescing.mu.Unlock()

}

func delCombineHandlers(key string) {
	// 外面已经加锁了， 再加就死锁了
	delete(requestCoalescing.inFL, key)
}

func initCombineHandlers(key string) *requestGroup {
	requestCoalescing.mu.Lock()
	g := &requestGroup{
		done:     make(chan []byte, 100), // 足够大即可
		connects: 0,
		callback: make(chan struct{}, 100),
	}
	requestCoalescing.inFL[key] = g
	requestCoalescing.mu.Unlock()
	return g
}

func OptimizerModule(w http.ResponseWriter, r *http.Request) bool {
	key := r.URL.String()
	g, ok := getCombineHandlers(key)
	if ok {
		// 已有进行中的请求，直接等待
		setCombineHandlers(key)
		rb := <-g.done
		g.callback <- struct{}{}
		w.Write(rb)
		return true
	} else {
		// 第一个请求，建组并启动
		g := initCombineHandlers(key)
		GetInstance(r).Set("xmux_comb", g)
		return false
	}

}

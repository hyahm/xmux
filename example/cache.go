package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hyahm/gocache"
	"github.com/hyahm/xmux"
)

func c(w http.ResponseWriter, r *http.Request) {
	fmt.Println("comming c")
	now := time.Now().String()
	xmux.GetInstance(r).Response.(*Response).Data = now
}

func noCache(w http.ResponseWriter, r *http.Request) {
	fmt.Println("update c")
	xmux.NeedUpdate("/aaa")
}

func noCache1(w http.ResponseWriter, r *http.Request) {
	fmt.Println("comming noCache1")
	now := time.Now().String()
	xmux.GetInstance(r).Response.(*Response).Data = now
}

func setKey(w http.ResponseWriter, r *http.Request) bool {
	xmux.GetInstance(r).SetCacheKey(r.URL.Path)
	fmt.Println(r.URL.Path + "    is cachedaaa")
	return false
}

type Response struct {
	Code int `json:"code"`

	Data interface{} `json:"data"`
}

func main() {

	r := &Response{
		Code: 0,
	}
	cth := gocache.NewCache[string, []byte](100, gocache.LFU)
	xmux.InitResponseCache(cth)

	router := xmux.NewRouter().AddModule(setKey, xmux.DefaultCacheTemplateCacheWithResponse).AddModule(xmux.DefaultCacheTemplateCacheWithResponse)
	router.BindResponse(r).Prefix("test")
	router.Get("/bbb", c)
	router.Get("/ccc", c).DelPrefix("test")
	router.IgnoreSlash = true
	g := xmux.NewRouteGroup().Prefix().DelPrefix()
	g.Get("/aaa", noCache).DelModule(setKey)
	g.Get("/no/cache1", noCache1).DelModule(setKey).DelPrefix("test")
	// router.Request("/aaa", noCache, http.MethodGet, http.MethodPost)
	router.AddGroup(g)
	router.DebugRoute()
	router.SetAddr(":7777")
	router.Run()
}

// func debug(paths ...string) {
// RELOAD:
// 	for _, path := range paths {
// 		fi, err := os.Stat(path)
// 		if err != nil {
// 			fmt.Println(err)
// 			continue
// 		}
// 		if fi.IsDir() {
// 		}
// 	}
// 	watch := make(map[string]*os.File)
// 	cache, err := os.Open("example\\cache.go")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fi, err := cache.Stat()
// 	if err != nil {
// 		panic(err)
// 	}
// 	this := fi.ModTime().Unix()

// 	exit := make(chan os.Signal)
// 	signal.Notify(exit, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
// 	reload := make(chan bool)

// 	ctx, cancel := context.WithCancel(context.Background())
// 	go func() {
// 		ticker := time.NewTicker(time.Second)
// 		for {
// 			select {
// 			case <-ticker.C:
// 				fi1, err := cache.Stat()
// 				if err != nil {
// 					panic(err)
// 				}
// 				if this != fi1.ModTime().Unix() {
// 					reload <- true
// 					return
// 				}

// 			}
// 		}

// 	}()
// 	r := &Response{
// 		Code: 0,
// 	}
// 	cth := gocache.NewCache[string, []byte](100, gocache.LFU)
// 	xmux.InitResponseCache(cth)
// 	router := xmux.NewRouter().AddModule(setKey, xmux.DefaultCacheTemplateCacheWithResponse)
// 	router.BindResponse(r)
// 	router.Get("/aaa", c)

// 	router.Get("/update/aaa", noCache).DelModule(setKey)
// 	router.Get("/no/cache1", noCache1).DelModule(setKey)
// 	router.AddGroup(xmux.Pprof().DelModule(setKey))
// 	go router.Debug(ctx)
// 	select {
// 	case <-reload:
// 		cancel()
// 		goto RELOAD
// 	case <-exit:
// 		fmt.Println("exit")
// 	}
// }

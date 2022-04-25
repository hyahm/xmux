package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hyahm/xmux"
	"github.com/hyahm/xmux/cache"
)

func c(w http.ResponseWriter, r *http.Request) {
	fmt.Println("comming c")

	now := time.Now().String()
	xmux.GetInstance(r).Response.(*Response).Data = now
}

func noCache(w http.ResponseWriter, r *http.Request) {
	fmt.Println("update c")
	cache.NeedUpdate("/aaa")
}

func noCache1(w http.ResponseWriter, r *http.Request) {
	fmt.Println("comming noCache1")
	now := time.Now().String()
	xmux.GetInstance(r).Response.(*Response).Data = now
}

func setKey(w http.ResponseWriter, r *http.Request) bool {
	xmux.GetInstance(r).Set(xmux.CacheKey, r.URL.Path)
	fmt.Print(r.URL.Path + " is cached")
	return false
}

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func main() {
	r := &Response{
		Code: 0,
	}
	cache.InitResponseCache()
	router := xmux.NewRouter().AddModule(setKey, xmux.DefaultCacheTemplateCacheWithResponse) // 设置所有路由都缓存
	router.BindResponse(r)
	router.Get("/aaa", c)                                // 缓存了
	router.Get("/update/aaa", noCache).DelModule(setKey) // 更新/aaa缓存
	router.Get("/no/cache1", noCache1).DelModule(setKey) // 没缓存
	router.Run()
}

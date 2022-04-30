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
	cth := gocache.NewCache[string, []byte](100, gocache.LFU)
	xmux.InitResponseCache(cth)
	router := xmux.NewRouter().AddModule(setKey, xmux.DefaultCacheTemplateCacheWithResponse)
	// router.Cache = cache.NewCache(10000, cache.LRU)
	router.BindResponse(r)
	router.Get("/aaa", c)
	router.Get("/update/aaa", noCache).DelModule(setKey)
	router.Get("/no/cache1", noCache1).DelModule(setKey)
	router.Run()
}

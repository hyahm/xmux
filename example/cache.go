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

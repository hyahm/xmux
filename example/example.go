package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/hyahm/xmux"
)

func home1(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Connection", "Close")
	xmux.GetInstance(r).Set("aaaa", "bbb")
	fmt.Println("home1")
	return false
}

func home(w http.ResponseWriter, r *http.Request) {
	xmux.GetInstance(r).Set("aaaa", "bbb")
	fmt.Println("home")
}

func v1Group() *xmux.GroupRoute {
	global := &Global{
		Code: 100,
		Msg:  "ok",
	}
	v1 := xmux.NewGroupRoute().BindResponse(global)
	v1.Get("/v1/login", home)
	// v1.Get("/v1/22", home)
	return v1
}

// func v2Group() *xmux.GroupRoute {
// 	v2 := xmux.NewGroupRoute().DelModule(home1)
// 	v2.Get("/v2/login", home)
// 	v2.Get("/v2/22", home)
// 	v2.AddGroup(v1Group())
// 	return v2
// }

type Global struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func midware() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		return
	})
}

func main() {
	global := &Global{
		Code: 200,
		Msg:  "ok",
	}
	router := xmux.NewRouter().AddModule(xmux.DefaultPermissionTemplate, home1).BindResponse(global)
	router.AddGroup(v1Group())
	router.Run(":8888")
}

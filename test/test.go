package main

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/hyahm/xmux"
)

func home1(w http.ResponseWriter, r *http.Request) bool {
	return false
}

func home(w http.ResponseWriter, r *http.Request) {
	xmux.GetInstance(r).Set("aaaa", "bbb")
}

type Global struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type User struct {
	UserName string `json:"username,require" form:"username,require"`
	PassWord string `json:"password" form:"password"`
	Gender   bool   `json:"form" form:"gender"`
}

func subgroup() *xmux.GroupRoute {
	sub := xmux.NewGroupRoute()
	sub.Get("/sub/get", home)
	sub.Post("/sub/post", home)
	sub.Any("/sub/any", home)
	sub.AddGroup(sub1group())
	return sub
}

func sub1group() *xmux.GroupRoute {
	sub1 := xmux.NewGroupRoute()
	sub1.Get("/sub1/get", home)
	sub1.Post("/sub1/post", home)
	sub1.Any("/sub1/any", home)
	return sub1
}

func main() {
	// global := &Global{
	// 	Code: 200,
	// }

	router := xmux.NewRouter().AddModule(home1)
	group := xmux.NewGroupRoute()
	group.Post("/post", home).AddModule()
	group.AddGroup(subgroup())
	router.Get("/get", home)
	router.Post("/", home).DelModule(home1)
	router.DebugAssignRoute("/")
	router.Run(":8888")
}

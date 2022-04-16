package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/hyahm/xmux"
)

func home1(w http.ResponseWriter, r *http.Request) bool {
	return false
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header.Get("Content-Length"))
	fmt.Printf("%T\n", r.Body)
}

type Global struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type User struct {
	UserName string `json:"username" form:"username"`
	PassWord string `json:"password" form:"password"`
	Gender   bool   `json:"form" form:"gender"`
}

func subgroup() *xmux.GroupRoute {
	sub := xmux.NewGroupRoute()
	sub.Get("/sub/get", home)
	sub.Post("/sub/post", home)
	sub.Any("/sub/any", home)
	sub.Any("/get", home)
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
	router := xmux.NewRouter().AddModule(home1)
	router.AddGroup(subgroup())
	router.PrintRequestStr = true
	router.Post("/get", home)
	router.Post("/", home).DelModule(home1).BindForm(User{})
	router.DebugAssignRoute("/")
	router.Run(":8888")
}

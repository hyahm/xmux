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

func homeargs(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header.Get("Content-Length"))
		fmt.Println("name", name)
	}
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
	sub.Any("/any/get", home)
	sub.AddGroup(sub1group())
	return sub
}

func sub1group() *xmux.GroupRoute {
	sub1 := xmux.NewGroupRoute()
	sub1.Get("/sub1/get", home)
	sub1.Post("/sub1/post", home)
	sub1.Any("/sub1/any", homeargs("sub1"))
	return sub1
}

func main() {
	g := &Global{
		Code: 200,
	}
	router := xmux.NewRouter().AddModule(home1).BindResponse(g)
	// router.SetHeader("Access-Control-Allow-Origin", "*")
	// router.SetHeader("Access-Control-Allow-Methods", "*")
	router.AddGroup(subgroup())
	router.Post("/get", home)
	router.Post("/", home).DelModule(home1).BindForm(User{})
	router.AddGroup(router.ShowSwagger("/docs", "127.0.0.1:8888"))

	router.Run(":8888")
}

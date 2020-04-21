package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hyahm/xmux"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println(xmux.Var[r.URL.Path]["test"])
	w.Write([]byte("hello world home"))
	return
}

func name(w http.ResponseWriter, r *http.Request) {
	fmt.Println(xmux.Var[r.URL.Path]["name"])
	fmt.Println(xmux.Ctx[r.URL.Path].Value("conf"))
	w.Write([]byte("hello world name"))
	return
}

func me(w http.ResponseWriter, r *http.Request) {
	fmt.Println(xmux.Var[r.URL.Path]["me"])
	fmt.Println(xmux.Ctx[r.URL.Path].Value("conf"))
	w.Write([]byte("hello world me"))
	return
}

func all(w http.ResponseWriter, r *http.Request) {
	fmt.Println(xmux.Var[r.URL.Path]["all"])
	fmt.Println(xmux.Var[r.URL.Path]["oid"])
	w.Write([]byte("hello world all"))
	return
}

func login(w http.ResponseWriter, r *http.Request) bool {
	fmt.Println("login mw")

	fmt.Println(r.Header.Get("bbb"))
	return false
}

func filter(w http.ResponseWriter, r *http.Request) bool {
	fmt.Println("login mw")
	r.Header.Set("bbb", "ccc")

	xmux.Ctx[r.URL.Path] = context.WithValue(context.Background(), "conf", "body")
	return false
}

func main() {
	router := xmux.NewRouter()
	router.IgnoreIco = false

	fmt.Println(router.Slash)

	router.Pattern("/home").Get(home)
	router.Pattern("/aaa/{name}").Get(name).AddMidware(filter).AddMidware(login)
	router.Pattern("/aaa/bbbb/{path:me}").Get(me)
	router.Pattern("/bbb/ccc/{int:oid}/{string:all}").Get(all)
	if err := http.ListenAndServe(":9000", router); err != nil {
		log.Fatal(err)
	}

}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hyahm/xmux"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world home"))
	return
}

func name(w http.ResponseWriter, r *http.Request) {
	fmt.Println(xmux.Var[r.URL.Path]["name"])
	w.Write([]byte("hello world name"))
	return
}

func me(w http.ResponseWriter, r *http.Request) {
	fmt.Println(xmux.Var[r.URL.Path]["me"])
	w.Write([]byte("hello world me"))
	return
}

func all(w http.ResponseWriter, r *http.Request) {
	fmt.Println(xmux.Var[r.URL.Path]["all"])
	fmt.Println(xmux.Var[r.URL.Path]["oid"])
	w.Write([]byte("hello world all"))
	return
}

func main() {
	router := xmux.NewRouter()
	router.Slash = true
	router.IgnoreIco = false
	fmt.Println(router.Slash)
	router.Pattern("/home").Get(home)
	router.Pattern("/aaa/{name}").Get(name)
	router.Pattern("/aaa/bbbb/{path:me}").Get(me)
	router.Pattern("/bbb/ccc/{int:oid}/{string:all}").Get(all)
	if err := http.ListenAndServe(":9000", router); err != nil {
		log.Fatal(err)
	}

}

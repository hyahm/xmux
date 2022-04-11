package main

import (
	"fmt"
	"net/http"

	"github.com/hyahm/xmux"
)

func home1(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Connection", "Close")
	xmux.GetInstance(r).Set("aaaa", "bbb")
	return false
}

func main() {
	router := xmux.NewRouter().AddModule(xmux.DefaultPermissionTemplate)
	router.Get("/aaa/adf{re:([a-z]{1,4})sf([0-9]{0,10})sd:name,age}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Connection", "Close")
		fmt.Println(xmux.GetInstance(r).Get("aaaa"))
		fmt.Println(xmux.GetConnents())
		xmux.GetInstance(r).Append(xmux.RESPONSEBODY, []byte("ok"))
		w.Write([]byte("ok"))
	}).AddModule(home1)

	router.Run(":8888")
}

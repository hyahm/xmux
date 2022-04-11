package main

import (
	"fmt"
	"net/http"

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
}

func v1Group() *xmux.GroupRoute {
	v1 := xmux.NewGroupRoute()
	v1.Get("/v1/login", home)
	v1.Get("/v1/22", home)
	return v1
}

func v2Group() *xmux.GroupRoute {
	v2 := xmux.NewGroupRoute().DelModule(home1)
	v2.Get("/v2/login", home)
	v2.Get("/v2/22", home)
	v2.AddGroup(v1Group())
	return v2
}

func main() {
	router := xmux.NewRouter().AddModule(xmux.DefaultPermissionTemplate, home1)
	router.Get("/aaa/adf{re:([a-z]{1,4})sf([0-9]{0,10})sd:name,age}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Connection", "Close")
		fmt.Println(xmux.GetInstance(r).Get("aaaa"))
		fmt.Println(xmux.GetConnents())
		xmux.GetInstance(r).Append(xmux.RESPONSEBODY, []byte("ok"))
		w.Write([]byte("ok"))
	})
	router.AddGroup(v2Group())
	router.DebugAssignRoute("/v1/login")
	router.Run(":8888")
}

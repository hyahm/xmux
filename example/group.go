package main

import (
	"fmt"
	"net/http"

	"github.com/hyahm/xmux"
)

func test(w http.ResponseWriter, r *http.Request) bool {
	fmt.Println("test")
	return false
}

func Admin() *xmux.GroupRoute {
	admin := xmux.NewGroupRoute()

	admin.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("admin"))
	})
	return admin
}

func Home() *xmux.GroupRoute {
	home := xmux.NewGroupRoute().AddModule(test)
	home.Get("/home", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("home"))
	})
	home.AddGroup(Admin())
	return home
}

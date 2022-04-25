package main

import "github.com/hyahm/xmux"

func main() {
	router := xmux.NewRouter()
	router.Get("/", nil)
	router.AddGroup(router.ShowSwagger("/docs", "localhost:8080"))
	router.Run()
}

package main

import (
	"log"
	"net/http"
	"xmux"
	"xmux/example/aritclegroup"
)

func show(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("show me!!!!"))
	return
}

func main() {
	router := xmux.NewRouter()
	router.HandleFunc("/get").Get(show)
	router.AddGroup(aritclegroup.Article())

	log.Fatal(http.ListenAndServe(":8080", router))
}

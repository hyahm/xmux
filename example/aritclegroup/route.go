package aritclegroup

import (
	"net/http"
	"xmux"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world!!!!"))
	return
}

func Article() *xmux.GroupRoute {
	article := xmux.NewGroupRoute("/article")
	article.HandleFunc("name").Get(hello)
	return article
}

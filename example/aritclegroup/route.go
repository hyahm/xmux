package aritclegroup

import (
	"fmt"
	"net/http"
	"xmux"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println(xmux.Var[r.URL.Path]["id"])
	w.Write([]byte("hello world!!!!"))
	return
}

func Article() *xmux.GroupRoute {
	article := xmux.NewGroupRoute("/article")
	article.HandleFunc("{int:id}").Get(hello)
	return article
}

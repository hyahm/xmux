package aritclegroup

import (
	"net/http"
	"xmux"
)

func hello(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(xmux.Var[r.URL.Path]["id"])
	w.Write([]byte("hello world!!!!"))
	return
}

var Article *xmux.GroupRoute

func init() {
	Article = xmux.NewGroupRoute()
	Article.Pattern("/{int:id}").Get(hello)

}

//func Article() *xmux.GroupRoute {
//	article := xmux.NewGroupRoute()
//	article.Pattern("/{int:id}").Get(hello)
//	return article
//}

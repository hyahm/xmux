package main

import (
	"net/http"
	"xmux"
)

func see(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(xmux.Var[r.URL.Path]["name"])
	w.Write([]byte("me"))
	return
}

func main() {
	router := xmux.NewRouter()
	router.SetHeader("momo", "bibi")
	router.Pattern("/name/age").SetHeader("one", "this").Get(see)
	article := xmux.NewGroupRoute()
	article.SetHeader("two", "self").Pattern("/my/{name}").Get(see).SetHeader("four", "func")
	article.SetHeader("three", "def").Pattern("/my/home/{name}").Get(see)
	router.AddGroup(article)
	http.ListenAndServe(":2000", router)

}

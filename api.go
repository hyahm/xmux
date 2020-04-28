package xmux

import (
	"fmt"
	"net/http"
)

func css(w http.ResponseWriter, r *http.Request) {
	filename := GetData(r).Var["name"]
	fmt.Println(filename)
	switch filename {
	case "style":
		w.Write([]byte(style))
		return
	case "left":
		w.Write([]byte(left))
		return
	case "font":
		w.Write([]byte(font))
		return
	}
}

func js(w http.ResponseWriter, r *http.Request) {
	filename := GetData(r).Var["name"]
	fmt.Println(filename)
	switch filename {
	case "jquery":
		w.Write([]byte(jqueryMin))
		return
	case "left":
		w.Write([]byte(left))
		return
	case "slimscroll":
		w.Write([]byte(slimscroll))
		return
	case "click":
		w.Write([]byte(click))
		return
	}
}

func NewApiDoc(name string) *GroupRoute {
	api := NewGroupRoute(name)
	api.Pattern("/-/js/{name}.js").Get(js).SetHeader("Content-Type", "application/javascript; charset=utf8")
	api.Pattern("/-/css/{name}.css").Get(css).SetHeader("Content-Type", "text/css; charset=utf8")
	return api
}

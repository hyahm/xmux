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

func ShowApi(name string, pattern string, r *Router) *GroupRoute {
	api := NewGroupRoute(name)
	api.Pattern("/-/js/{name}.js").Get(js).SetHeader("Content-Type", "application/javascript; charset=utf8")
	api.Pattern("/-/css/{name}.css").Get(css).SetHeader("Content-Type", "text/css; charset=utf8")
	api.Pattern(pattern).Get(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")

		doc := &Doc{
			Api:   make([]Document, 0),
			Title: "xmux docs",
		}

		t := NewTemplate()
		// 单路由
		r.route.AppendTo(pattern, doc)
		r.tpl.AppendTo(pattern, doc)

		// 组路由

		for _, g := range r.group {
			g.route.AppendTo(pattern, doc)
			g.tpl.AppendTo(pattern, doc)

		}
		err := t.Execute(w, *doc)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		return
	}))
	return api
}

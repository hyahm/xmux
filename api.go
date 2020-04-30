package xmux

import (
	"net/http"
	"strconv"
)

func css(w http.ResponseWriter, r *http.Request) {
	filename := Var(r)["name"]
	switch filename {
	case "style":
		w.Write([]byte(style))
		return
	case "left":
		w.Write([]byte(cssleft))
		return
	case "font":
		w.Write([]byte(font))
		return
	}
}

func js(w http.ResponseWriter, r *http.Request) {
	filename := Var(r)["name"]
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

func showThisDoc(w http.ResponseWriter, r *http.Request) {
	id := Var(r)["id"]
	t := NewTemplate()
	intid, _ := strconv.Atoi(id)
	if api, ok := ApiDocument[intid]; ok {
		api.Sidebar = sidebar

		err := t.Execute(w, api)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		return
	} else {
		ApiDocument[0] = Doc{
			Title:   "this is document home page",
			Sidebar: sidebar,
		}

		http.Redirect(w, r, "/-/api/0.html", 302)
	}

}

func homeDoc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	t := NewHomeTemplate()
	ApiDocument[0] = Doc{
		Title:   "this is document home page",
		Sidebar: sidebar,
	}
	err := t.Execute(w, ApiDocument[0])
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	return

}

func ShowApi(name string, pattern string, r *Router) *GroupRoute {
	api := NewGroupRoute(name)
	NewDocs(name, r)

	api.Pattern("/-/js/{name}.js").Get(js).SetHeader("Content-Type", "application/javascript; charset=utf8")
	api.Pattern("/-/css/{name}.css").Get(css).SetHeader("Content-Type", "text/css; charset=utf8")
	api.Pattern("/-/api/{int:id}.html").Get(showThisDoc).SetHeader("Content-Type", "text/html; charset=UTF-8")
	api.Pattern("/-/api/0.html").Get(homeDoc).SetHeader("Content-Type", "text/html; charset=UTF-8")
	api.Pattern(pattern).Get(homeDoc)
	return api
}

package xmux

// func css(w http.ResponseWriter, r *http.Request) {
// 	filename := Var(r)["name"]
// 	switch filename {
// 	case "style":
// 		w.Write([]byte(style))
// 		return
// 	case "left":
// 		w.Write([]byte(cssleft))
// 		return
// 	case "font":
// 		w.Write([]byte(font))
// 		return
// 	}
// }

// func js(w http.ResponseWriter, r *http.Request) {
// 	filename := Var(r)["name"]
// 	switch filename {
// 	case "jquery":
// 		w.Write([]byte(jqueryMin))
// 		return
// 	case "left":
// 		w.Write([]byte(left))
// 		return
// 	case "slimscroll":
// 		w.Write([]byte(slimscroll))
// 		return
// 	case "click":
// 		w.Write([]byte(click))
// 		return
// 	}
// }

// func showThisDoc(w http.ResponseWriter, r *http.Request) {
// 	id := Var(r)["id"]
// 	t := NewTemplate()
// 	intid, _ := strconv.Atoi(id)
// 	if api, ok := apiDocument[intid]; ok {
// 		api.Sidebar = sidebar
// 		err := t.Execute(w, api)
// 		if err != nil {
// 			w.Write([]byte(err.Error()))
// 		}
// 		return
// 	} else {
// 		apiDocument[0] = Doc{
// 			Title:   "this is document home page",
// 			Sidebar: sidebar,
// 		}
// 		http.Redirect(w, r, "/-/api/0.html", 302)
// 	}

// }

// func homeDoc(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
// 	t := NewHomeTemplate()
// 	apiDocument[0] = Doc{
// 		Title:   "this is document home page",
// 		Sidebar: sidebar,
// 	}
// 	err := t.Execute(w, apiDocument[0])
// 	if err != nil {
// 		w.Write([]byte(err.Error()))
// 	}
// 	w.WriteHeader(http.StatusOK)
// }

// func testdoc(w http.ResponseWriter, r *http.Request) {
// 	send, err := json.Marshal(apiDocument)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	w.Write(send)
// }

// func (r *Router) ShowApi(pattern string) *RouteGroup {
// 	if !r.new {
// 		panic("must be use get router by NewRouter()")
// 	}
// 	api := NewRouteGroup()
// 	NewDocs(r)
// 	api.Get("/-/js/{name}.js", js).SetHeader("Content-Type", "application/javascript; charset=utf8")
// 	api.Get("/-/css/{name}.css", css).SetHeader("Content-Type", "text/css; charset=utf8")
// 	api.Get("/-/api/{int:id}.html", showThisDoc).SetHeader("Content-Type", "text/html; charset=UTF-8")
// 	api.Get("/-/api/0.html", homeDoc).SetHeader("Content-Type", "text/html; charset=UTF-8")
// 	api.Get("/-/api/help", testdoc).SetHeader("Content-Type", "text/html; charset=UTF-8")
// 	api.Get(pattern, homeDoc)
// 	return api
// }

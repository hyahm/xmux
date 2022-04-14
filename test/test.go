package main

import (
	"net/http"
	_ "net/http/pprof"
	"strings"

	"github.com/hyahm/xmux"
)

func home1(w http.ResponseWriter, r *http.Request) bool {
	ct := r.Header.Get("Content-Type")
	ct = strings.ToLower(ct)
	if r.Method == http.MethodGet {

	}
	if r.Method == http.MethodPost {
		// b, err := io.ReadAll(r.Body)
		// if err != nil {
		// 	golog.Info(err)
		// 	w.WriteHeader(404)
		// 	return true
		// }
		// golog.Info(string(b))
		// err = json.Unmarshal(b, xmux.GetInstance(r).Data)
		// if err != nil {
		// 	w.WriteHeader(404)
		// 	return true
		// }

	}

	return false

}

func home(w http.ResponseWriter, r *http.Request) {
	xmux.GetInstance(r).Set("aaaa", "bbb")
}

type Global struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type User struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

func main() {
	global := &Global{
		Code: 200,
	}

	router := xmux.NewRouter()
	group := xmux.NewGroupRoute().BindResponse(global)
	group.Post("/post", home)
	router.AddGroup(group)
	router.Get("/get", home)
	router.Any("/", home).AddModule(home1).Bind(User{})
	router.DebugAssignRoute("/")
	router.Run(":8888")
}

package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/hyahm/xmux"
)

func home1(w http.ResponseWriter, r *http.Request) bool {
	// b, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	w.WriteHeader(404)
	// 	return true
	// }
	// fmt.Println(string(b))
	fmt.Println(r.FormValue("username"))
	// err := r.ParseForm()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(r.Form.Get("username"))
	// fmt.Println(r.Form.Get("username"))
	// ct := r.Header.Get("Content-Type")
	// ct = strings.ToLower(ct)

	// b, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	w.WriteHeader(404)
	// 	return true
	// }
	// fmt.Println(string(b))
	// err = json.Unmarshal(b, xmux.GetInstance(r).Data)
	// if err != nil {
	// 	w.WriteHeader(404)
	// 	return true
	// }

	return false

}

func home(w http.ResponseWriter, r *http.Request) {
	xmux.GetInstance(r).Set("aaaa", "bbb")

	// user := xmux.GetInstance(r).Data.(*User)
	// fmt.Printf("%#v\n", *user)
}

type Global struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type User struct {
	UserName string `json:"username,require" form:"username,require"`
	PassWord string `json:"password" form:"password"`
	Gender   bool   `json:"form" form:"gender"`
}

func main() {
	// global := &Global{
	// 	Code: 200,
	// }

	router := xmux.NewRouter()
	group := xmux.NewGroupRoute()
	group.Post("/post", home)
	router.Get("/get", home)
	router.Post("/", home).AddModule(home1)
	router.Run(":8888")
}

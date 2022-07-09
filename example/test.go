package main

import (
	"fmt"
	"net/http"

	"github.com/hyahm/xmux"
)

type DataFoo struct {
	UserName string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func AddFoo(w http.ResponseWriter, r *http.Request) {

	// r.ParseForm()
	// fmt.Println(r.PostFormValue("username"))
	// fmt.Println(r.Form.Get("username"))
	// fmt.Println(r.FormValue("username"))
	fmt.Println(string(xmux.GetInstance(r).Body))
	df := xmux.GetInstance(r).Data.(*DataFoo)
	fmt.Printf("%#v\n", df)
	w.Write([]byte("hello world"))
}

func UnmarshalError(err error, w http.ResponseWriter, r *http.Request) bool {
	fmt.Println(err)
	return false
}

func main() {
	router := xmux.NewRouter()
	router.Post("/bind/form", AddFoo)
	router.UnmarshalError = UnmarshalError
	// 也可以直接使用内置的
	router.Post("/bind/json", AddFoo).BindByContentType(&DataFoo{}) // 如果是json格式的可以直接 BindJson 与上面是类似的效果
	router.RunTLS("pri.key", "ca.crt")

}

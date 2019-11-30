package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"xmux"
	"xmux/example/aritclegroup"
)

func show(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("show me!!!!"))
	return
}

func postme(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("post me!!!!"))
	return
}

func Who(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(xmux.Var[r.URL.Path]["name"])
	//fmt.Println(xmux.Var[r.URL.Path]["age"])
	w.Write([]byte("yes is mine"))
	return
}

func testbool(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.FormValue("username"))
	fmt.Println(r.FormValue("password"))
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	fmt.Println(string(b))
	w.Write([]byte("yes is mine"))
	return
}

// 默认已经是这样的了，  如果有其他的请自定义
func Options() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})
}

func main() {
	router := xmux.NewRouter()
	router.Options = Options()                    // 这个是全局的options 请求处理， 前端预请求免除每次都要写个预请求的处理
	router.Pattern("/get").Get(show).Post(postme) // 不同请求分别处理
	//router.Pattern("/get/").Get(show).Post(postme) // 不同请求分别处理

	router.AddGroup(aritclegroup.Article)

	router.Pattern("/{string:age}").Get(Who).SetHeader("Host", "two")
	router.Pattern("/home/id").SetHeader("Host", "two").Post(testbool)
	router.Pattern("/home/{re:([a-z]{1,3})AAA([0-9]{1,3}):ch,zz}").SetHeader("Host", "two").Get(re)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func re(w http.ResponseWriter, r *http.Request) {
	fmt.Println(xmux.Var[r.URL.Path]["ch"])
	fmt.Println(xmux.Var[r.URL.Path]["id"])
	w.Write([]byte("hello"))
	return
}
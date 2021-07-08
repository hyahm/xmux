package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(time.Now().String()))
	return
}

func name(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("hello world name"))
	return
}

func me(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world me"))
	return
}

func all(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	w.Write([]byte("hello world all"))
}

func login(w http.ResponseWriter, r *http.Request) bool {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("not found data"))
		return true
	}
	err = json.Unmarshal(b, xmux.GetInstance(r).Data)
	return false
}

func Start(w http.ResponseWriter, r *http.Request) bool {
	r.Header.Set("bbb", "ccc")
	xmux.GetInstance(r).Set("start", time.Now())
	return false
}

type bbb struct {
	Name string
}

type Call struct {
	Code int    `json:"code" type:"int" need:"是" information:"错误返回码"`
	Msg  string `json:"msg" type:"string" need:"是" information:"错误信息"`
}

func (c *Call) Marshal() []byte {
	send, _ := json.Marshal(c)
	return send
}

func (c *Call) Error(msg string) []byte {
	c.Code = 1
	c.Msg = msg
	return c.Marshal()
}

func api(w http.ResponseWriter, r *http.Request) {
	c := &Call{}
	w.Write(c.Marshal())
}

type aaa struct {
	A int
	B int
}

// func NoHandleModule(w http.ResponseWriter, r *http.Request) bool {
// 	a := &aaa{
// 		A: 10,
// 		B: 20,
// 	}
// 	if err := xmux.HTML(w, `<html><head><title>{{ .A }}</title></head><body><h1>{{ .B }}</h1></body></html>`, a); err != nil {
// 		log.Fatal(err)
// 	}
// 	// w.Write([]byte("hello world"))
// 	return true
// }
// func HandleNotFound() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Add("Access-Control-Allow-Origin", "*")
// 		w.WriteHeader(http.StatusNotFound)
// 	})
// }

func hook(start time.Time, w http.ResponseWriter, r *http.Request) {
	golog.Info(time.Since(start).Seconds(), r.URL.Path)
}

func main() {
	router := xmux.NewRouter()
	router.Hook = hook
	// router.HandleNotFound = HandleNotFound()
	router.SetHeader("Content-Type", "aaa")
	router.Post("/home", home)
	// router.HandleNotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("not found this url in server, url: " + r.URL.Path))
	// })
	router.Get("/home", home).ApiCreateGroup("home", "showthis home", "hometest").SetHeader("Content-Type", "bbbb").
		ApiDescribe("这是home接口的测试").
		ApiReqHeader("content-type", "application/json").
		ApiRequestTemplate(`{"addr": "shenzhen", "people": 5}`).
		ApiResStruct(Call{}).
		ApiResponseTemplate(`{"code": 0, "msg": ""}`).
		ApiCodeField("133").ApiCodeMsg("1", "56").ApiCodeMsg("3", "akhsdklfhl").ApiDelReqHeader("aaaa").ApiCodeMsg("78", "").MiddleWare(xmux.GetExecTime)

	user := xmux.NewGroupRoute().ApiReqHeader("aaaa", "bbbb")
	user.ApiCodeMsg("98", "你是98")
	user.ApiCodeMsg("78", "你是78")
	user.ApiCodeMsg("0", "成功")
	user.Post("/api", api)
	user.Get("/bbb/ccc/{int:oid}/{string:all}", all)

	router.AddGroup(user)
	router.AddGroup(router.Pprof())
	doc := router.ShowApi("/docs")
	router.AddGroup(doc) // 开启文档， 一般都是写在路由的最后, 后面的api不会显示
	fmt.Println("-----------------------")
	// router.DebugRoute()
	router.Run(":8888")

}

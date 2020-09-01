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
	fmt.Println(xmux.GetData(r).Data)
	w.Write([]byte("hello world home"))
	return
}

func name(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home")
	fmt.Println(xmux.Var(r)["name"])

	w.Write([]byte("hello world name"))
	return
}

func me(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world me"))
	return
}

func all(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world all"))
	return
}

func login(w http.ResponseWriter, r *http.Request) bool {
	fmt.Println("login mw")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("not found data"))
		return true
	}
	err = json.Unmarshal(b, xmux.GetData(r).Data)

	fmt.Println(xmux.GetData(r).Data)
	return false
}

func Start(w http.ResponseWriter, r *http.Request) bool {
	r.Header.Set("bbb", "ccc")
	xmux.GetData(r).Set("start", time.Now())
	return false
}

type aaa struct {
	Age int
}

type bbb struct {
	Name string
}

type Home struct {
	Addr   string `json:"addr" type:"string" need:"是" default:"深圳" information:"家庭住址"`
	People int    `json:"people" type:"int" need:"是" default:"1" information:"有多少个人"`
}

type Call struct {
	Code int    `json:"code" type:"int" need:"是" information:"错误返回码"`
	Msg  string `json:"msg" type:"string" need:"是" information:"错误信息"`
}

func End(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("time: %f \n", time.Since(xmux.GetData(r).Get("start").(time.Time)).Seconds())
}

func main() {
	router := xmux.NewRouter()
	router.SetHeader("Content-Type", "aaa")
	router.AddModule(Start).EndModule(End)
	router.Get("/asdf/{name}", all)
	router.HandleNotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not found this url in server, url: " + r.URL.Path))
		return
	})
	user := xmux.NewGroupRoute().ApiReqHeader("aaaa", "bbbb")
	user.ApiCodeMsg("98", "你是98")
	user.ApiCodeMsg("78", "你是78")
	user.ApiCodeMsg("0", "成功")

	router.Get("/home", home).ApiCreateGroup("home", "showthis home", "hometest").SetHeader("Content-Type", "bbbb").
		ApiDescribe("这是home接口的测试").
		ApiReqHeader("content-type", "application/json").
		ApiReqStruct(&Home{}).
		ApiRequestTemplate(`{"addr": "shenzhen", "people": 5}`).
		ApiResStruct(Call{}).
		ApiResponseTemplate(`{"code": 0, "msg": ""}`).
		ApiSupplement("这个是接口的说明补充， 没补充就不填").Bind(&Home{}).AddModule(login).
		ApiCodeField("133").ApiCodeMsg("1", "56").ApiCodeMsg("3", "akhsdklfhl").ApiDelReqHeader("aaaa").ApiCodeMsg("78", "")

	user.Post("/aaa/{name}", name).ApiCreateGroup("test", "这是一个大写的测试组", "testaaa").
		ApiReqHeader("content-type", "application/json").
		ApiReqStruct(&Home{}).
		ApiRequestTemplate(`{"addr": "shenzhen", "people": 5}`).
		ApiResStruct(Call{}).
		ApiResponseTemplate(`{"code": 0, "msg": ""}`).
		ApiSupplement("这个是接口的说明补充， 没补充就不填").Bind(&Home{}).AddModule(login).
		ApiCodeField("133").ApiCodeMsg("1", "56").ApiCodeMsg("3", "akhsdklfhl").ApiDelReqHeader("aaaa").ApiCodeMsg("78", "")
	user.Post("/aaa/bbbb/{path:me}", me).
		ApiReqHeader("content-type", "application/json").
		ApiReqStruct(&Home{}).
		ApiRequestTemplate(`{"addr": "shenzhen", "people": 5}`).
		ApiResStruct(Call{}).
		ApiResponseTemplate(`{"code": 0, "msg": ""}`).
		ApiSupplement("这个是接口的说明补充， 没补充就不填").Bind(&Home{}).AddModule(login).
		ApiCodeField("133").ApiCodeMsg("1", "56").ApiCodeMsg("3", "akhsdklfhl")

	user.Get("/bbb/ccc/{int:oid}/{string:all}", all)
	router.Post("/cut", func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadAll(r.Body)
		golog.Info(string(b))
		a := struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}{}
		send, _ := json.Marshal(a)
		w.Write(send)
	})
	router.AddGroup(user)
	router.DebugRoute()
	router.AddGroup(router.Pprof())
	doc := router.ShowApi("/docs")
	router.AddGroup(doc) // 开启文档， 一般都是写在路由的最后, 后面的api不会显示

	router.Run()

}

func Midware(midware func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request) bool {
	midware(w, r)
	return true
}

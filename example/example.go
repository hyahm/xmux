package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hyahm/xmux"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println(xmux.Var[r.URL.Path]["test"])
	fmt.Println(xmux.GetData(r).Data)
	w.Write([]byte("hello world home"))
	return
}

func name(w http.ResponseWriter, r *http.Request) {
	fmt.Println(xmux.Var[r.URL.Path]["name"])
	w.Write([]byte("hello world name"))
	return
}

func me(w http.ResponseWriter, r *http.Request) {
	fmt.Println(xmux.Var[r.URL.Path]["me"])
	w.Write([]byte("hello world me"))
	return
}

func all(w http.ResponseWriter, r *http.Request) {
	fmt.Println(xmux.Var[r.URL.Path]["all"])
	fmt.Println(xmux.Var[r.URL.Path]["oid"])
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

func filter(w http.ResponseWriter, r *http.Request) bool {
	fmt.Println("-----------------------")
	fmt.Println("login filter")
	r.Header.Set("bbb", "ccc")

	xmux.Ctx[r.URL.Path] = context.WithValue(context.Background(), "conf", "body")
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
	Code int    `json:"code" type:"int" information:"错误返回码"`
	Msg  string `json:"msg" type:"string" information:"错误信息"`
}

func main() {

	router := xmux.NewRouter()
	router.IgnoreIco = true
	// fmt.Println(router.Slash)
	router.AddMidware(filter)
	router.Pattern("/home").Post(home).ApiDescribe("这是home接口的测试").
		ApiReqHeader(map[string]string{"content-type": "application/json"}).
		ApiReqStruct(&Home{}).
		ApiRequestTemplate(`{"addr": "shenzhen", "people": 5}`).
		ApiResStruct(Call{}).
		ApiResponseTemplate(`{"code": 0, "msg": ""}`).
		ApiSupplement("这个是接口的说明补充， 没补充就不填").Bind(&Home{}).AddMidware(login).Get(home)
	router.Pattern("/aaa/{name}").Post(name).DelMidware(filter).Get(name)
	router.Pattern("/aaa/bbbb/{path:me}").Post(me)
	router.Pattern("/bbb/ccc/{int:oid}/{string:all}").Get(all)

	router.ShowApi("/doc") // 开启文档， 一般都是写在路由的最后, 后面的api不会显示
	if err := http.ListenAndServe(":9000", router); err != nil {
		log.Fatal(err)
	}

}

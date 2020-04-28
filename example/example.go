package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hyahm/xmux"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println(xmux.GetData(r).Data)
	w.Write([]byte("hello world home"))
	return
}

func name(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home")
	fmt.Println(xmux.GetData(r).Var["name"])

	w.Write([]byte("hello world name"))
	return
}

func me(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world me"))
	return
}

func all(w http.ResponseWriter, r *http.Request) {
	xmux.GetData(r).End = "13333"
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
	fmt.Println("login filter.............")
	r.Header.Set("bbb", "ccc")

	return false
}

func end(end interface{}) {
	fmt.Println("-----------------------")
	fmt.Println(end)
	fmt.Println("end function ")

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

func main() {

	router := xmux.NewRouter()
	router.IgnoreIco = true
	router.AddMidware(filter)
	router.Pattern("/home").Post(home).
		ApiDescribe("这是home接口的测试").
		ApiReqHeader(map[string]string{"content-type": "application/json"}).
		ApiReqStruct(&Home{}).
		ApiRequestTemplate(`{"addr": "shenzhen", "people": 5}`).
		ApiResStruct(Call{}).
		ApiResponseTemplate(`{"code": 0, "msg": ""}`).
		ApiSupplement("这个是接口的说明补充， 没补充就不填").Bind(&Home{}).AddMidware(login).
		ApiCodeField("133").ApiCodeMsg("1", "56").ApiCodeMsg("3", "akhsdklfhl")
	router.Pattern("/aaa/{name}").Post(name).DelMidware(filter).Get(name)
	router.Pattern("/aaa/bbbb/{path:me}").Post(me).Get(me)
	router.Pattern("/bbb/ccc/{int:oid}/{string:all}").Get(all).End(end)

	router.AddGroup(xmux.ShowApi("doc", "/doc", router)) // 开启文档， 一般都是写在路由的最后, 后面的api不会显示
	if err := http.ListenAndServe(":9000", router); err != nil {
		log.Fatal(err)
	}

}

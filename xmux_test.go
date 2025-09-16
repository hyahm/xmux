package xmux

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func home(w http.ResponseWriter, r *http.Request) {
	// name := Var(r)["name"]
	// time.Sleep(time.Millisecond * 30)
	// fmt.Println(1111)
	// GetInstance(r).Set(RESPONSEBYTES, []byte("asjdlfjalsdf"))

}

func grouphome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("mmmmmmmm group")
	w.Write([]byte("grouphome" + Var(r)["name"] + "-" + Var(r)["age"]))
}

func adminhandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("mmmmmmmm admin")
}

func adminGroup() *RouteGroup {
	admin := NewRouteGroup().Prefix("test")
	admin.Get("/admin/{bbb}", home)
	admin.Get("/admin", adminhandle)
	admin.Get("/aaa/adf{re:([a-z]{1,4})sf([0-9]{0,10})sd: name, age}", grouphome)
	return admin
}

func userGroup() *RouteGroup {
	user := NewRouteGroup()
	// user.Get("/group", home).Use(CombineHandlers())
	user.Get("/group", home)
	user.AddGroup(adminGroup()).DelPostModule(postModule)
	return user
}

func TestMain(t *testing.T) {

	router := NewRouter()
	// router.HandleAll = LimitFixedWindowCounterTemplate
	router.SetHeader("Access-Control-Allow-Origin", "*")
	router.SetHeader("Content-Type", "application/x-www-form-urlencoded,application/json; charset=UTF-8")
	router.SetHeader("Access-Control-Allow-Headers", "Content-Type")
	router.SetHeader("Access-Control-Max-Age", "1728000")
	// router.SetHeader("Access-Control-Allow-Origin", "*").
	// 	SetHeader("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	router.AddGroup(Pprof())

	// router.Prefix("/api")
	// router.EnableConnect = true
	router.Get("/test", nil)
	router.Get("/post", pp).AddPostModule(postModule)
	router.HandleAll = nil
	// router.SetAddr(":8080")
	router.AddGroup(userGroup())
	log.Fatal(router.SetAddr(":9999").Run())
}

type Binding struct {
	ID   int64  `json:"id,required"`
	Name string `json:"name,required"`
}

func postModule(w http.ResponseWriter, r *http.Request) bool {

	fmt.Println("这是一个后置模块")
	return false
}

func pp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("这是一个处理函数")
}

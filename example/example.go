package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hyahm/xmux"
)

func home(w http.ResponseWriter, r *http.Request) {
	// pages := xmux.GetInstance(r).Data.(*AAA)
	// fmt.Println(*pages)
	fmt.Println(r.Method, r.URL.Path)
	// atomic.AddInt64(&count, 1)
	// // time.Sleep(1 * time.Second)
	// fmt.Fprintf(w, "hello home")
	w.Write([]byte(r.Method))
}

func Start(w http.ResponseWriter, r *http.Request) bool {
	r.Header.Set("bbb", "ccc")
	xmux.GetInstance(r).Set("start", time.Now())
	return false
}

func PermMudule(w http.ResponseWriter, r *http.Request) bool {

	// 通过类似 token 获取到用户的uid
	// uid := "from token" //
	// 直接写在token 验证路由里面或者 单独在加一个module 都可以
	// 根据uid 获取用户的CURD
	// 获取一个给定的结构存储函数对应的细致权限
	// 根据uid 判断自己有什么权限， 假如是retrieve 权限
	// ------------    这是开发给定的固定值  ---------------------
	// create := []string{"Create"} // 注意大小写方便判断
	retrieve := []string{"List", "Get"}
	// ------------    这是开发给定的固定值  ---------------------
	// 获取执行函数的方法名

	currFun := xmux.GetInstance(r).Get(xmux.CURRFUNCNAME) // module 或 handle 中都必定有此 值
	// pages := xmux.GetInstance(r).Get(xmux.PAGES)          // 页面的权限， 一般都是假如到路由组 必定有此 值
	// 假如路由匹配到这里 func List(w w http.ResponseWriter, r *http.Request) {}  currFun = "List"
	// router.Post("/home", List)
	// 因为 currFun = "List" 所以 retrieve 中 包含了 List 也就是 权限符合
	// 增加switch 来判断多条件即可
	for _, v := range retrieve {
		if v == currFun {
			// 符合条件， 放行
			return false
		}
	}
	w.Write([]byte("no permission"))
	// 认证失败， 直接返回
	return true
}

type AAA struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func GetExecTime(handle func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Set example variable
	xmux.GetInstance(r).Set("example", "12345")

	// before request

	handle(w, r)

	fmt.Printf("url: %s -- addr: %s -- method: %s -- exectime: %f\n", r.URL.Path, r.RemoteAddr, r.Method, time.Since(start).Seconds())
}

func TestM(rw http.ResponseWriter, r *http.Request) bool {
	fmt.Println("module M")
	return false
}
func TestN(rw http.ResponseWriter, r *http.Request) bool {
	fmt.Println("module N")
	return false
}
func TestJ(rw http.ResponseWriter, r *http.Request) bool {
	fmt.Println("module J")
	return false
}

func TestK(rw http.ResponseWriter, r *http.Request) bool {
	fmt.Println("module K")
	return false
}
func TestL(rw http.ResponseWriter, r *http.Request) bool {
	fmt.Println("module L")
	return false
}

func TestO(rw http.ResponseWriter, r *http.Request) bool {
	fmt.Println("module O")
	return false
}

func initv1() *xmux.GroupRoute {
	v1 := xmux.NewGroupRoute().DelModule(CheckToken).DelModule(CheckPage)
	// 更新状态
	// 由key的url 多出来的一个url
	v1.Get("/json/key", nil)
	// v1.Get("/iv/{string:filename}", nil)
	// v1.Get("/mysource/{type}/{string:filename}", nil)
	// v1.Post("/post/url/{int:id}", nil)
	v1.AddGroup(nil)
	return v1
}

func CheckToken(w http.ResponseWriter, r *http.Request) bool {
	a := r.Header.Get("Authorization")
	auths := strings.Split(a, " ")
	return len(auths) <= 1
}

func CheckPage(w http.ResponseWriter, r *http.Request) bool {
	return false
}

func printTime(ws *xmux.BaseWs, mt byte) {
	for {
		if ws.Conn == nil {
			return
		}
		err := ws.SendMessage([]byte(time.Now().String()), mt)
		if err != nil {
			log.Println("write:", err)
			break
		}
		time.Sleep(1 * time.Second)
	}

}

func main() {
	router := xmux.NewRouter().MiddleWare(xmux.DefaultMidwareTemplate)

	router.SetHeader("Access-Control-Allow-Origin", "*")
	router.SetHeader("Content-Type", "application/x-www-form-urlencoded,application/json; charset=UTF-8")
	router.SetHeader("Access-Control-Allow-Headers", "Content-Type,smail,authorization")
	// router.Slash = true
	// router.MiddleWare(GetExecTime)
	router.Request("/aaaa", home, http.MethodGet, http.MethodPost)
	router.Get("/{user}/{info}", func(rw http.ResponseWriter, r *http.Request) {
		ws, err := xmux.UpgradeWebSocket(rw, r)
		os.OpenFile("", os.O_APPEND, 0755)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer ws.Close()
		for {
			mt, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", message)
			err = ws.SendMessage([]byte(message), mt)
			if err != nil {
				log.Println("write:", err)
				break
			}
			go printTime(ws, mt)
		}
	})
	router.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("ok"))
	}).DelMiddleWare(xmux.DefaultMidwareTemplate)

	router.AddGroup(initv1())
	router.DebugAssignRoute("/json/key")
	// router.Get("/stop", func(rw http.ResponseWriter, r *http.Request) {
	// 	xmux.StopService()
	// 	r.BasicAuth()
	// })
	// router.AddGroup(xmux.Pprof())
	// router.SetHeader("Content-Type", "aaa")
	// // router.HandleNotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// // 	w.Write([]byte("not found this url in server, url: " + r.URL.Path))
	// // })
	// router.Get("/home", home).ApiCreateGroup("home", "showthis home", "hometest").SetHeader("Content-Type", "bbbb").
	// 	ApiDescribe("这是home接口的测试").
	// 	ApiReqHeader("content-type", "application/json").
	// 	ApiRequestTemplate(`{"addr": "shenzhen", "people": 5}`).
	// 	ApiResStruct(Call{}).
	// 	ApiResponseTemplate(`{"code": 0, "msg": ""}`).
	// 	ApiCodeField("133").ApiCodeMsg("1", "56").ApiCodeMsg("3", "akhsdklfhl").
	// 	ApiDelReqHeader("aaaa").ApiCodeMsg("78", "")

	// user := xmux.NewGroupRoute().ApiCreateGroup("user", "用户相关", "用户相关").
	// 	ApiReqHeader("aaaa", "bbbb").DelPageKeys("me")
	// user.ApiCodeMsg("98", "你是98")
	// user.ApiCodeMsg("78", "你是78")
	// user.ApiCodeMsg("0", "成功")

	// user.Get("/api", home).ApiDescribe("yonghu api")
	// user.Post("/api", home).BindJson(AAA{})

	// user.Get("/user/info", home).AddPageKeys("roles")
	// user.Get("/bbb/ccc/{int:oid}/{string:all}", all)

	// router.AddGroup(user)
	// router.AddGroup(router.Pprof())
	// doc := router.ShowApi("/docs")
	// router.AddGroup(doc) // 开启文档， 一般都是写在路由的最后, 后面的api不会显示
	// router.DebugRoute()
	router.Run(":8888")
}

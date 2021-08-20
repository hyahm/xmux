package xmux

import (
	"net/http"
	"strings"
	"testing"
)

func mw1(handle func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("mw1"))
	handle(w, r)
}

func mw2(handle func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("mw2"))
	handle(w, r)
}

func mw3(handle func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("mw3"))
	handle(w, r)
}

func mw4(handle func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("mw4"))
	handle(w, r)
}

var g1 *GroupRoute
var g2 *GroupRoute
var g3 *GroupRoute
var g4 *GroupRoute
var g5 *GroupRoute
var g6 *GroupRoute

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("home"))
}

func init() {
	g2 = NewGroupRoute().MiddleWare(mw4)
	g2.Get("/g2", home)
}

func init() {
	g3 = NewGroupRoute()
	g3.Get("/g3", home).MiddleWare(mw1)
}

func init() {
	g4 = NewGroupRoute().DelMiddleWare(mw1)
	g4.Get("/g4", home)
	g4.Get("/g4_1", home)
}

func init() {
	g5 = NewGroupRoute().MiddleWare(mw3)
	g5.Get("/g5", home).MiddleWare(mw2)
}

func init() {
	g6 = NewGroupRoute().MiddleWare(mw3)
	g6.Get("/g6", home).DelMiddleWare(mw3)
	g6.Get("/g6_1", home).DelMiddleWare(mw1)
}

func init() {
	g1 = NewGroupRoute().MiddleWare(mw2)
	g1.AddGroup(g2)
}

// 嵌套带中间件的组路由，组路由的中间件覆盖最外层的router的中间件
func TestGroupMidwareCoveredRouterMidware(t *testing.T) {
	router := NewRouter().MiddleWare(mw1)
	router.AddGroup(g2)
	t.Log(router.GetAssignRoute("/g2")["GET"].GetMidwareName())
	if !strings.Contains(router.GetAssignRoute("/g2")["GET"].GetMidwareName(), "mw4") {
		t.Fatal("get error midware")
	}
}

// 路由的中间件覆盖路由组的中间件
func TestRouteMidwareCoveredGroupMidware(t *testing.T) {
	router := NewRouter()
	router.AddGroup(g5)
	t.Log(router.GetAssignRoute("/g5")["GET"].GetMidwareName())
	if !strings.Contains(router.GetAssignRoute("/g5")["GET"].GetMidwareName(), "mw2") {
		t.Fatal("get error midware")
	}
}

// 组路由删除中间，此组的中间件都将删除，但是外层的没影响
func TestDeleteGroupMidwareWillDeleteRouterMidware(t *testing.T) {
	router := NewRouter().MiddleWare(mw1)
	router.Get("/r1", nil)
	router.AddGroup(g4)
	router.Get("/r2", nil)
	t.Log(router.GetAssignRoute("/g4")["GET"].GetMidwareName())
	t.Log(router.GetAssignRoute("/g4_1")["GET"].GetMidwareName())
	if router.GetAssignRoute("/g4")["GET"].GetMidwareName() != "" {
		t.Fatal("get error midware")
	}
	if router.GetAssignRoute("/g4_1")["GET"].GetMidwareName() != "" {
		t.Fatal("get error midware")
	}
}

// router和组里面都带了中间件，组里面的中间件会覆盖最外层
// 删除的话，只有删除组里面的才会删除
// 路由删除中间件，将会删除对应Router或者组的中间件， 如果外层的的不一样
func TestOnlyOneMidwareBeStayRouteDelMidwareWillDeleteTheLastestMidware(t *testing.T) {
	router := NewRouter().MiddleWare(mw1)
	router.AddGroup(g6)
	t.Log(router.GetAssignRoute("/g6")["GET"].GetMidwareName())
	t.Log(router.GetAssignRoute("/g6_1")["GET"].GetMidwareName())
	if router.GetAssignRoute("/g6")["GET"].GetMidwareName() != "" {
		t.Fatal("get error midware")
	}
	if !strings.Contains(router.GetAssignRoute("/g6_1")["GET"].GetMidwareName(), "mw3") {
		t.Fatal("get error midware")
	}
}

// router里面挂载了全局的中间件， 测试下面的url是否都挂载了
func TestOnlyInRouter(t *testing.T) {
	router := NewRouter().MiddleWare(mw1)
	router.Get("/r1", nil)
	router.Get("/r2", nil)
	t.Log(router.GetAssignRoute("/r1")["GET"].GetMidwareName())
	t.Log(router.GetAssignRoute("/r2")["GET"].GetMidwareName())
	if !strings.Contains(router.GetAssignRoute("/r1")["GET"].GetMidwareName(), "mw1") {
		t.Fatal("get error midware")
	}
	if !strings.Contains(router.GetAssignRoute("/r2")["GET"].GetMidwareName(), "mw1") {
		t.Fatal("get error midware")
	}
}

// func TestHttp(t *testing.T) {
// 	handler := func(w http.ResponseWriter, r *http.Request) {
// 		io.WriteString(w, "<html><body>Hello World!</body></html>")
// 	}

// 	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
// 	w := httptest.NewRecorder()
// 	handler(w, req)

// 	resp := w.Result()
// 	body, _ := ioutil.ReadAll(resp.Body)

// 	fmt.Println(resp.StatusCode)
// 	fmt.Println(resp.Header.Get("Content-Type"))
// 	fmt.Println(string(body))
// }

// func TestHealthCheckHandler(t *testing.T) {
// 	//创建一个请求
// 	req, err := http.NewRequest("GET", "/health-check", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// 我们创建一个 ResponseRecorder (which satisfies http.ResponseWriter)来记录响应
// 	rr := httptest.NewRecorder()

// 	//直接使用HealthCheckHandler，传入参数rr,req
// 	home(rr, req)

// 	// 检测返回的状态码
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}
// 	t.Log(rr.Body.String())
// 	// 检测返回的数据
// 	expected := `{"alive": true}`
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}
// }

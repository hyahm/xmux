package xmux

import (
	"net/http"
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

var g1 *GroupRoute
var g2 *GroupRoute

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("home"))
}

func init() {
	g2 = NewGroupRoute().MiddleWare(mw3)
	g2.Get("/", home)
}

func init() {
	g1 = NewGroupRoute().MiddleWare(mw2)
	g1.AddGroup(g2)
}

func Test_Module(t *testing.T) {
	router := NewRouter().MiddleWare(mw1)
	router.AddGroup(g1)
	router.DebugAssignRoute("/")
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

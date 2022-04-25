package xmux

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hyahm/xmux/helper"
)

func BenchmarkOneRoute(B *testing.B) {
	router := NewRouter()
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {})
	runRequest(B, router, "Get", "/ping")
}

func BenchmarkRecoveryMiddleware(B *testing.B) {
	router := NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {})
	runRequest(B, router, "Get", "/")
}

func BenchmarkLoggerMiddleware(B *testing.B) {
	router := NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {})
	runRequest(B, router, "Get", "/")
}

func BenchmarkManyHandlers(B *testing.B) {
	router := NewRouter()
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {})
	runRequest(B, router, "Get", "/ping")
}

func Benchmark5Params(B *testing.B) {
	router := NewRouter()
	router.Get("/param/{param1}/{params2}/{param3}/{param4}/{param5}", func(w http.ResponseWriter, r *http.Request) {})
	runRequest(B, router, "Get", "/param/path/to/parameter/john/12345")
}

func BenchmarkOneRouteJSON(B *testing.B) {
	router := NewRouter()
	data := struct {
		Status string `json:"status"`
	}{"ok"}
	router.Get("/json", func(w http.ResponseWriter, r *http.Request) {
		send, _ := json.Marshal(data)
		// w.Write([]byte(`{"status": "ok}`))
		w.Write(send)
	})
	runRequest(B, router, "Get", "/json")
}

// func BenchmarkOneRouteHTML(B *testing.B) {
// 	router := NewRouter()
// 	t := template.Must(template.NewRouter("index").Parse(`
// 		<html><body><h1>{{.}}</h1></body></html>`))
// 	router.SetHTMLTemplate(t)

// 	router.Get("/html", func(w http.ResponseWriter, r *http.Request) {
// 		c.HTML(http.StatusOK, "index", "hola")
// 	})
// 	runRequest(B, router, "Get", "/html")
// }

// func BenchmarkOneRouteSet(B *testing.B) {
// 	router := NewRouter()
// 	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
// 		c.Set("key", "value")
// 	})
// 	runRequest(B, router, "Get", "/ping")
// }

func BenchmarkOneRouteString(B *testing.B) {
	router := NewRouter()
	router.Get("/text", func(w http.ResponseWriter, r *http.Request) {
		w.Write(helper.StringToBytes("this is a plain text"))
	})
	runRequest(B, router, "Get", "/text")
}

func BenchmarkManyRoutesFist(B *testing.B) {
	router := NewRouter()
	router.Any("/ping", func(w http.ResponseWriter, r *http.Request) {})
	runRequest(B, router, "Get", "/ping")
}

func BenchmarkManyRoutesLast(B *testing.B) {
	router := NewRouter()
	router.Any("/ping", func(w http.ResponseWriter, r *http.Request) {})
	runRequest(B, router, "OPTIONS", "/ping")
}

func Benchmark404(B *testing.B) {
	router := NewRouter()
	router.Any("/something", func(w http.ResponseWriter, r *http.Request) {})
	router.HandleNotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
	// router.NoRoute(func(w http.ResponseWriter, r *http.Request) {})
	runRequest(B, router, "Get", "/ping")
}

func Benchmark404Many(B *testing.B) {
	router := NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {})
	router.Get("/path/to/something", func(w http.ResponseWriter, r *http.Request) {})
	router.Get("/post/:id", func(w http.ResponseWriter, r *http.Request) {})
	router.Get("/view/:id", func(w http.ResponseWriter, r *http.Request) {})
	router.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	router.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {})
	router.Get("/delete/:id", func(w http.ResponseWriter, r *http.Request) {})
	router.Get("/user/:id/:mode", func(w http.ResponseWriter, r *http.Request) {})

	router.HandleNotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
	runRequest(B, router, "Get", "/viewfake")
}

type mockWriter struct {
	headers http.Header
}

func NewRouterMockWriter() *mockWriter {
	return &mockWriter{
		http.Header{},
	}
}

func (m *mockWriter) Header() (h http.Header) {
	return m.headers
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockWriter) WriteHeader(int) {}

func runRequest(B *testing.B, r *Router, method, path string) {
	// create fake request
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}
	w := NewRouterMockWriter()
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		r.ServeHTTP(w, req)
	}
}

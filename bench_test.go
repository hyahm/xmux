package xmux

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkOneRoute(B *testing.B) {
	router := NewRouter()
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {})
	runRequest(B, router, "GET", "/ping")
}

// func BenchmarkRecoveryMiddleware(B *testing.B) {
// 	router := New()
// 	router.Use(Recovery())
// 	router.GET("/", func(c *Context) {})
// 	runRequest(B, router, "GET", "/")
// }

// func BenchmarkLoggerMiddleware(B *testing.B) {
// 	router := New()
// 	router.Use(LoggerWithWriter(newMockWriter()))
// 	router.GET("/", func(c *Context) {})
// 	runRequest(B, router, "GET", "/")
// }

// func BenchmarkManyHandlers(B *testing.B) {
// 	router := New()
// 	router.Use(Recovery(), LoggerWithWriter(newMockWriter()))
// 	router.Use(func(c *Context) {})
// 	router.Use(func(c *Context) {})
// 	router.GET("/ping", func(c *Context) {})
// 	runRequest(B, router, "GET", "/ping")
// }

// func Benchmark5Params(B *testing.B) {
// 	DefaultWriter = os.Stdout
// 	router := New()
// 	router.Use(func(c *Context) {})
// 	router.GET("/param/:param1/:params2/:param3/:param4/:param5", func(c *Context) {})
// 	runRequest(B, router, "GET", "/param/path/to/parameter/john/12345")
// }

// func BenchmarkOneRouteJSON(B *testing.B) {
// 	router := New()
// 	data := struct {
// 		Status string `json:"status"`
// 	}{"ok"}
// 	router.GET("/json", func(c *Context) {
// 		c.JSON(http.StatusOK, data)
// 	})
// 	runRequest(B, router, "GET", "/json")
// }

// func BenchmarkOneRouteHTML(B *testing.B) {
// 	router := New()
// 	t := template.Must(template.New("index").Parse(`
// 		<html><body><h1>{{.}}</h1></body></html>`))
// 	router.SetHTMLTemplate(t)

// 	router.GET("/html", func(c *Context) {
// 		c.HTML(http.StatusOK, "index", "hola")
// 	})
// 	runRequest(B, router, "GET", "/html")
// }

// func BenchmarkOneRouteSet(B *testing.B) {
// 	router := New()
// 	router.GET("/ping", func(c *Context) {
// 		c.Set("key", "value")
// 	})
// 	runRequest(B, router, "GET", "/ping")
// }

// func BenchmarkOneRouteString(B *testing.B) {
// 	router := New()
// 	router.GET("/text", func(c *Context) {
// 		c.String(http.StatusOK, "this is a plain text")
// 	})
// 	runRequest(B, router, "GET", "/text")
// }

// func BenchmarkManyRoutesFist(B *testing.B) {
// 	router := New()
// 	router.Any("/ping", func(c *Context) {})
// 	runRequest(B, router, "GET", "/ping")
// }

// func BenchmarkManyRoutesLast(B *testing.B) {
// 	router := New()
// 	router.Any("/ping", func(c *Context) {})
// 	runRequest(B, router, "OPTIONS", "/ping")
// }

// func Benchmark404(B *testing.B) {
// 	router := New()
// 	router.Any("/something", func(c *Context) {})
// 	router.NoRoute(func(c *Context) {})
// 	runRequest(B, router, "GET", "/ping")
// }

func Benchmark404Many(B *testing.B) {
	router := NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {})
	router.Get("/path/to/something", func(w http.ResponseWriter, r *http.Request) {})
	router.Get("/post/{id}", func(w http.ResponseWriter, r *http.Request) {})
	router.Get("/view/:id", func(w http.ResponseWriter, r *http.Request) {})
	router.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	router.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {})
	router.Get("/delete/:id", func(w http.ResponseWriter, r *http.Request) {})
	router.Get("/user/:id/:mode", func(w http.ResponseWriter, r *http.Request) {})

	// router.NoRoute(func(c *Context) {})
	runRequest(B, router, "GET", "/viewfake")
}

type mockWriter struct {
	headers http.Header
}

func newMockWriter() *mockWriter {
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
	w := newMockWriter()
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		r.ServeHTTP(w, req)
	}
}
func BenchmarkMux(b *testing.B) {
	router := NewRouter()
	handler := func(w http.ResponseWriter, r *http.Request) {}
	router.Get("/v1/{v1}", handler)
	request, _ := http.NewRequest("GET", "/v1/anything", nil)
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(nil, request)
	}
}

func BenchmarkMuxAlternativeInRegexp(b *testing.B) {
	router := NewRouter()
	handler := func(w http.ResponseWriter, r *http.Request) {}
	router.Get("/v1/{v1}", handler)

	requestA, _ := http.NewRequest("GET", "/v1/a", nil)
	requestB, _ := http.NewRequest("GET", "/v1/b", nil)
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(nil, requestA)
		router.ServeHTTP(nil, requestB)
	}
}

func BenchmarkManyPathVariables(b *testing.B) {
	router := NewRouter()
	handler := func(w http.ResponseWriter, r *http.Request) {}
	router.Get("/v1/{v1}/{v2}/{v3}/{v4}/{v5}", handler)

	matchingRequest, _ := http.NewRequest("GET", "/v1/1/2/3/4/5", nil)
	notMatchingRequest, _ := http.NewRequest("GET", "/v1/1/2/3/4", nil)
	recorder := httptest.NewRecorder()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(nil, matchingRequest)
		router.ServeHTTP(recorder, notMatchingRequest)
	}
}

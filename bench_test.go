package xmux

import (
	"net/http"
	"testing"
)

func BenchmarkOneRoute(B *testing.B) {
	router := NewRouter()
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {})
	runRequest(B, router, "GET", "/ping")
}

func Benchmark404Many(B *testing.B) {
	router := NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {})
	router.Get("/path/to/something", func(w http.ResponseWriter, r *http.Request) {})
	router.Get("/post/{int:id}", func(w http.ResponseWriter, r *http.Request) {})
	router.Get("/view/{int:id}", func(w http.ResponseWriter, r *http.Request) {})
	router.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	router.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {})
	router.Get("/delete/{int:id}", func(w http.ResponseWriter, r *http.Request) {})
	router.Get("/user/{int:id}/{word:mode}", func(w http.ResponseWriter, r *http.Request) {})

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
	runRequest(b, router, "GET", "/v1/anything")

}

func BenchmarkMuxAlternativeInRegexp(b *testing.B) {
	router := NewRouter()
	handler := func(w http.ResponseWriter, r *http.Request) {}
	router.Get("/v1/{v1}", handler)

	runRequest(b, router, "GET", "/v1/a")
	runRequest(b, router, "GET", "/v1/b")
}

func BenchmarkManyPathVariables(b *testing.B) {
	router := NewRouter()
	handler := func(w http.ResponseWriter, r *http.Request) {}
	router.Get("/v1/{v1}/{v2}/{v3}/{v4}/{v5}", handler)

	runRequest(b, router, "GET", "/v1/1/2/3/4/5")
	runRequest(b, router, "GET", "/v1/1/2/3/4")
}

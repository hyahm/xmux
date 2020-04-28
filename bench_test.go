package xmux

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkMux(b *testing.B) {
	router := NewRouter()
	handler := func(w http.ResponseWriter, r *http.Request) {}
	router.Pattern("/v1/{v1}").Get(handler)
	request, _ := http.NewRequest("GET", "/v1/anything", nil)
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(nil, request)
	}
}

func BenchmarkMuxAlternativeInRegexp(b *testing.B) {
	router := NewRouter()
	handler := func(w http.ResponseWriter, r *http.Request) {}
	router.Pattern("/v1/{v1}").Get(handler)

	requestA, _ := http.NewRequest("GET", "/v1/a", nil)
	requestB, _ := http.NewRequest("GET", "/v1/b", nil)
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(nil, requestA)
		router.ServeHTTP(nil, requestB)
	}
}

func BenchmarkManyPathVariables(b *testing.B) {
	fmt.Println(style)
	router := NewRouter()
	handler := func(w http.ResponseWriter, r *http.Request) {}
	router.Pattern("/v1/{v1}/{v2}/{v3}/{v4}/{v5}").Get(handler)

	matchingRequest, _ := http.NewRequest("GET", "/v1/1/2/3/4/5", nil)
	notMatchingRequest, _ := http.NewRequest("GET", "/v1/1/2/3/4", nil)
	recorder := httptest.NewRecorder()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(nil, matchingRequest)
		router.ServeHTTP(recorder, notMatchingRequest)
	}
}

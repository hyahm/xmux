package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hyahm/xmux"
)

func hf(w http.ResponseWriter, r *http.Request) bool {
	fmt.Println("44444444444444444444444444")
	r.Header.Set("name", "cander")
	w.Write([]byte("return"))
	return true
}

func hf1(w http.ResponseWriter, r *http.Request) bool {
	fmt.Println("66666")
	fmt.Println(r.Header.Get("name"))
	return true
}

func TestHome(t *testing.T) {
	router := xmux.NewRouter()
	router.Get("/home/{test}", home).AddModule(hf).SetHeader("name", "cander").AddModule(hf1)
	router.Get("/home", home)
	router.Post("/home", home)
	var a string
	// client := http.Client{}
	r, err := http.NewRequest("GET", "/home/asdf", strings.NewReader(a))
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	t.Log(w.Code)

	t.Log(w.Body.String())
}

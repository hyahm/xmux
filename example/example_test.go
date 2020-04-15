package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hyahm/xmux"
)

func TestHome(t *testing.T) {
	router := xmux.NewRouter()
	router.Pattern("/home").Get(home)
	var a string
	// client := http.Client{}
	r, err := http.NewRequest("GET", "/home", strings.NewReader(a))
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	t.Log(w.Code)

	t.Log(w.Body.String())
}

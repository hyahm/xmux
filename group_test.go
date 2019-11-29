package xmux

import (
	"net/http"
)

func getHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello get"))
	return
}

func oneHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello one"))
	return
}

func twoHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello two"))
	return
}

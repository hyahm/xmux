package xmux

import "net/http"

func Cors(w http.ResponseWriter, r *http.Request) (ok bool) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded,application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Max-Age", "1728000")
	return
}

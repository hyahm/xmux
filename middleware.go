package xmux

import (
	"fmt"
	"net/http"
	"time"
)

func DefaultMidwareTemplate(handle func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	handle(w, r)
	fmt.Printf("url: %s -- addr: %s -- method: %s -- exectime: %f\n", r.URL.Path, r.RemoteAddr, r.Method, time.Since(start).Seconds())
}

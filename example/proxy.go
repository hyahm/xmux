package main

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

func conn(w http.ResponseWriter, r *http.Request) {
	golog.Info("7777")
	hj, ok := w.(http.Hijacker)
	if !ok {
		golog.Error("error")
		return
	}
	lconn, _, err := hj.Hijack()
	if err != nil {
		golog.Error(err)
		return
	}
	_, err = lconn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	if err != nil {
		golog.Error(err)
		return
	}
	golog.Info(r.URL.Host)

	rconn, err := net.Dial("tcp", r.URL.Host)
	if err != nil {
		golog.Error(err)
		return
	}
	go func() {
		io.Copy(lconn, rconn)
	}()

	io.Copy(rconn, lconn)
	// n, err := io.Copy(lconn, rconn)
	// if err != nil {
	// 	golog.Error(err)
	// 	return
	// }
	// golog.Info(n)
}

func tttt(w http.ResponseWriter, r *http.Request) {
	cli := &http.Client{}
	req, err := http.NewRequest(r.Method, r.URL.Path, r.Body)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Info(r.URL.String())
	for k, v := range r.Header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}
	resp, err := cli.Do(req)
	if err != nil {
		golog.Error(err)
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		golog.Error(err)
		return
	}
	w.Write(b)
	resp.Body.Close()
}

func main() {
	router := xmux.NewRouter()

	router.Connect("{all:path}", conn)
	router.All("{all:path}", tttt)
	router.Run(":8990")
}

package xmux

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

func unmarshalError(err error, w http.ResponseWriter, r *http.Request) bool {
	fmt.Println(err)
	return false
}

func notFoundRequireField(key string, w http.ResponseWriter, r *http.Request) bool {
	fmt.Println("required field not found", key)
	return false
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	// w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusNotFound)
}

func handleAll(w http.ResponseWriter, r *http.Request) bool {
	// w.Header().Add("Access-Control-Allow-Origin", "*")
	log.Printf("method: %s\turl: %s\n", r.Method, r.URL.Path)
	return false
}

func HandleConnect(w http.ResponseWriter, r *http.Request) {
	destConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer destConn.Close()

	// 向客户端返回成功响应
	w.WriteHeader(http.StatusOK)

	// 使用 Hijacker 获取客户端的 TCP 连接
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer clientConn.Close()

	// 在客户端和目标服务器之间建立双向隧道
	go transfer(destConn, clientConn)
	transfer(clientConn, destConn)
}

// 数据传输函数
func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func handleOptions(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func handleFavicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
}

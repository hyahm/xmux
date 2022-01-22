package xmux

// import (
// 	"fmt"
// 	"io"
// 	"net"
// 	"net/http"
// 	"strings"
// 	"sync/atomic"
// 	"time"

// 	"github.com/ouqiang/goproxy/cert"
// )

// var tunnelEstablishedResponseLine = []byte("HTTP/1.1 200 Connection established\r\n\r\n")

// const (
// 	// 连接目标服务器超时时间
// 	defaultTargetConnectTimeout = 5 * time.Second
// 	// 目标服务器读写超时时间
// 	defaultTargetReadWriteTimeout = 30 * time.Second
// 	// 客户端读写超时时间
// 	defaultClientReadWriteTimeout = 30 * time.Second
// )

// type Proxy struct {
// 	clientConnNum int32
// 	decryptHTTPS  bool
// 	cert          *cert.Certificate
// 	transport     *http.Transport
// }

// func NewProxy() *Proxy {
// 	return &Proxy{}

// }

// // ServeHTTP 实现了http.Handler接口
// func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Host == "" {
// 		r.URL.Host = r.Host
// 	}
// 	atomic.AddInt32(&p.clientConnNum, 1)
// 	defer func() {
// 		atomic.AddInt32(&p.clientConnNum, -1)
// 	}()
// 	switch {
// 	case r.Method == http.MethodConnect:
// 		p.forwardTunnel(w, r)
// 	default:
// 		p.proxy(w, r)
// 	}
// }

// // var (
// // 	// searcher是协程安全的
// // 	searcher = engine.Engine{}
// // )

// func hijacker(rw http.ResponseWriter) (net.Conn, error) {
// 	hijacker, ok := rw.(http.Hijacker)
// 	if !ok {
// 		return nil, fmt.Errorf("web server不支持Hijacker")
// 	}
// 	conn, _, err := hijacker.Hijack()
// 	if err != nil {
// 		return nil, fmt.Errorf("hijacker错误: %s", err)
// 	}

// 	return conn, nil
// }

// func removeConnectionHeaders(h http.Header) {
// 	if c := h.Get("Connection"); c != "" {
// 		for _, f := range strings.Split(c, ",") {
// 			if f = strings.TrimSpace(f); f != "" {
// 				h.Del(f)
// 			}
// 		}
// 	}
// }

// // CloneHeader 深拷贝Header
// func cloneHeader(h http.Header) http.Header {
// 	h2 := make(http.Header, len(h))
// 	for k, vv := range h {
// 		vv2 := make([]string, len(vv))
// 		copy(vv2, vv)
// 		h2[k] = vv2
// 	}
// 	return h2
// }

// var hopHeaders = []string{
// 	"Connection",
// 	"Proxy-Connection",
// 	"Keep-Alive",
// 	"Proxy-Authenticate",
// 	"Proxy-Authorization",
// 	"Te",
// 	"Trailer",
// 	"Transfer-Encoding",
// 	"Upgrade",
// }

// func (p *Proxy) forwardTunnel(w http.ResponseWriter, r *http.Request) {

// 	clientConn, err := hijacker(w)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadGateway)
// 		return
// 	}
// 	defer clientConn.Close()
// 	parentProxyURL, err := http.ProxyFromEnvironment(r)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadGateway)
// 		return
// 	}
// 	targetAddr := r.URL.Host
// 	if parentProxyURL != nil {
// 		targetAddr = parentProxyURL.Host
// 	}
// 	targetConn, err := net.Dial("tcp", targetAddr)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadGateway)
// 		return
// 	}
// 	defer targetConn.Close()
// 	if parentProxyURL == nil {
// 		_, err = clientConn.Write(tunnelEstablishedResponseLine)
// 		if err != nil {
// 			return
// 		}
// 	} else {
// 		tunnelRequestLine := makeTunnelRequestLine(r.URL.Host)
// 		targetConn.Write([]byte(tunnelRequestLine))
// 	}

// 	p.transfer(clientConn, targetConn)
// }

// func makeTunnelRequestLine(addr string) string {
// 	return fmt.Sprintf("CONNECT %s HTTP/1.1\r\n\r\n", addr)
// }

// func (p *Proxy) transfer(src net.Conn, dst net.Conn) {
// 	go func() {
// 		for {
// 			_, err := io.Copy(src, dst)
// 			if err != nil {
// 				break
// 			}
// 			src.Close()
// 			dst.Close()
// 		}

// 	}()
// 	for {
// 		_, err := io.Copy(dst, src)
// 		if err != nil {
// 			break
// 		}
// 		dst.Close()
// 		src.Close()
// 	}

// }
// func (p *Proxy) proxy(w http.ResponseWriter, r *http.Request) {
// 	newReq := &http.Request{}
// 	*newReq = *r
// 	newReq.Header = cloneHeader(newReq.Header)
// 	removeConnectionHeaders(newReq.Header)
// 	for _, item := range hopHeaders {
// 		if newReq.Header.Get(item) != "" {
// 			newReq.Header.Del(item)
// 		}
// 	}
// 	ts := &http.Transport{}
// 	resp, err := ts.RoundTrip(newReq)
// 	if err == nil {
// 		removeConnectionHeaders(resp.Header)
// 		for _, h := range hopHeaders {
// 			resp.Header.Del(h)
// 		}
// 	}
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadGateway)
// 		return
// 	}
// 	defer resp.Body.Close()
// 	copyHeader(w.Header(), resp.Header)
// 	w.WriteHeader(resp.StatusCode)
// 	io.Copy(w, resp.Body)
// }

// // CopyHeader 浅拷贝Header
// func copyHeader(dst, src http.Header) {
// 	for k, vv := range src {
// 		for _, v := range vv {
// 			dst.Add(k, v)
// 		}
// 	}
// }

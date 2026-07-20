package xmux

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/quic-go/quic-go/http3"
)

// 只是参考， 没有做什么特别的
func QuicGetClient(url string) ([]byte, error) {
	client := http.Client{
		Transport: &http3.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 仅测试用
			},
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// 可选的辅助函数：判断是否是私有内网 IP（防止外部恶意伪造 X-Forwarded-For）
func isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	// 检查是否是内网段: 10.x.x.x, 172.16.x.x-172.31.x.x, 192.168.x.x, 127.0.0.1
	return ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsPrivate()
}

func GetClientIP(r *http.Request) string {
	// 1. 优先尝试 X-Forwarded-For (常用于有多级代理或 CDN 的场景)
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// X-Forwarded-For 的格式可能是: "1.2.3.4, 5.6.7.8, 9.10.11.12"
		// 第一个非空的 IP 就是原始客户端 IP
		ips := strings.Split(xForwardedFor, ",")
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			if ip != "" && !isPrivateIP(ip) { // 可选：过滤掉可能伪造的内网 IP
				return ip
			}
		}
		// 如果全都是内网 IP，退而求其次返回第一个
		if len(ips) > 0 {
			firstIP := strings.TrimSpace(ips[0])
			if firstIP != "" {
				return firstIP
			}
		}
	}

	// 2. 尝试 X-Real-IP (通常是 Nginx 配置的)
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return strings.TrimSpace(xRealIP)
	}

	// 3. 如果是 Cloudflare CDN，还可以通过这个专属 Header 获取
	cfIP := r.Header.Get("CF-Connecting-IP")
	if cfIP != "" {
		return cfIP
	}

	// 4. 最后兜底：直接从 RemoteAddr 获取
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

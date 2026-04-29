package xmux

import (
	"crypto/tls"
	"io"
	"net/http"

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

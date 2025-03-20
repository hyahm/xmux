package xmux

import "testing"

func TestMain(t *testing.T) {
	router := NewRouter()
	router.EnableConnect = true
	router.SetAddr(":9000")
	router.RunTLS("server.crt", "server.key")
}

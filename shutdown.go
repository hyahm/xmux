package xmux

import (
	"fmt"
	"time"
)

// 平滑停止http
func ShutDown() {
	for GetConnections() > 0 {
		time.Sleep(time.Millisecond * 10)
	}
	fmt.Println("server had been safty stoped")
}

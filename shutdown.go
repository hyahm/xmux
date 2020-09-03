package xmux

import (
	"fmt"
	"time"
)

func ShutDown() {
	Stop = true
	for GetConnents() > 0 {
		time.Sleep(time.Millisecond * 10)
	}
	fmt.Println("server had been safty stoped")
}

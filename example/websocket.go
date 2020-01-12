package main

import (
	"fmt"

	"github.com/hyahm/xmux"
)

type player struct {
	*xmux.BaseWs
}

func (p *player) ws() {
	for {
		// 封包
		msgType, msg, err := p.ReadMessage()
		if err != nil {
			fmt.Println(err)
			if err == xmux.ErrorConnect {
				break
			}
		}
		fmt.Println("msg", msg)
		// 发送的msg的长度不能超过 1<<31, 否则掉内容， 建议分包

		p.SendMessage([]byte("g的长度不能超过"), msgType)
	}
}

// func main() {
// 	router := xmux.NewRouter()
// 	player := &player{
// 		xmux.NewWebSocket(),
// 	}
// 	player.Handle = player.ws
// 	router.SetHeader("Access-Control-Allow-Origin", "*")
// 	router.Pattern("/player").WebSocket(player)
// 	//router.Pattern("/ws").WebSocket(ws)
// 	if err := http.ListenAndServe(":7000", router); err != nil {
// 		log.Fatal(err)
// 	}

// }

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hyahm/xmux"
)

// type player struct {
// 	*xmux.BaseWs
// }
var msgchan chan string

func ws(w http.ResponseWriter, r *http.Request) {
	p := xmux.NewWebsocket(w, r)
	p.SendMessage([]byte("hello"), xmux.TypeMsg)
	for {
		// 封包
		msgType, msg, err := p.ReadMessage()
		if err != nil {
			if err == xmux.ErrorConnect {
				// 连接断开
				break
			}
		}
		fmt.Println("msg", msg)
		// 发送的msg的长度不能超过 1<<31, 否则掉内容， 建议分包

		p.SendMessage([]byte(msg), msgType)
	}
}

func main() {
	router := xmux.NewRouter()
	msgchan = make(chan string, 100)
	router.SetHeader("Access-Control-Allow-Origin", "*")
	router.Pattern("/ws").Get(ws)
	if err := http.ListenAndServe(":7000", router); err != nil {
		log.Fatal(err)
	}

}

package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/hyahm/xmux"
)

type client struct {
	msg string
	c   *xmux.BaseWs
}

var msgchan chan client
var wsmu sync.RWMutex
var ps map[*xmux.BaseWs]byte

func sendMsg() {
	for {
		c := <-msgchan
		for p := range ps {
			if c.c == p {
				// 不发给自己
				continue
			}
			fmt.Println(c.msg)
			// 发送的msg的长度不能超过 1<<31, 否则掉内容， 建议分包
			p.SendMessage([]byte(c.msg), ps[p])
		}
	}
}

func ws(w http.ResponseWriter, r *http.Request) {
	websocket, err := xmux.UpgradeWebSocket(w, r)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	websocket.SendMessage([]byte("hello"), xmux.TypeMsg)
	websocket.SendMessage([]byte("hello"), xmux.TypeMsg)
	wsmu.Lock()
	ps[websocket] = xmux.TypeMsg
	wsmu.Unlock()
	tt := time.NewTicker(time.Second * 2)
	go func() {
		for {
			<-tt.C
			if err := websocket.SendMessage([]byte(time.Now().String()), xmux.TypeMsg); err != nil {
				break
			}
		}
	}()
	for {

		if websocket.Conn == nil {
			return
		}
		// 封包
		msgType, msg, err := websocket.ReadMessage()
		if err != nil {
			fmt.Println(err.Error())
			// 连接断开
			wsmu.Lock()
			delete(ps, websocket)
			wsmu.Unlock()
			break
		}
		log.Println(msg)
		websocket.SendMessage([]byte(msg+r.RemoteAddr), msgType)
		ps[websocket] = msgType
		c := client{
			msg: msg,
			c:   websocket,
		}
		msgchan <- c
	}
}

func main() {
	router := xmux.NewRouter()
	wsmu = sync.RWMutex{}
	msgchan = make(chan client, 100)
	ps = make(map[*xmux.BaseWs]byte)
	router.SetHeader("Access-Control-Allow-Origin", "*")
	router.Get("/{int:uid}", ws)

	go sendMsg()
	if err := http.ListenAndServe(":8888", router); err != nil {
		log.Fatal(err)
	}

}

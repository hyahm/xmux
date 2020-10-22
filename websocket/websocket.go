package main

import (
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
			// 发送的msg的长度不能超过 1<<31, 否则掉内容， 建议分包
			if err := p.SendMessage([]byte(c.msg), ps[p]); err != nil {
				return
			}
		}
	}
}

func ws(w http.ResponseWriter, r *http.Request) {
	p, err := xmux.NewWebsocket(w, r)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	p.SendMessage([]byte("hello"), xmux.TypeMsg)
	wsmu.Lock()
	ps[p] = xmux.TypeMsg
	wsmu.Unlock()
	tt := time.NewTicker(time.Second * 2)
	go func() {
		for {
			<-tt.C
			// 发送的msg的长度不能超过 1<<31, 否则掉内容， 建议分包
			if err := p.SendMessage([]byte(time.Now().String()), xmux.TypeMsg); err != nil {
				break
			}
		}
	}()
	for {
		// 封包
		msgType, msg, err := p.ReadMessage()
		if err != nil {
			if err == xmux.ConnectClose || err == xmux.ErrorConnect {
				// 连接断开
				wsmu.Lock()
				delete(ps, p)
				wsmu.Unlock()
				break
			}
		}
		ps[p] = msgType
		c := client{
			msg: msg,
			c:   p,
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
	router.Get("/ws", ws)

	go sendMsg()
	if err := http.ListenAndServe(":7000", router); err != nil {
		log.Fatal(err)
	}

}

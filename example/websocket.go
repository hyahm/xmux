package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"xmux"
)

var ErrorConnect = errors.New("connect error")
var ErrorReadHeader = errors.New("read header error")
var ErrorGetLenth = errors.New("get length error")
var ErrorGetMsg = errors.New("read data error")
var ErrorMsgNotEnough = errors.New("data length not enough")

func ReadMessage(conn net.Conn) (byte, string, error) {
	lpack := make([]byte, 2)
	_, err := io.ReadFull(conn, lpack)
	if err != nil {
		conn.Write([]byte("websocket: client sent data before handshake is complete"))
		return byte(0), "", ErrorConnect
	}
	start := uint64(lpack[0] << 1)
	if start != 1 && start != 2 {
		return byte(0), "", ErrorReadHeader
	}
	var length int32
	playload := int32(lpack[1]) - 128
	if playload < 126 {
		length = playload
	} else if playload == 126 {
		bb := make([]byte, 2)
		bit2 := make([]byte, 2)
		_, err = io.ReadFull(conn, bit2)
		if err != nil {
			return byte(0), "", ErrorGetLenth
		}
		bb = append(bb, bit2...)
		binary.Read(bytes.NewReader(bb), binary.BigEndian, &length)
	} else {
		bit8 := make([]byte, 8)
		_, err = io.ReadFull(conn, bit8)
		if err != nil {
			return byte(0), "", ErrorGetLenth
		}
		binary.Read(bytes.NewReader(bit8), binary.BigEndian, &length)
	}
	mask := make([]byte, 4)
	_, err = io.ReadFull(conn, mask)
	if err != nil {
		return byte(0), "", ErrorGetLenth
	}
	data := make([]byte, length)
	n, err := io.ReadFull(conn, data)
	if err != nil {
		return byte(0), "", ErrorGetMsg
	}
	if int32(n) != length {
		// 数据不对
		return byte(0), "", ErrorMsgNotEnough
	}

	for i, v := range data {
		data[i] = v ^ mask[i%4]
	}
	//fmt.Print("mask", mKey)
	//解包
	return lpack[0], string(data), nil
}

func SendMessage(conn net.Conn, header byte, msg []byte) {
	var send []byte
	send = append(send, header)
	l := len(msg)
	if l < 126 {
		send = append(send, byte(l))
	} else if 126 <= l && l < 1<<16 {
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, int32(l))
		send = append(send, byte(126))
		send = append(send, bytesBuffer.Bytes()[2:4]...)
	} else {
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, uint64(l))
		send = append(send, byte(127))
		send = append(send, bytesBuffer.Bytes()...)
	}
	send = append(send, msg...)
	conn.Write(send)
}

func websocket(w http.ResponseWriter, r *http.Request) {
	// show num of goroutine
	w.Header().Set("Content-Type", "text/plain")

	key := r.Header.Get("Sec-WebSocket-Key")
	if key == "" {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	sha := sha1.New()

	sha.Write([]byte(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	key = base64.StdEncoding.EncodeToString(sha.Sum(nil))
	h, ok := w.(http.Hijacker)
	if !ok {
		w.Write([]byte("websocket: response does not implement http.Hijacker"))
		return
	}
	netConn, brw, err := h.Hijack()
	if err != nil {
		netConn.Write([]byte(err.Error()))
		return
	}

	if brw.Reader.Buffered() > 0 {
		netConn.Write([]byte("websocket: client sent data before handshake is complete"))
		return
	}

	header := "HTTP/1.1 101 Switching Protocols\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Version: 13\r\n" +
		"Sec-WebSocket-Accept: " + key + "\r\n" +
		"Upgrade: websocket\r\n\r\n"
	// 升级为websocket
	netConn.Write([]byte(header))
	quit := make(chan int)
	go ws(netConn, quit)

	select {
	case <-quit:
		netConn.Close()
	}
}

func ws(conn net.Conn, q chan int) {
	for {
		// 封包
		header, msg, err := ReadMessage(conn)
		if err != nil && err == ErrorConnect {
			q <- 1
		}
		fmt.Println(msg)
		// 发送的msg的长度不能超过 1<<31, 否则掉内容， 建议分包
		SendMessage(conn, header, []byte("张三里斯， 利斯"))
	}
}

func main() {
	router := xmux.NewRouter()
	router.SetHeader("Access-Control-Allow-Origin", "*")
	router.Pattern("/ws").Get(websocket)
	if err := http.ListenAndServe(":7000", router); err != nil {
		log.Fatal(err)
	}

}

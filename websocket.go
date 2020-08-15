package xmux

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

var ErrorConnect = errors.New("connect error")
var ErrorType = errors.New("type error")
var ErrorGetLenth = errors.New("get length error")
var ErrorGetMsg = errors.New("read data error")
var ErrorMsgNotEnough = errors.New("data length not enough")

type WsHandler interface {
	Websocket(w http.ResponseWriter, r *http.Request)
}

type BaseWs struct {
	Conn   net.Conn
	Handle func()
	Err    error
}

func NewWebSocket() *BaseWs {
	return &BaseWs{}
}

func (xws *BaseWs) HandleWsFunc() {
	if xws.Handle == nil {
		xws.Err = errors.New("please write a Handle")
		return
	}
	xws.Handle()
}

// 对应的 Type
// x0表示连续消息片断  128
// x1表示文本消息片断//表示传输文本型数据  129
// x2表未二进制消息片断//表示传输Blob以及二进制数据 130
// x3-7为将来非控制消息片断保留地操作码 131-135
// x8表示连接关闭  136
// x9表示心跳检查的ping  137
// xA表示心跳检查的pong  138
// xB-F为将来控制消息片断的保留地操作码

const (
	TypeMsg    = byte(129)
	TypeBinary = byte(130)
	TypeClose  = byte(136)
	TypePing   = byte(137)
	TypePong   = byte(138)
)

func (xws *BaseWs) ReadMessage() (byte, string, error) {
	//解包
	lpack := make([]byte, 2)
	_, err := io.ReadFull(xws.Conn, lpack)
	if err != nil {
		xws.Conn.Write([]byte("websocket: client sent data before handshake is complete"))
		return byte(0), "", ErrorConnect
	}
	if lpack[0] == TypePing {
		xws.SendMessage([]byte(""), TypePong)
		return TypePing, "", nil
	}
	// start := uint64(lpack[0] << 1)
	// if start != 1 && start != 2 {
	// 	return byte(0), "", ErrorType
	// }
	var length int32
	playload := int32(lpack[1]) - 128
	if playload < 126 {
		length = playload
	} else if playload == 126 {
		bb := make([]byte, 2)
		bit2 := make([]byte, 2)
		_, err = io.ReadFull(xws.Conn, bit2)
		if err != nil {
			return lpack[0], "", ErrorGetLenth
		}
		bb = append(bb, bit2...)
		binary.Read(bytes.NewReader(bb), binary.BigEndian, &length)
	} else {
		bit8 := make([]byte, 8)
		_, err = io.ReadFull(xws.Conn, bit8)
		if err != nil {
			return lpack[0], "", ErrorGetLenth
		}
		binary.Read(bytes.NewReader(bit8), binary.BigEndian, &length)
	}

	mask := make([]byte, 4)
	_, err = io.ReadFull(xws.Conn, mask)
	if err != nil {
		return lpack[0], "", ErrorGetLenth
	}
	data := make([]byte, length)
	n, err := io.ReadFull(xws.Conn, data)
	if err != nil {
		return lpack[0], "", ErrorGetMsg
	}

	if int32(n) != length {
		// 数据长度不对
		return lpack[0], "", ErrorMsgNotEnough
	}

	for i, v := range data {
		data[i] = v ^ mask[i%4]
	}

	return lpack[0], string(data), nil
}

func (xws *BaseWs) SendMessage(msg []byte, typ ...byte) {
	var send []byte
	var header byte

	if len(typ) == 0 || header == byte(0) {
		header = TypeMsg
	} else {
		fmt.Println("typ: ", typ[0])
		header = typ[0]
	}
	send = append(send, header)
	l := len(msg)
	if l < 126 {
		send = append(send, byte(l))
	} else if l >= 126 && l < 1<<16 {
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
	xws.Conn.Write(send)
}

func NewWebsocket(w http.ResponseWriter, r *http.Request) (xws *BaseWs) {
	// show num of goroutine
	xws = &BaseWs{}
	w.Header().Set("Content-Type", "text/plain")
	key := r.Header.Get("Sec-WebSocket-Key")
	if key == "" {
		w.WriteHeader(http.StatusBadGateway)
		xws.Err = errors.New("not found Sec-WebSocket-Key")
		return
	}
	sha := sha1.New()

	sha.Write([]byte(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	key = base64.StdEncoding.EncodeToString(sha.Sum(nil))
	h, ok := w.(http.Hijacker)
	if !ok {
		w.Write([]byte("websocket: response does not implement http.Hijacker"))
		xws.Err = errors.New("websocket: response does not implement http.Hijacker")
		return
	}
	netConn, brw, err := h.Hijack()
	if err != nil {
		netConn.Write([]byte(err.Error()))
		xws.Err = err
		return
	}

	if brw.Reader.Buffered() > 0 {
		netConn.Write([]byte("websocket: client sent data before handshake is complete"))
		xws.Err = errors.New("websocket: client sent data before handshake is complete")
		return
	}

	header := "HTTP/1.1 101 Switching Protocols\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Version: 13\r\n" +
		"Sec-WebSocket-Accept: " + key + "\r\n" +
		"Upgrade: websocket\r\n\r\n"
	// 升级为websocket
	netConn.Write([]byte(header))
	xws.Conn = netConn
	go xws.HandleWsFunc()
	return
}

type WebSocket interface {
	WS(*BaseWs)
	SendMessage(msg []byte)
	websocket(w http.ResponseWriter, r *http.Request)
	ReadMessage() (string, error)
}

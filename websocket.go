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

var ErrConnectClosed = errors.New("connect closed")
var ErrorType = errors.New("type error")
var ErrorProtocol = errors.New("protocol undefined")
var ErrorGetLenth = errors.New("get length error")
var ErrorGetMsg = errors.New("read data error")
var ErrorMsgNotEnough = errors.New("data length not enough")
var ErrorNotFoundHandle = errors.New("please write a Handle")
var ErrorRespose = errors.New("websocket: response does not implement http.Hijacker")
var ErrorHandshake = errors.New("websocket: client sent data before handshake is complete")
var ErrorNoWebsocketKey = errors.New("not found Sec-WebSocket-Key")

type WsHandler interface {
	Websocket(w http.ResponseWriter, r *http.Request)
}

type BaseWs struct {
	Conn       net.Conn
	RemoteAddr string
	IsExtras   bool
}

func (bw *BaseWs) Close() error {
	return bw.Conn.Close()
}

// 包结构如下

// Frame format:
// ​​
//       0                   1                   2                   3
//       0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//      +-+-+-+-+-------+-+-------------+-------------------------------+
//      |F|R|R|R| opcode|M| Payload len |    Extended payload length    |
//      |I|S|S|S|  (4)  |A|     (7)     |             (16/64)           |
//      |N|V|V|V|       |S|             |   (if payload len==126/127)   |
//      | |1|2|3|       |K|             |                               |
//      +-+-+-+-+-------+-+-------------+ - - - - - - - - - - - - - - - +
//      |     Extended payload length continued, if payload len == 127  |
//      + - - - - - - - - - - - - - - - +-------------------------------+
//      |                               |Masking-key, if MASK set to 1  |
//      +-------------------------------+-------------------------------+
//      | Masking-key (continued)       |          Payload Data         |
//      +-------------------------------- - - - - - - - - - - - - - - - +
//      :                     Payload Data continued ...                :
//      + - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - +
//      |                     Payload Data continued ...                |
//      +---------------------------------------------------------------+
// 分别解释一下各种字段的意思

// FIN： 数据是否是已结束。
// 当我们发送数据过长的时候，可以分段发送，这时候就需要FIN字段来判断是否发送完全。
// 在浏览器客户端发送大量数据的时候，浏览器会自动帮忙我们截取分段发送，所以不需要我们考虑这个字段。服务端在实现的时候，应当考虑。
// RSV1 RSV2 RSV3： 如果未定义扩展协议，则为0。

// opcode ：操作码

// 0x0表示附加数据帧
// 0x1表示文本数据帧
// 0x2表示二进制数据帧
// 0x3 - 7暂时无定义，为以后的非控制帧保留
// 0x8表示连接关闭
// 0x9表示ping
// 0xA表示pong
// 0xB - F暂时无定义，为以后的控制帧保留

// 这里可以看出，其实本身的websocket中已经定义了ping和pong的操作。
// 遗憾的是，浏览器客户端并没有提供直接操作包的能力，也没有提供类似于ws.ping() ws.pong()的api， 导致我们如果需要从客户端发起ping操作，只能在传输数据中定义。
// 但是如果是从服务端发起的ping，
// 浏览器客户端会自动响应pong操作，且这个交互对于客户端来说是无感知的。
// 所以在使用websocket服务的时候还是尽量使用服务端发起ping操作，这样可以减少传输内容（不需要传数据来表明），且减少客户端的操作。

// Mask： 是否含有掩码
// 一般客户端发送给服务端的数据都需要使用32位掩码进行处理，服务端则不需要。
// Payload len 发送数据的长度。

// 如果有掩码，则需要使用掩码，掩码处理逻辑如下

// const data = payloadData.map((item, index) => {
//  const j = index % 4;
//  return item ^ maskKey[j];
// })

const (
	TypeMsg    = byte(129)
	TypeBinary = byte(130)
	TypeClose  = byte(136)
	TypePing   = byte(137)
	TypePong   = byte(138)
)

func (xws *BaseWs) Ping(ping []byte) error {
	err := xws.SendMessage(ping, TypePing)
	fmt.Println(err)
	return err
}

// TODO: 没有处理附加数据
func (xws *BaseWs) ReadBytes() (byte, []byte, error) {
	//解包
	if xws.Conn == nil {
		return byte(0), nil, ErrConnectClosed
	}
	lpack := make([]byte, 2)
	_, err := io.ReadFull(xws.Conn, lpack)
	if err != nil {
		return byte(0), nil, ErrorGetMsg
	}
	//
	if lpack[0] <= 128 {
		xws.IsExtras = true
		// 这里是附加数据
	} else {
		xws.IsExtras = false
		switch lpack[0] {
		case TypePing:
			xws.SendMessage([]byte(""), TypePong)
			return TypePing, nil, nil
		case TypeClose:
			xws.SendMessage([]byte(""), TypePong)
			xws.Conn.Close()
			xws.Conn = nil
			return 0, nil, ErrConnectClosed
		case TypeMsg, TypeBinary:
			// 正常消息，就继续

		default:
			return 0, nil, ErrorProtocol
		}
	}
	var length int64
	// 长度遵循以下规则
	// 若 0-125，则直接表示Payload长度
	// 若等于126 ，则使用后续的16位的值为Payload长度
	// 若等于127 ，则使用后续的64位作为Payload的长度
	// Masking-key： Mask为1时，Mask的掩码值
	// Payload Data： 发送的真实数据。
	playload := int64(lpack[1]) % 128 // 如果存在减去最高位mask位
	if playload < 126 {
		length = playload
	} else if playload == 126 {
		bb := make([]byte, 2)
		bit2 := make([]byte, 2)
		_, err = io.ReadFull(xws.Conn, bit2)
		if err != nil {
			return lpack[0], nil, ErrorGetLenth
		}
		bb = append(bb, bit2...)
		binary.Read(bytes.NewReader(bb), binary.BigEndian, &length)
	} else {
		bit8 := make([]byte, 8)
		_, err = io.ReadFull(xws.Conn, bit8)
		if err != nil {
			return lpack[0], nil, ErrorGetLenth
		}
		binary.Read(bytes.NewReader(bit8), binary.BigEndian, &length)
	}

	// 读取mask位， 服务器不做处理
	mask := make([]byte, 4)
	_, err = io.ReadFull(xws.Conn, mask)
	if err != nil {
		return lpack[0], nil, ErrorGetLenth
	}
	// 获取数据
	data := make([]byte, length)
	n, err := io.ReadFull(xws.Conn, data)
	if err != nil {
		return lpack[0], nil, ErrorGetMsg
	}

	if int64(n) != length {
		// 数据长度不对
		return lpack[0], nil, ErrorMsgNotEnough
	}

	for i, v := range data {
		data[i] = v ^ mask[i%4]
	}
	return lpack[0], data, nil
}

func (xws *BaseWs) ReadMessage() (byte, string, error) {
	t, d, err := xws.ReadBytes()
	return t, string(d), err
}

func (xws *BaseWs) SendMessage(msg []byte, typ byte) error {
	if xws.Conn == nil {
		return ErrConnectClosed
	}
	var send []byte

	send = append(send, typ)
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
	_, err := xws.Conn.Write(send)
	if err != nil {
		xws.Conn.Close()
		xws.Conn = nil
	}
	return err
}

func UpgradeWebSocket(w http.ResponseWriter, r *http.Request) (*BaseWs, error) {

	// show num of goroutine
	xws := &BaseWs{}
	w.Header().Set("Content-Type", "text/plain")
	key := r.Header.Get("Sec-WebSocket-Key")
	if key == "" {
		w.WriteHeader(http.StatusBadGateway)
		return nil, ErrorNoWebsocketKey
	}
	sha := sha1.New()

	sha.Write([]byte(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	key = base64.StdEncoding.EncodeToString(sha.Sum(nil))
	h, ok := w.(http.Hijacker)
	if !ok {
		w.Write([]byte(ErrorRespose.Error()))
		return nil, ErrorRespose
	}
	netConn, brw, err := h.Hijack()
	if err != nil {
		netConn.Write([]byte(err.Error()))
		return nil, err
	}

	if brw.Reader.Buffered() > 0 {
		return nil, ErrorHandshake
	}

	header := "HTTP/1.1 101 Switching Protocols\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Version: 13\r\n" +
		"Sec-WebSocket-Accept: " + key + "\r\n" +
		"Upgrade: websocket\r\n\r\n"
	// 升级为websocket
	netConn.Write([]byte(header))
	xws.Conn = netConn
	xws.RemoteAddr = r.RemoteAddr

	return xws, nil
}

type WebSocket interface {
	WS(*BaseWs)
	SendMessage(msg []byte)
	websocket(w http.ResponseWriter, r *http.Request)
	ReadMessage() (string, error)
}

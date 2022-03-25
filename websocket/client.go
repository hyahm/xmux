package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/sacOO7/gowebsocket"
)

func main() {
	fmt.Println([]byte("asldjf;ljasld;fj;"))
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	socket := gowebsocket.New("ws://127.0.0.1:8888/23")

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server")
	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Recieved connect error ", err)
	}

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		log.Println("Recieved message " + message)
	}

	socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
		log.Println("Recieved binary data ", data)
	}

	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved ping " + data)
	}

	socket.OnPongReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved pong " + data)
	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from server ")
		return
	}

	socket.Connect()
	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			socket.Close()
			return
		}
	}
}

package main

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/rostekus/go-game-server/pkg/client"
)

const wsServerEndpoint = "ws://localhost:4000/ws"

func main() {
	dialer := websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, _, err := dialer.Dial(wsServerEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	c := client.NewGameClient(conn, "James")
	if err := c.Login(); err != nil {
		log.Fatal(err)
	}
	c.Start()
	select {}

}

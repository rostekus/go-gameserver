package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rostekus/go-game-server/pkg/gameserver"
)

const wsServerEndpoint = "ws://localhost:4000/ws"

type GameClient struct {
	conn     *websocket.Conn
	clientID int
	username string
}

func newGameClient(conn *websocket.Conn, username string) *GameClient {
	return &GameClient{
		clientID: rand.Intn(math.MaxInt),
		username: username,
		conn:     conn,
	}
}

func (c *GameClient) login() error {
	b, err := json.Marshal(gameserver.Login{
		ClientID: c.clientID,
		Username: c.username,
	})
	if err != nil {
		return err
	}
	msg := gameserver.WSMessage{
		Type: "login",
		Data: b,
	}
	return c.conn.WriteJSON(msg)
}

func main() {
	dialer := websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, _, err := dialer.Dial(wsServerEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	c := newGameClient(conn, "James")
	if err := c.login(); err != nil {
		log.Fatal(err)
	}

	go func() {
		var msg gameserver.WSMessage
		for {
			if err := conn.ReadJSON(&msg); err != nil {
				fmt.Println("WS read error", err)
				continue
			}
			switch msg.Type {
			case "state":
				var state gameserver.PlayerState
				if err := json.Unmarshal(msg.Data, &state); err != nil {
					fmt.Println("WS read error", err)
					continue
				}
				fmt.Println("need to update the state of player", state)
			default:
				fmt.Println("receiving message we dont know")
			}
		}
	}()

	for {
		x := rand.Intn(1000)
		y := rand.Intn(1000)
		state := gameserver.PlayerState{
			Health:   100,
			Position: gameserver.Position{X: x, Y: y},
		}
		b, err := json.Marshal(state)
		if err != nil {
			log.Fatal(err)
		}

		msg := gameserver.WSMessage{
			Type: "playerState",
			Data: b,
		}
		if err := conn.WriteJSON(msg); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
	}
}

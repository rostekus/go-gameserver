package client

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

type GameClient struct {
	conn     *websocket.Conn
	clientID int
	username string
}

func NewGameClient(conn *websocket.Conn, username string) *GameClient {
	return &GameClient{
		clientID: rand.Intn(math.MaxInt),
		username: username,
		conn:     conn,
	}
}

func (c *GameClient) Login() error {
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

func (gc *GameClient) poolState() {
	var msg gameserver.WSMessage
	for {
		if err := gc.conn.ReadJSON(&msg); err != nil {
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

}

func (gc *GameClient) sendPos(x, y int) {
	for {
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
		if err := gc.conn.WriteJSON(msg); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
	}

}

func (gc *GameClient) Start() {
	go gc.poolState()
	go gc.sendPos(1, 1)
}

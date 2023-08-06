package client

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rostekus/go-game-server/pkg/game"
	"github.com/rostekus/go-game-server/pkg/gameserver"
)

type GameClient struct {
	conn     *websocket.Conn
	clientID int
	username string
	game     *game.Game
	gameClientConfig
}

type gameClientConfig struct {
	wsServerEndpoint string
	h, w             int
}

func NewGameConfig(wsServerEndpoint string, w, h int) gameClientConfig {
	return gameClientConfig{
		wsServerEndpoint: wsServerEndpoint,
		w:                w,
		h:                h,
	}
}

func NewGameClient(username string, cfg gameClientConfig) *GameClient {
	return &GameClient{
		clientID:         rand.Intn(math.MaxInt),
		username:         username,
		gameClientConfig: cfg,
	}
}

func (c *GameClient) Login() error {
	if err := c.initConnToServer(c.wsServerEndpoint); err != nil {
		log.Println("cannot connect to server")
		return err
	}
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
			gc.game.AddOnlinePlayer(state.Position.X, state.Position.Y, "player")
		}
	}

}

func (gc *GameClient) initGame(w, h int) {

	gc.game = game.NewGame(w, h)
}

func (gc *GameClient) initConnToServer(wsServerEndpoint string) error {

	dialer := websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	//
	conn, _, err := dialer.Dial(wsServerEndpoint, nil)
	if err != nil {
		return err
	}
	gc.conn = conn
	return nil
}

func (gc *GameClient) sendPos() {
	for {
		x, y := gc.game.GetPosSnake()
		state := gameserver.PlayerState{
			Health:   100,
			Position: gameserver.Position{X: x, Y: y},
		}
		b, err := json.Marshal(state)
		if err != nil {
			log.Println(err)
		}

		msg := gameserver.WSMessage{
			Type: "playerState",
			Data: b,
		}
		if err := gc.conn.WriteJSON(msg); err != nil {
			log.Fatal(err)
		}
		// TODO(@rostekus) add channel
		time.Sleep(time.Millisecond * 16)
	}

}

func (gc *GameClient) Start() {
	gc.initGame(gc.w, gc.h)
	go gc.game.Start()
	go gc.poolState()
	go func() {
		//TODO(@rostekus)
		//block, do not send all the time pos
		gc.sendPos()

	}()
}

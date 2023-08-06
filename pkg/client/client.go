package client

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rostekus/go-game-server/pkg/game"
	pr "github.com/rostekus/go-game-server/proto"
	"google.golang.org/protobuf/proto"
)

type GameClient struct {
	conn     *websocket.Conn
	clientID int32
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
		clientID:         int32(rand.Intn(math.MaxInt)),
		username:         username,
		gameClientConfig: cfg,
	}
}

func (c *GameClient) Login() error {
	if err := c.initConnToServer(c.wsServerEndpoint); err != nil {
		log.Println("cannot connect to server")
		return err
	}
	loginData := pr.Login{
		ClientID: c.clientID,
		Username: c.username,
	}
	b, err := proto.Marshal(&loginData)
	if err != nil {
		return err
	}
	msg := pr.WSMessage{
		Type: "login",
		Data: b,
	}
	dataBytes, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.BinaryMessage, dataBytes)
}

func (gc *GameClient) poolState() {
	var msg pr.WSMessage
	for {
		_, dataBytes, err := gc.conn.ReadMessage()
		if err != nil {
			fmt.Println("WS read error", err)
			continue
		}
		err = proto.Unmarshal(dataBytes, &msg)
		if err != nil {
			fmt.Println("WS read error", err)
			continue
		}
		switch msg.Type {
		case "state":
			var state pr.PlayerState
			if err := proto.Unmarshal(msg.Data, &state); err != nil {
				fmt.Println("WS read error", err)
				continue
			}
			gc.game.AddOnlinePlayer(state.Position.Pos, "player")
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
		pos := make([]int32, 2)
		pos[0] = int32(x)
		pos[1] = int32(y)
		position := pr.Position{
			Pos: pos,
		}
		state := pr.PlayerState{
			Position: &position,
		}
		b, err := proto.Marshal(&state)
		if err != nil {
			log.Println(err)
			continue
		}

		msg := pr.WSMessage{
			Type: "playerState",
			Data: b,
		}
		dataBytes, err := proto.Marshal(&msg)
		if err != nil {
			log.Println(err)
			continue
		}
		if err := gc.conn.WriteMessage(websocket.BinaryMessage, dataBytes); err != nil {
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

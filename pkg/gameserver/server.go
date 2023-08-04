package gameserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type PlayerState struct {
	Health   int      `json:"health"`
	Position Position `json:"position"`
}

type WSMessage struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}
type GameServer struct {
	GameServerConfig
}

type GameServerConfig struct {
	Port   string
	Logger *logrus.Logger
}

func DefaultGameServer() actor.Receiver {
	logger := logrus.New()

	return &GameServer{
		GameServerConfig{
			Port:   ":4000",
			Logger: logger,
		},
	}
}

func (s *GameServer) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		s.Logger.Println("GameServer has started", msg)
		s.startHTTPServer()
	}

}

func (s *GameServer) handleWS(w http.ResponseWriter, r *http.Request) {
	s.Logger.Println("New Connection")
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		s.Logger.Errorf("Cannot upgrade connetion: %v\n", err)
	}
	msg := []byte("hello")
	conn.WriteMessage(websocket.BinaryMessage, msg)
	for {
		var msg WSMessage
		for {
			if err := conn.ReadJSON(&msg); err != nil {
				fmt.Println("read error", err)
				return
			}
			var pos PlayerState
			json.Unmarshal(msg.Data, &pos)
			fmt.Printf("pos %+v\n", pos)
		}
	}
}

func (s *GameServer) startHTTPServer() {
	s.Logger.Printf("Listening on port: %s", s.Port)
	go func() {
		http.HandleFunc("/ws", s.handleWS)
		log.Panic(http.ListenAndServe(s.Port, nil))
	}()
}

package gameserver

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"

	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

type GameServer struct {
	GameServerConfig
	ctx      *actor.Context
	sessions map[*actor.PID]struct{}
}

type GameServerConfig struct {
	Port string
}

func DefaultGameServer() actor.Receiver {

	return &GameServer{
		GameServerConfig: GameServerConfig{
			Port: ":4000",
		}, sessions: make(map[*actor.PID]struct{}),
	}
}

func (s *GameServer) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		s.ctx = ctx
		_ = msg
		s.startHTTPServer()
	}

}

func (s *GameServer) handleWS(w http.ResponseWriter, r *http.Request) {
	logger.Println("New Connection")
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		logger.Errorf("Cannot upgrade connetion: %v\n", err)
	}
	sid := rand.Intn(math.MaxInt)
	pid := s.ctx.SpawnChild(newPlayerSession(sid, conn), fmt.Sprintf("session_%d", sid))
	s.sessions[pid] = struct{}{}
}

func (s *GameServer) startHTTPServer() {
	logger.Printf("Listening on port: %s", s.Port)
	go func() {
		http.HandleFunc("/ws", s.handleWS)
		log.Panic(http.ListenAndServe(s.Port, nil))
	}()
}

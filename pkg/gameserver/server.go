package gameserver

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"

	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

type GameServerConfig struct {
	Port string
}

type GameServer struct {
	GameServerConfig
	ctx      *actor.Context
	sessions map[int]*actor.PID
}

func DefaultGameServer() actor.Receiver {
	return &GameServer{
		GameServerConfig: GameServerConfig{
			Port: ":4000"},
		sessions: make(map[int]*actor.PID),
	}
}

func (s *GameServer) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case *PlayerState:
		s.bcast(c.Sender(), msg)
	case actor.Started:
		s.startHTTPServer()
		s.ctx = c
		_ = msg
	default:
		logger.Println("recv", msg)
	}
}

func (s *GameServer) bcast(from *actor.PID, state *PlayerState) {
	for _, pid := range s.sessions {
		if !pid.Equals(from) {
			s.ctx.Send(pid, state)
		}
	}
}

func (s *GameServer) startHTTPServer() {
	logger.Printf("Listening on port: %s", s.Port)
	go func() {
		http.HandleFunc("/ws", s.handleWS)
		http.ListenAndServe(s.Port, nil)
	}()
}

// handles the upgrade of the websocket
func (s *GameServer) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		logger.Errorf("Cannot upgrade connetion: %v\n", err)
		return
	}
	logger.Println("New Connection")
	sid := rand.Intn(math.MaxInt)
	pid := s.ctx.SpawnChild(newPlayerSession(s.ctx.PID(), sid, conn), fmt.Sprintf("session_%d", sid))
	s.sessions[sid] = pid
}

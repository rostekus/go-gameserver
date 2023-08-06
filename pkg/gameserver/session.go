package gameserver

import (
	"encoding/json"
	"fmt"

	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
)

type PlayerSession struct {
	sessionID int
	clientID  int
	username  string
	inLobby   bool
	conn      *websocket.Conn
	ctx       *actor.Context
	serverPID *actor.PID
}

func newPlayerSession(serverPID *actor.PID, sid int, conn *websocket.Conn) actor.Producer {
	return func() actor.Receiver {
		return &PlayerSession{
			conn:      conn,
			sessionID: sid,
			serverPID: serverPID,
		}
	}
}

func (s *PlayerSession) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Started:
		s.ctx = c
		go s.readLoop()
	case *PlayerState:
		s.sendPlayerState(msg)
	default:
		fmt.Println("recv", msg)
	}
}

func (s *PlayerSession) sendPlayerState(state *PlayerState) {
	b, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	msg := WSMessage{
		Type: "state",
		Data: b,
	}
	if err := s.conn.WriteJSON(msg); err != nil {
		panic(err)
	}
}

func (s *PlayerSession) readLoop() {
	var msg WSMessage
	for {
		if err := s.conn.ReadJSON(&msg); err != nil {
			fmt.Println("read error", err)
			return
		}
		go s.handleMessage(msg)
	}
}

func (s *PlayerSession) handleMessage(msg WSMessage) {
	switch msg.Type {
	case "login":
		var loginMsg Login
		if err := json.Unmarshal(msg.Data, &loginMsg); err != nil {
			panic(err)
		}
		s.clientID = loginMsg.ClientID
		s.username = loginMsg.Username
	case "playerState":
		var ps PlayerState
		if err := json.Unmarshal(msg.Data, &ps); err != nil {
			panic(err)
		}
		ps.SessionID = s.sessionID
		if s.ctx != nil {
			s.ctx.Send(s.serverPID, &ps)
		}
	}
}

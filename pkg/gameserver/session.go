package gameserver

import (
	"fmt"

	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
	pr "github.com/rostekus/go-game-server/proto"
	"google.golang.org/protobuf/proto"
)

type PlayerSession struct {
	sessionID int32
	clientID  int32
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
			sessionID: int32(sid),
			serverPID: serverPID,
		}
	}
}

func (s *PlayerSession) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Started:
		s.ctx = c
		go s.readLoop()
	case *pr.PlayerState:
		s.sendPlayerState(msg)
	default:
		fmt.Println("recv", msg)
	}
}

func (s *PlayerSession) sendPlayerState(state *pr.PlayerState) {
	b, err := proto.Marshal(state)
	if err != nil {
		panic(err)
	}
	msg := pr.WSMessage{
		Type: "state",
		Data: b,
	}
	dataBytes, err := proto.Marshal(&msg)
	if err != nil {
		logger.Errorln("cannot encode msg")
	}
	if err := s.conn.WriteMessage(websocket.BinaryMessage, dataBytes); err != nil {
		panic(err)
	}
}

func (s *PlayerSession) readLoop() {
	var msg pr.WSMessage
	for {
		_, dataBytes, err := s.conn.ReadMessage()
		if err != nil {
			logger.Println("read error", err)
			return
		}
		if err := proto.Unmarshal(dataBytes, &msg); err != nil {
			logger.Println("unmarshal error", err)
		}
		go s.handleMessage(&msg)
	}
}

func (s *PlayerSession) handleMessage(msg *pr.WSMessage) {
	switch msg.Type {
	case "login":
		var loginMsg pr.Login
		if err := proto.Unmarshal(msg.Data, &loginMsg); err != nil {
			panic(err)
		}
		s.clientID = loginMsg.ClientID
		s.username = loginMsg.Username
	case "playerState":
		var ps pr.PlayerState
		if err := proto.Unmarshal(msg.Data, &ps); err != nil {
			panic(err)
		}
		ps.SessionID = s.sessionID
		if s.ctx != nil {
			s.ctx.Send(s.serverPID, &ps)
		}
	}
}

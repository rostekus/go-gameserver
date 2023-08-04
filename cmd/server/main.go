package main

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/rostekus/go-game-server/pkg/gameserver"
)

func main() {
	e := actor.NewEngine()
	_ = e.Spawn(gameserver.DefaultGameServer, "game_server")
	select {}
}

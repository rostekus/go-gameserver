package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nsf/termbox-go"
	"github.com/rostekus/go-game-server/pkg/client"
)

const wsServerEndpoint = "ws://localhost:4000/ws"

func main() {
	err := termbox.Init()
	if err != nil {
		log.Fatalf("failed to init termbox: %v", err)
	}
	defer termbox.Close()
	w, h := termbox.Size()

	cfg := client.NewGameConfig(wsServerEndpoint, w, h)
	c := client.NewGameClient("James", cfg)
	if err := c.Login(); err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	c.Start()

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	<-interruptChan

	fmt.Println("Received interrupt signal. Exiting gracefully...")

	os.Exit(0)
}

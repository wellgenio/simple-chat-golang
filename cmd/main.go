package main

import (
	"github.com/wellgenio/simple-chat-golang/internal/ws"
	"github.com/wellgenio/simple-chat-golang/router"
)

func main() {
	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub)

	go hub.Run()

	router.InitRouter(wsHandler)
	router.Start("0.0.0.0:8080")
}

package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader *websocket.Upgrader
}

func (server *Server) handle(w http.ResponseWriter, r *http.Request) {
	conn, err := server.upgrader.Upgrade(w, r, nil)

	if err != nil {
		return
	}
	defer conn.Close()

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if err := conn.WriteMessage(messageType, data); err != nil {
			return
		}
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const MaxConnections int64 = 128

type Server struct {
	ctx         context.Context
	upgrader    *websocket.Upgrader
	presets     map[string][]byte
	connections atomic.Int64
}

func newServer() *Server {
	return &Server{
		ctx:      context.Background(),
		upgrader: &websocket.Upgrader{HandshakeTimeout: 10 * time.Second},
		presets:  makePresets(),
	}
}

func main() {
	server := newServer()

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/", http.FileServer(http.Dir("web")))
	mux.HandleFunc("/shell", server.websocket)

	httpServer := &http.Server{
		Addr:              ":4242",
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("portfolio available at http://localhost:4242")
	log.Fatal(httpServer.ListenAndServe())
}

func makePresets() map[string][]byte {
	presets := make(map[string][]byte)

	projects, err := projectsAscii()
	if err != nil {
		panic(fmt.Sprintf("failed to load ascii file: %s", err.Error()))
	}
	presets[AboutCommand] = profileAscii()
	presets[ContactCommand] = contactAscii()
	presets[ProjectsCommand] = []byte(projects)
	presets[HelpCommand] = helpMenu()
	return presets
}

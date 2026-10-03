package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	ctx      context.Context
	upgrader *websocket.Upgrader
	presets  map[string][]byte
}

func newServer() *Server {
	return &Server{
		ctx:      context.Background(),
		upgrader: &websocket.Upgrader{},
		presets:  makePresets(),
	}
}

func main() {
	server := newServer()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/", http.FileServer(http.Dir("web")))
	http.HandleFunc("/shell", server.websocket)

	log.Println("portfolio available at http://localhost:4242")
	log.Fatal(http.ListenAndServe(":4242", nil))
}

func makePresets() map[string][]byte {
	presets := make(map[string][]byte)

	projects, err := projectsAscii()
	if err != nil {
		panic(fmt.Sprintf("failed to load ascii file: %s", err.Error()))
	}
	presets[AboutCommand] = profileAscii()
	presets[ProjectsCommand] = []byte(projects)
	presets[HelpCommand] = helpMenu()
	return presets
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader *websocket.Upgrader
	presets  map[string][]byte
}

const (
	About    string = "about"
	Projects string = "projects"
)

func newServer() *Server {
	return &Server{
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

	ascii, err := os.ReadFile("static/mr_robot/profile.ascii")
	if err != nil {
		panic(fmt.Sprintf("failed to load ascii file: %s", err.Error()))
	}
	projects, err := projectsAscii()
	if err != nil {
		panic(fmt.Sprintf("failed to load ascii file: %s", err.Error()))
	}
	presets[AboutCommand] = ascii
	presets[ProjectsCommand] = []byte(projects)
	return presets
}

func projectsAscii() (string, error) {
	ascii := strings.Builder{}

	b, err := os.ReadFile("static/banners/perkins.ascii")
	if err != nil {
		return "", err
	}
	ascii.WriteString(string(b))

	b, err = os.ReadFile("static/banners/tinyrenderer.ascii")
	if err != nil {
		return "", err
	}
	ascii.WriteString("\n")
	ascii.WriteString(string(b))

	b, err = os.ReadFile("static/banners/wireframe.ascii")
	if err != nil {
		return "", err
	}
	ascii.WriteString("\n")
	ascii.WriteString(string(b))
	return ascii.String(), nil
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Mr-Robot-err-404/portfolio/pkg/ascii"
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
	presets[HelpCommand] = helpMenu()
	return presets
}

func helpMenu() []byte {
	menu := ascii.Table([]ascii.Stat{
		{Key: "about", Value: "Meet the person behind the terminal"},
		{Key: "contact", Value: "Find me elsewhere"},
		{Key: "projects", Value: "Explore what I've built"},
		{Key: "clear", Value: "Clear the terminal"},
	}, 60, ascii.TableStyle{
		Primary:   ascii.Amber,
		Secondary: ascii.Ocean,
		Border:    ascii.Ocean,
	})
	menu = appendLine(menu)
	menu = appendLine(menu)

	add := []byte(fmt.Sprintf(
		"      %s -> %s",
		ascii.Color("shell", ascii.Amber),
		ascii.Color("connect to a real unix shell!", ascii.Amber),
	))
	menu = append(menu, add...)
	return menu
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

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

const (
	InputEvent  string = "input"
	ResizeEvent string = "resize"
)

const (
	AboutCommand string = "about"
	ClearCommand string = "clear"
)

const (
	Green    string = "\x1b[32m"
	Blue     string = "\x1b[34m"
	Reset    string = "\x1b[0m"
	Shell    string = "❯"
	NextLine string = "\n\r"
	Clear    string = "\x1b[2J\x1b[H"
)

type ClientMessage struct {
	Event   string  `json:"event"`
	Command *string `json:"command"`
	Resize  *Resize `json:"resize"`
}
type Resize struct {
	Rows int `json:"rows"`
	Cols int `json:"cols"`
}

func (server *Server) parseInput(input *string) []byte {
	if input == nil {
		return []byte("Received corrupt payload")
	}
	command := strings.TrimSpace(strings.ToLower(*input))
	if response, ok := server.presets[command]; ok {
		return response
	}
	switch command {
	case ClearCommand:
		return clearAll()
	default:
		return unknown(command)
	}
}

func flush(conn *websocket.Conn, response []byte) {
	if err := conn.WriteMessage(websocket.TextMessage, response); err != nil {
		fmt.Println(err)
		return
	}
}

func (server *Server) websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := server.upgrader.Upgrade(w, r, nil)

	if err != nil {
		return
	}
	defer conn.Close()
	flush(conn, startup())

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if messageType != websocket.TextMessage {
			fmt.Printf("received unexpected messageType: %d", messageType)
			continue
		}
		var message ClientMessage

		if err = json.Unmarshal(data, &message); err != nil {
			fmt.Printf("unmarshal failure: %s", data)
			continue
		}
		var response []byte

		switch message.Event {
		case ResizeEvent:
		case InputEvent:
			response = server.parseInput(message.Command)
			flush(conn, response)
		}
	}
}

func startup() []byte {
	return fmt.Appendf(nil, "%s\n\n\r%s", color("PORTFOLIO / SYSTEM ONLINE", Green), prompt())
}
func clearAll() []byte {
	return fmt.Appendf(nil, "%s%s", Clear, prompt())
}
func color(msg string, color string) string {
	return fmt.Sprintf("%s%s%s", color, msg, Reset)
}
func prompt() []byte {
	user := color("visitor@portfolio", Green)
	return fmt.Appendf(nil, "%s:~ %s ", user, color(Shell, Blue))
}
func unknown(command string) []byte {
	return fmt.Appendf(nil, "%sCommand not found: %s%s%s", NextLine, command, NextLine, prompt())
}

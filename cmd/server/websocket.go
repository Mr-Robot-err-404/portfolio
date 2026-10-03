package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Mr-Robot-err-404/portfolio/pkg/ascii"
	"github.com/gorilla/websocket"
)

const (
	InputEvent    string = "input"
	ResizeEvent   string = "resize"
	KeypressEvent string = "keypress"
)

const (
	AboutCommand    string = "about"
	ClearCommand    string = "clear"
	HelpCommand     string = "help"
	ProjectsCommand string = "projects"
)
const (
	UpArrow   string = "up"
	DownArrow string = "down"
	Backspace string = "backspace"
)

type ClientMessage struct {
	Event    string  `json:"event"`
	Keypress *string `json:"keypress"`
	Command  *string `json:"command"`
	Resize   *Resize `json:"resize"`
}
type Resize struct {
	Rows int `json:"rows"`
	Cols int `json:"cols"`
}

func (server *Server) websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := server.upgrader.Upgrade(w, r, nil)

	if err != nil {
		return
	}
	defer conn.Close()
	flush(conn, server.startup())

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
		switch message.Event {
		case ResizeEvent:
		case KeypressEvent:
		case InputEvent:
			flush(conn, server.parseInput(message.Command))
		}
	}
}

func (server *Server) parseInput(input *string) []byte {
	if input == nil {
		return []byte("Received corrupt payload")
	}
	command := strings.TrimSpace(strings.ToLower(*input))

	if len(command) == 0 {
		return []byte{}
	}
	if response, ok := server.presets[command]; ok {
		return batch(sandwich(response), ascii.Prompt())
	}
	if isShellCommand(command) {
		return batch(sandwich(connectHint(command)), ascii.Prompt())
	}
	switch command {
	case ClearCommand:
		return ascii.ClearAll()
	default:
		return ascii.Unknown(command)
	}
}

func flush(conn *websocket.Conn, response []byte) {
	if err := conn.WriteMessage(websocket.BinaryMessage, response); err != nil {
		fmt.Println(err)
		return
	}
}

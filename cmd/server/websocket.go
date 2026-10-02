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
type History struct {
	list []string
	idx  int
}

func (server *Server) websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := server.upgrader.Upgrade(w, r, nil)

	if err != nil {
		return
	}
	defer conn.Close()
	flush(conn, server.startup())

	history := History{idx: -1}

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
			flush(conn, server.parseKeypress(message.Keypress, &history))
		case InputEvent:
			history.append(message.Command)
			flush(conn, server.parseInput(message.Command))
		}
	}
}

func (server *Server) parseKeypress(input *string, history *History) []byte {
	if input == nil {
		return []byte("Received corrupt payload")
	}
	keypress := strings.TrimSpace(strings.ToLower(*input))

	switch keypress {
	case UpArrow:
		command := history.previous()
		if len(command) == 0 {
			return nil
		}
		return batch(ascii.ClearLine(), []byte(command))
	case DownArrow:
		command := history.next()
		if len(command) == 0 {
			return nil
		}
		return batch(ascii.ClearLine(), []byte(command))
	default:
		return nil
	}
}

func (h *History) append(command *string) {
	if command == nil || len(*command) == 0 {
		return
	}
	if len(h.list) != 0 && h.list[len(h.list)-1] == *command {
		return
	}
	h.list = append(h.list, *command)
	h.idx++
}

func (h *History) previous() string {
	if len(h.list) == 0 {
		return ""
	}
	h.idx = max(0, h.idx-1)
	return h.list[h.idx]
}

func (h *History) next() string {
	if len(h.list) == 0 {
		return ""
	}
	h.idx = min(len(h.list)-1, h.idx+1)
	return h.list[h.idx]
}

func (server *Server) parseInput(input *string) []byte {
	if input == nil {
		return []byte("Received corrupt payload")
	}
	command := strings.TrimSpace(strings.ToLower(*input))
	if response, ok := server.presets[command]; ok {
		return batch(sandwich(response), ascii.Prompt())
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

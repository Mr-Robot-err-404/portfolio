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
	InputEvent  string = "input"
	ResizeEvent string = "resize"
)

const (
	AboutCommand string = "about"
	ClearCommand string = "clear"
	StatsCommand string = "stats"
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

func (server *Server) websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := server.upgrader.Upgrade(w, r, nil)

	if err != nil {
		return
	}
	defer conn.Close()
	flush(conn, ascii.Startup())

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

func batch(payloads ...[]byte) []byte {
	var total int

	for _, payload := range payloads {
		total += len(payload)
	}
	result := make([]byte, 0, total)

	for _, payload := range payloads {
		result = append(result, payload...)
	}
	return result
}
func prependLine(payload []byte) []byte {
	return append([]byte("\n"), payload...)
}
func appendLine(payload []byte) []byte {
	return append(payload, []byte("\n")...)
}
func sandwich(payload []byte) []byte {
	return appendLine(prependLine(payload))
}

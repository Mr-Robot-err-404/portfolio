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
	ConnectCommand  string = "connect"
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
type Dimensions struct {
	width  int
	height int
}
type ClientState struct {
	shell      *Shell
	conn       *websocket.Conn
	dimensions Dimensions
}

func (server *Server) websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := server.upgrader.Upgrade(w, r, nil)

	if err != nil {
		return
	}
	defer conn.Close()
	flush(conn, server.startup())

	client := ClientState{conn: conn}
	defer client.cleanup()

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
			if message.Resize == nil {
				continue
			}
			client.dimensions.width = message.Resize.Cols
			client.dimensions.height = message.Resize.Rows

		case KeypressEvent:
			if client.shell == nil {
				continue
			}
			if message.Keypress == nil {
				continue
			}
			client.shell.write([]byte(*message.Keypress))

		case InputEvent:
			flush(conn, server.parseInput(message.Command, &client))
		}
	}
}

func (server *Server) parseInput(input *string, client *ClientState) []byte {
	if input == nil {
		return []byte("Received corrupt payload")
	}
	if client.shell != nil {
		return batch(sandwich([]byte("Use Keypress events for shell sessions")), ascii.Prompt())
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
	case ConnectCommand:
		flush(client.conn, sandwich([]byte("connecting")))

		shell, err := spawnShell()
		if err != nil {
			fmt.Println(err)
			return nil
		}
		flush(client.conn, batch(ascii.ClearLine(), appendLine([]byte("connected"))))

		go shell.receiveShellOutput(client.conn)
		client.shell = shell
		return nil
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

func (client *ClientState) cleanup() {
	if client.shell == nil {
		return
	}
	try(client.shell.removeContainer)
}

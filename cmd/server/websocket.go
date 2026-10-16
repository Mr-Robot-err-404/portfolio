package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Mr-Robot-err-404/portfolio/pkg/ascii"
	"github.com/gorilla/websocket"
)

const (
	ReadLimit    int64         = 4096
	PongWait     time.Duration = 90 * time.Second
	PingInterval time.Duration = 30 * time.Second
	WriteWait    time.Duration = 10 * time.Second
	SessionLimit time.Duration = 15 * time.Minute
)

const (
	InputEvent    string = "input"
	ResizeEvent   string = "resize"
	KeypressEvent string = "keypress"
)

const (
	AboutCommand    string = "about"
	ClearCommand    string = "clear"
	ContactCommand  string = "contact"
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
	writeMu    sync.Mutex
	id         uint64
	expired    atomic.Bool
	commands   atomic.Uint64
}

func (client *ClientState) closeReason() string {
	if client.expired.Load() {
		return "timeout"
	}
	return "disconnected"
}

var session atomic.Uint64

func nextSessionID() uint64 {
	return session.Add(1)
}

func clientAddr(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		if comma := strings.IndexByte(forwarded, ','); comma != -1 {
			return strings.TrimSpace(forwarded[:comma])
		}
		return strings.TrimSpace(forwarded)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func (server *Server) websocket(w http.ResponseWriter, r *http.Request) {
	addr := clientAddr(r)

	if server.connections.Add(1) > MaxConnections {
		server.connections.Add(-1)
		slog.Warn("connection rejected", "addr", addr, "reason", "limit reached")
		http.Error(w, "too many connections", http.StatusServiceUnavailable)
		return
	}
	defer server.connections.Add(-1)

	if !server.acquire(addr) {
		slog.Warn("connection rejected", "addr", addr, "reason", "per ip limit reached")
		http.Error(w, "too many connections", http.StatusTooManyRequests)
		return
	}
	defer server.release(addr)

	conn, err := server.upgrader.Upgrade(w, r, nil)

	if err != nil {
		slog.Warn("upgrade failed", "addr", addr, "error", err)
		return
	}
	defer conn.Close()

	id := nextSessionID()
	started := time.Now()

	slog.Info("session opened",
		"id", id,
		"addr", addr,
		"active", server.connections.Load(),
	)

	client := ClientState{conn: conn, id: id}
	defer client.cleanup()

	conn.SetReadLimit(ReadLimit)
	conn.SetReadDeadline(time.Now().Add(PongWait))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(PongWait))
	})

	done := make(chan struct{})
	defer close(done)
	go client.keepalive(done)

	deadline := time.After(SessionLimit)
	go func() {
		select {
		case <-deadline:
			client.expired.Store(true)
			conn.Close()
		case <-done:
		}
	}()

	defer func() {
		slog.Info("session closed",
			"id", id,
			"addr", addr,
			"reason", client.closeReason(),
			"duration", time.Since(started).Round(time.Second).String(),
			"commands", client.commands.Load(),
		)
	}()

	client.flush(server.startup())

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if messageType != websocket.TextMessage {
			slog.Debug("unexpected message type", "id", id, "type", messageType)
			continue
		}
		var message ClientMessage

		if err = json.Unmarshal(data, &message); err != nil {
			slog.Debug("unmarshal failure", "id", id, "bytes", len(data))
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
			client.commands.Add(1)
			client.flush(server.parseInput(message.Command, &client))
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
		// client.flush(sandwich([]byte("connecting")))
		//
		// shell, err := spawnShell(client.dimensions)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return nil
		// }
		// client.flush(batch(ascii.ClearLine(), appendLine([]byte("connected"))))
		//
		// go shell.receiveShellOutput(client)
		// client.shell = shell
		// return nil
		msg := fmt.Appendf(nil, "Real shell sessions coming soon, with a CTF challenge included!")
		return batch(sandwich(msg), ascii.Prompt())
	default:
		return ascii.Unknown(command)
	}
}

func (client *ClientState) flush(response []byte) {
	client.writeMu.Lock()
	defer client.writeMu.Unlock()

	client.conn.SetWriteDeadline(time.Now().Add(WriteWait))

	if err := client.conn.WriteMessage(websocket.BinaryMessage, response); err != nil {
		slog.Debug("write failed", "id", client.id, "error", err)
		return
	}
}

func (client *ClientState) keepalive(done <-chan struct{}) {
	ticker := time.NewTicker(PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			client.writeMu.Lock()
			client.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			err := client.conn.WriteMessage(websocket.PingMessage, nil)
			client.writeMu.Unlock()

			if err != nil {
				return
			}
		}
	}
}

func (client *ClientState) cleanup() {
	if client.shell == nil {
		return
	}
	try(client.shell.removeContainer)
}

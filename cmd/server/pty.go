package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

func receiveShellOutput(ptmx *os.File, conn *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ptmx.Read(buf)
		if err != nil {
			return
		}
		conn.WriteMessage(websocket.BinaryMessage, buf[:n])
	}
}

func writeToShell(ctx context.Context, ptmx *os.File, ch <-chan []byte) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			_, err := ptmx.Write(msg)
			if err != nil {
				fmt.Println("failed to write ptmx: %w", err)
			}
		}
	}
}

func spawnShell() (*os.File, error) {
	cmd := exec.Command("podman", "run", "--rm", "-i",
		"--network=none",
		"--memory=64m",
		"--cpus=0.5",
		"--pids-limit=50",
		"--cap-drop=ALL",
		"--security-opt", "no-new-privileges",
		"--read-only",
		"--tmpfs", "/tmp:rw,size=16m,mode=1777",
		"--hostname", "sandbox",
		"--user", "nobody",
		"alpine:latest",
		"/bin/sh",
	)
	return pty.Start(cmd)
}

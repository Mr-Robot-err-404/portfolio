package main

import (
	"fmt"
	"os"
	"os/exec"
	"sync/atomic"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

type Shell struct {
	ptmx *os.File
	cmd  *exec.Cmd
	name string
}

func (shell *Shell) receiveShellOutput(conn *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := shell.ptmx.Read(buf)
		if err != nil {
			return
		}
		conn.WriteMessage(websocket.BinaryMessage, buf[:n])
	}
}

func (shell *Shell) write(chunk []byte) error {
	_, err := shell.ptmx.Write(chunk)
	return err
}

func (shell *Shell) removeContainer() error {
	try(shell.ptmx.Close)
	cmd := exec.Command("podman", "rm", "--force", shell.name)
	return cmd.Run()
}

func spawnShell(dimensions Dimensions) (*Shell, error) {
	name := nextContainerName()

	cmd := exec.Command("podman", "run", "--rm", "-it",
		fmt.Sprintf("--name=%s", name),
		"--network=none",
		"--memory=64m",
		"--cpus=0.5",
		"--pids-limit=50",
		"--cap-drop=ALL",
		"--security-opt", "no-new-privileges",
		"--read-only",
		"--pull=never",
		"--tmpfs",
		"/tmp:rw,noexec,nosuid,nodev,size=16m,mode=1777",
		"--hostname", "sandbox",
		"--user", "nobody",
		"alpine:latest",
		"/bin/sh",
	)
	ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{
		Rows: uint16(dimensions.height),
		Cols: uint16(dimensions.width),
	})
	if err != nil {
		return nil, err
	}
	return &Shell{
		ptmx: ptmx,
		cmd:  cmd,
		name: name,
	}, nil
}

var container atomic.Uint64

func nextContainerName() string {
	return fmt.Sprintf(
		"visitor_%d_%d",
		os.Getpid(),
		container.Add(1),
	)
}

func try(fn func() error) {
	if err := fn(); err != nil {
		fmt.Println(err)
	}
}

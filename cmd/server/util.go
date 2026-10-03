package main

import (
	"fmt"
	"strings"

	"github.com/Mr-Robot-err-404/portfolio/pkg/ascii"
)

var ShellCommands = map[string]struct{}{
	"ls": {}, "cat": {}, "pwd": {}, "cd": {},
	"echo": {}, "grep": {}, "find": {}, "head": {},
	"tail": {}, "wc": {}, "sort": {}, "uniq": {},
	"mkdir": {}, "touch": {}, "rm": {}, "cp": {},
	"mv": {}, "chmod": {}, "whoami": {}, "date": {},
	"ps": {}, "top": {}, "vi": {}, "vim": {},
}

func connectHint(command string) []byte {
	var b []byte
	b = fmt.Appendf(
		b,
		"'%s' requires a shell session.\n Run %s to start your own shell!",
		command,
		ascii.ColorWithAnsi("connect", ascii.Amber),
	)
	return b
}

func isShellCommand(command string) bool {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return false
	}
	_, ok := ShellCommands[fields[0]]
	return ok
}

func (server *Server) startup() []byte {
	return fmt.Appendf(
		nil,
		"%s\n\r%s%s",
		ascii.Color("PORTFOLIO / SYSTEM ONLINE", ascii.Green),
		server.presets[HelpCommand],
		ascii.Prompt(),
	)
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
func title(msg string) []byte {
	p := []byte{}
	p = appendLine(p)
	p = append(p, []byte(msg)...)
	p = appendLine(p)
	return p
}
func prependLine(payload []byte) []byte {
	return append([]byte("\n"), payload...)
}
func appendLine(payload []byte) []byte {
	return append(payload, '\n')
}
func sandwich(payload []byte) []byte {
	return appendLine(prependLine(payload))
}

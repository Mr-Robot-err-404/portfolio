package main

import "bytes"

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
func addNote(p []byte, note []byte) []byte {
	p = append(p, bytes.Repeat([]byte("-"), 50)...)
	p = appendLine(p)
	p = append(p, note...)
	p = appendLine(p)
	p = append(p, bytes.Repeat([]byte("-"), 50)...)
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

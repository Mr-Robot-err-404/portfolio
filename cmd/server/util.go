package main

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

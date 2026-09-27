package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func notification(content string, width, height int, fg, bg string) []byte {
	if width < 2 || height < 1 || strings.ContainsAny(content, "\r\n") {
		return nil
	}
	inner := width - 2
	contentWidth := utf8.RuneCountInString(content)
	if contentWidth > inner {
		return nil
	}

	var out strings.Builder
	for row := range height {
		text := ""
		if row == height/2 {
			text = content
		}

		padding := inner - utf8.RuneCountInString(text)
		left := padding / 2
		right := padding - left

		fmt.Fprintf(&out, "\x1b[%s;%sm▏%s%s%s▕\x1b[0m",
			fg, bg,
			strings.Repeat(" ", left),
			text,
			strings.Repeat(" ", right),
		)

		if row < height-1 {
			out.WriteString("\r\n")
		}
	}

	return []byte(out.String())
}

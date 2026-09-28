package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type Stat struct {
	Key   string
	Value string
}

func stats(rows []Stat, width int, primary, secondary string) []byte {
	if width < 7 || len(rows) == 0 {
		return nil
	}
	maxKeyWidth := 0
	maxValueWidth := 0

	for _, row := range rows {
		if strings.ContainsAny(row.Key+row.Value, "\r\n") {
			return nil
		}

		maxKeyWidth = max(maxKeyWidth, utf8.RuneCountInString(row.Key))
		maxValueWidth = max(maxValueWidth, utf8.RuneCountInString(row.Value))
	}

	keyWidth := maxKeyWidth + 2
	valueWidth := width - keyWidth - 3

	if valueWidth < maxValueWidth+2 {
		return nil
	}
	var out strings.Builder

	write := func(text, color string) {
		fmt.Fprintf(&out, "\x1b[%s;%sm%s\x1b[0m", color, StatsBG, text)
	}
	border := func(left, middle, right string) {
		write(
			left+
				strings.Repeat("─", keyWidth)+
				middle+
				strings.Repeat("─", valueWidth)+
				right,
			primary,
		)
		out.WriteString("\r\n")
	}
	border("┌", "┬", "┐")

	for i, row := range rows {
		color := primary
		if i%2 != 0 {
			color = secondary
		}
		keyPadding := keyWidth - utf8.RuneCountInString(row.Key) - 1
		valuePadding := valueWidth - utf8.RuneCountInString(row.Value) - 1

		line := "│ " +
			row.Key +
			strings.Repeat(" ", keyPadding) +
			"│ " +
			row.Value +
			strings.Repeat(" ", valuePadding) +
			"│"

		write(line, color)
		out.WriteString("\r\n")

		if i < len(rows)-1 {
			border("├", "┼", "┤")
		}
	}
	write(
		"└"+
			strings.Repeat("─", keyWidth)+
			"┴"+
			strings.Repeat("─", valueWidth)+
			"┘",
		primary,
	)
	return []byte(out.String())
}

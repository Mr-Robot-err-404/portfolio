package ascii

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	Green   string = "\x1b[32m"
	Blue    string = "\x1b[34m"
	Reset   string = "\x1b[0m"
	Shell   string = "❯"
	Clear   string = "\x1b[2J\x1b[H"
	Amber   string = "38;2;230;195;132"
	Ocean   string = "38;2;126;156;216"
	StatsBG string = "48;2;13;12;12"
)

type Stat struct {
	Key   string
	Value string
}

type TableStyle struct {
	Primary    string
	Secondary  string
	Background string
}

func Table(rows []Stat, width int, style TableStyle) []byte {
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

	styled := func(text, foreground string) string {
		return fmt.Sprintf(
			"\x1b[%s;%sm%s%s",
			foreground,
			style.Background,
			text,
			Reset,
		)
	}
	border := func(left, middle, right string, newline bool) {
		out.WriteString(styled(
			left+
				strings.Repeat("─", keyWidth)+
				middle+
				strings.Repeat("─", valueWidth)+
				right,
			style.Primary,
		))

		if newline {
			out.WriteString("\r\n")
		}
	}
	border("┌", "┬", "┐", true)

	for i, row := range rows {
		rowColor := style.Primary
		if i%2 != 0 {
			rowColor = style.Secondary
		}
		keyPadding := keyWidth - utf8.RuneCountInString(row.Key) - 1
		valuePadding := valueWidth - utf8.RuneCountInString(row.Value) - 1

		out.WriteString(styled("│ ", style.Primary))
		out.WriteString(styled(
			row.Key+strings.Repeat(" ", keyPadding),
			rowColor,
		))
		out.WriteString(styled("│ ", style.Primary))
		out.WriteString(styled(
			row.Value+strings.Repeat(" ", valuePadding),
			rowColor,
		))
		out.WriteString(styled("│", style.Primary))
		out.WriteString("\r\n")

		if i < len(rows)-1 {
			border("├", "┼", "┤", true)
		}
	}
	border("└", "┴", "┘", false)

	return []byte(out.String())
}

package ascii

import "fmt"

func ClearLine() []byte {
	return fmt.Appendf(nil, "\r\x1b[2K%s", Prompt())
}

func ClearAll() []byte {
	return fmt.Appendf(nil, "%s%s", Clear, Prompt())
}
func Color(msg string, color string) string {
	return fmt.Sprintf("%s%s%s", color, msg, Reset)
}
func ColorWithAnsi(msg string, color string) string {
	return fmt.Sprintf("\x1b[%sm%s%s", color, msg, Reset)
}
func OpenAnsi(color string) string {
	return fmt.Sprintf("\x1b[%sm", color)
}
func Prompt() []byte {
	user := Color("visitor@portfolio", Green)
	return fmt.Appendf(nil, "%s:~ %s ", user, Color(Shell, Blue))
}
func Unknown(command string) []byte {
	return fmt.Appendf(nil, "\nCommand not found: %s\n%s", command, Prompt())
}

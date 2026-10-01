package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Mr-Robot-err-404/portfolio/pkg/ascii"
)

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	switch os.Args[1] {
	case "stats":
		stats()
	case "bg":
		if len(os.Args) != 4 {
			usage()
		}
		background(os.Args[2], os.Args[3])
	case "stitch":
		if len(os.Args) != 6 {
			usage()
		}
		gap, err := strconv.Atoi(os.Args[5])
		if err != nil || gap < 0 {
			usage()
		}
		stitch(os.Args[2], os.Args[3], os.Args[4], gap)
	default:
		usage()
	}
}

func stats() {
	stats := ascii.Table([]ascii.Stat{
		{Key: "Name", Value: "Harry Lawton"},
		{Key: "Role", Value: "Software Engineer"},
		{Key: "Languages", Value: "Go, Odin, Typescript"},
		{Key: "Work", Value: "Backend, Infrastructure, Systems"},
		{Key: "Domains of interest", Value: "Graphics, Game Development"},
		{Key: "Approach", Value: "Generalist"},
	}, 60, ascii.Ocean, ascii.Amber)

	if err := os.WriteFile("stats.ascii", stats, 0o644); err != nil {
		log.Fatal(err)
	}
}

func background(inputPath, outputPath string) {
	input, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatal(err)
	}

	background := []byte("\x1b[" + ascii.StatsBG + "m")
	resetBackground := append([]byte(ascii.Reset), background...)
	output := make([]byte, 0, len(input)+len(background))
	output = append(output, background...)
	output = append(output, bytes.ReplaceAll(input, []byte(ascii.Reset), resetBackground)...)
	output = bytes.ReplaceAll(output, []byte("⠀"), []byte(" "))
	output = append(output, []byte(ascii.Reset)...)

	if err := os.WriteFile(outputPath, output, 0o644); err != nil {
		log.Fatal(err)
	}
}

func stitch(firstPath, secondPath, outputPath string, gap int) {
	first, err := os.ReadFile(firstPath)
	if err != nil {
		log.Fatal(err)
	}
	second, err := os.ReadFile(secondPath)
	if err != nil {
		log.Fatal(err)
	}

	firstLines := lines(first)
	secondLines := lines(second)
	firstStates := sgrStates(firstLines)
	firstWidth := 0
	for _, line := range firstLines {
		firstWidth = max(firstWidth, visibleWidth(line))
	}
	secondWidth := 0
	for _, line := range secondLines {
		secondWidth = max(secondWidth, visibleWidth(line))
	}

	lineCount := max(len(firstLines), len(secondLines))
	firstOffset := (lineCount - len(firstLines)) / 2
	secondOffset := (lineCount - len(secondLines)) / 2
	background := "\x1b[" + ascii.StatsBG + "m"
	var output strings.Builder
	for i := range lineCount {
		output.WriteString(ascii.Reset)
		output.WriteString(background)
		firstIndex := i - firstOffset
		if firstIndex >= 0 && firstIndex < len(firstLines) {
			output.WriteString(firstStates[firstIndex])
			output.WriteString(firstLines[firstIndex])
			output.WriteString(background)
			output.WriteString(strings.Repeat(" ", firstWidth-visibleWidth(firstLines[firstIndex])))
		} else {
			output.WriteString(strings.Repeat(" ", firstWidth))
		}
		output.WriteString(strings.Repeat(" ", gap))
		secondIndex := i - secondOffset
		if secondIndex >= 0 && secondIndex < len(secondLines) {
			output.WriteString(secondLines[secondIndex])
			output.WriteString(background)
			output.WriteString(strings.Repeat(" ", secondWidth-visibleWidth(secondLines[secondIndex])))
		} else {
			output.WriteString(strings.Repeat(" ", secondWidth))
		}
		if i < lineCount-1 {
			output.WriteString("\r\n")
		}
	}

	if err := os.WriteFile(outputPath, []byte(output.String()), 0o644); err != nil {
		log.Fatal(err)
	}
}

func sgrStates(lines []string) []string {
	states := make([]string, len(lines))
	state := ""
	for lineIndex, line := range lines {
		states[lineIndex] = state
		for i := 0; i < len(line); {
			if line[i] != '\x1b' || i+1 >= len(line) || line[i+1] != '[' {
				i++
				continue
			}
			end := i + 2
			for end < len(line) && (line[end] < '@' || line[end] > '~') {
				end++
			}
			if end >= len(line) {
				break
			}
			end++
			sequence := line[i:end]
			if sequence[len(sequence)-1] == 'm' {
				if sequence == ascii.Reset || sequence == "\x1b[m" {
					state = sequence
				} else {
					state += sequence
				}
			}
			i = end
		}
	}
	return states
}

func lines(input []byte) []string {
	normalized := strings.ReplaceAll(string(input), "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")
	return strings.Split(strings.TrimSuffix(normalized, "\n"), "\n")
}

func visibleWidth(line string) int {
	width := 0
	for i := 0; i < len(line); {
		if line[i] == '\x1b' && i+1 < len(line) && line[i+1] == '[' {
			i += 2
			for i < len(line) && (line[i] < '@' || line[i] > '~') {
				i++
			}
			if i < len(line) {
				i++
			}
			continue
		}
		_, size := utf8.DecodeRuneInString(line[i:])
		width++
		i += size
	}
	return width
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintln(os.Stderr, "  go run ./dev stats")
	fmt.Fprintln(os.Stderr, "  go run ./dev bg <input> <output>")
	fmt.Fprintln(os.Stderr, "  go run ./dev stitch <first> <second> <output> <gap>")
	os.Exit(2)
}

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
	case "colors":
		if len(os.Args) != 3 {
			usage()
		}
		colors(os.Args[2])
	case "replace":
		if len(os.Args) != 5 || os.Args[2] == "" {
			usage()
		}
		replace(os.Args[2], os.Args[3], os.Args[4])
	default:
		usage()
	}
}

func colors(path string) {
	input, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	seen := make(map[string]bool)
	for i := 0; i < len(input); i++ {
		if input[i] != '\x1b' || i+1 >= len(input) || input[i+1] != '[' {
			continue
		}
		end := i + 2
		for end < len(input) && (input[end] < '@' || input[end] > '~') {
			end++
		}
		if end >= len(input) {
			break
		}
		if input[end] != 'm' {
			i = end
			continue
		}

		params := strings.Split(string(input[i+2:end]), ";")
		for j := 0; j < len(params); j++ {
			color := ""
			switch {
			case params[j] == "38" || params[j] == "48":
				if j+1 < len(params) {
					size := 0
					switch params[j+1] {
					case "2":
						size = 5
					case "5":
						size = 3
					}
					if size > 0 {
						if params[j] == "38" && j+size <= len(params) {
							color = strings.Join(params[j:j+size], ";")
						}
						j += min(size-1, len(params)-j-1)
					}
				}
			case foregroundCode(params[j]):
				color = params[j]
			}
			if color != "" && !seen[color] {
				seen[color] = true
				fmt.Printf("%s  \x1b[%sm██\x1b[0m\n", color, color)
			}
		}
		i = end
	}
}

func foregroundCode(value string) bool {
	code, err := strconv.Atoi(value)
	return err == nil && (code >= 30 && code <= 37 || code >= 90 && code <= 97)
}

func replace(old, new, path string) {
	input, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	if !bytes.Contains(input, []byte(old)) {
		log.Fatalf("color %q not found in %s", old, path)
	}
	info, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(path, bytes.ReplaceAll(input, []byte(old), []byte(new)), info.Mode().Perm()); err != nil {
		log.Fatal(err)
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
	}, 60, ascii.TableStyle{
		Primary:    ascii.Ocean,
		Secondary:  ascii.Amber,
		Background: ascii.StatsBG,
	})

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
	fmt.Fprintln(os.Stderr, "  go run ./dev colors <file>")
	fmt.Fprintln(os.Stderr, "  go run ./dev replace <old> <new> <file>")
	os.Exit(2)
}

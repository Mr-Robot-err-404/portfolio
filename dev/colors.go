package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

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

package main

import (
	"fmt"
	"log"
	"os"
)

func inspect(path string, count int) {
	input, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	characters := []rune(string(input))
	fmt.Println("     BEFORE │ NEWLINE │ AFTER")
	newline := 0
	for i, character := range characters {
		if character != '\n' {
			continue
		}
		newline++
		before := characters[max(0, i-count):i]
		after := characters[i+1 : min(len(characters), i+1+count)]
		fmt.Printf("%4d %q │ \\n │ %q\n", newline, string(before), string(after))
	}
}

package main

import (
	"bytes"
	"log"
	"os"

	"github.com/Mr-Robot-err-404/portfolio/pkg/ascii"
)

func moveReset(inputPath, outputPath string) {
	input, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatal(err)
	}

	reset := []byte(ascii.Reset)
	output := bytes.ReplaceAll(input, []byte("\r"), nil)
	output = bytes.ReplaceAll(output, append([]byte("\n"), reset...), append(reset, '\n'))

	if err := os.WriteFile(outputPath, output, 0o644); err != nil {
		log.Fatal(err)
	}
}

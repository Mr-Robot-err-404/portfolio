package main

import (
	"bytes"
	"log"
	"os"

	"github.com/Mr-Robot-err-404/portfolio/pkg/ascii"
)

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

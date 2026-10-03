package main

import (
	"bytes"
	"log"
	"os"
)

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

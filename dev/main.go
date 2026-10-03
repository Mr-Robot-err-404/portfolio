package main

import (
	"fmt"
	"os"
	"strconv"
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
	case "inspect":
		if len(os.Args) != 4 {
			usage()
		}
		count, err := strconv.Atoi(os.Args[3])
		if err != nil || count < 0 {
			usage()
		}
		inspect(os.Args[2], count)
	case "move-reset":
		if len(os.Args) != 4 {
			usage()
		}
		moveReset(os.Args[2], os.Args[3])
	default:
		usage()
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintln(os.Stderr, "  go run ./dev stats")
	fmt.Fprintln(os.Stderr, "  go run ./dev bg <input> <output>")
	fmt.Fprintln(os.Stderr, "  go run ./dev stitch <first> <second> <output> <gap>")
	fmt.Fprintln(os.Stderr, "  go run ./dev colors <file>")
	fmt.Fprintln(os.Stderr, "  go run ./dev replace <old> <new> <file>")
	fmt.Fprintln(os.Stderr, "  go run ./dev inspect <file> <n>")
	fmt.Fprintln(os.Stderr, "  go run ./dev move-reset <input> <output>")
	os.Exit(2)
}

package main

import (
	"log"
	"os"

	"github.com/Mr-Robot-err-404/portfolio/pkg/ascii"
)

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
		Border:     ascii.Ocean,
	})

	if err := os.WriteFile("stats.ascii", stats, 0o644); err != nil {
		log.Fatal(err)
	}
}

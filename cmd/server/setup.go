package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Mr-Robot-err-404/portfolio/pkg/ascii"
)

func helpMenu() []byte {
	menu := ascii.Table([]ascii.Stat{
		{Key: "about", Value: "Meet the person behind the terminal"},
		{Key: "contact", Value: "Find me elsewhere"},
		{Key: "projects", Value: "Explore what I've built"},
		{Key: "clear", Value: "Clear the terminal"},
	}, 60, ascii.TableStyle{
		Primary:   ascii.Amber,
		Secondary: ascii.Ocean,
		Border:    ascii.Ocean,
	})
	menu = appendLine(menu)
	menu = appendLine(menu)

	note := []byte(fmt.Sprintf(
		"  \x1b[%s  → %s%s",
		ascii.ColorWithAnsi("connect", ascii.Orange),
		ascii.ColorWithAnsi("Gain access to your own unix shell!", ascii.Amber),
		ascii.Reset,
	))
	menu = addNote(menu, note)
	return menu
}

func projectsAscii() (string, error) {
	ascii := strings.Builder{}

	b, err := os.ReadFile("static/banners/perkins.ascii")
	if err != nil {
		return "", err
	}
	ascii.WriteString(string(b))

	b, err = os.ReadFile("static/banners/tinyrenderer.ascii")
	if err != nil {
		return "", err
	}
	ascii.WriteString("\n")
	ascii.WriteString(string(b))

	b, err = os.ReadFile("static/banners/wireframe.ascii")
	if err != nil {
		return "", err
	}
	ascii.WriteString("\n")
	ascii.WriteString(string(b))
	return ascii.String(), nil
}

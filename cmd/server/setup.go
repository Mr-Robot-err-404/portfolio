package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Mr-Robot-err-404/portfolio/pkg/ascii"
)

func helpMenu() []byte {
	menu := title(ascii.ColorWithAnsi("  COMMAND    DESCRIPTION", ascii.Ocean))

	table := ascii.Table([]ascii.Stat{
		{Key: "about", Value: "Meet the person behind the terminal"},
		{Key: "contact", Value: "Find me elsewhere"},
		{Key: "projects", Value: "Explore what I've built"},
		{Key: "clear", Value: "Clear the terminal"},
		{Key: "help", Value: "Display this help menu"},
	}, 60, ascii.TableStyle{
		Primary:   ascii.Amber,
		Secondary: ascii.Ocean,
		Border:    ascii.Ocean,
	})
	menu = append(menu, table...)
	menu = appendLine(menu)
	menu = appendLine(menu)

	highlight := ascii.Table(
		[]ascii.Stat{{Key: "connect ", Value: "Gain access to your own unix shell!"}},
		60,
		ascii.TableStyle{
			Primary: ascii.Ocean,
			Border:  ascii.Amber,
		})
	menu = append(menu, highlight...)
	menu = appendLine(menu)
	return menu
}

func profileAscii() []byte {
	b, err := os.ReadFile("static/mr_robot/profile.ascii")
	if err != nil {
		panic(fmt.Sprintf("failed to load ascii file: %s", err.Error()))
	}
	b = append(b, []byte(ascii.Reset)...)
	return b
}

func projectsAscii() (string, error) {
	str := strings.Builder{}

	b, err := os.ReadFile("static/banners/perkins.ascii")
	if err != nil {
		return "", err
	}
	str.WriteString(string(b))

	b, err = os.ReadFile("static/banners/tinyrenderer.ascii")
	if err != nil {
		return "", err
	}
	str.WriteString("\n")
	str.WriteString(string(b))

	b, err = os.ReadFile("static/banners/wireframe.ascii")
	if err != nil {
		return "", err
	}
	str.WriteString("\n")
	str.WriteString(string(b))
	return str.String(), nil
}

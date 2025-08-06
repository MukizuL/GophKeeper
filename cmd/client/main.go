package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	var serverAddr string
	flag.StringVar(&serverAddr, "s", "localhost:8080", "Server address")

	flag.Parse()

	s := newStart()
	a := newAbout()

	m := model{
		start:  s,
		about:  a,
		window: "start",
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

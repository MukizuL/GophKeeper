package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type about struct{}

func newAbout() about {
	return about{}
}

func updateAbout(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			m.window = "start"
		}
	}

	return m, nil
}

func viewAbout(m model) string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(itemStyle.Render("Version: "))
	b.WriteString(itemStyle.Render(buildVersion))
	b.WriteString("\n")
	b.WriteString(itemStyle.Render("Build commit: "))
	b.WriteString(itemStyle.Render(buildCommit))
	b.WriteString("\n")
	b.WriteString(itemStyle.Render("Build date: "))
	b.WriteString(itemStyle.Render(buildDate))
	b.WriteString("\n")

	button := &backButton

	fmt.Fprintf(&b, "\n%s\n\n", *button)

	return b.String()
}

package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type viewText struct {
	text text
}

func newViewText(t text) viewText {
	return viewText{text: t}
}

func updateViewText(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			m.window = "storage-text"
		}
	}

	return m, nil
}

func viewViewText(m model) string {
	var b strings.Builder

	lines, err := WrapNoSplitWords(m.viewText.text.Text, 20)
	if err != nil {
		b.WriteString(err.Error())
		b.WriteString("\n\n")

		button := &backButtonFocused

		fmt.Fprintf(&b, "\n%s\n\n", *button)

		return b.String()
	}

	b.WriteString("\n")
	b.WriteString(itemStyle.Render(m.viewText.text.Name))
	b.WriteString("\n\n")
	for _, line := range lines {
		b.WriteString(itemStyle.Render(line))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	button := &backButtonFocused

	fmt.Fprintf(&b, "\n%s\n\n", *button)

	return b.String()
}

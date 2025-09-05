package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type viewBank struct {
	card card
}

func newViewBank(c card) viewBank {
	return viewBank{card: c}
}

func updateViewBank(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			m.window = "storage-bank"
		}
	}

	return m, nil
}

func viewViewBank(m model) string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(itemStyle.Render("CCN: "))
	b.WriteString(itemStyle.Render(m.viewBank.card.CCN))
	b.WriteString("\n")
	b.WriteString(itemStyle.Render("EXP: "))
	b.WriteString(itemStyle.Render(m.viewBank.card.EXP))
	b.WriteString("\n")
	b.WriteString(itemStyle.Render("CVV: "))
	b.WriteString(itemStyle.Render(m.viewBank.card.CVV))
	b.WriteString("\n")
	b.WriteString(itemStyle.Render("Name: "))
	b.WriteString(itemStyle.Render(m.viewBank.card.Name))
	b.WriteString("\n")

	button := &backButtonFocused

	fmt.Fprintf(&b, "\n%s\n\n", *button)

	return b.String()
}

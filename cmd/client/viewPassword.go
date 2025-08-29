package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type viewPassword struct {
	password password
}

func newViewPassword(p password) viewPassword {
	return viewPassword{password: p}
}

func updateViewPassword(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			m.window = "storage-password"
		}
	}

	return m, nil
}

func viewViewPassword(m model) string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(itemStyle.Render("Name: "))
	b.WriteString(itemStyle.Render(m.viewPassword.password.Name))
	b.WriteString("\n")
	b.WriteString(itemStyle.Render("Login: "))
	b.WriteString(itemStyle.Render(m.viewPassword.password.Login))
	b.WriteString("\n")
	b.WriteString(itemStyle.Render("Password: "))
	b.WriteString(itemStyle.Render(m.viewPassword.password.Password))
	b.WriteString("\n")
	b.WriteString(itemStyle.Render("Description: "))
	b.WriteString(itemStyle.Render(m.viewPassword.password.Description))
	b.WriteString("\n")

	button := &backButtonFocused

	fmt.Fprintf(&b, "\n%s\n\n", *button)

	return b.String()
}

package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type start struct {
	list list.Model
}

func newStart() start {
	startItems := []list.Item{
		item("Login"),
		item("Create new account"),
		item("About"),
	}

	l := list.New(startItems, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Welcome to GophKeeper"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = listTitleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return start{
		list: l,
	}
}

func updateStart(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.start.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := m.start.list.SelectedItem().(item)
			if ok {
				switch i {
				case "Login":
					m.window = "login"
					m.login.inputs[0].Focus()
				case "Create new account":
					m.window = "register"
					m.register.inputs[0].Focus()
				case "About":
					m.window = "about"
				}
			}
		}
	}

	var cmd tea.Cmd
	m.start.list, cmd = m.start.list.Update(msg)
	return m, cmd
}

func viewStart(m model) string {
	var b strings.Builder

	if m.register.success {
		b.WriteString(successStyle.Render("Account created"))
	}
	b.WriteString("\n")
	b.WriteString(m.start.list.View())

	return b.String()
}

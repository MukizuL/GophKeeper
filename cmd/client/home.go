package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type home struct {
	list list.Model
}

func newHome() home {
	homeItems := []list.Item{
		item("View"),
		item("Create"),
		item("Logout"),
	}

	l := list.New(homeItems, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Home"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = listTitleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return home{
		list: l,
	}
}

func updateHome(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.home.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := m.home.list.SelectedItem().(item)
			if ok {
				switch i {
				case "View":
					m.window = "storage"
				case "Create":
					m.window = "create"
				case "Logout":

					m.token = ""
					m.window = "start"
				}
			}
		}
	}

	var cmd tea.Cmd
	m.home.list, cmd = m.home.list.Update(msg)
	return m, cmd
}

func viewHome(m model) string {
	var b strings.Builder

	if m.createPassword.success || m.createBank.success || m.createText.success {
		b.WriteString(successStyle.Render("Entry created"))
	}
	b.WriteString("\n")
	b.WriteString(m.home.list.View())

	return b.String()
}

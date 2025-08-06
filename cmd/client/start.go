package main

import (
	listB "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type start struct {
	list listB.Model
}

func newStart() start {
	startItems := []listB.Item{
		item("Login"),
		item("Create new account"),
		item("About"),
	}

	l := listB.New(startItems, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Welcome to GophKeeper"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
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
				case "Create new account":
					m.window = "create-account"
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
	return "\n" + m.start.list.View()
}

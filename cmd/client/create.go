package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type create struct {
	list list.Model
}

func newCreate() create {
	createItems := []list.Item{
		item("Passwords"),
		item("Bank details"),
		item("Texts"),
		item("Data"),
		item("Back"),
	}

	l := list.New(createItems, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Create"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = listTitleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return create{
		list: l,
	}
}

func updateCreate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.create.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := m.create.list.SelectedItem().(item)
			if ok {
				switch i {
				case "Passwords":
					m.window = "create-password"
					m.createPassword.inputs[0].Focus()
				case "Bank details":
					m.window = "create-bank"
					m.createBank.inputs[0].Focus()
				case "Texts":
					m.window = "create-text"
					m.createText.name.Focus()
				case "Data":
					m.window = "create-data"
					if !m.createData.once {
						cmds = append(cmds, m.createData.fp.Init())
						m.createData.once = true
					}
				case "Back":
					m.window = "home"
				}
			}
		}
	}

	var cmd tea.Cmd
	m.create.list, cmd = m.create.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func viewCreate(m model) string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(m.create.list.View())

	return b.String()
}

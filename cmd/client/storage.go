package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type storage struct {
	list list.Model
}

func newStorage() storage {
	items := []list.Item{
		item("Passwords"),
		item("Bank details"),
		item("Texts"),
		item("Data"),
		item("Back"),
	}

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Storage"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = listTitleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return storage{
		list: l,
	}
}

func updateStorage(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.storage.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := m.storage.list.SelectedItem().(item)
			if ok {
				switch i {
				case "Passwords":
					m.window = "storage-password"
					m.storagePasswords = newStoragePasswords(m)
				case "Bank details":
					m.window = "storage-bank"
				case "Texts":
					m.window = "storage-text"
				case "Data":
					m.window = "storage-data"
				case "Back":
					m.window = "home"
				}
			}
		}
	}

	var cmd tea.Cmd
	m.storage.list, cmd = m.storage.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func viewStorage(m model) string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(m.storage.list.View())

	return b.String()
}

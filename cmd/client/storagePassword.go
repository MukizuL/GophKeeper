package main

import (
	"encoding/json"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type storagePasswords struct {
	list      list.Model
	passwords map[string]password
	err       error
}

type password struct {
	Name        string `json:"name"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	Description string `json:"description"`
}

func newStoragePasswords(m model) storagePasswords {
	var items []list.Item

	data, err := GetPasswords(m.token, m.dk)
	if err != nil {
		return storagePasswords{err: err}
	}

	out := storagePasswords{
		passwords: make(map[string]password, len(data)),
	}

	var passwords []password
	for _, v := range data {
		var temp password
		err = json.Unmarshal(v, &temp)
		if err != nil {
			return storagePasswords{err: err}
		}

		out.passwords[temp.Name] = temp

		items = append(items, item(temp.Name))

		passwords = append(passwords, temp)
	}

	items = append(items, item("Back"))

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Passwords"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = listTitleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	out.list = l

	return out
}

func updateStoragePasswords(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.storagePasswords.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := m.storagePasswords.list.SelectedItem().(item)
			if ok {
				switch i {
				case "Back":
					m.window = "storage"
				default:
					m.window = "view-password"
					m.viewPassword = newViewPassword(m.storagePasswords.passwords[string(i)])
				}
			}
		}
	}

	var cmd tea.Cmd
	m.storagePasswords.list, cmd = m.storagePasswords.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func viewStoragePasswords(m model) string {
	var b strings.Builder

	b.WriteString("\n")
	if m.storagePasswords.err != nil {
		b.WriteString("An error occurred. Try again")
	} else {
		b.WriteString(m.storagePasswords.list.View())
	}

	return b.String()
}

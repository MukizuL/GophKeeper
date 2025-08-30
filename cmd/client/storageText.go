package main

import (
	"encoding/json"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type storageText struct {
	list  list.Model
	texts map[string]text
	err   error
}

type text struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

func newStorageText(m model) storageText {
	var items []list.Item

	data, err := GetText(m.token, m.dk)
	if err != nil {
		return storageText{err: err}
	}

	out := storageText{
		texts: make(map[string]text, len(data)),
	}

	for _, v := range data {
		var temp text
		err = json.Unmarshal(v, &temp)
		if err != nil {
			return storageText{err: err}
		}

		out.texts[temp.Name] = temp

		items = append(items, item(temp.Name))
	}

	items = append(items, item("Back"))

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Texts"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = listTitleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	out.list = l

	return out
}

func updateStorageText(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.storageText.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := m.storageText.list.SelectedItem().(item)
			if ok {
				switch i {
				case "Back":
					m.window = "storage"
				default:
					m.window = "view-text"
					m.viewText = newViewText(m.storageText.texts[string(i)])
				}
			}
		}
	}

	var cmd tea.Cmd
	m.storageText.list, cmd = m.storageText.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func viewStorageText(m model) string {
	var b strings.Builder

	b.WriteString("\n")
	if m.storageText.err != nil {
		b.WriteString("An error occurred. Try again")
	} else {
		b.WriteString(m.storageText.list.View())
	}

	return b.String()
}

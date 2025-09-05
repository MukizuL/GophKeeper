package main

import (
	"encoding/json"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type storageBank struct {
	list  list.Model
	cards map[string]card
	err   error
}

type card struct {
	Name string `json:"name"`
	CCN  string `json:"ccn"`
	EXP  string `json:"exp"`
	CVV  string `json:"cvv"`
}

func newStorageBank(m model) storageBank {
	var items []list.Item

	data, err := GetBank(m.token, m.dk)
	if err != nil {
		return storageBank{err: err}
	}

	out := storageBank{
		cards: make(map[string]card, len(data)),
	}

	for _, v := range data {
		var temp card
		err = json.Unmarshal(v, &temp)
		if err != nil {
			return storageBank{err: err}
		}

		out.cards[temp.Name] = temp

		items = append(items, item(temp.Name))
	}

	items = append(items, item("Back"))

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Cards"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = listTitleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	out.list = l

	return out
}

func updateStorageBank(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.storageBank.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := m.storageBank.list.SelectedItem().(item)
			if ok {
				switch i {
				case "Back":
					m.window = "storage"
				default:
					m.window = "view-bank"
					m.viewBank = newViewBank(m.storageBank.cards[string(i)])
				}
			}
		}
	}

	var cmd tea.Cmd
	m.storageBank.list, cmd = m.storageBank.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func viewStorageBank(m model) string {
	var b strings.Builder

	b.WriteString("\n")
	if m.storageBank.err != nil {
		b.WriteString("An error occurred. Try again")
	} else {
		b.WriteString(m.storageBank.list.View())
	}

	return b.String()
}

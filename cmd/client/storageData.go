package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type storageData struct {
	list            list.Model
	files           map[string]file
	downloadSuccess bool
	err             error
}

type file struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
}

func newStorageData(m model) storageData {
	var items []list.Item

	files, err := GetData(m.token, m.dk)
	if err != nil {
		return storageData{err: err}
	}

	out := storageData{
		files: make(map[string]file, len(files)),
	}

	for _, v := range files {
		out.files[v.Filename] = v

		items = append(items, item(v.Filename))
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

func updateStorageData(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.storageData.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := m.storageData.list.SelectedItem().(item)
			if ok {
				switch i {
				case "Back":
					m.window = "storage"
				default:
					m.window = "view-data"
					m.storageData.downloadSuccess = false
					m.storageData.err = nil
					m.viewData = newViewData(m.storageData.files[string(i)])
					cmds = append(cmds, m.viewData.spinner.Tick)
				}
			}
		}
	}

	var cmd tea.Cmd
	m.storageData.list, cmd = m.storageData.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func viewStorageData(m model) string {
	var b strings.Builder

	if m.storageData.downloadSuccess {
		b.WriteString(successStyle.Render("Download successful"))
		b.WriteString("\n")
	}
	if m.storageData.err != nil {
		b.WriteString(errorStyle.Render(m.storageData.err.Error()))
		b.WriteString("\n")
	} else {
		b.WriteString("\n")
	}

	b.WriteString(m.storageData.list.View())

	return b.String()
}

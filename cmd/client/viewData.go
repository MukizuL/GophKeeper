package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type viewData struct {
	file        file
	downloading bool
	spinner     spinner.Model
	err         error
}

func newViewData(f file) viewData {
	s := spinner.New()
	s.Spinner = spinner.Line

	return viewData{
		file:    f,
		spinner: s,
	}
}

func updateViewData(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	if !m.viewData.downloading {
		m.viewData.downloading = true
		return m, DownloadFile(m.token, m.dk, m.viewData.file.ID, m.viewData.file.Filename)
	}

	switch msg := msg.(type) {
	case Done:
		m.storageData.downloadSuccess = true
		m.window = "storage-data"
		return m, nil
	case errMsg:
		m.storageData.err = msg.error
		m.window = "storage-data"
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.viewData.spinner, cmd = m.viewData.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func viewViewData(m model) string {
	var b strings.Builder

	if m.viewData.err != nil {
		b.WriteString(m.viewData.err.Error())
		b.WriteString("\n\n")
	} else {
		b.WriteString("\n")
	}

	b.WriteString("Downloading:")
	b.WriteString("\n")
	b.WriteString(m.viewData.spinner.View())
	b.WriteString("\n")

	return b.String()
}

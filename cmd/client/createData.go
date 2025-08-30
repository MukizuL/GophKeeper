package main

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type createData struct {
	fp           filepicker.Model
	selectedFile string
	once         bool
	success      bool
	error        error
	progress     progress.Model
	uploading    bool
	percent      float64
	ch           chan tea.Msg
}

func newCreateData() createData {
	fp := filepicker.New()
	fp.CurrentDirectory, _ = os.UserHomeDir()
	fp.ShowPermissions = false
	fp.AutoHeight = false
	fp.SetHeight(14)

	p := progress.New(progress.WithDefaultGradient())

	return createData{
		fp:       fp,
		progress: p,
		ch:       make(chan tea.Msg),
	}
}

func updateCreateData(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case uploadProgressMsg:
		m.createData.percent = float64(msg)
		if m.createData.percent >= 1.0 {
			m.createData.success = true
			m.createData.uploading = false
			m.window = "home"
		}
		return m, updateProgressBar(m.createData.ch)

	case errMsg:
		m.createData.error = msg.error
		m.createData.uploading = false
		return m, nil
	}

	var cmd tea.Cmd
	m.createData.fp, cmd = m.createData.fp.Update(msg)

	if didSelect, path := m.createData.fp.DidSelectFile(msg); didSelect {
		m.createData.selectedFile = path
		m.createData.uploading = true

		return m, CreateData(m.token, m.dk, m.createData.selectedFile, m.createData.ch)
	}

	return m, tea.Batch(cmd, updateProgressBar(m.createData.ch))
}

func viewCreateData(m model) string {
	var b strings.Builder

	if m.createData.error != nil {
		b.WriteString(errorStyle.Render(m.createData.error.Error()))
		b.WriteString("\n\n")
	} else {
		b.WriteString("\n")
	}

	if m.createData.uploading {
		b.WriteString("Uploading:")
		b.WriteString("\n")
		b.WriteString(m.createData.progress.ViewAs(m.createData.percent))
		b.WriteString("\n")

		return b.String()
	}

	if m.createData.selectedFile != "" {
		b.WriteString("Select a file")
		b.WriteString("\n\n")
	}

	b.WriteString(m.createData.fp.View())

	return b.String()
}

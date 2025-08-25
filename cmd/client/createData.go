package main

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

type createData struct {
	fp           filepicker.Model
	selectedFile string
	once         bool
	success      bool
	error        error
}

func newCreateData() createData {
	fp := filepicker.New()
	fp.CurrentDirectory, _ = os.UserHomeDir()
	fp.ShowPermissions = false
	fp.AutoHeight = false
	fp.SetHeight(10)

	return createData{
		fp: fp,
	}
}

func updateCreateData(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.createData.fp, cmd = m.createData.fp.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.createData.fp.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		m.createData.selectedFile = path
		// TODO: GRPC request to upload a file.
	}

	return m, cmd
}

func viewCreateData(m model) string {
	var b strings.Builder

	if m.createData.error != nil {
		b.WriteString(errorStyle.Render(m.createData.error.Error()))
		b.WriteString("\n\n")
	} else {
		b.WriteString("\n")
	}

	if m.createData.selectedFile != "" {
		b.WriteString("Select a file")
		b.WriteString("\n\n")
	}

	b.WriteString(m.createData.fp.View())

	return b.String()
}

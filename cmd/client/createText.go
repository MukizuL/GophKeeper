package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type createText struct {
	focusIndex int
	name       textinput.Model
	text       textarea.Model
	success    bool
	error      error
}

func newCreateText() createText {
	ti := textinput.New()
	ti.Cursor.Style = selectedItemStyle
	ti.Placeholder = "Name"
	ti.PromptStyle = formStyle
	ti.TextStyle = formSelectedStyle
	ti.CharLimit = 20
	ti.Width = 30

	ta := textarea.New()

	return createText{
		name: ti,
		text: ta,
	}
}

func (c *createText) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 2)

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	c.name, cmds[0] = c.name.Update(msg)
	c.text, cmds[1] = c.text.Update(msg)

	return tea.Batch(cmds...)
}

func updateCreateText(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.createText.focusIndex == 3 {
				if m.createText.name.Value() == "" {
					m.createText.error = fmt.Errorf("name cannot be empty")
					return m, nil
				}
				m.createText.error = nil
				// TODO: GRPC request to create text. Each name should be unique and not empty.

				m.window = "home"
				m.createText.success = true

				return m, nil
			}

			// If user hits Back, return him to Create
			if s == "enter" && m.createText.focusIndex == 4 {
				m.window = "create"
				resetCreateText(&m)
				return m, nil
			}

			if s == "up" || s == "shift+tab" {
				m.createText.focusIndex--
			} else {
				m.createText.focusIndex++
			}

			if m.createText.focusIndex > 3 {
				m.createText.focusIndex = 0
			} else if m.createText.focusIndex < 0 {
				m.createText.focusIndex = 3
			}

			cmds := make([]tea.Cmd, 2)
			for i := range 2 {
				switch i {
				case 0:
					if m.createText.focusIndex == 0 {
						cmds[0] = m.createText.name.Focus()
						m.createText.name.PromptStyle = formSelectedStyle
						m.createText.name.TextStyle = formSelectedStyle
					} else {
						m.createText.name.Blur()
						m.createText.name.PromptStyle = formStyle
						m.createText.name.TextStyle = formStyle
					}
				case 1:
					if m.createText.focusIndex == 1 {
						cmds[1] = m.createText.text.Focus()
					} else {
						m.createText.text.Blur()
					}
				}
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.createText.updateInputs(msg)

	return m, cmd
}

func viewCreateText(m model) string {
	var b strings.Builder

	if m.createText.error != nil {
		b.WriteString(errorStyle.Render(m.createText.error.Error()))
		b.WriteString("\n")
	} else {
		b.WriteString("\n")
	}

	b.WriteString(titleStyle.Render("Create text entry"))
	b.WriteString("\n\n")

	//b.WriteString(titleStyle.Render("Name"))
	//b.WriteString("\n")
	b.WriteString(m.createText.name.View())
	b.WriteString("\n\n")

	//b.WriteString(titleStyle.Render("Text"))
	//b.WriteString("\n")
	b.WriteString(m.createText.text.View())
	b.WriteString("\n\n")

	ok := okButton
	if m.createText.focusIndex == 2 {
		ok = okButtonFocused
	}

	back := backButton
	if m.createText.focusIndex == 3 {
		back = backButtonFocused
	}

	fmt.Fprintf(&b, "\n\n%s\n%s\n\n", ok, back)

	return b.String()
}

func resetCreateText(m *model) {
	m.createText.error = nil
	m.createText.success = false
	m.createText.focusIndex = 0
	m.createText.name.Reset()
	m.createText.text.Reset()
}

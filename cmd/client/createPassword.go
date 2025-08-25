package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type createPassword struct {
	focusIndex int
	inputs     []textinput.Model
	success    bool
	error      error
}

func newCreatePassword() createPassword {
	l := createPassword{
		inputs: make([]textinput.Model, 4),
	}

	var t textinput.Model
	for i := range l.inputs {
		t = textinput.New()
		t.Cursor.Style = selectedItemStyle
		t.CharLimit = 255
		t.Width = 20

		switch i {
		case 0:
			t.Placeholder = "Site/Name"
			t.PromptStyle = formStyle
			t.TextStyle = formSelectedStyle
		case 1:
			t.Placeholder = "Login"
			t.PromptStyle = formStyle
			t.TextStyle = formSelectedStyle
		case 2:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
			t.PromptStyle = formStyle
			t.TextStyle = formStyle
		case 3:
			t.Placeholder = "Description"
			t.PromptStyle = formStyle
			t.TextStyle = formSelectedStyle
		}

		l.inputs[i] = t
	}

	return l
}

func (c *createPassword) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(c.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range c.inputs {
		c.inputs[i], cmds[i] = c.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func updateCreatePassword(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.createPassword.focusIndex == len(m.createPassword.inputs) {
				m.createPassword.error = nil

				if m.createPassword.inputs[0].Value() == "" {
					m.createPassword.error = errors.New("name/site cannot be empty")
					return m, nil
				}
				// TODO: GRPC request to create password. Should reject same login/password combination within
				// TODO: same site.
				err := CreatePassword(m.token, m.dk,
					m.createPassword.inputs[0].Value(),
					m.createPassword.inputs[1].Value(),
					m.createPassword.inputs[2].Value(),
					m.createPassword.inputs[3].Value())
				if err != nil {
					m.createPassword.error = err
					return m, nil
				}

				m.window = "home"
				m.createPassword.success = true

				return m, nil
			}

			// If user hits Back, return him to Create
			if s == "enter" && m.createPassword.focusIndex == len(m.createPassword.inputs)+1 {
				m.window = "create"
				resetCreatePassword(&m)
				return m, nil
			}

			if s == "up" || s == "shift+tab" {
				m.createPassword.focusIndex--
			} else {
				m.createPassword.focusIndex++
			}

			if m.createPassword.focusIndex > len(m.createPassword.inputs)+1 {
				m.createPassword.focusIndex = 0
			} else if m.createPassword.focusIndex < 0 {
				m.createPassword.focusIndex = len(m.createPassword.inputs) + 1
			}

			cmds := focusOrBlur(m.createPassword.inputs, m.createPassword.focusIndex)

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.createPassword.updateInputs(msg)

	return m, cmd
}

func viewCreatePassword(m model) string {
	var b strings.Builder

	if m.createPassword.error != nil {
		b.WriteString(errorStyle.Render(m.createPassword.error.Error()))
		b.WriteString("\n")
	} else {
		b.WriteString("\n")
	}

	b.WriteString(titleStyle.Render("Create password"))
	b.WriteString("\n")

	for i := range m.createPassword.inputs {
		b.WriteString(m.createPassword.inputs[i].View())
		if i < len(m.createPassword.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	ok := okButton
	if m.createPassword.focusIndex == len(m.createPassword.inputs) {
		ok = okButtonFocused
	}

	back := backButton
	if m.createPassword.focusIndex == len(m.createPassword.inputs)+1 {
		back = backButtonFocused
	}

	fmt.Fprintf(&b, "\n\n%s\n%s\n\n", ok, back)

	return b.String()
}

func resetCreatePassword(m *model) {
	m.createPassword.error = nil
	m.createPassword.success = false
	m.createPassword.focusIndex = 0
	for i := range m.createPassword.inputs {
		m.createPassword.inputs[i].Reset()
	}
}

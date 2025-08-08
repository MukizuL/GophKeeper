package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type login struct {
	focusIndex int
	inputs     []textinput.Model
	success    bool
	error      error
}

func newLogin() login {
	l := login{
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range l.inputs {
		t = textinput.New()
		t.Cursor.Style = selectedItemStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Login"
			//t.Focus()
			t.PromptStyle = formStyle
			t.TextStyle = formSelectedStyle
		case 1:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
			t.PromptStyle = formStyle
			t.TextStyle = formStyle
		}

		l.inputs[i] = t
	}

	return l
}

func (l *login) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(l.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range l.inputs {
		l.inputs[i], cmds[i] = l.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func updateLogin(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.login.focusIndex == len(m.login.inputs) {
				m.register.error = nil
				token, err := Login(m.login.inputs[0].Value(), m.login.inputs[1].Value())
				if err != nil {
					m.login.error = err
					return m, nil
				}

				m.window = "home"
				m.login.success = true

				m.token = token

				return m, nil
			}

			// If user hits Back, return him to Start
			if s == "enter" && m.login.focusIndex == len(m.login.inputs)+1 {
				m.window = "start"
				resetLogin(&m)
				return m, nil
			}

			if s == "up" || s == "shift+tab" {
				m.login.focusIndex--
			} else {
				m.login.focusIndex++
			}

			if m.login.focusIndex > len(m.login.inputs)+1 {
				m.login.focusIndex = 0
			} else if m.login.focusIndex < 0 {
				m.login.focusIndex = len(m.login.inputs) + 1
			}

			cmds := make([]tea.Cmd, len(m.login.inputs))
			for i := 0; i < len(m.login.inputs); i++ {
				if i == m.login.focusIndex {
					// Set focused state
					cmds[i] = m.login.inputs[i].Focus()
					m.login.inputs[i].PromptStyle = formSelectedStyle
					m.login.inputs[i].TextStyle = formSelectedStyle
					continue
				}
				// Remove focused state
				m.login.inputs[i].Blur()
				m.login.inputs[i].PromptStyle = formStyle
				m.login.inputs[i].TextStyle = formStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.login.updateInputs(msg)

	return m, cmd
}

func viewLogin(m model) string {
	var b strings.Builder

	if m.login.error != nil {
		b.WriteString(errorStyle.Render(m.login.error.Error()))
		b.WriteString("\n")
	} else {
		b.WriteString("\n")
	}

	b.WriteString(titleStyle.Render("Create new account"))
	b.WriteString("\n")

	for i := range m.login.inputs {
		b.WriteString(m.login.inputs[i].View())
		if i < len(m.login.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	ok := okButton
	if m.login.focusIndex == len(m.login.inputs) {
		ok = okButtonFocused
	}

	back := backButton
	if m.login.focusIndex == len(m.login.inputs)+1 {
		back = backButtonFocused
	}

	fmt.Fprintf(&b, "\n\n%s\n%s\n\n", ok, back)

	return b.String()
}

func resetLogin(m *model) {
	m.login.error = nil
	m.login.success = false
	m.login.focusIndex = 0
	m.login.inputs[0].Reset()
	m.login.inputs[1].Reset()
}

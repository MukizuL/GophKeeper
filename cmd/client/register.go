package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type register struct {
	focusIndex int
	inputs     []textinput.Model
	success    bool
	error      error
}

func newRegister() register {
	r := register{
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range r.inputs {
		t = textinput.New()
		t.Cursor.Style = selectedItemStyle
		t.CharLimit = 32
		t.Width = 20

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

		r.inputs[i] = t
	}

	return r
}

func (r *register) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(r.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range r.inputs {
		r.inputs[i], cmds[i] = r.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func updateRegister(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.register.focusIndex == len(m.register.inputs) {
				m.register.error = nil
				err := Register(m.register.inputs[0].Value(), m.register.inputs[1].Value())
				if err != nil {
					m.register.error = err
					return m, nil
				}

				m.window = "start"
				m.register.success = true

				return m, nil
			}

			// If user hits Back, return him to Start
			if s == "enter" && m.register.focusIndex == len(m.register.inputs)+1 {
				m.window = "start"
				resetRegister(&m)
				return m, nil
			}

			if s == "up" || s == "shift+tab" {
				m.register.focusIndex--
			} else {
				m.register.focusIndex++
			}

			if m.register.focusIndex > len(m.register.inputs)+1 {
				m.register.focusIndex = 0
			} else if m.register.focusIndex < 0 {
				m.register.focusIndex = len(m.register.inputs) + 1
			}

			cmds := focusOrBlur(m.register.inputs, m.register.focusIndex)

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.register.updateInputs(msg)

	return m, cmd
}

func viewRegister(m model) string {
	var b strings.Builder

	if m.register.error != nil {
		b.WriteString(errorStyle.Render(m.register.error.Error()))
		b.WriteString("\n")
	} else {
		b.WriteString("\n")
	}

	b.WriteString(titleStyle.Render("Create new account"))
	b.WriteString("\n")

	for i := range m.register.inputs {
		b.WriteString(m.register.inputs[i].View())
		if i < len(m.register.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	ok := okButton
	if m.register.focusIndex == len(m.register.inputs) {
		ok = okButtonFocused
	}

	back := backButton
	if m.register.focusIndex == len(m.register.inputs)+1 {
		back = backButtonFocused
	}

	fmt.Fprintf(&b, "\n\n%s\n%s\n\n", ok, back)

	return b.String()
}

func resetRegister(m *model) {
	m.register.error = nil
	m.register.success = false
	m.register.focusIndex = 0
	m.register.inputs[0].Reset()
	m.register.inputs[1].Reset()
}

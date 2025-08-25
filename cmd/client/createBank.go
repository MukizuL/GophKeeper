package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type createBank struct {
	focusIndex int
	inputs     []textinput.Model
	success    bool
	error      error
}

func newCreateBank() createBank {
	l := createBank{
		inputs: make([]textinput.Model, 4),
	}

	var t textinput.Model
	for i := range l.inputs {
		t = textinput.New()
		t.Cursor.Style = selectedItemStyle

		switch i {
		case 0:
			t.Placeholder = "4505 **** **** 1234"
			t.PromptStyle = formStyle
			t.TextStyle = formSelectedStyle
			t.CharLimit = 20
			t.Width = 30
		case 1:
			t.Placeholder = "MM/YY "
			t.PromptStyle = formStyle
			t.TextStyle = formSelectedStyle
			t.CharLimit = 5
			t.Width = 5
		case 2:
			t.Placeholder = "XXX"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = 'X'
			t.PromptStyle = formStyle
			t.TextStyle = formStyle
			t.CharLimit = 3
			t.Width = 3
		case 3:
			t.Placeholder = "John Jonson"
			t.PromptStyle = formStyle
			t.TextStyle = formSelectedStyle
			t.CharLimit = 255
			t.Width = 30
		}

		l.inputs[i] = t
	}

	return l
}

func ccnValidator(s string) error {
	// Credit Card Number should a string less than 20 digits
	// It should include 16 integers and 3 spaces
	if len(s) > 16+3 {
		return fmt.Errorf("CCN is too long")
	}

	if len(s) < 16+3 {
		return fmt.Errorf("CCN is too short")
	}

	if len(s) == 0 || len(s)%5 != 0 && (s[len(s)-1] < '0' || s[len(s)-1] > '9') {
		return fmt.Errorf("CCN is invalid")
	}

	// The last digit should be a number unless it is a multiple of 4 in which
	// case it should be a space
	if len(s)%5 == 0 && s[len(s)-1] != ' ' {
		return fmt.Errorf("CCN must separate groups with spaces")
	}

	// The remaining digits should be integers
	c := strings.ReplaceAll(s, " ", "")
	_, err := strconv.ParseInt(c, 10, 64)

	return err
}

func expValidator(s string) error {
	// The 3 character should be a slash (/)
	// The rest should be numbers
	e := strings.ReplaceAll(s, "/", "")
	_, err := strconv.ParseInt(e, 10, 64)
	if err != nil {
		return fmt.Errorf("EXP is invalid")
	}

	// There should be only one slash, and it should be in the 2nd index (3rd character)
	if len(s) != 5 && (strings.Index(s, "/") != 2 || strings.LastIndex(s, "/") != 2) {
		return fmt.Errorf("EXP is invalid")
	}

	return nil
}

func cvvValidator(s string) error {
	// The CVV should be a number of 3 digits
	// Since the input will already ensure that the CVV is a string of length 3,
	// All we need to do is check that it is a number
	if len(s) != 3 {
		return fmt.Errorf("CVV is invalid")
	}
	_, err := strconv.ParseInt(s, 10, 64)
	return err
}

func (c *createBank) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(c.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range c.inputs {
		c.inputs[i], cmds[i] = c.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func updateCreateBank(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.createBank.focusIndex == len(m.createBank.inputs) {
				err := ccnValidator(m.createBank.inputs[0].Value())
				if err != nil {
					m.createBank.error = err
					return m, nil
				}
				err = expValidator(m.createBank.inputs[1].Value())
				if err != nil {
					m.createBank.error = err
					return m, nil
				}
				err = cvvValidator(m.createBank.inputs[2].Value())
				if err != nil {
					m.createBank.error = err
					return m, nil
				}
				m.createBank.error = nil
				// TODO: GRPC request to create bank card. Should error if it is a duplicate number.

				m.window = "home"
				m.createBank.success = true

				return m, nil
			}

			// If user hits Back, return him to Create
			if s == "enter" && m.createBank.focusIndex == len(m.createBank.inputs)+1 {
				m.window = "create"
				resetCreateBank(&m)
				return m, nil
			}

			if s == "up" || s == "shift+tab" {
				m.createBank.focusIndex--
			} else {
				m.createBank.focusIndex++
			}

			if m.createBank.focusIndex > len(m.createBank.inputs)+1 {
				m.createBank.focusIndex = 0
			} else if m.createBank.focusIndex < 0 {
				m.createBank.focusIndex = len(m.createBank.inputs) + 1
			}

			cmds := focusOrBlur(m.createBank.inputs, m.createBank.focusIndex)

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.createBank.updateInputs(msg)

	return m, cmd
}

func viewCreateBank(m model) string {
	var b strings.Builder

	if m.createBank.error != nil {
		b.WriteString(errorStyle.Render(m.createBank.error.Error()))
		b.WriteString("\n")
	} else {
		b.WriteString("\n")
	}

	b.WriteString(titleStyle.Render("Create bank card"))
	b.WriteString("\n\n")

	b.WriteString(titleStyle.Render("Card number"))
	b.WriteString("\n")
	b.WriteString(m.createBank.inputs[0].View())
	b.WriteString("\n\n")

	b.WriteString(titleStyle.Render("EXP         CVV"))
	b.WriteString("\n")
	b.WriteString(m.createBank.inputs[1].View())
	b.WriteString(m.createBank.inputs[2].View())
	b.WriteString("\n\n")

	b.WriteString(titleStyle.Render("Name"))
	b.WriteString("\n")
	b.WriteString(m.createBank.inputs[3].View())

	ok := okButton
	if m.createBank.focusIndex == len(m.createBank.inputs) {
		ok = okButtonFocused
	}

	back := backButton
	if m.createBank.focusIndex == len(m.createBank.inputs)+1 {
		back = backButtonFocused
	}

	fmt.Fprintf(&b, "\n\n%s\n%s\n\n", ok, back)

	return b.String()
}

func resetCreateBank(m *model) {
	m.createBank.error = nil
	m.createBank.success = false
	m.createBank.focusIndex = 0
	for i := range m.createBank.inputs {
		m.createBank.inputs[i].Reset()
	}
}

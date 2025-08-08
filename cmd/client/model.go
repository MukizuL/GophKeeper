package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	listHeight   = 14
	defaultWidth = 20
)

var (
	listTitleStyle    = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	blurredStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)

	formStyle         = lipgloss.NewStyle().PaddingLeft(4)
	formSelectedStyle = lipgloss.NewStyle().PaddingLeft(4).Foreground(lipgloss.Color("170"))

	titleStyle          = lipgloss.NewStyle().PaddingLeft(4)
	errorStyle          = lipgloss.NewStyle().PaddingLeft(4).Foreground(lipgloss.Color("9"))
	successStyle        = lipgloss.NewStyle().PaddingLeft(4).Foreground(lipgloss.Color("2"))
	selectedButtonStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))

	okButton        = itemStyle.Render("[ OK ]")
	okButtonFocused = itemStyle.Render(selectedButtonStyle.Render("[ OK ]"))

	backButton        = itemStyle.Render("[ Back ]")
	backButtonFocused = itemStyle.Render(selectedButtonStyle.Render("[ Back ]"))
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	token    string
	start    start
	about    about
	register register
	login    login
	window   string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	switch m.window {
	case "start":
		return updateStart(msg, m)
	case "about":
		return updateAbout(msg, m)
	case "register":
		return updateRegister(msg, m)
	case "login":
		return updateLogin(msg, m)
	default:
		return m, nil
	}
}

func (m model) View() string {
	switch m.window {
	case "start":
		return viewStart(m)
	case "about":
		return viewAbout(m)
	case "register":
		return viewRegister(m)
	case "login":
		return viewLogin(m)
	default:
		return fmt.Sprintf("Unknown window. token = %s", m.token)
	}
}

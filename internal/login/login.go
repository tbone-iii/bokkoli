package login

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	cursor   int
	choices  []string
	selected map[int]struct{}
}

func New() Model {
	return Model{
		choices:  []string{"Start chatting", "option 2", "option 3"},
		selected: make(map[int]struct{}),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.SetWindowTitle("start here")
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			if m.cursor == 0 {
				return m, func() tea.Msg { return "start chatting" }
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	s := strings.Builder{}
	s.WriteString("Bokkoli\n\n")

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("34")).
		Bold(true)

	defaultStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	for i := 0; i < len(m.choices); i++ {
		itemStyle := defaultStyle
		if m.cursor == i {
			itemStyle = selectedStyle
		}

		cursor := " "
		if m.cursor == i {
			cursor = "(•)"
		} else {
			cursor = "( )"
		}

		s.WriteString(fmt.Sprintf("%s %s\n", cursor, itemStyle.Render(m.choices[i])))
	}

	return s.String()
}

func StartChatting() tea.Cmd {
	return func() tea.Msg {
		return "start chatting"
	}
}

package main

import (
	"fmt"
	"log"

	"bokkoli/internal/control"
	"bokkoli/internal/instruction"
	"bokkoli/internal/login"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedModelStyle = lipgloss.NewStyle().
				Width(15).
				Height(5).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

// used to track which model is focused
type sessionState uint

const (
	//default??
	loginView sessionState = iota //
	controlView
	instructionView
)

type mainModel struct {
	state       sessionState //
	instruction instruction.Model
	login       login.Model
	control     control.Model
}

func newModel() mainModel {
	m := mainModel{state: loginView}
	m.login = login.New()
	m.control = control.New()
	m.instruction = instruction.New()
	return m
}

func (m mainModel) Init() tea.Cmd {
	return nil //start views on program start
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //At some point maybe arrow keys instead of tab and press enter to actually enter the view
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		//what is the key pressed?
		switch msg.String() {
		//exit program
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.state == loginView {
				m.state = controlView
			} else if m.state == controlView {
				m.state = instructionView
			} else {
				m.state = loginView
			}
		}

		switch m.state {
		// update whichever model is focused
		case loginView:
			m.login, cmd = m.login.Update(msg)
			cmds = append(cmds, cmd)
		case controlView:
			m.control, cmd = m.control.Update(msg)
			cmds = append(cmds, cmd)
		case instructionView:
			m.instruction, cmd = m.instruction.Update(msg)
			cmds = append(cmds, cmd)
		default:
			fmt.Println("Oops. Defaulted.")
		}
	}
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	var s string
	model := m.currentFocusedModel()
	if m.state == loginView {
		s += lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render(m.login.View()), m.control.View(), m.instruction.View())
	} else if m.state == controlView {
		s += lipgloss.JoinHorizontal(lipgloss.Top, m.login.View(), focusedModelStyle.Render(m.control.View()), m.instruction.View())
	} else {
		s += lipgloss.JoinHorizontal(lipgloss.Top, m.login.View(), m.control.View(), focusedModelStyle.Render(m.instruction.View()))
	}
	s += helpStyle.Render(fmt.Sprintf("\ntab: focus next • n: new %s • q: exit\n", model))
	return s
}

func (m mainModel) currentFocusedModel() string {
	if m.state == loginView {
		return "login"
	} else if m.state == controlView {
		return "control"
	} else {
		return "instruction"
	}
}

func main() {
	p := tea.NewProgram(newModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

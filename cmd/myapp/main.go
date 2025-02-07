package main

import (
	"bokkoli/internal/login"
	"bokkoli/internal/message"
	"bokkoli/internal/setup"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedModelStyle = lipgloss.NewStyle().
				Width(30).
				Height(8).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

type sessionState uint

const (
	loginView sessionState = iota
	chatView
	setupView
)

type mainModel struct {
	state sessionState
	login login.Model
	chat  *message.ChatModel
	setup *setup.Model
}

func newModel() mainModel {
	m := mainModel{state: loginView}
	m.login = login.New()
	m.chat = message.New()
	m.setup = setup.New()
	return m
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		if msg.String() == "enter" && m.state == loginView {
			m.state = chatView
			// Here you'll want to start up the chat room
			m.chat = message.New()
		}

	}

	switch m.state {
	case loginView:
		m.login, cmd = m.login.Update(msg)
		cmds = append(cmds, cmd)
	case chatView:
		var updatedChat tea.Model
		updatedChat, cmd = m.chat.Update(msg)

		if chatModel, ok := updatedChat.(*message.ChatModel); ok {
			m.chat = chatModel
		} else {
			log.Println("Unexpected type assertion failure for ChatModel")
		}
		cmds = append(cmds, cmd)
	case setupView:
		var updatedSetup tea.Model
		updatedSetup, cmd = m.chat.Update(msg)

		if Model, ok := updatedSetup.(*setup.Model); ok {
			m.setup = Model
		} else {
			log.Println("Unexpected type assertion failure for Model")
		}
	}

	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	if m.state == chatView {
		return m.chat.View()
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		focusedModelStyle.Render(m.login.View()),
		helpStyle.Render("\nPress ↑/↓ to navigate • Press Enter to select • Press Q to quit"),
	)
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	fmt.Println("Welcome to Bokkoli!")
	p := tea.NewProgram(newModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

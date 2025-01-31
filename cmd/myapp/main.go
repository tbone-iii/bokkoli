package main

import (
	"bokkoli/internal/login"
	"bokkoli/internal/message"
	"fmt"
	"log"

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
)

type mainModel struct {
	state sessionState
	login login.Model
	chat  *message.ChatModel // Use pointer type here
}

func newModel() mainModel {
	m := mainModel{state: loginView}
	m.login = login.New()
	m.chat = message.New() // Corrected to use pointer
	return m
}

func (m mainModel) Init() tea.Cmd {
	return nil // Start views when the program begins
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

		// Check if "start chatting" was typed in the login view
		if msg.String() == "enter" && m.state == loginView {
			// Transition to chat view when "start chatting" is selected
			m.state = chatView
			// Start the server when chat view is entered
			go func() {
				// Simulate server start when entering the chat view
				// Call your method to start the server here
				message.RunChat(tea.NewProgram(m.chat)) // Pass chat as pointer to tea program
			}()
		}

	}

	// Handle different views
	switch m.state {
	case loginView:
		m.login, cmd = m.login.Update(msg)
		cmds = append(cmds, cmd)
	case chatView:
		var updatedChat tea.Model
		updatedChat, cmd = m.chat.Update(msg) // Calling on pointer

		if chatModel, ok := updatedChat.(*message.ChatModel); ok { // Corrected type assertion to pointer
			m.chat = chatModel
		} else {
			log.Println("Unexpected type assertion failure for ChatModel")
		}
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	if m.state == chatView {
		return m.chat.View() // We can call View on pointer type as well
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		focusedModelStyle.Render(m.login.View()),
		helpStyle.Render("\nPress ↑/↓ to navigate • Press Enter to select • Press Q to quit"),
	)
}

func main() {
	fmt.Println("Welcome to Bokkoli!")
	p := tea.NewProgram(newModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

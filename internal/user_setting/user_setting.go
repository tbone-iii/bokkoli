package user_setting

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

var (
	username   string
	portNumber string
	confirm    bool
)

type Model struct {
	form *huh.Form
}

func isEmpty(input string) error {
	if input == "" {
		return errors.New("username cannot be empty")
	}
	return nil
}

func NewModel() Model {
	return Model{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("username").
					Title("Username").
					Prompt("?").
					Validate(isEmpty).
					Value(&username),
			),
			huh.NewGroup(
				huh.NewInput().
					Title("Enter new port number").
					Validate(isEmpty).
					Value(&portNumber),
			),
			//save NewInput port values in NewSelect options
			huh.NewGroup(
				huh.NewSelect[string]().
					Key("port").
					Options(
						huh.NewOption("user input 1", "1"),
						huh.NewOption("user input 2", "2"),
						huh.NewOption("user input 3", "3"),
					).
					Value(&portNumber),
			),
			huh.NewGroup(
				huh.NewConfirm().
					Title("Please confirm username and port number").
					Affirmative("Save").
					Value(&confirm),
			),
		),
	}
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	return m, cmd
}

func (m Model) View() string {
	if m.form.State == huh.StateCompleted {
		username := m.form.GetString("class")
		port := m.form.GetString("level")
		return fmt.Sprintf("You selected: username. %s, port. %s", username, port)
	}
	return m.form.View()
}

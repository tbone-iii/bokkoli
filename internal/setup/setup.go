package setup

import (
	"errors"
	"fmt"
	"strconv"

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
		return errors.New("input cannot be empty")
	}
	return nil
}

func New() *Model {
	return &Model{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Input username").
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
					Title("Select port").
					OptionsFunc(func() []huh.Option[string] {
						var options []huh.Option[string]
						for i := 1; i <= 4; i++ {
							options = append(options, huh.NewOption("value-"+strconv.Itoa(i), "value-"+strconv.Itoa(i)))
						}
						return options
					}, nil). // Corrected parenthesis here
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
		username := m.form.GetString("username")
		port := m.form.GetString("port")
		return fmt.Sprintf("You selected: username. %s, port. %s", username, port)
	}
	return m.form.View()
}

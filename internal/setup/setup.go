package setup

import (
	"bokkoli/internal/db"
	"fmt"
	"log"
	"reflect"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// TODO: Another file in setup, same package, create functions to pull data from database settings table
// Then load that into here upon start

const (
	LOWERBOUND_PORT_NUMBER int = 1024
	UPPERBOUND_PORT_NUMBER int = 49151
)

var (
	username   string = "Username-read-from-db" // db.readUsername()
	portNumber string = "8080"                  // db.readPortNumber()
	confirm    bool
)

type SetupModel struct {
	form                    *huh.Form
	dbHandler               *db.DbHandler
	isValidDataAndCompleted bool
}

func New() *SetupModel {
	return &SetupModel{
		form: huh.NewForm(
			huh.NewGroup(
				// TODO: Skipping to new port field by default for some reason, shift-tab is required to return to first input field
				huh.NewInput().
					Key("username").
					Title("Input username").
					Prompt("> ").
					Placeholder("<username>").
					Value(&username),
				huh.NewInput().
					Key("port").
					Title("Enter new port number").
					Suggestions([]string{"8080", "8081"}).
					Value(&portNumber),
				huh.NewConfirm().
					Title("Please confirm username and port number").
					Affirmative("Save").
					Value(&confirm),
			),
		),
	}
}

func (m SetupModel) Init() tea.Cmd {
	dbHandler, err := db.NewDbHandler(db.DefaultDbFilePath)
	if err != nil {
		log.Fatal("DB failed to open in setup model.")
	}

	m.dbHandler = dbHandler
	return m.form.Init()
}

func (m SetupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	} else {
		log.Fatal("Wrong type assertion, expected *huh.Form, got ", reflect.TypeOf(form))
	}

	// TODO: Escape to return to main menu

	if m.form.State == huh.StateCompleted && !m.isValidDataAndCompleted {
		tempUsername := m.form.GetString("username")
		tempPort := m.form.GetString("port")

		if !validateUsername(tempUsername) {
			log.Println("Bad username, clearing it out.")
			username = ""
		}

		if !validatePort(tempPort) {
			log.Println("Bad port, clearing it out.")
			portNumber = ""
		}

		if validateUsername(tempUsername) && validatePort(tempPort) {
			err := saveSetup(m.dbHandler, tempPort, tempUsername)
			// TODO: Improve error message
			if err != nil {
				log.Panic("PANICKKKKK!!11, didn't save to DB properly for ", tempPort, tempUsername)
			}
			m.isValidDataAndCompleted = true
		}

	}

	return m, cmd
}

func validateUsername(username string) bool {
	// TODO: Add extra restrictions based on mood of the day
	return username != ""
}

func validatePort(port string) bool {
	portNumber, err := strconv.Atoi(port)
	if err != nil {
		log.Printf("%s is not valid port number due to not being integer: ", port)
		return false
	}

	if !(LOWERBOUND_PORT_NUMBER <= portNumber && portNumber <= UPPERBOUND_PORT_NUMBER) {
		log.Printf("%d is out of range [%d, %d]", portNumber, LOWERBOUND_PORT_NUMBER, UPPERBOUND_PORT_NUMBER)
		return false
	}

	return true
}

func (m SetupModel) View() string {
	if m.isValidDataAndCompleted {
		return m.form.View() +
			fmt.Sprintf("\n\nSaved successfully, you selected username: %s, port: %s", username, portNumber) +
			lipgloss.NewStyle().Blink(true).Faint(true).Render("\nPress 'esc' to return back to main menu.")
	}
	return m.form.View()
}

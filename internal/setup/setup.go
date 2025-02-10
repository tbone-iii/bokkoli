package setup

import (
	"bokkoli/internal/db"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const (
	LOWERBOUND_PORT_NUMBER int = 1024
	UPPERBOUND_PORT_NUMBER int = 49151
)

var (
	username   string //= "Username-read-from-db" // db.readUsername()
	portNumber string //= "8080"                  // db.readPortNumber()
	confirm    bool
)

type SetupModel struct {
	Form                    *huh.Form
	dbHandler               *db.DbHandler
	isValidDataAndCompleted bool
}

func New() *SetupModel {
	// Reset global variables

	dbHandler, err := db.NewDbHandler(db.DefaultDbFilePath)
	if err != nil {
		log.Fatal("DB failed to open in setup model.")
	}

	err = dbHandler.SetupSchemas()
	if err != nil {
		log.Fatal("DB failed to set up schema for setup.")
	}

	// Create the form with proper validation
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(""), // bubbletea Huh bug, possibly remove this field if not needed
			huh.NewInput().
				Key("username").
				Title("Input username").
				Placeholder("<username>").
				Value(&username).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("sorry, username cannot be left empty")
					}
					return nil
				}),
			huh.NewInput().
				Key("port").
				Title("Enter new port number").
				Placeholder("<port number>").
				Value(&portNumber).
				Validate(func(str string) error {
					// Try to convert the port string to an integer
					portInt, err := strconv.Atoi(str)
					if err != nil {
						return errors.New("invalid port number, ports are numbers")
					}

					// Port validation logic: must be in the range [1024, 49151]
					if portInt < LOWERBOUND_PORT_NUMBER || portInt > UPPERBOUND_PORT_NUMBER {
						return fmt.Errorf("sorry, only ports in the range %d-%d are allowed", LOWERBOUND_PORT_NUMBER, UPPERBOUND_PORT_NUMBER)
					}
					return nil
				}),
			huh.NewConfirm().
				Title("Please confirm username and port number").
				Validate(func(v bool) error {
					if !v {
						return fmt.Errorf("no isn't actually an option, press 'Save' ~(^-^~)")
					}
					return nil
				}).
				Affirmative("Save").
				Value(&confirm),
		),
	)
	return &SetupModel{
		Form:      form,
		dbHandler: dbHandler,
	}
}

func (m SetupModel) Init() tea.Cmd {
	return m.Form.Init()
}

func (m SetupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.Form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.Form = f
	} else {
		log.Fatal("Wrong type assertion, expected *huh.Form, got ", reflect.TypeOf(form))
	}

	// TODO: Even when no is selected on confirm field, form is considered "complete", make sure to check the condition
	if m.Form.State == huh.StateCompleted && !m.isValidDataAndCompleted {
		tempUsername := m.Form.GetString("username")
		tempPort := m.Form.GetString("port")

		if !validateUsername(tempUsername) {
			log.Println("Bad username, clearing it out.")
			username = ""
		}

		if !validatePort(tempPort) {
			log.Println("Bad port, clearing it out.")
			portNumber = ""
		}

		if validateUsername(tempUsername) && validatePort(tempPort) {
			err := m.dbHandler.SaveSetup(tempPort, tempUsername)

			if err != nil {
				log.Panicf("DB did not save record properly to settings.\nPort: %s\nUsername: %s", tempPort, tempUsername)
			}
			m.isValidDataAndCompleted = true
			if m.Form.State == huh.StateCompleted && !m.isValidDataAndCompleted {
				fmt.Printf("Form State: %v\n", m.Form.State)                           // Add debug output here
				fmt.Printf("isValidDataAndCompleted: %v\n", m.isValidDataAndCompleted) // Add debug output here
			}
		}

	}

	return m, cmd
}

func validateUsername(username string) bool {
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
		return m.Form.View() +
			fmt.Sprintf("\n\nSaved successfully, you selected username: %s, port: %s", username, portNumber) +
			lipgloss.NewStyle().Faint(true).Render("\nPress 'esc' to return back to main menu.")
	}
	return m.Form.View()
}

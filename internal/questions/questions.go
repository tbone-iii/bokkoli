package questions

import (
	"bokkoli/internal/db"
	"log"
	"reflect"

	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	question string
	confirm  bool
)

type QuestionsModel struct {
	Form                    *huh.Form
	Handler                 *db.DbHandler
	isValidDataAndCompleted bool
}

func New() *QuestionsModel {
	dbHandler, err := db.NewDbHandler(db.DefaultDbFilePath)
	if err != nil {
		log.Fatal("DB failed to open debug model.")
	}

	err = dbHandler.bugsSchemas()
	if err != nil {
		log.Fatal("DB failed to set up schema for bugs")
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Any questions or bugs that you found?").
				Validate(validateQuestion).
				Value(&question),
			huh.NewConfirm().
				Title("Would you like to submit your question?").
				Affirmative("Yes!").
				Negative("No.").
				Value(&confirm),
		),
	)

	return &QuestionsModel{
		Form:      form,
		dbHandler: dbHandler,
	}
}

func (m QuestionsModel) Init() tea.Cmd {
	return m.Form.Init()
}

func (m QuestionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.Form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.Form = f
	} else {
		log.Fatal("Wrong type assertion, expected *huh.Form, got ", reflect.TypeOf(form))
	}
}

func (m QuestionsModel) View() string {
	if m.isValidDataAndCompleted {
		return m.Form.View() +
			fmt.Sprintf("\n\nSaved successfully, you have submitted a bug, questions, comment")
		lipgloss.NewStyle().Blink(true).Faint(true).Render("\nPress 'esc' to return back to main menu.")
	}
}

func validateQuestion(question string) bool {
	return question != ""
}

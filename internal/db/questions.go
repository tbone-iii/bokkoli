package db

import (
	_ "github.com/mattn/go-sqlite3"
)

type Questions struct {
	Question string
}

func (handler *DbHandler) setupQuestionsSchema() error {
	query := `
	`
}

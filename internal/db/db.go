package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

const DefaultDbFilePath string = "./bokkoli.db"

type DbHandler struct {
	db *sql.DB
}

// Execute a query without returning any rows, includes a SQL result
func (handler DbHandler) ExecuteQuery(query string, args ...any) (sql.Result, error) {
	var result sql.Result
	var err error

	if len(args) == 0 {
		result, err = handler.db.Exec(query)
	} else {
		result, err = handler.db.Exec(query, args...)
	}

	if err != nil {
		log.Printf("failed to execute query: %v", err)
		return result, err
	}

	return result, nil
}

// Run a query that returns records/rows
func (handler DbHandler) Query(query string) (*sql.Rows, error) {
	result, err := handler.db.Query(query)
	if err != nil {
		log.Println("failed to query: ", err)
		return result, err
	}

	return result, nil
}

func NewDbHandler(filePath string) (*DbHandler, error) {
	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	return &DbHandler{db: db}, nil
}

func (handler DbHandler) SetupSchemas() error {
	schemas := []func() error{
		handler.setupMessageSchema,
		handler.setupSetupSchema,
	}

	for _, setupFn := range schemas {
		if err := setupFn(); err != nil {
			return err
		}
	}
	return nil
}

// Close the DB connection
func (handler DbHandler) Close() error {
	return handler.db.Close()
}

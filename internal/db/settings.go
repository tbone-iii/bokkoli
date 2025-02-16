package db

import (
	"errors"
	"log"

	_ "modernc.org/sqlite"
)

type Setup struct {
	Port     string
	Username string
}

func (handler *DbHandler) setupSetupSchema() error {
	query := `
    CREATE TABLE IF NOT EXISTS setup (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        port TEXT NOT NULL,
		username TEXT NOT NULL
    );`

	_, err := handler.ExecuteQuery(query)
	return err
}

// Save user settings to the DB; new record is created if one doesn't exist. Otherwise, previous record is overwritten.
func (handler *DbHandler) SaveSetup(port string, username string) error {
	query := `
	SELECT id
	FROM setup	
	`

	rows, err := handler.Query(query)
	if err != nil {
		return err
	}

	var id int
	count := 0
	for rows.Next() {
		count++
		if count > 1 {
			log.Fatal("More than 1 row inserted into setup, please identify issue.")
		}

		rows.Scan(&id)
	}

	if count == 0 {
		query := `
		INSERT INTO setup (port, username)
		VALUES (?, ?);
		`

		_, err := handler.ExecuteQuery(query, port, username)
		return err
	}

	query = `
	UPDATE setup
	SET port = ?, username = ?
	WHERE id = ?
	`
	result, err := handler.ExecuteQuery(query, port, username, id)
	log.Println(result.LastInsertId())
	return err
}

func (handler *DbHandler) ReadSetup() (Setup, error) {
	query := `
	SELECT id, port, username
	FROM setup
	LIMIT 1
	`

	var setup Setup

	rows, err := handler.Query(query)
	if err != nil {
		return setup, err
	}

	var id int
	for rows.Next() {
		rows.Scan(&id, &setup.Port, &setup.Username)
	}

	if setup.Port == "" || setup.Username == "" {
		log.Println("Port and/or Username are empty", setup.Port, setup.Username)
		return setup, errors.New("no results in the read setup query")
	}

	log.Println("Port and Username values read from DB: ", setup.Port, setup.Username)

	return setup, nil

}

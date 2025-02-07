package setup

import (
	"bokkoli/internal/db"

	_ "github.com/mattn/go-sqlite3"
)

// TODO: Consider receiver methods for ALL things DBhandler related (e.g. chatdb methods)
func setupSchema(handler *db.DbHandler) error {
	query := `
    CREATE TABLE IF NOT EXISTS setup (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        port TEXT NOT NULL,
		username TEXT NOT NULL,
    );`

	_, err := handler.ExecuteQuery(query)
	return err
}

func saveSetup(handler *db.DbHandler, port string, username string) error {
	query := `
	INSERT INTO setup (port, username)
	VALUES (?, ?);
	`

	_, err := handler.ExecuteQuery(query, port, username)
	return err
}

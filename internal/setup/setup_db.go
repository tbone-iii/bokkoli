package setup

import (
	"bokkoli/internal/db"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// TODO: Consider receiver methods for ALL things DBhandler related (e.g. chatdb methods)
func setupSchema(handler *db.DbHandler) error {
	query := `
    CREATE TABLE IF NOT EXISTS setup (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        port TEXT NOT NULL,
		username TEXT NOT NULL
    );`

	_, err := handler.ExecuteQuery(query)
	return err
}

func saveSetup(handler *db.DbHandler, port string, username string) error {

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

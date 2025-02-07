package message

import (
	"bokkoli/internal/db"

	_ "github.com/mattn/go-sqlite3"
)

// TODO: Consider receiver methods for ALL things DBhandler related (e.g. chatdb methods)
func setupSchema(handler *db.DbHandler) error {
	query := `
    CREATE TABLE IF NOT EXISTS messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        text TEXT NOT NULL,
        sender TEXT NOT NULL,
        receiver TEXT NOT NULL,
		direction TEXT NOT NULL,
        timestamp DATETIME NOT NULL
    );`

	_, err := handler.ExecuteQuery(query)
	return err
}

func saveMessage(handler *db.DbHandler, msg Message) error {
	query := `
	INSERT INTO messages (text, sender, receiver, direction, timestamp)
	VALUES (?, ?, ?, ?, ?);
	`

	_, err := handler.ExecuteQuery(query, msg.Text, msg.Sender, msg.Receiver, msg.Direction, msg.Timestamp)
	return err
}

package message

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type ChatDb struct {
	db *sql.DB
}

const DefaultDbFilePath string = "./bokkoli.db"

func NewChatDB(filePath string) (*ChatDb, error) {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	if err := setupSchema(db); err != nil {
		return nil, fmt.Errorf("failed fot setup schema: %v", err)
	}
	return &ChatDb{db: db}, nil //work iwtha reference for ChatDB rather than copy
}

func setupSchema(db *sql.DB) error {
	schema := `
    CREATE TABLE IF NOT EXISTS messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        text TEXT NOT NULL,
        sender TEXT NOT NULL,
        receiver TEXT NOT NULL,
		direction TEXT NOT NULL,
        timestamp DATETIME NOT NULL
    );`
	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to execute schema: %v", err)
	}
	return nil
}

func (chatdb *ChatDb) saveMessage(msg Message) error {
	query := `
	INSERT INTO messages (text, sender, receiver, direction, timestamp)
	VALUES (?, ?, ?, ?, ?)
	`

	_, err := chatdb.db.Exec(query, msg.Text, msg.Sender, msg.Receiver, msg.Direction, msg.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to save message: %v", err)
	}
	return nil
}

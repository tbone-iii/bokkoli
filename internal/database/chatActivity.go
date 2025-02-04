package chatActivity

import (
	"database/sql"
	"fmt"
	"sync"
)

type ChatActivity struct {
	mu sync.Mutex
	db *sql.DB
}

const file string = "./chatActivity.db"

func NewChatActivity() (*ChatActivity, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	if err := setupSchema(db); err != nil {
		return nil, fmt.Errorf("failed fot setup schema: %v", err)
	}
	return &ChatActivity{db: db}, nil
}

func setupSchema(db *sql.DB) error {
	schema := `
    CREATE TABLE IF NOT EXISTS messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        text TEXT NOT NULL,
        sender TEXT NOT NULL,
        receiver TEXT NOT NULL,
        message_type TEXT NOT NULL,
        timestamp DATETIME NOT NULL
    );`
	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to execute schema: %v", err)
	}
	return nil
}

func (ca *ChatActivity) saveMessage(msg Message) error {
	ca.mu.Lock()

	query := `
	INSERT INTO messages (text, sender, receiver, message_type, timestamp)
	VALUES (?, ?, ?, ?, ?)`

	_, err := ca.db.Exec(query, msg.Text, msg.Sender, msg.Receiver, msg.Type, msg.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to save message: %v", err)
	}
	return nil
}

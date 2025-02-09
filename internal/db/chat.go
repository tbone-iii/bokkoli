package db

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// TODO: Find import Message struct solution to avoid circular imports
type Direction string

const (
	Outgoing Direction = "outgoing"
	Incoming Direction = "incoming"
)

type Message struct {
	Text      string    `json:"text"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Direction Direction `json:"direction"`
	Timestamp time.Time `json:"timestamp"`
}

// TODO: End todo ^^

func (handler *DbHandler) setupMessageSchema() error {
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

func (handler *DbHandler) SaveMessage(msg Message) error {
	query := `
	INSERT INTO messages (text, sender, receiver, direction, timestamp)
	VALUES (?, ?, ?, ?, ?);
	`

	_, err := handler.ExecuteQuery(query, msg.Text, msg.Sender, msg.Receiver, msg.Direction, msg.Timestamp)
	return err
}

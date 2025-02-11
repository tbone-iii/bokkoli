package db

import (
	"log"
	"os"
	"slices"
	"testing"
	"time"
)

var dbHandler *DbHandler

// Setup and teardown function for entire test suite
func TestMain(m *testing.M) {
	// Setup
	filePath := "./test-bokkoli.db"
	handler, err := NewDbHandler(filePath)
	if err != nil {
		log.Fatal("Got an error on DB creation: ", err)
	}
	dbHandler = handler
	code := m.Run()

	// Clean up
	if err := dbHandler.Close(); err != nil {
		log.Fatal("DB failed to close.", err)
	}
	removeFile(filePath)

	os.Exit(code)
}

// Useful for cleanup, especially for temporary files
func removeFile(filePath string) {
	if err := os.Remove(filePath); err != nil {
		log.Printf("ERROR: Removing file %q: %v", filePath, err)
	}
}

func TestSetupMessageSchema(t *testing.T) {
	err := dbHandler.setupMessageSchema()
	if err != nil {
		t.Error("Got an error on DB schema setup: ", err)
	}

	// Read the DB file that was created and check that the schema is correct
	rows, err := dbHandler.Query("SELECT * FROM messages")
	if err != nil {
		t.Error("Got error on select statement: ", err)
	}

	cols, err := rows.Columns()
	if err != nil {
		t.Error("Got error on retreiving columns: ", err)
	}

	if !slices.Contains(cols, "text") {
		t.Error("'text' field NOT in list of columns. Database malformed.")
	}
}

func TestSaveMessage(t *testing.T) {
	testText := "Working on it"
	testMessage := Message{
		Text:      testText,
		Sender:    "Sender",
		Direction: Outgoing,
		Timestamp: time.Now(),
	}

	err := dbHandler.setupMessageSchema()
	if err != nil {
		t.Error("Got an error on DB schema setup: ", err)
	}

	err = dbHandler.SaveMessage(testMessage)
	if err != nil {
		t.Error("Saving message produced an error: ", err)
	}

	// Read the DB file that was created and make sure values were added
	rows, err := dbHandler.Query("SELECT text, sender, direction, timestamp FROM messages")
	if err != nil {
		t.Error("Got error on select statement: ", err)
	}

	var tempMessages Message
	for rows.Next() {
		rows.Scan(&tempMessages.Text, &tempMessages.Sender, &tempMessages.Direction, &tempMessages.Timestamp)
		if tempMessages.Text != testText {
			t.Errorf("Saved message for 'text' is incorrect, expected %s, got %s", testText, tempMessages.Text)
		}
	}
}

package message

import (
	"bokkoli/internal/db"
	"log"
	"os"
	"slices"
	"testing"
	"time"
)

// TODO: Use setup/teardown functions to reduce code duplication

func removeTempFile(filePath string) {
	if err := os.Remove(filePath); err != nil {
		log.Printf("ERROR: Removing temp file %q: %v", filePath, err)
	}
}

func TestSetupSchema(t *testing.T) {
	filePath := "./test-bokkoli.db"
	dbHandler, err := db.NewDbHandler(filePath)
	if err != nil {
		t.Error("Got an error on DB creation: ", err)
	}

	err = setupSchema(dbHandler)
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

	// Clean up temp file
	removeTempFile(filePath) // TODO: Get this working
}

func TestSaveMessage(t *testing.T) {
	testText := "Working on it"
	testMessage := Message{
		Text:      testText,
		Sender:    "Somy",
		Receiver:  "Pickle",
		Direction: Outgoing,
		Timestamp: time.Now(),
	}

	filePath := "./test-bokkoli.db"
	dbHandler, err := db.NewDbHandler(filePath)
	if err != nil {
		t.Error("Got an error on DB creation: ", err)
	}

	err = setupSchema(dbHandler)
	if err != nil {
		t.Error("Got an error on DB schema setup: ", err)
	}

	err = saveMessage(dbHandler, testMessage)
	if err != nil {
		t.Error("Saving message produced an error: ", err)
	}

	// Read the DB file that was created and make sure values were added
	rows, err := dbHandler.Query("SELECT text, sender, receiver, direction, timestamp FROM messages")
	if err != nil {
		t.Error("Got error on select statement: ", err)
	}

	var tempMessages Message
	for rows.Next() {
		rows.Scan(&tempMessages.Text, &tempMessages.Sender, &tempMessages.Receiver, &tempMessages.Direction, &tempMessages.Timestamp)
		if tempMessages.Text != testText {
			t.Errorf("Saved message for 'text' is incorrect, expected %s, got %s", testText, tempMessages.Text)
		}
	}

	// Clean up temp file
	removeTempFile(filePath) // TODO: Get this working
}

// package chatActivity

// import (
// 	"database/sql"
// 	"fmt"
// 	"sync"
// )

// type ChatActivity struct {
// 	mu sync.Mutex
// 	db *sql.DB
// }

// const file string = "chatActivity.db"

// func NewChatActivity() (*ChatActivity, error) {
// 	db, err := sql.Open("sqlite3", file)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open database: %v", err)
// 	}
// 	if err := setupSchema(db); err != nil {
// 		return nil, fmt.Errorf("failed fot setup schema: %v", err)
// 	}
// 	return &ChatActivity{db: db}, nil
// }

// func setupSchema(db *sql.DB) error {

// }

package db

import (
	"database/sql"
	"errors"
	"fmt"
)

var Conn *sql.DB

// Initialize all tables and triggers
func InitializeSchema(db *sql.DB) error {
	// Enable foreign-key enforcement
	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return errors.New(fmt.Sprintf("enable foreign_keys: %v\n", err))
	}

	// Analytics
	if _, err := db.Exec(CreateAnalyticsTable); err != nil {
		return errors.New(fmt.Sprintf("create analytics table: %v\n", err))
	}

	// Blogs
	if _, err := db.Exec(CreateBlogsTable); err != nil {
		return errors.New(fmt.Sprintf("create blogs table: %v\n", err))
	}

	return nil
}

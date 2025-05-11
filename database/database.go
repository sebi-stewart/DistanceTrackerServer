package database

import (
	"database/sql"
	"fmt"
)

func InitDatabase(dbConn *sql.DB) error {
	// Create the users table if it doesn't exist
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email VARCHAR(50) NOT NULL UNIQUE,
		name VARCHAR(20) NOT NULL,
		password TEXT NOT NULL,
	    linked_account INTEGER NULL,
	    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
	    modified_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
	    
	    CONSTRAINT fk_linked_account FOREIGN KEY(linked_account) REFERENCES users(id)
	);
	`

	createLinkCodeTable := `
	CREATE TABLE IF NOT EXISTS link_code (
		user_id INTEGER PRIMARY KEY,
		code UUID NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
		
		CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
	)
	`
	_, err := dbConn.Exec(createUsersTable)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	_, err = dbConn.Exec(createLinkCodeTable)
	if err != nil {
		return fmt.Errorf("failed to create link code table: %w", err)
	}

	return nil
}

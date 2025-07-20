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

	createUsersTableIndex := `
	CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
	`

	createLinkCodeTable := `
	CREATE TABLE IF NOT EXISTS link_code (
		user_id INTEGER PRIMARY KEY,
		code UUID NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
		
		CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
	)
	`

	createLinkCodeTableIndex := `
	CREATE INDEX IF NOT EXISTS idx_link_code_code ON link_code (code);
	`

	createLocationsTable := `
	CREATE TABLE IF NOT EXISTS locations (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    user_id INTEGER NOT NULL,
	    latitude REAL NOT NULL,
	    longitude REAL NOT NULL,
	    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
	    is_valid BOOLEAN DEFAULT TRUE NOT NULL,
	    validation_reason TEXT DEFAULT '' NOT NULL,
	    
	    CONSTRAINT fk_user_location FOREIGN KEY(user_id) REFERENCES users(id)
	)
	`

	createLocationsTableIndex := `
	CREATE INDEX IF NOT EXISTS idx_user_id_valid ON locations (user_id, is_valid);
	`

	createRejectedRequestsTable := `
	CREATE TABLE IF NOT EXISTS rejected_requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_email VARCHAR(50) NOT NULL,
		status_code INTEGER NOT NULL,
		reason TEXT NOT NULL,
	    ip_address VARCHAR(45) NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
	)
	`

	createRejectedRequestsTableIndex := `
	CREATE INDEX IF NOT EXISTS idx_rejected_requests_ip_created ON rejected_requests (ip_address, created_at);
	`

	createBannedIpTable := `
	CREATE TABLE IF NOT EXISTS banned_ips (
	 	id INTEGER PRIMARY KEY AUTOINCREMENT,
	 	ip_address VARCHAR(45) NOT NULL UNIQUE,
	    reason TEXT NOT NULL,
	    last_banned_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
	    banned_length REAL NOT NULL,
	    banned_until DATETIME NOT NULL,
	    banned_times INTEGER DEFAULT 1 NOT NULL,
	    
	    CONSTRAINT chk_banned_until CHECK (banned_until > last_banned_at)
	)
	`

	createBannedIpTableIndex := `
	CREATE INDEX IF NOT EXISTS idx_banned_ips_ip ON banned_ips (ip_address);
	`

	_, err := dbConn.Exec(createUsersTable)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	_, err = dbConn.Exec(createUsersTableIndex)
	if err != nil {
		return fmt.Errorf("failed to create index on users table: %w", err)
	}

	_, err = dbConn.Exec(createLinkCodeTable)
	if err != nil {
		return fmt.Errorf("failed to create link code table: %w", err)
	}

	_, err = dbConn.Exec(createLinkCodeTableIndex)
	if err != nil {
		return fmt.Errorf("failed to create index on link code table: %w", err)
	}

	_, err = dbConn.Exec(createLocationsTable)
	if err != nil {
		return fmt.Errorf("failed to create locations table: %w", err)
	}
	_, err = dbConn.Exec(createLocationsTableIndex)
	if err != nil {
		return fmt.Errorf("failed to create index on locations table: %w", err)
	}

	_, err = dbConn.Exec(createRejectedRequestsTable)
	if err != nil {
		return fmt.Errorf("failed to create request rejection table: %w", err)
	}
	_, err = dbConn.Exec(createRejectedRequestsTableIndex)
	if err != nil {
		return fmt.Errorf("failed to create index on rejected requests table: %w", err)
	}

	_, err = dbConn.Exec(createBannedIpTable)
	if err != nil {
		return fmt.Errorf("failed to create banned IP table: %w", err)
	}

	_, err = dbConn.Exec(createBannedIpTableIndex)
	if err != nil {
		return fmt.Errorf("failed to create index on banned IP table: %w", err)
	}

	return nil
}

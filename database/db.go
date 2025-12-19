package database

import (
	"database/sql"
	"fmt"
	"mangadb/config"

	_ "github.com/lib/pq"
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	connStr := cfg.GetConnectionString()

	// Ensure client_encoding is UTF8
	if connStr[len(connStr)-1] != ' ' {
		connStr += " "
	}
	// make suure go understands that these strings must be in utf8 format to suppot Japanese / Korean / Chinese
	// characters
	connStr += "client_encoding=UTF8"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil

}

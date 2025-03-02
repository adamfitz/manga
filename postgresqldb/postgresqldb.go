package postgresqldb

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func OpenDatabase(dbURL string) (*sql.DB, error) {
	/*
		Open a connection to the specified PostgreSQL database and return the *sql.DB instance.
	*/
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL database: %w", err)
	}
	
	// Verify the connection is successful
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL database: %w", err)
	}
	
	return db, nil
}


func InsertDb(pgDB *sql.DB, tableName string, rows []map[string]interface{}) error {
	if len(rows) == 0 {
		return nil
	}

	// Extract column names from the first row
	columns := make([]string, 0, len(rows[0]))
	for col := range rows[0] {
		columns = append(columns, col)
	}

	// Create INSERT query with placeholders
	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, 
		strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	// Insert each row
	for _, row := range rows {
		values := make([]interface{}, len(columns))
		for i, col := range columns {
			values[i] = row[col]
		}

		if _, err := pgDB.Exec(query, values...); err != nil {
			return fmt.Errorf("failed to insert row: %w", err)
		}
	}

	return nil
}

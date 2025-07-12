// manga table code
package postgresqldb

import (
	"database/sql"
	"fmt"
	"log"
	"main/auth"
)

// Add row to manga table
func AddMangaRow(db *sql.DB, name, altTitle, url, mangadexID string, completed, ongoing, hiatus, cancelled *bool) (int64, error) {
	query := `
		INSERT INTO manga (name, alt_name, url, completed, ongoing, hiatus, cancelled)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	// Ensure all 7 parameters are passed, using `nil` for unchecked fields
	var newID int64
	err := db.QueryRow(query, name, altTitle, url, nullableBool(completed), nullableBool(ongoing), nullableBool(hiatus), nullableBool(cancelled)).Scan(&newID)
	if err != nil {
		log.Printf("PG AddMangaRow - failed to insert new row entry %v", err)
		return 0, fmt.Errorf("failed to insert new row entry: %w", err)
	}

	return newID, nil
}

// Search query on manga table for column string.  Return all row/s data if found
func QuerySearchMangaSubstring(db *sql.DB, tableName, columnName, subString string) ([]map[string]any, error) {
	// ILIKE is case insensitive LIKE (search)
	query := fmt.Sprintf("SELECT id, name, alt_name, url, ongoing, completed, hiatus, cancelled FROM %s WHERE %s ILIKE $1", tableName, columnName)
	rows, err := db.Query(query, "%"+subString+"%")
	if err != nil {
		log.Printf("PG QuerySearchMangaSubstring - failed to execute query %v", err)
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var results []map[string]any

	for rows.Next() {
		var id int
		var name string
		var altName, url sql.NullString // Handle NULL values
		var ongoing, completed, hiatus, cancelled sql.NullBool

		err := rows.Scan(&id, &name, &altName, &url, &ongoing, &completed, &hiatus, &cancelled)
		// Check for errors during scanning
		if err != nil {
			log.Printf("PG QuerySearchMangaSubstring - failed to scan row %v", err)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		result := map[string]any{
			"id":        id,
			"name":      name,                    // Guaranteed to be non-NULL
			"alt_name":  altName.String,          // Returns "" if NULL
			"url":       url.String,              // Returns "" if NULL
			"ongoing":   boolToString(ongoing),   // Returns "" if NULL
			"completed": boolToString(completed), // Returns "" if NULL
			"hiatus":    boolToString(hiatus),    // Returns "" if NULL
			"cancelled": boolToString(cancelled), // Returns "" if NULL
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		log.Printf("PG QuerySearchMangaSubstring - row iteration error %v", err)
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

// Lookup and return all rows in manga table
func AllMangaTableRows() ([]map[string]any, error) {

	var allRows []map[string]any

	//load db connection config
	config, _ := auth.LoadConfig()

	// Connect to postgresql db
	pgDb, err := OpenDatabase(
		config.PgServer,
		config.PgPort,
		config.PgUser,
		config.PgPassword,
		config.PgDbName)
	if err != nil {
		return nil, fmt.Errorf("AllMangaTableRows() Database conenction failure, %s", err)
	}

	allRows, err = LookupAllRows(pgDb, "manga")
	if err != nil {
		return nil, fmt.Errorf("AllMangaTableRows() table lookup failure, %s", err)
	}

	return allRows, nil
}

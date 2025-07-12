// mangadex table code
package postgresqldb

import (
	"database/sql"
	"fmt"
	"log"
)

// Add row to mangadex table
func AddMangadexRow(db *sql.DB, name, altTitle, url, mangadexID string, completed, ongoing, hiatus, cancelled *bool) (int64, error) {
	query := `
		INSERT INTO mangadex (name, alt_name, url, mangadex_id, completed, ongoing, hiatus, cancelled)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	// Ensure all 8 parameters are passed, using `nil` for unchecked fields
	var newID int64
	err := db.QueryRow(query, name, altTitle, url, mangadexID,
		nullableBool(completed), nullableBool(ongoing), nullableBool(hiatus), nullableBool(cancelled),
	).Scan(&newID)
	if err != nil {
		log.Printf("PG AddMangadexRow - failed to insert new row entry %v", err)
		return 0, fmt.Errorf("failed to insert new row entry: %w", err)
	}

	return newID, nil
}

// Search query on mangadex table for column string.  Return all row data if found
func QuerySearchMangadexSubstring(db *sql.DB, tableName, columnName, subString string) ([]map[string]any, error) {
	// ILIKE is case insensitive LIKE (search)
	query := fmt.Sprintf("SELECT id, name, alt_name, url, mangadex_id, ongoing, completed, hiatus, cancelled FROM %s WHERE %s ILIKE $1", tableName, columnName)
	rows, err := db.Query(query, "%"+subString+"%")
	if err != nil {
		log.Printf("PG QuerySearchMangadexSubstring - failed to execute query %v", err)
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var results []map[string]any

	for rows.Next() {
		var id int
		var name string
		var altName, url, mangadexID sql.NullString // Handle NULL values
		var ongoing, completed, hiatus, cancelled sql.NullBool

		err := rows.Scan(&id, &name, &altName, &url, &mangadexID, &ongoing, &completed, &hiatus, &cancelled)
		// Check for errors during scanning
		if err != nil {
			log.Printf("PG QuerySearchMangadexSubstring - failed to scan row %v", err)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		result := map[string]any{
			"id":          id,
			"name":        name,                    // Guaranteed to be non-NULL
			"alt_name":    altName.String,          // Returns "" if NULL or false
			"url":         url.String,              // Returns "" if NULL or false
			"mangadex_id": mangadexID.String,       // Returns "" if NULL or false
			"ongoing":     boolToString(ongoing),   // Returns "" if NULL or false
			"completed":   boolToString(completed), // Returns "" if NULL or false
			"hiatus":      boolToString(hiatus),    // Returns "" if NULL or false
			"cancelled":   boolToString(cancelled), // Returns "" if NULL or false
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		log.Printf("PG QuerySearchMangadexSubstring - row iteration error %v", err)
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

// Update applicable mangadex status column
func InsertMangaDexStatus(pgDB *sql.DB, tableName, status, mangadexId string) error {
	// Map statuses to valid column names
	validStatuses := map[string]bool{
		"completed": true,
		"ongoing":   true,
		"cancelled": true,
		"hiatus":    true,
	}

	// Ensure the given status is valid
	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}

	// Construct the SQL query to update the correct column
	query := fmt.Sprintf("UPDATE %s SET %s = TRUE WHERE mangadex_id = $1", tableName, status)

	// Execute the update
	_, err := pgDB.Exec(query, mangadexId)
	if err != nil {
		log.Printf("Error updating manga status: %v", err)
		return fmt.Errorf("failed to update status: %w", err)
	}

	return nil
}

// lightnovel table code
package postgresqldb

import (
	"database/sql"
	"fmt"
	"log"
)

// Light novel substring search (specific table columns are specified in the query)
func LightNovelSearchSubstring(db *sql.DB, tableName, columnName, subString string) ([]map[string]any, error) {
	// ILIKE is case insensitive LIKE (search)
	query := fmt.Sprintf("SELECT id, name, alt_name, url, volumes, completed FROM %s WHERE %s ILIKE $1", tableName, columnName)
	rows, err := db.Query(query, "%"+subString+"%")
	if err != nil {
		log.Printf("PG LightNovelSearchSubstring - failed to execute query %v", err)
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var results []map[string]any

	for rows.Next() {
		var id int
		var name string
		var altName, url sql.NullString // Handle NULL values
		var volumes int64
		var completed sql.NullBool

		err := rows.Scan(&id, &name, &altName, &url, &volumes, &completed)
		// Check for errors during scanning
		if err != nil {
			log.Printf("PG LightNovelSearchSubstring - failed to scan row %v", err)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		result := map[string]any{
			"id":        id,
			"name":      name,           // Guaranteed to be non-NULL
			"alt_name":  altName.String, // Returns "" if NULL
			"url":       url.String,     // Returns "" if NULL
			"volumes":   volumes,        // Returns "" if NULL
			"completed": completed.Bool, // Returns "" if NULL
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		log.Printf("PG LightNovelSearchSubstring - row iteration error %v", err)
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

// Add new row to lightnovel TABLE
func AddLightNovelRow(db *sql.DB, name, altTitle, url string, volumes int, completed *bool) (int64, error) {
	query := `
		INSERT INTO lightnovel (name, alt_name, url, volumes, completed)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	// Ensure all 6 parameters are passed, using `nil` for unchecked fields
	var newID int64
	err := db.QueryRow(query, name, altTitle, url, volumes, nullableBool(completed)).Scan(&newID)
	if err != nil {
		log.Printf("PG AddLightNovelRow - failed to insert new row entry %v", err)
		return 0, fmt.Errorf("failed to insert new row entry: %w", err)
	}

	return newID, nil
}

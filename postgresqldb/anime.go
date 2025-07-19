// anime table code
package postgresqldb

import (
	"database/sql"
	"fmt"
	"log"
)

// Add new row to anime TABLE
func AddAnimeRow(db *sql.DB, name, altTitle, url string, completed, watched *bool) (int64, error) {
	query := `
		INSERT INTO anime (name, alt_name, url, completed, watched)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	// Ensure all 5 parameters are passed, using `nil` for unchecked fields
	var newID int64
	err := db.QueryRow(query, name, altTitle, url,
		nullableBool(completed), nullableBool(watched)).Scan(&newID)
	if err != nil {
		log.Printf("PG AddAnimeRow - failed to insert new row entry %v", err)
		return 0, fmt.Errorf("failed to insert new row entry: %w", err)
	}

	return newID, nil
}

// Anime substring search (specific table columns are specified in the query)
func AnimeSearchSubstring(db *sql.DB, tableName, columnName, subString string) ([]map[string]any, error) {
	// ILIKE is case insensitive LIKE (search)
	query := fmt.Sprintf("SELECT id, name, alt_name, url, completed, watched FROM %s WHERE %s ILIKE $1", tableName, columnName)
	rows, err := db.Query(query, "%"+subString+"%")
	if err != nil {
		log.Printf("PG AnimeSearchSubstring - failed to execute query %v", err)
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var results []map[string]any

	for rows.Next() {
		var id int
		var name string
		var altName, url sql.NullString // Handle NULL values
		var completed, watched sql.NullBool

		err := rows.Scan(&id, &name, &altName, &url, &completed, &watched)
		// Check for errors during scanning
		if err != nil {
			log.Printf("PG AnimeSearchSubstring - failed to scan row %v", err)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		result := map[string]any{
			"id":        id,
			"name":      name,           // Guaranteed to be non-NULL
			"alt_name":  altName.String, // Returns "" if NULL
			"url":       url.String,     // Returns "" if NULL
			"completed": completed.Bool, // Returns "" if NULL
			"watched":   watched.Bool,   // Returns "" if NULL
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		log.Printf("PG QuerySearchSubstring - row iteration error %v", err)
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

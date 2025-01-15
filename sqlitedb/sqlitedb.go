package sqlitedb

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// OpenDatabase opens a connection to the SQLite database and returns the *sql.DB instance.
func OpenDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return db, nil
}

// QueryRow retrieves a single row from the database.
func QueryRow(db *sql.DB, query string, args ...interface{}) (map[string]interface{}, error) {
	row := db.QueryRow(query, args...)

	// Prepare a statement to get column names
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	columns, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer columns.Close()

	columnNames, err := columns.Columns()
	if err != nil {
		return nil, err
	}

	// Create a container for the values
	values := make([]interface{}, len(columnNames))
	valuePtrs := make([]interface{}, len(columnNames))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Scan the row
	err = row.Scan(valuePtrs...)
	if err != nil {
		return nil, err
	}

	// Map the column names to their values
	result := make(map[string]interface{})
	for i, colName := range columnNames {
		result[colName] = values[i]
	}

	return result, nil
}

func QueryWithCondition(db *sql.DB, tableName, columnName, condition string) (map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", tableName, columnName)
	rows, err := db.Query(query, condition)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Retrieve column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	// Create a slice to hold the results
	var results []map[string]interface{}

	// Iterate through the rows
	for rows.Next() {
		// Create a container for the values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the row into the values container
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Map the column names to their values
		result := make(map[string]interface{})
		for i, colName := range columns {
			result[colName] = values[i]
		}

		// Append the result to the list
		results = append(results, result)
	}

	// Check if any error occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed during row iteration: %w", err)
	}

	// Return only the first result (map) if there is at least one
	if len(results) > 0 {
		return results[0], nil
	}

	// If no rows were found, return nil
	return nil, nil
}

// MangaDexLookupChapterList retrieves the JSON string array from the "mangadex_ch_list" column
// for a given name in the SQLite database.
func MangaDexLookupChapterList(db *sql.DB, name string) (string, error) {
	// Query to select the "mangadex_ch_list" column based on the "name"
	query := `SELECT mangadex_ch_list FROM chapters WHERE name = ?`
	var jsonString string

	// Execute the query
	err := db.QueryRow(query, name).Scan(&jsonString)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no entry found for name: %s", name)
		}
		return "", err
	}

	return jsonString, nil
}

// Initial function to add the mangadex json string to the correct  database column
// UpdateMangaDexChapterList updates the mangadex_ch_list column in the chapters table
// with the provided JSON string for the corresponding name.
func MangaDexInitialDbChapterListUpdate(db *sql.DB, name string, jsonString string) error {
	// SQL query to update the mangadex_ch_list column
	query := `UPDATE chapters SET mangadex_ch_list = ? WHERE name = ?`

	// Execute the query
	result, err := db.Exec(query, jsonString, name)
	if err != nil {
		return fmt.Errorf("failed to update mangadex_ch_list: %v", err)
	}

	// Check if any row was affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated; name '%s' may not exist", name)
	}

	return nil
}

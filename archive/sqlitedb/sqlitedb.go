package sqlitedb

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func OpenDatabase(dbPath string) (*sql.DB, error) {
	/*
		Open a connection to the sprecified SQLite database and return the *sql.DB instance.
	*/
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return db, nil
}

func QueryRow(db *sql.DB, query string, args ...interface{}) (map[string]interface{}, error) {
	/*
		Perform a DB lookup for specific row and return the result as a map.
	*/
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
	/*
		Perform a gerneric DB lookup for specific row and return the result as a map, based on the provided column name
		and condition eg: name
	*/
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

func MangadexDbLookupChapterList(db *sql.DB, name string) (string, error) {
	/*
		Perform a DB lookup for a given string in the name column (saerch by name) and return the JSON string array
		(if it exists) from the mangadex_ch_list column
	*/

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

func MangadexInitialDbChapterListUpdate(db *sql.DB, name string, jsonString string) error {
	/*
		func will insert (update/overwrite) the mangadex_ch_list column in the chapters table with the provided JSON
		string for the provided name.
	*/

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

func MangadexIdDbLookup(db *sql.DB, name, tableName string) (string, error) {
	/*
		func performs a lookup in the provided table for the provided name and returns the mangadex_id in string format.
	*/

	// Use a parameterized query to avoid SQL injection
	query := fmt.Sprintf("SELECT mangadex_id FROM %s WHERE name = ?", tableName)

	stmt, err := db.Prepare(query)
	if err != nil {
		return "", fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	// Execute the query with the provided name
	var mangadexId string
	err = stmt.QueryRow(name).Scan(&mangadexId)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("no matching record found for name: %s", name)
	} else if err != nil {
		return "", fmt.Errorf("failed to execute query: %w", err)
	}

	return mangadexId, nil
}

func MangaNameDbLookup(db *sql.DB, name, tableName string) (bool, error) {
	/*
		This function performs a lookup in the provided table for the provided name.

		NOTE: This function exists to perform a comparison with bookmark names and database names.
	*/

	// check if the record exists in DB dont retrieve the name column
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE name = ?", tableName)

	// Prepare the statement
	stmt, err := db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	// Execute the query with the provided name
	var exists int
	err = stmt.QueryRow(name).Scan(&exists)
	if err == sql.ErrNoRows {
		// No matching record found
		return false, nil
	} else if err != nil {
		// Other errors
		return false, fmt.Errorf("failed to execute query: %w", err)
	}

	// Record exists
	return true, nil
}

func AddMangaEntry(db *sql.DB, name, altTitle, url, mangadexID string) (int64, error) {
	/*
		Function to add a new entry to the `chapters` table and return the ID of the newly inserted row.
		- `name` goes into the `name` column.
		- `altTitle` goes into the `alt_name` column.
		- `url` goes into the `url` column.
		- `mangadexID` goes into the `mangadex_id` column.
		- Other columns are left empty.
	*/

	// SQL query to insert a new row into the chapters table
	query := `
		INSERT INTO chapters (name, alt_name, url, mangadex_id) 
		VALUES (?, ?, ?, ?)
	`

	// Execute the query
	result, err := db.Exec(query, name, altTitle, url, mangadexID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert new chapter entry: %v", err)
	}

	// Retrieve the ID of the newly inserted row
	newID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get the ID of the new chapter entry: %v", err)
	}

	return newID, nil
}

func QueryAllMangadexNames(db *sql.DB) []string {
	/*
		Function to query for all the names of mangas that have a Null or empty mangadex_ch_list and contains the
		substring "mangadex" in the url column.

		Returns a slice of strings containing the names matching the above conditions.
	*/

	// SQL query to select all names
	query := `SELECT name FROM chapters WHERE (mangadex_ch_list IS NULL OR mangadex_ch_list = '') AND url LIKE '%mangadex%' ORDER BY name ASC;`

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("Error querying all manga names: %v\n", err)
		return nil
	}
	defer rows.Close()

	// Create a slice to hold the manga names
	var mangaNames []string

	// Iterate over the rows and append each name to the slice
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			continue
		}
		mangaNames = append(mangaNames, name)
	}

	// Check for any error that occurred during iteration
	if err := rows.Err(); err != nil {
		fmt.Printf("Error during row iteration: %v\n", err)
	}

	return mangaNames
}

func QuerySearchSubstring(db *sql.DB, tableName, columnName, subString string) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT id, name, alt_name, url, mangadex_id FROM %s WHERE %s LIKE ?", tableName, columnName)
	rows, err := db.Query(query, "%"+subString+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}

	for rows.Next() {
		var id int
		var name, altName, url, mangadexID string

		err := rows.Scan(&id, &name, &altName, &url, &mangadexID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		result := map[string]interface{}{
			"id":          id,
			"name":        name,
			"alt_name":    altName,
			"url":         url,
			"mangadex_id": mangadexID,
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

func QueryByID(db *sql.DB, tableName string, id int64) (map[string]interface{}, error) {
	/*
		Query the table by the specified ID and return the entry as a map[string]interface{}.
		This function is tailored to retrieve a single row by its ID.
	*/

	// Prepare the query
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", tableName)

	// Execute the query
	row := db.QueryRow(query, id)

	// Get column names
	columnNames, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", tableName))
	if err != nil {
		return nil, fmt.Errorf("failed to get column names for table %s: %v", tableName, err)
	}
	defer columnNames.Close()

	// Collect column names
	columns := []string{}
	for columnNames.Next() {
		var cid int
		var name string
		var ctype string
		var notnull int
		var dfltValue interface{}
		var pk int
		if err := columnNames.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); err != nil {
			return nil, fmt.Errorf("failed to parse column info: %v", err)
		}
		columns = append(columns, name)
	}

	// Prepare storage for the row values
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Scan the result
	if err := row.Scan(valuePtrs...); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no row found with id %d", id)
		}
		return nil, fmt.Errorf("failed to scan row: %v", err)
	}

	// Map the result
	result := make(map[string]interface{})
	for i, col := range columns {
		val := values[i]
		if b, ok := val.([]byte); ok {
			result[col] = string(b)
		} else {
			result[col] = val
		}
	}

	return result, nil
}

// extract all rows from the database
func QueryAllRows(db *sql.DB, tableName string) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		rowMap := make(map[string]interface{})
		for i, colName := range columns {
			rowMap[colName] = values[i]
		}

		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return results, nil
}

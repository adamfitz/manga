package postgresqldb

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// table name comes from an untrusted source (user input) so this map is used to validate the table name
var allowedTables = map[string]bool{"mangadex": true, "manga": true}

/*
Open a connection to a remote PostgreSQL database using host, port, user, password, and database name.
*/
func OpenDatabase(dbHost, dbPort, dbUser, dbPassword, dbName string) (*sql.DB, error) {
	dBSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	pgDb, err := sql.Open("postgres", dBSourceName)
	if err != nil {
		log.Printf("PG OpenDatabase - failed to open PostgreSQL database: %v", err)
		return nil, fmt.Errorf("failed to open PostgreSQL database: %w", err)
	}

	// Verify the connection is successful
	if err := pgDb.Ping(); err != nil {
		log.Printf("PG OpenDatabase - failed to ping PostgreSQL database: %v", err)
		return nil, fmt.Errorf("failed to ping PostgreSQL database: %w", err)
	}

	return pgDb, nil
}

// insert rows into database
func InsertRow(pgDB *sql.DB, tableName string, rows []map[string]interface{}) error {
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
			log.Printf("PG InsertRow - failed to insert row: %v", err)
			return fmt.Errorf("failed to insert row: %w", err)
		}
	}

	return nil
}

// Lookup and return row from DB table
func LookupRow(pgDB *sql.DB, tableName string, conditions map[string]any) ([]byte, error) {
	if len(conditions) == 0 {
		return nil, errors.New("no conditions provided for query")
	}

	// Build query condition string
	conditionClauses := make([]string, 0, len(conditions))
	values := make([]any, 0, len(conditions))
	i := 1

	for col, val := range conditions {
		conditionClauses = append(conditionClauses, fmt.Sprintf("%s = $%d", col, i))
		values = append(values, val)
		i++
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT 1", tableName, strings.Join(conditionClauses, " AND "))

	// Execute query
	row := pgDB.QueryRow(query, values...)

	// Retrieve column names
	cols, err := pgDB.Query(fmt.Sprintf("SELECT column_name FROM information_schema.columns WHERE table_name='%s'", tableName))
	if err != nil {
		log.Printf("PG LookupRow - failed to retrieve column names: %v", err)
		return nil, fmt.Errorf("failed to retrieve column names: %w", err)
	}
	defer cols.Close()

	columnNames := []string{}
	for cols.Next() {
		var colName string
		if err := cols.Scan(&colName); err != nil {
			log.Printf("PG LookupRow - failed to scan column names: %v", err)
			return nil, fmt.Errorf("failed to scan column names: %w", err)
		}
		columnNames = append(columnNames, colName)
	}

	// Prepare storage for scanned values
	columnPointers := make([]interface{}, len(columnNames))
	columnValues := make([]interface{}, len(columnNames))
	for i := range columnValues {
		columnPointers[i] = &columnValues[i]
	}

	// Scan row into column values
	if err := row.Scan(columnPointers...); err != nil {
		log.Printf("PG LookupRow - failed to scan row: %v", err)
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	// Convert to map
	result := make(map[string]interface{})
	for i, colName := range columnNames {
		result[colName] = columnValues[i]
	}

	// Convert result to JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Printf("PG LookupRow - failed to convert row to JSON: %v", err)
		return nil, fmt.Errorf("failed to convert row to JSON: %w", err)
	}

	return jsonData, nil
}

// Query all data from the specified PostgreSQL table and return the results as a slice of maps.
func LookupAllRows(db *sql.DB, tableName string) ([]map[string]any, error) {

	query := fmt.Sprintf("SELECT * FROM %s", tableName)

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("PG LookupAllRows - failed to execute query %v", err)
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.Printf("PG LookupAllRows - failed to get columns %v", err)
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var results []map[string]any

	// Iterate over the rows
	for rows.Next() {
		// Create a slice for values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the row into the value pointers
		if err := rows.Scan(valuePtrs...); err != nil {
			log.Printf("PG LookupAllRows - failed to scan row %v", err)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Map column names to values
		rowMap := make(map[string]interface{})
		for i, colName := range columns {
			rowMap[colName] = values[i]
		}

		results = append(results, rowMap)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		log.Printf("PG LookupAllRows - failed during row iteration %v", err)
		return nil, fmt.Errorf("failed during row iteration: %w", err)
	}

	// Sort the slice alphabetically based on the "name" key
	sort.Slice(results, func(i, j int) bool {
		return results[i]["name"].(string) < results[j]["name"].(string)
	})

	return results, nil
}

/*
Query the table by the specified ID and return the entry as a map[string]interface{}.
This function is tailored to retrieve a single row by its ID.
*/
func LookupByID(db *sql.DB, tableName string, id string) (map[string]any, error) {

	// Prepare the query
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", tableName)

	// Execute the query
	row := db.QueryRow(query, id)

	// Get column names from the PostgreSQL catalog
	columnNamesQuery := `
		SELECT column_name
		FROM information_schema.columns
		WHERE table_name = $1
	`
	columnNamesRows, err := db.Query(columnNamesQuery, tableName)
	if err != nil {
		log.Printf("PG LookupByID - failed to get column names for table %v", err)
		return nil, fmt.Errorf("failed to get column names for table %s: %v", tableName, err)
	}
	defer columnNamesRows.Close()

	// Collect column names
	columns := []string{}
	for columnNamesRows.Next() {
		var name string
		if err := columnNamesRows.Scan(&name); err != nil {
			log.Printf("PG LookupByID - failed to gparse column info %v", err)
			return nil, fmt.Errorf("failed to parse column info: %v", err)
		}
		columns = append(columns, name)
	}

	// Prepare storage for the row values
	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Scan the result
	if err := row.Scan(valuePtrs...); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("PG LookupByID - no row found with id %v", err)
			return nil, fmt.Errorf("no row found with id %v", id)
		}
		log.Printf("PG LookupByID - failed to scan row %v", err)
		return nil, fmt.Errorf("failed to scan row: %v", err)
	}

	// Map the result
	result := make(map[string]any)
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

// Perform a gerneric DB table lookup for specific row and return the result as a map, based on the provided column name
// and condition eg: name
func QueryWithCondition(db *sql.DB, tableName, columnName, condition string) (map[string]any, error) {

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", tableName, columnName)
	rows, err := db.Query(query, condition)
	if err != nil {
		log.Printf("PG QueryWithCondition - failed to execute query %v", err)
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Retrieve column names
	columns, err := rows.Columns()
	if err != nil {
		log.Printf("PG QueryWithCondition - failed to get columns %v", err)
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
			log.Printf("PG QueryWithCondition - failed to scan row %v", err)
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
		log.Printf("PG QueryWithCondition - failed during row iteration %v", err)
		return nil, fmt.Errorf("failed during row iteration: %w", err)
	}

	// Return only the first result (map) if there is at least one
	if len(results) > 0 {
		return results[0], nil
	}

	// If no rows were found, return nil
	return nil, nil
}

/*
This function performs a lookup in the provided table for the provided name.

NOTE: This function exists to perform a comparison with bookmark names and database names.
*/
func LookupByName(db *sql.DB, name, tableName string) (bool, error) {
	// Validate the table name (exists in allowed tables)
	if !allowedTables[tableName] {
		log.Printf("Illegal table name, validation failed: %s", tableName)
		return false, fmt.Errorf("invalid table name")
	}

	// check if the record exists in DB dont retrieve the name column
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE name = $1", tableName)

	// Prepare the statement
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("PG QueryByName - failed to prepare query %v", err)
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
		log.Printf("PG QueryByName - failed to execute query %v", err)
		return false, fmt.Errorf("failed to execute query: %w", err)
	}

	// Record exists
	return true, nil
}

// Helper function to handle *bool -> SQL NULL conversion
func nullableBool(b *bool) any {
	if b == nil {
		return nil // Store as NULL in database
	}
	return *b // Store TRUE if checked
}

// Perform a lookup for a specific column value based on the provided condition.
func LookupColumnValues(db *sql.DB, tableName, columnName string) ([]string, error) {

	// Validate the table name (exists in allowed tables)
	if !allowedTables[tableName] {
		log.Printf("Illegal table name, validation failed: %s", tableName)
		return nil, fmt.Errorf("invalid table name")
	}

	query := fmt.Sprintf("SELECT %s FROM %s", columnName, tableName)
	result, err := db.Query(query)
	if err != nil {
		log.Printf("PG LookupColumnValues - failed to execute query %v", err)
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer result.Close()

	var results []string

	// Iterate over rows and scan into the slice
	for result.Next() {
		var value string
		if err := result.Scan(&value); err != nil {
			log.Printf("PG LookupColumnValues - failed to scan row: %v", err)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, value)
	}

	// Check for iteration errors
	if err := result.Err(); err != nil {
		log.Printf("PG LookupColumnValues - error iterating rows: %v", err)
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return results, nil
}

// Extract values from multiple columns at the same time
/*
	Perform a lookup for multiple column values and return each row as a string.
*/
func LookupMultipleColumnValues(db *sql.DB, tableName string, columnNames ...string) ([]string, error) {

	// Validate the table name (exists in allowed tables)
	if !allowedTables[tableName] {
		log.Printf("Illegal table name, validation failed: %s", tableName)
		return nil, fmt.Errorf("invalid table name")
	}

	// Ensure at least one column name is provided
	if len(columnNames) == 0 {
		return nil, fmt.Errorf("no column names provided")
	}

	// Construct the query dynamically
	query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(columnNames, ", "), tableName)
	result, err := db.Query(query)
	if err != nil {
		log.Printf("PG LookupMultipleColumnValues - failed to execute query %v", err)
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer result.Close()

	var rows []string

	// Iterate over rows
	for result.Next() {
		// Create a slice of interface{} to hold scanned values
		values := make([]any, len(columnNames))
		valuePtrs := make([]any, len(columnNames))

		// Assign pointers to the interface slice
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the row
		if err := result.Scan(valuePtrs...); err != nil {
			log.Printf("PG LookupMultipleColumnValues - failed to scan row: %v", err)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Convert row data into a single string
		var rowStrings []string
		for _, val := range values {
			if val != nil {
				rowStrings = append(rowStrings, fmt.Sprintf("%v", val)) // Convert to string
			} else {
				rowStrings = append(rowStrings, "NULL") // Handle NULL values
			}
		}

		rows = append(rows, strings.Join(rowStrings, " ")) // Join values with a space
	}

	// Check for iteration errors
	if err := result.Err(); err != nil {
		log.Printf("PG LookupMultipleColumnValues - error iterating rows: %v", err)
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return rows, nil
}

// Perform DB table lookup by manga status
func LookupByStatus(db *sql.DB, tableName string, statusColumn string) ([]map[string]any, error) {
	// Allowlist of valid status columns
	validColumns := map[string]bool{
		"completed": true,
		"hiatus":    true,
		"ongoing":   true,
		"cancelled": true,
	}

	if !validColumns[statusColumn] {
		return nil, fmt.Errorf("invalid status column: %s", statusColumn)
	}

	// Construct the query safely
	query := fmt.Sprintf(`SELECT name, alt_name, mangadex_id FROM %s WHERE %s = TRUE`, tableName, statusColumn)

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("PG LookupByStatus - query execution failed: %v", err)
		return nil, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	var results []map[string]any

	for rows.Next() {
		var name, altName, mangadexID sql.NullString
		if err := rows.Scan(&name, &altName, &mangadexID); err != nil {
			log.Printf("PG LookupByStatus - row scan failed: %v", err)
			return nil, fmt.Errorf("row scan failed: %v", err)
		}

		entry := map[string]any{
			"name":        name.String,
			"alt_name":    altName.String,
			"mangadex_id": mangadexID.String,
		}
		results = append(results, entry)
	}

	if err := rows.Err(); err != nil {
		log.Printf("PG LookupByStatus - rows iteration error: %v", err)
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return results, nil
}

// Perform DB table lookup by name or alt_name and returns the status of the manga eg: ongoing, completed, hiatus or cancelled
func LookupByNameOrAltName(db *sql.DB, tableName string, searchColumn string, value string) (map[string]any, error) {
	// Allowlist to prevent SQL injection
	validColumns := map[string]bool{
		"name":     true,
		"alt_name": true,
	}

	if !validColumns[searchColumn] {
		return nil, fmt.Errorf("invalid search column: %s", searchColumn)
	}

	// Query to select needed fields
	query := fmt.Sprintf(`
		SELECT name, alt_name, mangadex_id, completed, ongoing, hiatus, cancelled
		FROM %s
		WHERE %s = $1
		LIMIT 1
	`, tableName, searchColumn)

	row := db.QueryRow(query, value)

	// Fields for scanning
	var name, altName, mangadexID sql.NullString
	var completed, ongoing, hiatus, cancelled sql.NullBool

	err := row.Scan(&name, &altName, &mangadexID, &completed, &ongoing, &hiatus, &cancelled)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no record found where %s = '%s'", searchColumn, value)
		}
		return nil, fmt.Errorf("failed to scan row: %v", err)
	}

	// Determine the status from the boolean flags
	status := ""
	switch {
	case completed.Valid && completed.Bool:
		status = "completed"
	case ongoing.Valid && ongoing.Bool:
		status = "ongoing"
	case hiatus.Valid && hiatus.Bool:
		status = "hiatus"
	case cancelled.Valid && cancelled.Bool:
		status = "cancelled"
	default:
		status = "unknown"
	}

	// Build result map
	result := map[string]any{
		"name":        name.String,
		"alt_name":    altName.String,
		"mangadex_id": mangadexID.String,
		"status":      status,
	}

	return result, nil
}

/*
Convert to emptry string if false returned (to dispaly nothing when the page is rendered
*/
func boolToString(b sql.NullBool) string {
	if b.Valid && b.Bool {
		return "true"
	}
	return ""
}

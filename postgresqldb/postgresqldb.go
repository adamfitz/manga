package postgresqldb

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func OpenDatabase(dbHost, dbPort, dbUser, dbPassword, dbName string) (*sql.DB, error) {
	/*
		Open a connection to a remote PostgreSQL database using host, port, user, password, and database name.
	*/
	dBSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	pgDb, err := sql.Open("postgres", dBSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL database: %w", err)
	}

	// Verify the connection is successful
	if err := pgDb.Ping(); err != nil {
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
			return fmt.Errorf("failed to insert row: %w", err)
		}
	}

	return nil
}

// Lookup and return row from database
func LookupRow(pgDB *sql.DB, tableName string, conditions map[string]interface{}) ([]byte, error) {
	if len(conditions) == 0 {
		return nil, errors.New("no conditions provided for query")
	}

	// Build query condition string
	conditionClauses := make([]string, 0, len(conditions))
	values := make([]interface{}, 0, len(conditions))
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
		return nil, fmt.Errorf("failed to retrieve column names: %w", err)
	}
	defer cols.Close()

	columnNames := []string{}
	for cols.Next() {
		var colName string
		if err := cols.Scan(&colName); err != nil {
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
		return nil, fmt.Errorf("failed to convert row to JSON: %w", err)
	}

	return jsonData, nil
}

// query all data in a table
func LookupAllRows(db *sql.DB, tableName string) ([]map[string]interface{}, error) {
	/*
		Query all data from the specified PostgreSQL table and return the results as a slice of maps.
	*/
	query := fmt.Sprintf("SELECT * FROM %s;", tableName)

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var results []map[string]interface{}

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
		return nil, fmt.Errorf("failed during row iteration: %w", err)
	}

	// Sort the slice alphabetically based on the "name" key
	sort.Slice(results, func(i, j int) bool {
		return results[i]["name"].(string) < results[j]["name"].(string)
	})

	return results, nil
}

func LookupByID(db *sql.DB, tableName string, id string) (map[string]interface{}, error) {
	/*
		Query the table by the specified ID and return the entry as a map[string]interface{}.
		This function is tailored to retrieve a single row by its ID.
	*/

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
		return nil, fmt.Errorf("failed to get column names for table %s: %v", tableName, err)
	}
	defer columnNamesRows.Close()

	// Collect column names
	columns := []string{}
	for columnNamesRows.Next() {
		var name string
		if err := columnNamesRows.Scan(&name); err != nil {
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
			return nil, fmt.Errorf("no row found with id %v", id)
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

package actions

import (
	"fmt"
	"log"
	"main/auth"
	"main/bookmarks"
	"main/parser"
	"main/postgresqldb"
	"sort"
)

func CompareNames() {
	/*
		compares the manga names in bookmarks to the names in the database, returns:
		- names from the bookmarks file that are not in the DB
		- names from the DB that are not in the bookmarks file
	*/

	//load db connection config
	config, _ := auth.LoadConfig()

	// 1 - Load bookmarks
	bookmarksFromFile, err := bookmarks.LoadBookmarks()
	if err != nil {
		log.Fatalf("Error loading bookmarks: %v", err)
	}
	// 2 - Get a list of the titles with "mangadex" connector from bookmarks
	bookmarkNames := bookmarks.MangadexBookmarks(bookmarksFromFile)

	// 3 - get all the DB names
	dbConnection, _ := postgresqldb.OpenDatabase(config.PgServer, config.PgPort, config.PgUser, config.PgPassword, config.PgDbName)
	dbNames, _ := postgresqldb.LookupColumnValues(dbConnection, "mangadex", "name")

	// Sort both slices
	sort.Strings(bookmarkNames)
	sort.Strings(dbNames)

	// slices to store missing names
	var missingInDB []string        // Bookmarks not in DB
	var missingInBookmarks []string // DB names not in bookmarks

	// Use two pointers to compare sorted slices
	i, j := 0, 0
	for i < len(bookmarkNames) && j < len(dbNames) {
		if bookmarkNames[i] == dbNames[j] {
			// Name exists in both
			i++
			j++
		} else if bookmarkNames[i] < dbNames[j] {
			// bookmarkNames[i] is missing in the DB
			missingInDB = append(missingInDB, bookmarkNames[i])
			i++
		} else {
			// dbNames[j] is missing in the bookmarks
			missingInBookmarks = append(missingInBookmarks, dbNames[j])
			j++
		}
	}

	// Add remaining unmatched elements
	for i < len(bookmarkNames) {
		missingInDB = append(missingInDB, bookmarkNames[i])
		i++
	}
	for j < len(dbNames) {
		missingInBookmarks = append(missingInBookmarks, dbNames[j])
		j++
	}

	// Print results
	if len(missingInDB) > 0 {
		fmt.Println("\n--- Bookmarks missing from the database ---")
		for _, name := range missingInDB {
			fmt.Println(name)
		}
	} else {
		fmt.Println("All bookmarks are present in the database.")
	}

	if len(missingInBookmarks) > 0 {
		fmt.Println("\n--- DB names missing from the bookmarks file ---")
		for _, name := range missingInBookmarks {
			fmt.Println(name)
		}
	} else {
		fmt.Println("All database entries are present in the bookmarks file.")
	}
}

func DumpPostgressTable(tableName string, columns []string) {
	/*
		Dumps the postgresql table.
	*/

	// Load db connection config
	config, _ := auth.LoadConfig()

	// Connect to PostgreSQL db
	pgDb, err := postgresqldb.OpenDatabase(
		config.PgServer,
		config.PgPort,
		config.PgUser,
		config.PgPassword,
		config.PgDbName)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer pgDb.Close()

	// Get all data in PostgreSQL table
	data, err := postgresqldb.LookupAllRows(pgDb, tableName)
	if err != nil {
		log.Fatalf("Error querying data: %v", err)
	}

	// Convert the columns slice to a map for quick lookup
	columnFilter := make(map[string]struct{}, len(columns))
	for _, col := range columns {
		columnFilter[col] = struct{}{}
	}

	// Iterate over each row in the data slice
	for _, row := range data {
		// Create a slice to hold the keys
		var keys []string
		for key := range row {
			if _, ok := columnFilter[key]; ok {
				keys = append(keys, key)
			}
		}

		// Sort the selected keys
		sort.Strings(keys)

		// Print the selected key-value pairs
		for _, key := range keys {
			fmt.Printf("\n%s:\t%v", key, row[key])
		}
		fmt.Println("\n----------------------------------------")
	}
}

// Get a list of all directories from the provided rootDir.
func DirList(rootDir string, exclusionList ...string) ([]string, error) {
	dirListing, err := parser.DirList(rootDir, exclusionList...) // Exclusion list is optional, indicated by the variadic parameter
	if err != nil {
		log.Fatalf("Error getting directory list: %v", err)
		return nil, err
	}

	return dirListing, nil
}

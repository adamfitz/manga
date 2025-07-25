package actions

import (
	"fmt"
	"log"
	"main/auth"
	"main/bookmarks"
	"main/mangadex"
	"main/postgresqldb"
)

// Compare the managa names in bookmarks to the names in the database, prints out the difference if the name does
// not exist in the DB.
func CompareBookmarksAndDB() {

	//load db connection config
	config, _ := auth.LoadConfig()

	// 1 - Load bookmarks
	bookmarksFromFile, err := bookmarks.LoadBookmarks()
	if err != nil {
		log.Fatalf("Error loading bookmarks: %v", err)
	}
	// 2 - Get a list of the titles with "mangadex" connector from bookmarks
	bookmarkNames := bookmarks.MangadexBookmarks(bookmarksFromFile)

	// open the database
	dbConnection, _ := postgresqldb.OpenDatabase(config.PgServer, config.PgPort, config.PgUser, config.PgPassword, config.PgDbName)
	// iterate of the names of the mangas in the bookmark list
	for _, name := range bookmarkNames {

		// a. extract the mangadex id from the mangadex table based on the manga name
		dbNames, _ := postgresqldb.LookupByName(dbConnection, name, "mangadex")

		if !dbNames {
			fmt.Printf("Bookmark not in DB: %s\n", name)
		}
	}
}

// Return the manga status attirbutes from the mangadex API and write them to the DB
func MangaStatusAttributes() {

	//load db connection config
	config, _ := auth.LoadConfig()

	// Connect to postgresql db
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

	// get all the mangadex ids from the (mangadex_id column)
	mangadexIds, _ := postgresqldb.LookupColumnValues(pgDb, "mangadex", "mangadex_id")

	// for each mangadex_id in the table, lookup the manga status and convert the API response to a string, write
	// the string into the table in the applicable column
	for _, id := range mangadexIds {
		response, _ := mangadex.MangaAttributes(id)
		status := mangadex.MangaStatus(response)
		postgresqldb.InsertMangaDexStatus(pgDb, "mangadex", status, id)
	}

	// slice of all the columns to perform the lookup on
	columns := []string{"name", "completed", "hiatus", "ongoing", "cancelled"}
	// get all the manga names from the (name column)
	outputList, _ := postgresqldb.LookupMultipleColumnValues(pgDb, "mangadex", columns...)

	// show me the DB table output for all statues columns by name
	for _, columnValues := range outputList {
		fmt.Println(columnValues)
	}
}

// Return all entries in name column
func MangadexNames() ([]string, error) {
	//load db connection config
	config, _ := auth.LoadConfig()

	// Connect to postgresql db
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

	// all entry names
	mangadexNames, err := postgresqldb.LookupColumnValues(pgDb, "mangadex", "name")
	if err != nil {
		log.Println("error retrieving mangadex table column n\"name\"", err)
		return nil, fmt.Errorf("error retrieving mangadex table column n\"name\" %v", err)
	}

	return mangadexNames, nil
}

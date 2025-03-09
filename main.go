package main

import (
	//"encoding/json"
	"fmt"
	//"strings"
	"log"
	"main/auth"
	"main/bookmarks"
	"main/mangadex"
	//"main/compare"
	//"main/mangadex"
	"main/postgresqldb"
	//"main/sqlitedb"
	//"main/webfrontend"
	"os"
	"sort"
)

func init() {
	logFile, err := os.OpenFile("manga.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)
}

func main() {

	MangaAttributes()
	//NewMangaDbUpdate()
	//CheckIfBookmarkInDb()
	//CompareNames()
	//BlanketUpdateDb()
	//ExtractMangasWithoutChapterList()
	//UpdateMangasWithoutChapterList()
	//webfrontend.StartServer("8080")
	//DumpPostgressDb()
	//PgQueryByID("21")
}

/*
func CheckForNewChapters() {


	//	func gets a list of chapters from mangadex and comnpares it to the list in the database, outputting the differences
	//	(if any).


	// 1 - Load bookmarks
	bookmarksFromFile, err := bookmarks.LoadBookmarks()
	if err != nil {
		log.Fatalf("Error loading bookmarks: %v", err)
	}

	// 2 - Get a list of the titles with "mangadex" connector
	names := bookmarks.MangadexMangaTitles(bookmarksFromFile)

	// 3 - performing chapter list comparison between mangadex.org and the local DB
	fmt.Println("== Performing chapter list comparison between local DB and mangadex.org ==:")

	// open the database
	dbConnection, _ := sqlitedb.OpenDatabase("database/mangaList.db")

	// iterate of the names of the mangas in the bookmark list
	for _, name := range names {

		// a. extract the mangadex id from the database based on the manga name
		mangadexId, _ := sqlitedb.MangadexIdDbLookup(dbConnection, name, "chapters")

		// b. get the list of chapters from mangadex
		chapterList, _ := mangadex.ChaptersSorted(mangadexId)

		// c. extract the list of chapters from the database
		chapterListDb, _ := sqlitedb.MangadexDbLookupChapterList(dbConnection, name)

		// d. perform a comparison of both of the json arrays, returning the difference (compares A to B) and returns
		// the elements from list A that are not present in list B
		differences, _ := compare.CompareJSONArrays(chapterList, chapterListDb)

		// Output the differences
		if len(differences) != 0 {
			fmt.Printf("New chapters available on mangadex for\t\t %s\n", name)
			for _, diff := range differences {
				fmt.Println(diff)
			}
		}
	}
}
*/

/*
func BlanketUpdateDb() {

	//	func gets a sorted list of all the chapters from mangadex and write them into the database


	// 1 - Load bookmarks
	bookmarksFromFile, err := bookmarks.LoadBookmarks()
	if err != nil {
		log.Fatalf("Error loading bookmarks: %v", err)
	}

	// 2 - Get a list of the titles with "mangadex" connector
	names := bookmarks.MangadexMangaTitles(bookmarksFromFile)

	// open the database
	dbConnection, _ := sqlitedb.OpenDatabase("database/mangaList.db")
	// iterate of the names of the mangas in the bookmark list
	for _, name := range names {

		// a. extract the mangadex id from the database based on the manga name
		mangadexId, _ := sqlitedb.MangadexIdDbLookup(dbConnection, name, "chapters")

		// b. get the list of chapters from mangadex
		chapterList, _ := mangadex.ChaptersSorted(mangadexId)

		// c. update the DB with the new chapter list
		sqlitedb.MangadexInitialDbChapterListUpdate(dbConnection, name, chapterList)

		fmt.Println("Updated DB for: ", name)
	}
}
*/

func CheckIfBookmarkInDb() {

	/*
		compares the managa names in bookmarks to the names in the database, prints out the difference if the name does
		not exist in the DB.
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

func MangaAttributes() {

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
		postgresqldb.InsertMangaStatus(pgDb, "mangadex", status, id)
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

/*
func NewMangaDbUpdate() {

	// func to make rest call for managa information and the update the database with the filters/processed information
	// for a specific list of mangas.


	// 1 - Load bookmarks
	bookmarksFromFile, err := bookmarks.LoadBookmarks()
	if err != nil {
		log.Fatalf("Error loading bookmarks: %v", err)
	}
	// 2 - Get a list of the titles with "mangadex" connector from bookmarks
	names := bookmarks.MangadexMangaTitles(bookmarksFromFile)

	// 3 - open the database
	dbConnection, _ := sqlitedb.OpenDatabase("database/mangaList.db")

	// 3a - declare list to hold the return dicts
	var mangaNotInDb []string

	// iterate of the names of the mangas in the bookmark list
	for _, name := range names {

		// 3b - extract the mangadex id from the database based on the manga name
		mangaNameDb, _ := sqlitedb.MangaNameDbLookup(dbConnection, name, "chapters")

		if !mangaNameDb {
			mangaData, _ := mangadex.TitleSearch(name)
			mangaNotInDb = append(mangaNotInDb, mangaData)
		}
	}
	// 4 - Iterate through the mangaNotInDb list
	for _, mangaDict := range mangaNotInDb {
		// create datamap to hold the json data
		var dataMap map[string]interface{}

		// 4a - Convert JSON string to a map
		err := json.Unmarshal([]byte(mangaDict), &dataMap)
		if err != nil {
			fmt.Printf("Error unmarshalling manga data: %v\n", err)
			continue
		}
		// 4b - if there is no error then update the DB with the new manga data
		sqlitedb.AddMangaEntry(dbConnection, dataMap["name"].(string), dataMap["altTitle"].(string), dataMap["url"].(string), dataMap["id"].(string))
		fmt.Println("Updated DB for: ", dataMap["name"].(string))
	}
}
*/

/*
func ExtractMangasWithoutChapterList() {

	// 1 - open the database
	dbConnection, _ := sqlitedb.OpenDatabase("database/mangaList.db")

	mangasWithoutChapterLists := sqlitedb.QueryAllMangadexNames(dbConnection)

	for _, manga := range mangasWithoutChapterLists {
		fmt.Println(manga)
	}
}
*/

/*
func UpdateMangasWithoutChapterList() {

	//	func to grab a list of chapters from mangadex and then add to the database.  The list of mangas is manually
	//	provided.


	// 1 - open the database
	dbConnection, _ := sqlitedb.OpenDatabase("database/mangaList.db")

	// 2 list of mangas to get the chapter list for
	mangasToUpdate := sqlitedb.QueryAllMangadexNames(dbConnection)

	// 3 - iterate over the list of mangas
	for _, manga := range mangasToUpdate {

		// a. extract the mangadex id from the database based on the manga name
		mangadexId, _ := sqlitedb.MangadexIdDbLookup(dbConnection, manga, "chapters")

		// b. get the list of chapters from mangadex
		chapterList, _ := mangadex.ChaptersSorted(mangadexId)

		// c. update the DB with the new chapter list
		sqlitedb.MangadexInitialDbChapterListUpdate(dbConnection, manga, chapterList)

		fmt.Println("Updated DB for: ", manga)
	}
}
*/

func DumpPostgressDb() {
	/*
		Dumps the postgresql db.
	*/

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

	// get all data in postgresql manga table
	data, err := postgresqldb.LookupAllRows(pgDb, "manga")
	if err != nil {
		log.Fatalf("Error querying data: %v", err)
	}

	// Iterate over each map in the data slice
	for _, row := range data {
		// Create a slice to hold the keys
		var keys []string
		for key := range row {
			keys = append(keys, key)
		}

		// Sort the keys slice
		sort.Strings(keys)

		// Iterate over the sorted keys and print the corresponding key-value pairs
		for _, key := range keys {
			value := row[key]
			switch key {
			case "id", "name", "mangadex_id", "url":
				fmt.Printf("%s:\t%v\n", key, value)
			}
			// Add more cases as needed for other keys
		}
	}

}

func PgQueryByID(id string) {
	/*
		Dumps the postgresql db.
	*/

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

	// get all data in postgresql manga table
	data, err := postgresqldb.LookupByID(pgDb, "manga", id)
	if err != nil {
		log.Fatalf("Error querying data: %v", err)
	}

	for key, value := range data {
		fmt.Printf("%s: %v\n", key, value)
	}
}

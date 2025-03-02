package main

import (
	"encoding/json"
	"fmt"
	//"strings"
	"log"
	"main/bookmarks"
	"main/compare"
	"main/mangadex"
	"main/sqlitedb"
	"main/webfrontend"
	"os"
)

func init() {
	logFile, err := os.OpenFile("manga.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)
}

func main() {

	//NewMangaDbUpdate()
	//CheckIfBookmarkInDb()
	//BlanketUpdateDb()
	//ExtractMangasWithoutChapterList()
	//UpdateMangasWithoutChapterList()
	webfrontend.StartServer("8080")
}

func CheckForNewChapters() {

	/*
		func gets a list of chapters from mangadex and comnpares it to the list in the database, outputting the differences
		(if any).
	*/

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
		chapterList, _ := mangadex.MangadexChaptersSorted(mangadexId)

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

func BlanketUpdateDb() {
	/*
		func gets a sorted list of all the chapters from mangadex and write them into the database
	*/

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
		chapterList, _ := mangadex.MangadexChaptersSorted(mangadexId)

		// c. update the DB with the new chapter list
		sqlitedb.MangadexInitialDbChapterListUpdate(dbConnection, name, chapterList)

		fmt.Println("Updated DB for: ", name)
	}
}

func CheckIfBookmarkInDb() {

	/*
		compares the managa names in bookmarks to the names in the database, prints out the difference if the name does
		not exist in the DB.
	*/

	// 1 - Load bookmarks
	bookmarksFromFile, err := bookmarks.LoadBookmarks()
	if err != nil {
		log.Fatalf("Error loading bookmarks: %v", err)
	}
	// 2 - Get a list of the titles with "mangadex" connector from bookmarks
	names := bookmarks.MangadexMangaTitles(bookmarksFromFile)

	// open the database
	dbConnection, _ := sqlitedb.OpenDatabase("database/mangaList.db")
	// iterate of the names of the mangas in the bookmark list
	for _, name := range names {

		// a. extract the mangadex id from the database based on the manga name
		mangaNameDb, _ := sqlitedb.MangaNameDbLookup(dbConnection, name, "chapters")

		if !mangaNameDb {
			fmt.Printf("Bookmark not in DB: %s\n", name)
		}
	}
}

func NewMangaDbUpdate() {
	/*
	 func to make rest call for managa information and the update the database with the filters/processed information
	 for a specific list of mangas.
	*/

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
			mangaData, _ := mangadex.MangadexTitleSearch(name)
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

func ExtractMangasWithoutChapterList() {

	// 1 - open the database
	dbConnection, _ := sqlitedb.OpenDatabase("database/mangaList.db")

	mangasWithoutChapterLists := sqlitedb.QueryAllMangadexNames(dbConnection)

	for _, manga := range mangasWithoutChapterLists {
		fmt.Println(manga)
	}
}

func UpdateMangasWithoutChapterList() {
	/*
		func to grab a list of chapters from mangadex and then add to the database.  The list of mangas is manually
		provided.
	*/

	// 1 - open the database
	dbConnection, _ := sqlitedb.OpenDatabase("database/mangaList.db")

	// 2 list of mangas to get the chapter list for
	mangasToUpdate := sqlitedb.QueryAllMangadexNames(dbConnection)

	// 3 - iterate over the list of mangas
	for _, manga := range mangasToUpdate {

		// a. extract the mangadex id from the database based on the manga name
		mangadexId, _ := sqlitedb.MangadexIdDbLookup(dbConnection, manga, "chapters")

		// b. get the list of chapters from mangadex
		chapterList, _ := mangadex.MangadexChaptersSorted(mangadexId)

		// c. update the DB with the new chapter list
		sqlitedb.MangadexInitialDbChapterListUpdate(dbConnection, manga, chapterList)

		fmt.Println("Updated DB for: ", manga)
	}
}

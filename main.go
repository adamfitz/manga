package main

import (
	//"encoding/json"
	"fmt"
	//"strings"
	"log"

	//"main/httprequests"
	//"main/sqlitedb" // importing custom code from 'sqlitedb' package in subdir
	"main/bookmarks"
	"main/compare"
	"main/httprequests"
	"main/sqlitedb"
)

func main() {

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
	dbConnection, _ := sqlitedb.OpenDatabase("database/mangaList_test.db")

	// iterate of the names of the mangas in the bookmark list
	for _, name := range names {

		// a. extract the mangadex id from the database based on the manga name
		mangadexId, _ := sqlitedb.MangadexIdDbLookup(dbConnection, name, "chapters")

		// b. get the list of chapters from mangadex
		chapterList, _ := httprequests.MangadexChaptersSorted(mangadexId)

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
	dbConnection, _ := sqlitedb.OpenDatabase("database/mangaList_test.db")
	// iterate of the names of the mangas in the bookmark list
	for _, name := range names {

		// a. extract the mangadex id from the database based on the manga name
		mangadexId, _ := sqlitedb.MangadexIdDbLookup(dbConnection, name, "chapters")

		// b. get the list of chapters from mangadex
		chapterList, _ := httprequests.MangadexChaptersSorted(mangadexId)

		// c. update the DB with the new chapter list
		sqlitedb.MangadexInitialDbChapterListUpdate(dbConnection, name, chapterList)

		fmt.Println("Updated DB for: ", name)
	}
}

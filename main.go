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

	/*
		These managas are a PITA due to the fact that there is a possibliity that the latest chapter is available but
		earlier chapters are not.

		So instead we will operate under this assumption and instead store a list of the chapters in a specific format
		in the DB and then compare that with the list of chapters from the mangadex API.

	*/


	// load bookmarks and iterate through each name if the connect ir mangadex

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
	for _, name := range names {

		// a. extract the mangadex id from the database based on the manga name
		mangadexId, _ := sqlitedb.MangadexIdDbLookup(dbConnection, name, "chapters")

		// b. get the list of chapters from mangadex
		chapterList, _ := httprequests.MangadexChaptersSorted(mangadexId)

		// c. extract the list of chapters from the database
		chapterListDb, _ := sqlitedb.MangadexDbLookupChapterList(dbConnection, name)

		// d. perform a comparison of both of the json arrays, returning teh difference (compares A to B) and returns
		// the elements from list A that are not present in list B
		differences, _ := compare.CompareJSONArrays(chapterList, chapterListDb)

		// Output the differences
		if len(differences) != 0 {
			fmt.Printf("Chapters available on mangadex for %s\n", name)
			for _, diff := range differences {
				fmt.Println(diff)
			}
		}
	}
}

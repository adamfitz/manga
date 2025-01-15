package main

import (
	//"encoding/json"
	"fmt"
	//"strings"
	"log"

	//"main/httprequests"
	//"main/sqlitedb" // importing custom code from 'sqlitedb' package in subdir
	"main/bookmarks"
	"main/httprequests"
	"main/sqlitedb"
)

func main() {

	/*
				These managas are a PITA due to the fact that there is a possibliity that the latest chapter is available but
				earlier chapters are not.


			manga_id := "05a56be4-26ab-4f50-8fc0-ab8304570258"

			response, err := httprequests.MangadexGetChapterList(manga_id)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			//fmt.Printf("Response: %+v\n", response)

			// If you want to print specific nested data:
			for volume, volumeData := range response.Volumes {
				fmt.Printf("Volume: %s, Chapter Count: %d\n", volume, volumeData.Count)
				for chapter, chapterData := range volumeData.Chapters {
					fmt.Printf("  Chapter %s: ID=%s, Count=%d\n", chapter, chapterData.ID, chapterData.Count)
				}
			}


		chapterID := "06fb0bd0-855a-44d6-ad5a-319c84513ce7"

		pagesResponse, err := httprequests.MangadexGetPagesList(chapterID)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Print the base URL and hash
		fmt.Printf("Base URL: %s\n", pagesResponse.BaseURL)
		fmt.Printf("Hash: %s\n", pagesResponse.Chapter.Hash)

		// Print each page in the "data" array
		fmt.Println("Pages:")
		for i, page := range pagesResponse.Chapter.Data {
			fmt.Printf("  Page %d: %s\n", i+1, page)
		}

		// Print each page in the "dataSaver" array
		fmt.Println("Data Saver Pages:")
		for i, page := range pagesResponse.Chapter.DataSaver {
			fmt.Printf("  Page %d: %s\n", i+1, page)
		}

		//printing out a list of all the pages
		for _, page := range pagesResponse.Chapter.Data {
			fmt.Printf("Page: %v\n", page)
		}

		/*
			for _, chapter := range response.Data {
				if chapter.Attributes.TranslatedLanguage == "en" {
					fmt.Printf("ID: %s\n", chapter.Id)
					fmt.Printf("Chapter: %s\n", chapter.Attributes.Chapter)
					fmt.Printf("Updated At: %s\n\n", chapter.Attributes.UpdatedAt)
				}
			}
	*/

	/*
		// Convert the response to JSON
		responseJSON, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling response:", err)
			return
		}

		// Print the JSON string
		fmt.Println(string(responseJSON))

			// load bookmarks and return struct to iterate
			mangaList, err := utils.LoadBookmarks()
			if err != nil {
				fmt.Println(err)
				return
			}

			// Iterate through the returned slice (each managa in the list)
			for _, entry := range mangaList {

				// fmt.Printf("Manga Id: %s\n", entry.Key.Manga)
				// fmt.Println()

				// for each entry in mangaList struct
				// if the id contains https then the manga is not at mangadex (and has its own URL / URI)
				if !strings.Contains(entry.Key.Manga, "https") {
					// get the mangadex id
					mangaId := entry.Key.Manga

					// make api call for the mangaId
					response, err := httprequests.GetResponseAsStruct(mangaId)
					if err != nil {
						fmt.Println("Error:", err)
						return
					}

					// print out the manga title
					fmt.Printf("Manga Title: %s\n", entry.Title.Manga)
					// print out the mangaId, chapter number and updated date
					for _, chapter := range response.Data {
						if chapter.Attributes.TranslatedLanguage == "en" {
							fmt.Printf("ID: %s\n", chapter.Id)
							fmt.Printf("Chapter: %s\n", chapter.Attributes.Chapter)
							fmt.Printf("Updated At: %s\n\n", chapter.Attributes.UpdatedAt)

						}
					}
				}
			}
	*/

	/*


				// open database
				mangaListDb, _ := sqlitedb.OpenDatabase("./database/mangaList.db")
				// query for a row
				databaseRow, _ := sqlitedb.QueryWithCondition(mangaListDb, "chapters", "name", "Absolute Dominion")
				// Print the extracted row (single map)
				fmt.Println("Extracted Row:", databaseRow)
				// print out selected fields
				fmt.Printf("Name: %s\n", databaseRow["name"])
				fmt.Printf("Latest Chapter: %d\n", databaseRow["current_dld_chapter"])


			// chapter list

			manga_id := "05a56be4-26ab-4f50-8fc0-ab8304570258"
			//jsonSAtringArray, err := httprequests.MangadexChaptersSorted(manga_id)
			//if err != nil {
			//	fmt.Println("Error:", err)
			//	return
			//}

			//fmt.Printf("%s", jsonSAtringArray)

			//Update the sqlite database with a list of the mangadex chapters
			// load bookmarks and return SORTED struct to iterate
			bookmarks, err := bookmarks.LoadBookmarks()
			if err != nil {
				fmt.Println(err)
				return
			}


			// test on the first 3 bookmarks before continuing
			for idx, name := range bookmarks {
				// print out the struct title as a string
				if idx >= 3 {
					break
				}
				fmt.Println(name.Title.Manga)

				jsonChapters, err := httprequests.MangadexChaptersSorted(manga_id)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				fmt.Println(jsonChapters)
			}


		// Testing the flow to update initially the db column with a json list of all the chapters from mangadex

		// open the database
		mangaListDb, _ := sqlitedb.OpenDatabase("database/mangaList_test.db")

		// A Sharp-Eyed Classmate mangaid - get a list of all the english chapters
		jsonChapters, err := httprequests.MangadexChaptersSorted("e5f13b1a-eabd-4752-863a-cc3930a20d24")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Json chapters requested from mangadex:")
		fmt.Println(jsonChapters)

		// update the DB with the json list of chapters
		sqlitedb.MangaDexInitialDbChapterListUpdate(mangaListDb, "A Sharp-Eyed Classmate", jsonChapters)

		// After populating the table column - lookup the current list in the database
		dbChapterList, err := sqlitedb.MangaDexLookupChapterList(mangaListDb, "A Sharp-Eyed Classmate")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Json chapters retrieved from the database:")
		fmt.Println(dbChapterList)
	*/

	// load bookmarks and iterate through each name if the connect ir mangadex

	// 1 - Load bookmarks
	bookmarksFromFile, err := bookmarks.LoadBookmarks()
	if err != nil {
		log.Fatalf("Error loading bookmarks: %v", err)
	}

	// 2 - Get a list of the titles with "mangadex" connector
	names := bookmarks.MangadexMangaTitles(bookmarksFromFile)

	// 3 - Updating the database
	fmt.Println("== Updating database with Mangadex chapter list ==:")

	// open the database
	dbConnection, _ := sqlitedb.OpenDatabase("database/mangaList_test.db")
	for _, name := range names {
		fmt.Printf("\nUpdating Chapter list for %s ...\n", name)

		// a. extract the mangadex id from the database based on the manga name
		mangadexId, _ := sqlitedb.LookupMangadexId(dbConnection, name, "chapters")

		// b. get the list of chapters from mangadex
		chapterList, _ := httprequests.MangadexChaptersSorted(mangadexId)

		// c. update the database with the list of chapters
		sqlitedb.MangaDexInitialDbChapterListUpdate(dbConnection, name, chapterList)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	//"strings"

	//"mangadex/sqlitedb" // importing custom code from 'sqlitedb' package in subdir
	"main/httprequests"
	//"main/utils"
)

func main() {

	/*
		These managas are a PITA due to the fact that there is a possibliity that the latest chapter is available but
		earlier chapters are not.
	*/

	manga_id := "05a56be4-26ab-4f50-8fc0-ab8304570258"

	response, err := httprequests.MangadexGetChapterList(manga_id)
	if err != nil {
		fmt.Println("Error:", err)
		return
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
	// Convert the response to JSON
	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling response:", err)
		return
	}

	// Print the JSON string
	fmt.Println(string(responseJSON))
	/*
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
	*/

}

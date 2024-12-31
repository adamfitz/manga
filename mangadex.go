package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	//"strings"

	"mangadex/sqlitedb" // importing custom code from package in subdir
)

// MangaResponse struct represents the API response from Mangadex
type MangaResponse struct {
	Result   string        `json:"message"`
	Response string        `json:"response"`
	Data     []ChapterData `json:"data"`
}

// Chapter struct is the nested data part of the response struct
type ChapterData struct {
	Type       string            `json:"type"`
	Chapter    string            `json:"chapter"`
	Title      string            `json:"title"`
	Id         string            `json:"id"`
	Attributes ChapterAttributes `json:"attributes"`
	// Add other fields as needed based on the API response
}

// Attributes represents the nested attributes of each chapter
type ChapterAttributes struct {
	Volume             string `json:"volume"`
	Chapter            string `json:"chapter"` // chapter number
	Title              string `json:"title"`
	TranslatedLanguage string `json:"translatedLanguage"`
	ExternalUrl        string `json:"externalUrl"`
	PublishAt          string `json:"publishAt"`
	ReadableAt         string `json:"readableAt"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
	Pages              int    `json:"pages"`
	Version            int    `json:"version"`
}

// struct for bookmarks json file
type MangaList struct {
	Title Title `json:"title"`
	Key   Key   `json:"key"`
}

// struct for Title key in bookmarks file
type Title struct {
	Connector string `json:"connector"`
	Manga     string `json:"manga"`
}

// struct for key in bookmarks file
type Key struct {
	Connector string `json:"connector"`
	Manga     string `json:"manga"`
}

func getResponseAsString(manga_id string) (map[string]interface{}, error) {
	response, err := http.Get("https://api.mangadex.org/chapter?manga=" + manga_id)
	if err != nil {
		return nil, fmt.Errorf("error making http request: %s", err)
	}
	// schedules the resource cleanup for when the block of code finishes
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %s", err)
	}

	// Decode JSON into a map
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %s", err)
	}

	return result, nil
}

func getResponseAsStruct(manga_id string) (MangaResponse, error) {

	// parsed response
	var structuredResponse MangaResponse

	response, err := http.Get("https://api.mangadex.org/chapter?manga=" + manga_id)
	if err != nil {
		return MangaResponse{}, fmt.Errorf("error making http request")
	}
	// cleanup when the block of code finishes
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return MangaResponse{}, fmt.Errorf("error reading response body: %v", err)
	}

	// convert the byte array eponse body into the target struct datatype
	parsingError := json.Unmarshal(body, &structuredResponse)
	if parsingError != nil {
		panic(parsingError)
	}
	return structuredResponse, nil
}

func LoadBookmarks() ([]MangaList, error) {
	// Read the JSON file
	bookmarks, err := os.ReadFile("bookmarks.json")
	if err != nil {
		return nil, fmt.Errorf("Error reading file: %w", err)
	}

	// Unmarshal the JSON data into a slice of MangaList structs
	var mangaList []MangaList
	err = json.Unmarshal(bookmarks, &mangaList)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return mangaList, nil

}

func main() {

	/*
			manga_id := "05a56be4-26ab-4f50-8fc0-ab8304570258"

			response, err := getResponseAsStruct(manga_id)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			for _, chapter := range response.Data {
				if chapter.Attributes.TranslatedLanguage == "en" {
					fmt.Printf("ID: %s\n", chapter.Id)
					fmt.Printf("Chapter: %s\n", chapter.Attributes.Chapter)
					fmt.Printf("Updated At: %s\n\n", chapter.Attributes.UpdatedAt)
				}
			}


		// load bookmarks and return struct to iterate
		mangaList, err := LoadBookmarks()
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
				response, err := getResponseAsStruct(mangaId)
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

	// open datbase
	chaptersDb, _ := sqlitedb.OpenDatabase("./database/mangaList.db")
	// query for a row
	mangaRow, _ := sqlitedb.QueryWithCondition(chaptersDb, "chapters", "name", "Absolute Dominion")
	// Print the extracted row (single map)
	fmt.Println("Extracted Row:", mangaRow)

}

package bookmarks

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

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

func LoadBookmarks() ([]MangaList, error) {
	/*
		Loads the bookmarks from the boomarks.json file and sorts them by Title.Manga alphabetically

	*/
	// Read the JSON file
	bookmarks, err := os.ReadFile("bookmarks/bookmarks.json")
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Unmarshal the JSON data into a slice of MangaList structs
	var mangaList []MangaList
	err = json.Unmarshal(bookmarks, &mangaList)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	// Sort the mangaList by the Manga title alphabetically
	sort.Slice(mangaList, func(i, j int) bool {
		return mangaList[i].Title.Manga < mangaList[j].Title.Manga
	})

	return mangaList, nil
}

func MangadexMangaTitles(bookmarks []MangaList) []string {
	/*
		filter the loaded bookmarks variable on the manga title and return all titles as string slice.

		Can be used as a list to iterate over the DB.
	*/

	var mangaDexTitles []string

	// Iterate over bookmarks and filter titles with the "mangadex" connector
	for _, bookmark := range bookmarks {

		// Check if the connector is "mangadex" (case-insensitive)
		if strings.ToLower(bookmark.Title.Connector) == "mangadex" {
			mangaDexTitles = append(mangaDexTitles, bookmark.Title.Manga)
		}
	}

	// Check if mangaDexTitles is empty after filtering
	if len(mangaDexTitles) == 0 {
		fmt.Println("No MangaDex titles found.")
	}

	return mangaDexTitles
}

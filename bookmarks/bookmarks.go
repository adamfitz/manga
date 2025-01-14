package bookmarks

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
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

// LoadBookmarks loads the bookmarks from the file and sorts them by Title.Manga alphabetically
func LoadBookmarks() ([]MangaList, error) {
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
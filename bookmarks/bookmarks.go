package bookmarks

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	//"bytes"
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

// Normalize ensures Key.Manga is always a string
func (k *Key) Normalize() {
	switch v := any(k.Manga).(type) {
	case json.Number:
		k.Manga = v.String() // Convert to string safely
	case float64:
		k.Manga = fmt.Sprintf("%.0f", v) // Convert to string
	case string:
		// Already a string, do nothing
	default:
		k.Manga = "" // Handle unexpected types gracefully
	}
}

func LoadBookmarks() ([]MangaList, error) {
	/*
		Loads the bookmarks from the bookmarks.json file and sorts them by Title.Manga alphabetically
	*/

	// Read the JSON file
	bookmarks, err := os.ReadFile("bookmarks/bookmarks.json")
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Unmarshal into generic slice of maps to detect incorrect types
	var rawData []map[string]interface{}
	err = json.Unmarshal(bookmarks, &rawData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling into map: %w", err)
	}

	// Fix data type issues (convert numbers to strings)
	for _, item := range rawData {
		if keyData, ok := item["key"].(map[string]interface{}); ok {
			if mangaValue, exists := keyData["manga"]; exists {
				switch v := mangaValue.(type) {
				case float64:
					keyData["manga"] = fmt.Sprintf("%.0f", v) // Convert number to string
				case json.Number:
					keyData["manga"] = v.String()
				}
			}
		}
	}

	// Re-marshal fixed data
	fixedJSON, err := json.Marshal(rawData)
	if err != nil {
		return nil, fmt.Errorf("error marshalling fixed JSON: %w", err)
	}

	// Unmarshal into final struct
	var mangaList []MangaList
	err = json.Unmarshal(fixedJSON, &mangaList)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling fixed JSON: %w", err)
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

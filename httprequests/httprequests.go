package httprequests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
)

var mangadexBaseUri string = "https://api.mangadex.org"

// -- httprequests structs --

// Nested struct - MangaResponse struct represents the API response from Mangadex
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

// Nested struct - MangadexChapterList represents the root structure of the API response for chapter information
type MangadexChapterList struct {
	Result  string                `json:"result"`
	Volumes map[string]VolumeData `json:"volumes"`
}

// VolumeData represents the data for a specific volume
type VolumeData struct {
	Volume   string                 `json:"volume"`
	Count    int                    `json:"count"`
	Chapters map[string]ChapterInfo `json:"chapters"`
}

// ChapterInfo represents the data for a specific chapter
type ChapterInfo struct {
	Chapter string   `json:"chapter"`
	Count   int      `json:"count"`
	ID      string   `json:"id"`
	Others  []string `json:"others"`
}

// ChapterPageData represents the top-level JSON structure
type ChapterPageData struct {
	Result  string         `json:"result"`
	BaseURL string         `json:"baseUrl"`
	Chapter ChapterDetails `json:"chapter"`
}

// ChapterDetails represents the "chapter" field in the JSON
type ChapterDetails struct {
	Hash      string   `json:"hash"`
	Data      []string `json:"data"`
	DataSaver []string `json:"dataSaver"`
}

// -- httprequests functions --

func MangadexHttpRespAsString(manga_id string) (map[string]interface{}, error) {
	/*
		func returns the chapter information for a specific manga by the manga id as as a string
	*/
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

func MangadexHttpRespAsStruct(manga_id string) (MangaResponse, error) {
	/*
		func returns the chapter information for a specific manga by the manga id as as struct (custom type)
	*/

	// parsed response
	var structuredResponse MangaResponse

	response, err := http.Get(mangadexBaseUri + "/chapter?manga=" + manga_id)
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

func MangadexChapters(mangaID string) (*MangadexChapterList, error) {
	/*
		func returns a list of all chapters for a specific manga

		NOTE: This func is different from MangadexChaptersSorted() becuase this func uses the /aggregate URI which
		provides only the volume, chapter and chapter info (not detailed info).
	*/

	// Make the HTTP GET request, NOTE the translated language is hard coded to english
	response, err := http.Get(fmt.Sprintf("%s/manga/%s/aggregate?translatedLanguage[]=en", mangadexBaseUri, mangaID))
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request for list of volumes and chapters: %s", err)
	}
	// schedules the resource cleanup for when the block of code finishes
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body for list of volumes and chapters: %s", err)
	}

	// Decode the JSON response into the struct
	var chapterList MangadexChapterList
	err = json.Unmarshal(body, &chapterList)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON for list of volumes and chapters: %s", err)
	}

	return &chapterList, nil
}

func MangadexChaptersSorted(mangaId string) (string, error) {
	/*
		This function grabs a list of all the chapters for a specific manga from mangadex.com and returns a JSON string,
		sorted and ordered by chapter number.

		NOTE: This func is different from MangadexChapters() becuase this func uses the /feed URI which provides
		detailed information about each chapter.

		Also this func returns chapter list sorted by chapter number.  THe URis probably also should be swapped and
		this func use /aggregate and the otrher /feed.
	*/

	// Make the HTTP GET request
	response, err := http.Get(fmt.Sprintf("%s/manga/%s/feed?translatedLanguage[]=en", mangadexBaseUri, mangaId))
	if err != nil {
		return "", fmt.Errorf("error making HTTP request: %w", err)
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Parse JSON into a generic map
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	// Extract chapters
	var chapters []map[string]interface{}
	if dataArray, ok := data["data"].([]interface{}); ok {
		for _, item := range dataArray {
			if chapterMap, ok := item.(map[string]interface{}); ok {
				if attributes, ok := chapterMap["attributes"].(map[string]interface{}); ok {
					chapters = append(chapters, attributes)
				}
			}
		}
	}

	// Sort chapters by the "chapter" field
	sort.Slice(chapters, func(i, j int) bool {
		// Convert chapter strings to numbers for comparison
		chapterI, _ := strconv.ParseFloat(fmt.Sprintf("%v", chapters[i]["chapter"]), 64)
		chapterJ, _ := strconv.ParseFloat(fmt.Sprintf("%v", chapters[j]["chapter"]), 64)
		return chapterI < chapterJ
	})

	// Build JSON string array
	var chapterLines []string
	for _, chapter := range chapters {
		// Ensure each field has a string value or defaults to an empty string
		volume := fmt.Sprintf("%v", chapter["volume"])
		if volume == "<nil>" || volume == "null" {
			volume = ""
		}

		chapterNumber := fmt.Sprintf("%v", chapter["chapter"])
		if chapterNumber == "<nil>" || chapterNumber == "null" {
			chapterNumber = ""
		}

		title := fmt.Sprintf("%v", chapter["title"])
		if title == "<nil>" || title == "null" {
			title = ""
		}

		// Create the formatted line
		line := fmt.Sprintf("Volume: %s Chapter: %s Title: %s", volume, chapterNumber, title)
		chapterLines = append(chapterLines, line)
	}

	// Convert the slice of lines into a JSON string array
	jsonArray, err := json.Marshal(chapterLines)
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON array: %w", err)
	}

	return string(jsonArray), nil
}

func MangadexChapterPages(chapterID string) (*ChapterPageData, error) {
	/*
		func returns a list of all pages and page information for a specific chapter
	*/
	// Make the HTTP GET request
	response, err := http.Get(fmt.Sprintf("%s/at-home/server/%s", mangadexBaseUri, chapterID))
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request for chapter page information: %s", err)
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body for chapter page information: %s", err)
	}

	// Decode the JSON response into the struct
	var chapterPageData ChapterPageData
	err = json.Unmarshal(body, &chapterPageData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON for chapter page information: %s", err)
	}

	return &chapterPageData, nil
}

func MangadexTitleSearch(name string) (string, error) {
	/*
		Function to search for a manga by name (title) and extract the id and altTitles to use to populate the database.
	*/
	// Create the URL and add the query parameters
	baseURL := mangadexBaseUri + "/manga"
	params := url.Values{}
	params.Add("title", name)

	// Create the full URL with parameters
	mangaUrl := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Make the HTTP GET request
	response, err := http.Get(mangaUrl)
	if err != nil {
		return "", fmt.Errorf("error making HTTP request for manga information: %w", err)
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Parse JSON response
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	// Extract relevant information
	var result map[string]interface{}
	if mangaData, ok := data["data"].([]interface{}); ok && len(mangaData) > 0 {
		// Only process the first result
		firstManga := mangaData[0].(map[string]interface{})
		id := firstManga["id"].(string)
		attributes := firstManga["attributes"].(map[string]interface{})
		altTitles := attributes["altTitles"]

		// Build the result
		result = map[string]interface{}{
			"id":        id,
			"altTitles": altTitles,
		}
	} else {
		return "", fmt.Errorf("no manga found for the title: %s", name)
	}

	// Convert result to JSON string
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling result to JSON: %w", err)
	}

	return string(jsonResult), nil
}

/*

// Get request: baseurl/hash/pageNumber

// Download chapter page
func DownloadPage(baseUrl string, hash string, pageName string, targetDir string) (Nil error) {

}

// Create CBZ file (zip file)
func CreateCbzFile(sourceDir string, outputFileName string) {

}
*/

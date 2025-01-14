package httprequests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
)

// httprequests structs start here

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

// httprequests functions below here

func GetResponseAsString(manga_id string) (map[string]interface{}, error) {
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

func GetResponseAsStruct(manga_id string) (MangaResponse, error) {

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

// Return a list of all chapters for a specific manga
func MangadexGetChapterList(mangaID string) (*MangadexChapterList, error) {

	// Make the HTTP GET request, NOTE the translated language is hard coded to english
	response, err := http.Get(fmt.Sprintf("https://api.mangadex.org/manga/%s/aggregate?translatedLanguage[]=en", mangaID))
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

// Return a list of all pages and information within a chapter
func MangadexGetPagesList(chapterID string) (*ChapterPageData, error) {
	// Make the HTTP GET request
	response, err := http.Get(fmt.Sprintf("https://api.mangadex.org/at-home/server/%s", chapterID))
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

/*

// Get request: baseurl/hash/pageNumber

// Download chapter page
func DownloadPage(baseUrl string, hash string, pageName string, targetDir string) (Nil error) {

}

// Create CBZ file (zip file)
func CreateCbzFile(sourceDir string, outputFileName string) {

}
*/

func MangadexChaptersSorted(mangaId string) error {
	// Make the HTTP GET request
	response, err := http.Get(fmt.Sprintf("https://api.mangadex.org/manga/%s/feed?translatedLanguage[]=en", mangaId))
	if err != nil {
		return fmt.Errorf("error making HTTP request: %w", err)
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	// Parse JSON into a generic map
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
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
		chapterI, _ := strconv.ParseFloat(chapters[i]["chapter"].(string), 64)
		chapterJ, _ := strconv.ParseFloat(chapters[j]["chapter"].(string), 64)
		return chapterI < chapterJ
	})

	// Print sorted chapters
	for _, chapter := range chapters {
		fmt.Printf("Chapter: %v, Title: %v, Volume: %v\n",
			chapter["chapter"], chapter["title"], chapter["volume"])
	}

	return nil
}

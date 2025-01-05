package httprequests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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

// Nested struct - Chapter page data represents the data for all the pages in a chapter
type ChapterPageData struct {
	Result          string                     `json:"result"`
	Baseurl         string                     `json:"baseurl"`
	ChapterPageInfo map[string]ChapterPageInfo `json:"chapterpageinfo"`
}

type ChapterPageInfo struct {
	Hash  string  `json:"hash"`
	Pages []Pages `json:"pages"`
}

type Pages struct {
	Page string `json:"page"`
}

// Return a list of all pages and information within a chapter
func MangadexGetPagesList(mangaID string) (*ChapterPageData, error) {

	// Make the HTTP GET request for the specific chapter data
	response, err := http.Get(fmt.Sprintf("https://api.mangadex.org/at-home/server/%s", mangaID))
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request for chapter page information: %s", err)
	}
	// schedules the resource cleanup for when the block of code finishes
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body for chapter page information %s", err)
	}

	// Decode the JSON response into the struct
	var chapterPageData ChapterPageData
	err = json.Unmarshal(body, &chapterPageData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON for chapter page information: %s", err)
	}

	return &chapterPageData, nil
}

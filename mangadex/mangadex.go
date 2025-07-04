package mangadex

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/parser"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

var mangadexApiBaseUri string = "https://api.mangadex.org"
var mangadexBaseUri string = "https://mangadex.org"

// -- mangadex structs --

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

// -- mangadex functions --

/*
func returns the chapter information for a specific manga by the manga id as a JSON string
*/
func HttpResponseAsString(manga_id string) (string, error) {

	response, err := http.Get("https://api.mangadex.org/chapter?manga=" + manga_id)
	if err != nil {
		return "", fmt.Errorf("error making http request: %s", err)
	}
	// schedules the resource cleanup for when the block of code finishes
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %s", err)
	}

	// Decode JSON into a map
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling JSON: %s", err)
	}

	// Marshal the map back to a JSON string
	jsonString, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON to string: %s", err)
	}

	return string(jsonString), nil
}

/*
func returns the chapter information for a specific manga by the manga id as as struct (custom type)
*/
func HttpResponseAsStruct(manga_id string) (MangaResponse, error) {

	// parsed response
	var structuredResponse MangaResponse

	response, err := http.Get(mangadexApiBaseUri + "/chapter?manga=" + manga_id)
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

/*
func returns a list of all chapters for a specific manga

NOTE: This func is different from ChaptersSorted() becuase this func uses the /aggregate URI which
provides only the volume, chapter and chapter info (not detailed info).
*/
func Chapters(mangaID string) (*MangadexChapterList, error) {

	// Make the HTTP GET request, NOTE the translated language is hard coded to english
	response, err := http.Get(fmt.Sprintf("%s/manga/%s/aggregate?translatedLanguage[]=en", mangadexApiBaseUri, mangaID))
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

/*
This function grabs a list of all the chapters for a specific manga from mangadex.com and returns a JSON string,
sorted and ordered by chapter number.

NOTE: This func is different from Chapters() becuase this func uses the /feed URI which provides
detailed information about each chapter.

Also this func returns chapter list sorted by chapter number.  The URIs probably also should be swapped and
this func use /aggregate and the otrher /feed.
*/
func ChaptersSorted(mangaId string) (string, error) {
	const limit = 100 // Define the limit for pagination
	baseURL := fmt.Sprintf("%s/manga/%s/feed", mangadexApiBaseUri, mangaId)
	var chapters []map[string]interface{}
	var totalChapters int

	offset := 0
	for {
		// Construct URL with pagination parameters
		url := fmt.Sprintf("%s?limit=%d&offset=%d&translatedLanguage[]=en", baseURL, limit, offset)

		// Make the HTTP GET request
		response, err := http.Get(url)
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

		// Extract total chapters on the first request
		if totalChapters == 0 {
			if total, ok := data["total"].(float64); ok {
				totalChapters = int(total)
			}
		}

		// Extract chapter data
		if dataArray, ok := data["data"].([]interface{}); ok {
			for _, item := range dataArray {
				if chapterMap, ok := item.(map[string]interface{}); ok {
					if attributes, ok := chapterMap["attributes"].(map[string]interface{}); ok {
						chapters = append(chapters, attributes)
					}
				}
			}
		}

		// Break the loop if there are no more chapters to fetch
		if len(chapters) >= totalChapters {
			break
		}
		offset += limit
	}

	// Sort chapters by the "chapter" field
	sort.Slice(chapters, func(i, j int) bool {
		chapterI, _ := strconv.ParseFloat(fmt.Sprintf("%v", chapters[i]["chapter"]), 64)
		chapterJ, _ := strconv.ParseFloat(fmt.Sprintf("%v", chapters[j]["chapter"]), 64)
		return chapterI < chapterJ
	})

	// Build JSON string array
	var chapterLines []string
	for _, chapter := range chapters {
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

		line := fmt.Sprintf("Volume: %s Chapter: %s Title: %s", volume, chapterNumber, title)
		chapterLines = append(chapterLines, line)
	}

	// Print the total number of chapters
	//fmt.Printf("Total chapters: %d\n", totalChapters)

	// Convert the slice of lines into a JSON string array
	jsonArray, err := json.Marshal(chapterLines)
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON array: %w", err)
	}

	return string(jsonArray), nil
}

/*
func ChapterPages(chapterID string) (*ChapterPageData, error) {

	//	func returns a list of all pages and page information for a specific chapter
	// Make the HTTP GET request
	response, err := http.Get(fmt.Sprintf("%s/at-home/server/%s", mangadexApiBaseUri, chapterID))
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
*/

/*
Function to search for a manga by name (title) and extract the id and a prioritized altTitle to populate the database.
*/
func TitleSearch(name string) (string, error) {
	// Create the URL and add the query parameters
	baseURL := mangadexApiBaseUri + "/manga"
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
		altTitles := attributes["altTitles"].([]interface{})

		// Find the prioritized altTitle
		var prioritizedAltTitle string
		for _, alt := range altTitles {
			altMap := alt.(map[string]interface{})
			if enTitle, ok := altMap["en"]; ok {
				prioritizedAltTitle = enTitle.(string)
				break
			} else if jaTitle, ok := altMap["ja"]; ok {
				prioritizedAltTitle = jaTitle.(string)
			} else if zhTitle, ok := altMap["zh"]; ok {
				prioritizedAltTitle = zhTitle.(string)
			} else if prioritizedAltTitle == "" {
				// Assign any if no prioritized language found
				for _, v := range altMap {
					prioritizedAltTitle = v.(string)
					break
				}
			}
		}

		// Build the result
		result = map[string]interface{}{
			"id":       id,
			"altTitle": prioritizedAltTitle,
			"name":     name,                                            // add name to the result
			"url":      fmt.Sprintf("%s/manga/%s", mangadexBaseUri, id), // build the url from the id
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
func returns the chapter information for a specific manga by the manga id as a map
*/
func MangaAttributes(manga_id string) (map[string]any, error) {
	response, err := http.Get("https://api.mangadex.org/manga/" + manga_id)
	if err != nil {
		log.Printf("error making http request: %s", err)
		return nil, fmt.Errorf("error making http request: %s", err)
	}
	// schedules the resource cleanup for when the block of code finishes
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		return nil, fmt.Errorf("error reading response body: %s", err)
	}

	// Decode JSON into a map
	var result map[string]any
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("error unmarshalling JSON: %s", err)
		return nil, fmt.Errorf("error unmarshalling JSON: %s", err)
	}

	return result, nil
}

// return the manga status from the reponse
func MangaStatus(response map[string]any) string {
	// get the manga status
	status, err := parser.NestedMapValue(response, "data", "attributes", "status")
	if err != nil {
		log.Println("Error: MangaStatus - Problem extracting status", err)
		return "" //return empty string if there is an error
	}

	// validate string before returning
	if strStatus, ok := status.(string); ok {
		return strStatus
	}

	log.Println("Error: MangaStatus - status is not a string or is nil")
	return "" // Return an empty string if the type is incorrect
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

func ChaptersWithDetails(mangaId string) ([]map[string]any, error) {
	const limit = 100
	baseURL := fmt.Sprintf("%s/manga/%s/feed", mangadexApiBaseUri, mangaId)

	chapterMap := make(map[string]map[string]any) // key = chapter number (string)
	offset := 0

	for {
		// Construct URL with pagination parameters
		url := fmt.Sprintf("%s?limit=%d&offset=%d&translatedLanguage[]=en", baseURL, limit, offset)

		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch chapters: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		var result map[string]any
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}

		dataArray, ok := result["data"].([]any)
		if !ok {
			break
		}

		for _, item := range dataArray {
			itemMap, ok := item.(map[string]any)
			if !ok {
				continue
			}

			// Extract chapter ID and attributes
			id, _ := itemMap["id"].(string)
			attr, ok := itemMap["attributes"].(map[string]any)
			if !ok {
				continue
			}

			chapterStr := fmt.Sprintf("%v", attr["chapter"])
			versionFloat, _ := attr["version"].(float64)
			version := int(versionFloat)

			// Build final chapter entry
			chapterEntry := map[string]any{"id": id}
			for k, v := range attr {
				chapterEntry[k] = v
			}

			// Deduplication logic: keep highest version per chapter
			if existing, exists := chapterMap[chapterStr]; exists {
				existingVersion := int(existing["version"].(float64))
				if version > existingVersion {
					chapterMap[chapterStr] = chapterEntry
				}
			} else {
				chapterMap[chapterStr] = chapterEntry
			}
		}

		if len(dataArray) < limit {
			break
		}
		offset += limit
	}

	// Convert map to slice
	var chapters []map[string]any
	for _, chapter := range chapterMap {
		chapters = append(chapters, chapter)
	}

	// Optional: sort by chapter number
	sort.Slice(chapters, func(i, j int) bool {
		chI, _ := strconv.ParseFloat(fmt.Sprintf("%v", chapters[i]["chapter"]), 64)
		chJ, _ := strconv.ParseFloat(fmt.Sprintf("%v", chapters[j]["chapter"]), 64)
		return chI < chJ
	})

	return chapters, nil
}

/*
Returns a map of all page and server information for a specific chapter
*/
func ChapterPages(chapterID string) (map[string]any, error) {

	// Make the HTTP GET request
	response, err := http.Get(fmt.Sprintf("%s/at-home/server/%s", mangadexApiBaseUri, chapterID))
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request for chapter page information: %s", err)
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body for chapter page information: %s", err)
	}

	// Decode JSON into a map
	var data map[string]any
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON for chapter page information: %s", err)
	}

	return data, nil
}

// DownloadPage downloads a single image page and saves it to targetDir
func DownloadPage(baseUrl, hash, pageName, targetDir string) error {
	url := fmt.Sprintf("%s/data/%s/%s", baseUrl, hash, pageName)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download page %s: %w", pageName, err)
	}
	defer resp.Body.Close()

	outPath := filepath.Join(targetDir, pageName)
	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outPath, err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write image %s: %w", pageName, err)
	}

	return nil
}

// CreateCBZ zips the files in tempDir and stores the archive in ./<mangaName>/<chapter>.cbz
func CreateCBZ(tempDir, mangaName, chapter string) (string, error) {
	destDir := filepath.Join(".", mangaName)
	err := os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	cbzPath := filepath.Join(destDir, chapter+".cbz")
	zipFile, err := os.Create(cbzPath)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if !info.Mode().IsRegular() {
			return nil
		}

		relPath := filepath.Base(path)
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		writer, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		return err
	})

	if err != nil {
		return "", err
	}

	return cbzPath, nil
}

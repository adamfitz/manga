package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"io"
)

// MangaResponse struct represents the API response from Mangadex
type MangaResponse struct {
	Result string `json:"message"`
	Response  string `json:"response"`
	Data []ChapterData	`json:"data"`
}

// Chapter struct is the nested data part of the response struct
type ChapterData struct {
    Type    string `json:"type"`
    Chapter string `json:"chapter"`
	Title string `json:"title"`
	Id string `json:"id"`
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



func getResponseAsStruct(manga_id string) (MangaResponse, error)  {

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


func main () {

	manga_id := "a920060c-7e39-4ac1-980c-f0e605a40ae4"

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
}


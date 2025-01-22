package webfrontend

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"main/sqlitedb"
	"net/http"
	"strings"
)

// StartServer initializes and starts the web server on the given port.
func StartServer(port string) {
	http.HandleFunc("/", homePageHandler)
	http.HandleFunc("/query", queryHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/add", addMangaEntryHandler)

	log.Printf("Web server running at http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	tmplParsed, err := template.ParseFiles("./webfrontend/index.html")
	if err != nil {
		http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error parsing template: %v", err)
		return
	}

	if err := tmplParsed.Execute(w, nil); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
	}
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract and clean input variables
	mangaName := strings.TrimSpace(r.FormValue("manga_name"))
	alternateName := strings.TrimSpace(r.FormValue("alternate_name"))
	dbId := strings.TrimSpace(r.FormValue("id"))

	// If empty, set to "Null"
	if mangaName == "" {
		mangaName = "Null"
	}
	if alternateName == "" {
		alternateName = "Null"
	}
	if dbId == "" {
		dbId = "Null"
	}

	// Open database connection
	dbConnection, _ := sqlitedb.OpenDatabase("database/mangaList.db")

	// Prepare the response
	var result string
	var queryResult map[string]interface{}

	// Query by mangaName
	if mangaName != "Null" {
		queryResult, _ = sqlitedb.QueryWithCondition(dbConnection, "chapters", "name", mangaName)
		result = fmt.Sprintf("Query Result for Manga Name: %s", mangaName)
	} else if alternateName != "Null" {
		queryResult, _ = sqlitedb.QueryWithCondition(dbConnection, "chapters", "alt_name", alternateName)
		result = fmt.Sprintf("Query Result for Alternate Name: %s", alternateName)
	} else if dbId != "Null" {
		queryResult, _ = sqlitedb.QueryWithCondition(dbConnection, "chapters", "id", dbId)
		result = fmt.Sprintf("Query Result for ID: %s", dbId)
	}

	// Check if mangadex_ch_list exists and is a string
	if mangadexChList, ok := queryResult["mangadex_ch_list"].(string); ok {
		// Parse mangadex_ch_list if it's a string representation of a JSON array
		var chapters []string
		err := json.Unmarshal([]byte(mangadexChList), &chapters)
		if err != nil {
			http.Error(w, "Error parsing chapter list", http.StatusInternalServerError)
			return
		}
		// Update the queryResult with the parsed chapter list
		queryResult["mangadex_ch_list"] = chapters
	}

	// Marshal queryResult to pretty-printed JSON
	queryResultJSON, err := json.MarshalIndent(queryResult, "", "    ")
	if err != nil {
		http.Error(w, "Error marshaling query result", http.StatusInternalServerError)
		return
	}

	// Prepare data for the template
	data := struct {
		Result      string
		QueryResult string
	}{
		Result:      result,
		QueryResult: string(queryResultJSON),
	}

	// Load the queryresult page template with the requested data
	tmpl, err := template.ParseFiles("./webfrontend/queryresult.html")
	if err != nil {
		// Print the error to the server logs for debugging
		log.Println("Error loading template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	// Send the response to the user
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

// Column substring search handler
func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract and clean input variables
	mangaName := strings.TrimSpace(r.FormValue("manga_name"))
	alternateName := strings.TrimSpace(r.FormValue("alternate_name"))

	// If empty, set to "Null"
	if mangaName == "" {
		mangaName = "Null"
	}
	if alternateName == "" {
		alternateName = "Null"
	}

	// Open database connection
	dbConnection, _ := sqlitedb.OpenDatabase("database/mangaList.db")

	// Prepare the response
	var result string
	var searchResult []map[string]interface{}

	// Declare error variable
	var err error

	// Query by mangaName or alternateName
	if mangaName != "Null" {
		searchResult, err = sqlitedb.QuerySearchSubstring(dbConnection, "chapters", "name", mangaName)
		result = fmt.Sprintf("Search Result for Manga Name: %s", mangaName)
	} else if alternateName != "Null" {
		searchResult, err = sqlitedb.QuerySearchSubstring(dbConnection, "chapters", "alt_name", alternateName)
		result = fmt.Sprintf("Search Result for Alternate Name: %s", alternateName)
	}
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}

	// Prepare data for the template
	data := struct {
		Result       string
		SearchResult []map[string]interface{}
	}{
		Result:       result,
		SearchResult: searchResult,
	}

	// Load the searchresult page template with the requested data
	tmpl, err := template.ParseFiles("./webfrontend/searchresult.html")
	if err != nil {
		// Print the error to the server logs for debugging
		log.Println("Error loading template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	// Send the response to the user
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

// Update handler
func updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	update := r.FormValue("update")
	// Call your database update logic here
	result := "Update executed: " + update

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

// Add Manga Entry Handler
func addMangaEntryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract and clean input variables
	mangaName := strings.TrimSpace(r.FormValue("manga_name"))
	alternateName := strings.TrimSpace(r.FormValue("alternate_name"))
	url := strings.TrimSpace(r.FormValue("url"))
	mangadexID := strings.TrimSpace(r.FormValue("mangadex_id"))

	// Validate input (ensure mangaName is provided)
	if mangaName == "" {
		http.Error(w, "Manga name is required", http.StatusBadRequest)
		return
	}

	// Open database connection
	dbConnection, err := sqlitedb.OpenDatabase("database/mangaList.db")
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		log.Println("Database connection error:", err)
		return
	}
	defer dbConnection.Close()

	// Add entry to the database and get the new ID
	newID, err := sqlitedb.AddMangaEntry(dbConnection, mangaName, alternateName, url, mangadexID)
	if err != nil {
		http.Error(w, "Error adding manga entry to the database", http.StatusInternalServerError)
		log.Println("Error adding entry:", err)
		return
	}
	fmt.Printf("New entry added with ID: %d\n", newID)

	// Query the database using the new ID
	//queryCondition := fmt.Sprintf("id = %d", newID)
	//fmt.Printf("Querying table 'chapters' with condition: %s\n", queryCondition)

	newEntry, err := sqlitedb.QueryByID(dbConnection, "chapters", newID)
	if err != nil {
		http.Error(w, "Error retrieving the added manga entry from the database", http.StatusInternalServerError)
		log.Println("Error querying for added entry:", err)
		return
	}

	// Load the addmangaentryresult.html template
	tmpl, err := template.ParseFiles("./webfrontend/adddbentryresult.html")
	if err != nil {
		log.Println("Error loading template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	// Prepare data for the template
	data := struct {
		Message string
		Entry   map[string]interface{}
	}{
		Message: fmt.Sprintf("Manga entry '%s' was added successfully!", mangaName),
		Entry:   newEntry,
	}

	// Send the response to the user
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

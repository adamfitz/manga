package webfrontend

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"main/auth"
	"main/postgresqldb"
	"net/http"
	"strconv"
	"strings"
)

// StartServer initializes and starts the web server on the given port.
func StartServer(port string) {
	// define page handlers
	http.HandleFunc("/", homePageHandler)
	http.HandleFunc("/manga", mangaPageHandler)
	http.HandleFunc("/anime", animePageHandler)
	http.HandleFunc("/lightnovel", lightNovelPageHandler)
	http.HandleFunc("/webnovel", webNovelPageHandler)
	http.HandleFunc("/webtoons", webtoonPageHandler)

	// define action handlers
	// manga actions
	http.HandleFunc("/queryManga", mangaQueryHandler) // this is the DB lookup, must be exact match
	http.HandleFunc("/searchManga", mangaSearchHandler)
	//http.HandleFunc("/updateManga", mangaUpdateHandler)
	http.HandleFunc("/addManga", addMangaEntryHandler)
	http.HandleFunc("/queryMangaAll", mangaLookupAllRows)
	http.HandleFunc("/queryMangadexAll", mangadexLookupAllRows)

	// anime actions
	http.HandleFunc("/queryAnime", animeQueryHandler)   // this is the DB lookup, must be exact match
	http.HandleFunc("/searchAnime", animeSearchHandler) // substring search case insensitive
	http.HandleFunc("/addAnime", addAnimeEntryHandler)

	// light novel actions
	http.HandleFunc("/queryLightNovel", lightNovelQueryHandler)   // this is the DB lookup, must be exact match
	http.HandleFunc("/searchLightNovel", lightNovelSearchHandler) // substring search case insensitive
	http.HandleFunc("/addLightNovel", addLightNovelEntryHandler)

	// webtoons actions
	http.HandleFunc("/queryWebtoon", webtoonQueryHandler)   // this is the DB lookup, must be exact match
	http.HandleFunc("/searchWebtoon", webtoonSearchHandler) // substring search case insensitive
	http.HandleFunc("/addWebtoon", addWebtoonEntryHandler)

	// webnovel actions
	http.HandleFunc("/queryWebNovel", webNovelQueryHandler)   // this is the DB lookup, must be exact match
	http.HandleFunc("/searchWebNovel", webNovelSearchHandler) // substring search case insensitive
	http.HandleFunc("/addWebNovel", addWebNovelEntryHandler)

	log.Printf("Web server running at http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

////////////////////////////////////////////////// PAGE HANDLERS  //////////////////////////////////////////////////

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

func mangaPageHandler(w http.ResponseWriter, r *http.Request) {
	tmplParsed, err := template.ParseFiles("./webfrontend/manga/manga.html")
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

func animePageHandler(w http.ResponseWriter, r *http.Request) {
	tmplParsed, err := template.ParseFiles("./webfrontend/anime/anime.html")
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

func lightNovelPageHandler(w http.ResponseWriter, r *http.Request) {
	tmplParsed, err := template.ParseFiles("./webfrontend/lightnovel/lightnovel.html")
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

func webNovelPageHandler(w http.ResponseWriter, r *http.Request) {
	tmplParsed, err := template.ParseFiles("./webfrontend/webnovel/webnovel.html")
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

func webtoonPageHandler(w http.ResponseWriter, r *http.Request) {
	tmplParsed, err := template.ParseFiles("./webfrontend/webtoons/webtoons.html")
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

//////////////////////////////////////////////////  ACTION HANDLERS  //////////////////////////////////////////////////

////////////// MANGA ACTION HANDLERS

func mangaQueryHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request on /query/Manga")
	// Load config
	config, _ := auth.LoadConfig()

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
	dbConnection, _ := postgresqldb.OpenDatabase(config.PgServer, config.PgPort, config.PgUser, config.PgPassword, config.PgDbName)

	// Prepare the response
	var result string
	var queryResult map[string]any

	// Query by mangaName
	if mangaName != "Null" {
		queryResult, _ = postgresqldb.QueryWithCondition(dbConnection, "mangadex", "name", mangaName)
		result = fmt.Sprintf("Query Result for Manga Name: %s", mangaName)
	} else if alternateName != "Null" {
		queryResult, _ = postgresqldb.QueryWithCondition(dbConnection, "mangadex", "alt_name", alternateName)
		result = fmt.Sprintf("Query Result for Alternate Name: %s", alternateName)
	} else if dbId != "Null" {
		queryResult, _ = postgresqldb.QueryWithCondition(dbConnection, "mangadex", "id", dbId)
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
	tmpl, err := template.ParseFiles("./webfrontend/manga/mangaQueryResult.html")
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
func mangaSearchHandler(w http.ResponseWriter, r *http.Request) {
	config, _ := auth.LoadConfig()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	mangaName := strings.TrimSpace(r.FormValue("manga_name"))
	alternateName := strings.TrimSpace(r.FormValue("alternate_name"))

	if mangaName == "" {
		mangaName = "Null"
	}
	if alternateName == "" {
		alternateName = "Null"
	}

	dbConnection, _ := postgresqldb.OpenDatabase(config.PgServer, config.PgPort, config.PgUser, config.PgPassword, config.PgDbName)

	var mangadexResults, mangaResults []map[string]any
	var err error

	if mangaName != "Null" {
		// query mangadex table
		mangadexResults, err = postgresqldb.QuerySearchMangadexSubstring(dbConnection, "mangadex", "name", mangaName)
		if err != nil {
			http.Error(w, "Error querying mangadex", http.StatusInternalServerError)
			return
		}
		// query manga table
		mangaResults, err = postgresqldb.QuerySearchMangaSubstring(dbConnection, "manga", "name", mangaName)
		if err != nil {
			http.Error(w, "Error querying manga", http.StatusInternalServerError)
			return
		}
	} else if alternateName != "Null" {
		// query mangadex table
		mangadexResults, err = postgresqldb.QuerySearchMangadexSubstring(dbConnection, "mangadex", "alt_name", alternateName)
		if err != nil {
			http.Error(w, "Error querying mangadex", http.StatusInternalServerError)
			return
		}
		// query manga table
		mangaResults, err = postgresqldb.QuerySearchMangaSubstring(dbConnection, "manga", "alt_name", alternateName)
		if err != nil {
			http.Error(w, "Error querying manga", http.StatusInternalServerError)
			return
		}
	}

	data := struct {
		Result          string
		MangadexResults []map[string]any
		MangaResults    []map[string]any
		HasResults      bool
	}{
		MangadexResults: mangadexResults,
		MangaResults:    mangaResults,
		HasResults:      len(mangadexResults) > 0 || len(mangaResults) > 0,
	}

	tmpl, err := template.ParseFiles("./webfrontend/manga/mangaSearchResult.html")
	if err != nil {
		log.Println("Error loading template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

// Return all rows in manag DB table (NOT mangadex)
func mangaLookupAllRows(w http.ResponseWriter, r *http.Request) {
	// im super lazy and since the target page will render both or EITHER manga results on their own I have just left
	// empty vars in each of the separate funcs that lookup and return all the table rows so teh target results template
	// does not have to change, becuase it is also used for single manga lookups as well (same as below func)
	var mangadexResults []map[string]any

	mangaResults, err := postgresqldb.AllMangaTableRows()
	if err != nil {
		log.Println("Error querying all rows manga table", http.StatusInternalServerError)
		http.Error(w, "Error querying all rows manga table", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("./webfrontend/manga/mangaSearchResult.html")
	if err != nil {
		log.Println("Error loading template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	data := struct {
		Result             string
		MangadexResults    []map[string]any
		MangaResults       []map[string]any
		HasResults         bool
		MangaTableRowCount int
		MangadexRowCount   int
	}{
		MangadexResults:    mangadexResults,
		MangaResults:       mangaResults,
		HasResults:         len(mangadexResults) > 0 || len(mangaResults) > 0,
		MangaTableRowCount: len(mangaResults),
		MangadexRowCount:   len(mangadexResults),
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)

}

// Return all rows in mangadex DB table (NOT manga table)
func mangadexLookupAllRows(w http.ResponseWriter, r *http.Request) {
	// im super lazy and since the target page will render both or EITHER manga results on their own I have just left
	// empty vars in each of the separate funcs that lookup and return all the table rows so teh target results template
	// does not have to change, becuase it is also used for single manga lookups as well
	var mangaResults []map[string]any

	mangadexResults, err := postgresqldb.AllMangaDexTableRows()
	if err != nil {
		log.Println("Error querying all rows mangadex table", http.StatusInternalServerError)
		http.Error(w, "Error querying all rows mangadex table", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("./webfrontend/manga/mangaSearchResult.html")
	if err != nil {
		log.Println("Error loading template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	data := struct {
		Result                string
		MangadexResults       []map[string]any
		MangaResults          []map[string]any
		HasResults            bool
		MangaTableRowCount    int
		MangadexTableRowCount int
	}{
		MangadexResults:       mangadexResults,
		MangaResults:          mangaResults,
		HasResults:            len(mangadexResults) > 0 || len(mangaResults) > 0,
		MangaTableRowCount:    len(mangaResults),
		MangadexTableRowCount: len(mangadexResults),
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)

}

/*
// Update handler
func mangaUpdateHandler(w http.ResponseWriter, r *http.Request) {
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
*/

// Add Manga Entry Handler
func addMangaEntryHandler(w http.ResponseWriter, r *http.Request) {
	// Load config
	config, _ := auth.LoadConfig()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract and clean input variables
	mangaName := strings.TrimSpace(r.FormValue("manga_name"))
	alternateName := strings.TrimSpace(r.FormValue("alternate_name"))
	url := strings.TrimSpace(r.FormValue("url"))
	mangadexID := strings.TrimSpace(r.FormValue("mangadex_id"))
	table := strings.TrimSpace(r.FormValue("table_select"))

	// Extract boolean fields (use pointers so NULL can be stored)
	var completed, ongoing, hiatus, cancelled *bool

	if r.FormValue("completed") == "on" {
		val := true
		completed = &val
	}
	if r.FormValue("ongoing") == "on" {
		val := true
		ongoing = &val
	}
	if r.FormValue("hiatus") == "on" {
		val := true
		hiatus = &val
	}
	if r.FormValue("cancelled") == "on" {
		val := true
		cancelled = &val
	}

	// Validate input
	if mangaName == "" {
		http.Error(w, "Manga name is required", http.StatusBadRequest)
		return
	}
	if table != "manga" && table != "mangadex" {
		http.Error(w, "Invalid table selected", http.StatusBadRequest)
		return
	}

	// Open database connection
	dbConnection, err := postgresqldb.OpenDatabase(config.PgServer, config.PgPort, config.PgUser, config.PgPassword, config.PgDbName)
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		log.Println("Database connection error:", err)
		return
	}
	defer dbConnection.Close()

	// Insert based on table selection
	var newID int64
	var newEntry map[string]any

	switch table {
	case "manga":
		newID, err = postgresqldb.AddMangaRow(dbConnection, mangaName, alternateName, url, mangadexID, completed, ongoing, hiatus, cancelled)
		if err != nil {
			http.Error(w, "Error adding manga entry to manga table", http.StatusInternalServerError)
			log.Println("Manga table error:", err)
			return
		}
		newEntry, err = postgresqldb.LookupByID(dbConnection, "manga", fmt.Sprintf("%d", newID))
		if err != nil {
			http.Error(w, "Error retrieving added entry from manga table", http.StatusInternalServerError)
			log.Println("Manga lookup error:", err)
			return
		}

	case "mangadex":
		newID, err = postgresqldb.AddMangadexRow(dbConnection, mangaName, alternateName, url, mangadexID, completed, ongoing, hiatus, cancelled)
		if err != nil {
			http.Error(w, "Error adding manga entry to mangadex table", http.StatusInternalServerError)
			log.Println("Mangadex table error:", err)
			return
		}
		newEntry, err = postgresqldb.LookupByID(dbConnection, "mangadex", fmt.Sprintf("%d", newID))
		if err != nil {
			http.Error(w, "Error retrieving added entry from mangadex table", http.StatusInternalServerError)
			log.Println("Mangadex lookup error:", err)
			return
		}
	}

	// Load the result template
	tmpl, err := template.ParseFiles("./webfrontend/manga/mangaAddDbEntryResult.html")
	if err != nil {
		log.Println("Error loading template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	// Prepare data for the template
	data := struct {
		Message string
		Entry   map[string]any
	}{
		Message: fmt.Sprintf("Manga entry '%s' was added to table '%s' successfully!", mangaName, table),
		Entry:   newEntry,
	}

	// Send the response
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

////////////// ANIME ACTION HANDLERS

// Anime Lookup Handler (query for exact match)
func animeQueryHandler(w http.ResponseWriter, r *http.Request) {
	// Load config
	config, _ := auth.LoadConfig()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract and clean input variables
	animeName := strings.TrimSpace(r.FormValue("anime_name"))
	alternateName := strings.TrimSpace(r.FormValue("alternate_name"))
	dbId := strings.TrimSpace(r.FormValue("id"))

	// If empty, set to "Null"
	if animeName == "" {
		animeName = "Null"
	}
	if alternateName == "" {
		alternateName = "Null"
	}
	if dbId == "" {
		dbId = "Null"
	}

	// Open database connection
	dbConnection, _ := postgresqldb.OpenDatabase(config.PgServer, config.PgPort, config.PgUser, config.PgPassword, config.PgDbName)

	// Prepare the response
	var result string
	var queryResult map[string]any

	// Query by animeName
	if animeName != "Null" {
		queryResult, _ = postgresqldb.QueryWithCondition(dbConnection, "anime", "name", animeName)
		result = fmt.Sprintf("Query Result for Anime Name: %s", animeName)
	} else if alternateName != "Null" {
		queryResult, _ = postgresqldb.QueryWithCondition(dbConnection, "anime", "alt_name", alternateName)
		result = fmt.Sprintf("Query Result for Anime Alternate name: %s", alternateName)
	} else if dbId != "Null" {
		queryResult, _ = postgresqldb.QueryWithCondition(dbConnection, "anime", "id", dbId)
		result = fmt.Sprintf("Query Result for Anime ID: %s", dbId)
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
	tmpl, err := template.ParseFiles("./webfrontend/anime/animeQueryResult.html")
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

// Anime search specified colmun for substring
func animeSearchHandler(w http.ResponseWriter, r *http.Request) {
	// Load config
	config, _ := auth.LoadConfig()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract and clean input variables
	animeName := strings.TrimSpace(r.FormValue("anime_name"))
	alternateName := strings.TrimSpace(r.FormValue("alternate_name"))

	// If empty, set to "Null"
	if animeName == "" {
		animeName = "Null"
	}
	if alternateName == "" {
		alternateName = "Null"
	}

	// Open database connection
	dbConnection, _ := postgresqldb.OpenDatabase(config.PgServer, config.PgPort, config.PgUser, config.PgPassword, config.PgDbName)

	// Prepare the response
	var result string
	var searchResult []map[string]any

	// Declare error variable
	var err error

	// Query by animeName or alternateName
	if animeName != "Null" {
		searchResult, err = postgresqldb.AnimeSearchSubstring(dbConnection, "anime", "name", animeName)
		result = fmt.Sprintf("Search Result for Anime Name: %s", animeName)
	} else if alternateName != "Null" {
		searchResult, err = postgresqldb.AnimeSearchSubstring(dbConnection, "anime", "alt_name", alternateName)
		result = fmt.Sprintf("Search Result for Anime Alternate Name: %s", alternateName)
	}
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}

	// Prepare data for the template
	data := struct {
		Result       string
		SearchResult []map[string]any
	}{
		Result:       result,
		SearchResult: searchResult,
	}

	// Load the searchresult page template with the requested data
	tmpl, err := template.ParseFiles("./webfrontend/anime/animeSearchResult.html")
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

// Add Anime Entry Handler
func addAnimeEntryHandler(w http.ResponseWriter, r *http.Request) {
	// Load config
	config, _ := auth.LoadConfig()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract and clean input variables
	animeName := strings.TrimSpace(r.FormValue("anime_name"))
	alternateName := strings.TrimSpace(r.FormValue("alternate_name"))
	url := strings.TrimSpace(r.FormValue("url"))

	// Extract boolean fields (use pointers so NULL can be stored)
	var completed, watched *bool

	if r.FormValue("completed") == "on" {
		val := true
		completed = &val
	}
	if r.FormValue("watched") == "on" {
		val := true
		watched = &val
	}

	// Validate input (ensure animeName is provided)
	if animeName == "" {
		http.Error(w, "Anime name is required", http.StatusBadRequest)
		return
	}

	// Open database connection
	dbConnection, err := postgresqldb.OpenDatabase(config.PgServer, config.PgPort, config.PgUser, config.PgPassword, config.PgDbName)
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		log.Println("Database connection error:", err)
		return
	}
	defer dbConnection.Close()

	// Add entry to the database and get the new ID
	newID, err := postgresqldb.AddAnimeRow(dbConnection, animeName, alternateName, url, completed, watched)
	if err != nil {
		http.Error(w, "Error adding anime entry to the database", http.StatusInternalServerError)
		log.Println("Error adding entry:", err)
		return
	}
	fmt.Printf("New entry added with ID: %d\n", newID)

	// Query the database using the new ID
	newEntry, err := postgresqldb.LookupByID(dbConnection, "anime", fmt.Sprintf("%d", newID))
	if err != nil {
		http.Error(w, "Error retrieving the added anime entry from the database", http.StatusInternalServerError)
		log.Println("Error querying for added entry:", err)
		return
	}

	// Load the add anime db entry result.html template
	tmpl, err := template.ParseFiles("./webfrontend/anime/animeAddDbEntryResult.html")
	if err != nil {
		log.Println("Error loading template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	// Prepare data for the template
	data := struct {
		Message string
		Entry   map[string]any
	}{
		Message: fmt.Sprintf("Anime entry '%s' was added successfully!", animeName),
		Entry:   newEntry,
	}

	// Send the response to the user
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

////////////// LIGHT NOVEL ACTION HANDLERS

// Light novel Lookup Handler (query for exact match)
func lightNovelQueryHandler(w http.ResponseWriter, r *http.Request) {
	// Load config
	config, _ := auth.LoadConfig()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract and clean input variables
	lightNovelName := strings.TrimSpace(r.FormValue("ln_name"))
	alternateName := strings.TrimSpace(r.FormValue("ln_alt_name"))
	dbId := strings.TrimSpace(r.FormValue("id"))

	// If empty, set to "Null"
	if lightNovelName == "" {
		lightNovelName = "Null"
	}
	if alternateName == "" {
		alternateName = "Null"
	}
	if dbId == "" {
		dbId = "Null"
	}

	// Open database connection
	dbConnection, _ := postgresqldb.OpenDatabase(config.PgServer, config.PgPort, config.PgUser, config.PgPassword, config.PgDbName)

	// Prepare the response
	var result string
	var queryResult map[string]any

	// Query by lightnovel
	if lightNovelName != "Null" {
		queryResult, _ = postgresqldb.QueryWithCondition(dbConnection, "lightnovel", "name", lightNovelName)
		result = fmt.Sprintf("Query Result for Light Novel Name: %s", lightNovelName)
	} else if alternateName != "Null" {
		queryResult, _ = postgresqldb.QueryWithCondition(dbConnection, "lightnovel", "alt_name", alternateName)
		result = fmt.Sprintf("Query Result for Light Novel Alternate name: %s", alternateName)
	} else if dbId != "Null" {
		queryResult, _ = postgresqldb.QueryWithCondition(dbConnection, "lightnovel", "id", dbId)
		result = fmt.Sprintf("Query Result for Light Novel ID: %s", dbId)
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
	tmpl, err := template.ParseFiles("./webfrontend/lightnovel/lightnovelQueryResult.html")
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

// Light Novel search specified colmun for substring
func lightNovelSearchHandler(w http.ResponseWriter, r *http.Request) {
	// Load config
	config, _ := auth.LoadConfig()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract and clean input variables
	lightNovelName := strings.TrimSpace(r.FormValue("ln_name"))
	alternateName := strings.TrimSpace(r.FormValue("ln_alt_name"))

	// If empty, set to "Null"
	if lightNovelName == "" {
		lightNovelName = "Null"
	}
	if alternateName == "" {
		alternateName = "Null"
	}

	// Open database connection
	dbConnection, _ := postgresqldb.OpenDatabase(config.PgServer, config.PgPort, config.PgUser, config.PgPassword, config.PgDbName)

	// Prepare the response
	var result string
	var searchResult []map[string]any

	// Declare error variable
	var err error

	// Search by lightNovelName or alternateName
	if lightNovelName != "Null" {
		searchResult, err = postgresqldb.LightNovelSearchSubstring(dbConnection, "lightnovel", "name", lightNovelName)
		result = fmt.Sprintf("Search Result for Light Novel Name: %s", lightNovelName)
	} else if alternateName != "Null" {
		searchResult, err = postgresqldb.LightNovelSearchSubstring(dbConnection, "lightnovel", "alt_name", alternateName)
		result = fmt.Sprintf("Search Result for Light Novel Alternate Name: %s", alternateName)
	}
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}

	// Prepare data for the template
	data := struct {
		Result       string
		SearchResult []map[string]any
	}{
		Result:       result,
		SearchResult: searchResult,
	}

	// Load the searchresult page template with the requested data
	tmpl, err := template.ParseFiles("./webfrontend/lightnovel/lightnovelSearchResult.html")
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

// Add Light Novel Entry Handler
func addLightNovelEntryHandler(w http.ResponseWriter, r *http.Request) {
	// Load config
	config, _ := auth.LoadConfig()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract and clean input variables
	lightNovelName := strings.TrimSpace(r.FormValue("ln_name"))
	alternateName := strings.TrimSpace(r.FormValue("ln_alt_name"))
	url := strings.TrimSpace(r.FormValue("url"))

	// typecast the value of the volumes field to an integer
	// the field type is a javescript number (string)
	volumesStr := r.FormValue("volumes")
	volumes, err := strconv.Atoi(volumesStr)
	if err != nil {
		http.Error(w, "addLightNovelEntryHandler - Invalid convesion of volumes value to integer, must be a number", http.StatusBadRequest)
		return
	}

	// Extract boolean fields (use pointers so NULL can be stored)
	var completed *bool

	if r.FormValue("completed") == "on" {
		val := true
		completed = &val
	}

	// Validate input (ensure lightNovelName is provided)
	if lightNovelName == "" {
		http.Error(w, "Light Novel name is required", http.StatusBadRequest)
		return
	}

	// Open database connection
	dbConnection, err := postgresqldb.OpenDatabase(config.PgServer, config.PgPort, config.PgUser, config.PgPassword, config.PgDbName)
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		log.Println("Database connection error:", err)
		return
	}
	defer dbConnection.Close()

	// Add entry to the database and get the new ID
	newID, err := postgresqldb.AddLightNovelRow(dbConnection, lightNovelName, alternateName, url, volumes, completed)
	if err != nil {
		http.Error(w, "Error adding light novel entry to the database", http.StatusInternalServerError)
		log.Println("Error adding entry:", err)
		return
	}
	fmt.Printf("New entry added with ID: %d\n", newID)

	// Query the database using the new ID
	newEntry, err := postgresqldb.LookupByID(dbConnection, "lightnovel", fmt.Sprintf("%d", newID))
	if err != nil {
		http.Error(w, "Error retrieving the added light novel entry from the database", http.StatusInternalServerError)
		log.Println("Error querying for added entry:", err)
		return
	}

	// Load the add anime db entry result.html template
	tmpl, err := template.ParseFiles("./webfrontend/lightnovel/lighnovelAddDbEntryResult.html")
	if err != nil {
		log.Println("Error loading template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	// Prepare data for the template
	data := struct {
		Message string
		Entry   map[string]any
	}{
		Message: fmt.Sprintf("Light Novel entry '%s' was added successfully!", lightNovelName),
		Entry:   newEntry,
	}

	// Send the response to the user
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

////////////// WEBTOONS ACTION HANDLERS

// Webtoons Lookup Handler (query for exact match)
func webtoonQueryHandler(w http.ResponseWriter, r *http.Request) {
	// Load config
	config, _ := auth.LoadConfig()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract and clean input variables
	webtoonName := strings.TrimSpace(r.FormValue("wt_name"))
	alternateName := strings.TrimSpace(r.FormValue("wt_alternate_name"))
	dbId := strings.TrimSpace(r.FormValue("id"))

	// If empty, set to "Null"
	if webtoonName == "" {
		webtoonName = "Null"
	}
	if alternateName == "" {
		alternateName = "Null"
	}
	if dbId == "" {
		dbId = "Null"
	}

	// Open database connection
	dbConnection, _ := postgresqldb.OpenDatabase(config.PgServer, config.PgPort, config.PgUser, config.PgPassword, config.PgDbName)

	// Prepare the response
	var result string
	var queryResult map[string]any

	// Query by webtoons
	if webtoonName != "Null" {
		queryResult, _ = postgresqldb.QueryWithCondition(dbConnection, "webtoons", "name", webtoonName)
		result = fmt.Sprintf("Query Result for Webtoon Name: %s", webtoonName)
	} else if alternateName != "Null" {
		queryResult, _ = postgresqldb.QueryWithCondition(dbConnection, "webtoons", "alt_name", alternateName)
		result = fmt.Sprintf("Query Result for Webtoon Alternate name: %s", alternateName)
	} else if dbId != "Null" {
		queryResult, _ = postgresqldb.QueryWithCondition(dbConnection, "webtoons", "id", dbId)
		result = fmt.Sprintf("Query Result for Webtoon ID: %s", dbId)
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
	tmpl, err := template.ParseFiles("./webfrontend/webtoons/webtoonsQueryResult.html")
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

// webtoon search specified colmun for substring
func webtoonSearchHandler(w http.ResponseWriter, r *http.Request) {
	// Load config
	config, _ := auth.LoadConfig()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract and clean input variables
	webtoonName := strings.TrimSpace(r.FormValue("wt_name"))
	alternateName := strings.TrimSpace(r.FormValue("wt_alternate_name"))

	// If empty, set to "Null"
	if webtoonName == "" {
		webtoonName = "Null"
	}
	if alternateName == "" {
		alternateName = "Null"
	}

	// Open database connection
	dbConnection, _ := postgresqldb.OpenDatabase(config.PgServer, config.PgPort, config.PgUser, config.PgPassword, config.PgDbName)

	// Prepare the response
	var result string
	var searchResult []map[string]any

	// Declare error variable
	var err error

	// Search by lightNovelName or alternateName
	if webtoonName != "Null" {
		searchResult, err = postgresqldb.WebtoonSearchSubstring(dbConnection, "webtoons", "name", webtoonName)
		result = fmt.Sprintf("Search Result for Webtoon Name: %s", webtoonName)
	} else if alternateName != "Null" {
		searchResult, err = postgresqldb.WebtoonSearchSubstring(dbConnection, "webtoons", "alt_name", alternateName)
		result = fmt.Sprintf("Search Result for Webtoon Alternate Name: %s", alternateName)
	}
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}

	// Prepare data for the template
	data := struct {
		Result       string
		SearchResult []map[string]any
	}{
		Result:       result,
		SearchResult: searchResult,
	}

	// Load the searchresult page template with the requested data
	tmpl, err := template.ParseFiles("./webfrontend/webtoons/webtoonsSearchResult.html")
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

// Add webtoon Entry Handler
func addWebtoonEntryHandler(w http.ResponseWriter, r *http.Request) {
	// Load config
	config, _ := auth.LoadConfig()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract and clean input variables
	webtoonName := strings.TrimSpace(r.FormValue("wt_name"))
	alternateName := strings.TrimSpace(r.FormValue("wt_alternate_name"))
	url := strings.TrimSpace(r.FormValue("url"))

	// Extract boolean fields (use pointers so NULL can be stored)
	var completed *bool

	if r.FormValue("completed") == "on" {
		val := true
		completed = &val
	}

	// Validate input (ensure webtoon is provided)
	if webtoonName == "" {
		http.Error(w, "Webtoon name is required", http.StatusBadRequest)
		return
	}

	// Open database connection
	dbConnection, err := postgresqldb.OpenDatabase(config.PgServer, config.PgPort, config.PgUser, config.PgPassword, config.PgDbName)
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		log.Println("Database connection error:", err)
		return
	}
	defer dbConnection.Close()

	// Add entry to the database and get the new ID
	newID, err := postgresqldb.AddWebtoonRow(dbConnection, webtoonName, alternateName, url, completed)
	if err != nil {
		http.Error(w, "Error adding webtoon entry to the database", http.StatusInternalServerError)
		log.Println("Error adding entry:", err)
		return
	}
	fmt.Printf("New entry added with ID: %d\n", newID)

	// Query the database using the new ID
	newEntry, err := postgresqldb.LookupByID(dbConnection, "webtoons", fmt.Sprintf("%d", newID))
	if err != nil {
		http.Error(w, "Error retrieving the added webtoon entry from the database", http.StatusInternalServerError)
		log.Println("Error querying for added entry:", err)
		return
	}

	// Load the add anime db entry result.html template
	tmpl, err := template.ParseFiles("./webfrontend/webtoons/webtoonsAddDbEntryResult.html")
	if err != nil {
		log.Println("Error loading template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	// Prepare data for the template
	data := struct {
		Message string
		Entry   map[string]any
	}{
		Message: fmt.Sprintf("Webtoon entry '%s' was added successfully!", webtoonName),
		Entry:   newEntry,
	}

	// Send the response to the user
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

////////////// WEBNOVEL ACTION HANDLERS

// webnovel Lookup Handler (query for exact match)
func webNovelQueryHandler(w http.ResponseWriter, r *http.Request) {
	// Load config
	config, _ := auth.LoadConfig()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract and clean input variables
	webnovelName := strings.TrimSpace(r.FormValue("wn_name"))
	alternateName := strings.TrimSpace(r.FormValue("wn_alternate_name"))
	dbId := strings.TrimSpace(r.FormValue("id"))

	// If empty, set to "Null"
	if webnovelName == "" {
		webnovelName = "Null"
	}
	if alternateName == "" {
		alternateName = "Null"
	}
	if dbId == "" {
		dbId = "Null"
	}

	// Open database connection
	dbConnection, _ := postgresqldb.OpenDatabase(config.PgServer, config.PgPort, config.PgUser, config.PgPassword, config.PgDbName)

	// Prepare the response
	var result string
	var queryResult map[string]any

	// Query by webtoons
	if webnovelName != "Null" {
		queryResult, _ = postgresqldb.QueryWithCondition(dbConnection, "webnovel", "name", webnovelName)
		result = fmt.Sprintf("Query Result for Webnovel Name: %s", webnovelName)
	} else if alternateName != "Null" {
		queryResult, _ = postgresqldb.QueryWithCondition(dbConnection, "webnovel", "alt_name", alternateName)
		result = fmt.Sprintf("Query Result for Webnovel Alternate name: %s", alternateName)
	} else if dbId != "Null" {
		queryResult, _ = postgresqldb.QueryWithCondition(dbConnection, "webnovel", "id", dbId)
		result = fmt.Sprintf("Query Result for Webnovel ID: %s", dbId)
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
	tmpl, err := template.ParseFiles("./webfrontend/webnovel/webnovelQueryResult.html")
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

func webNovelSearchHandler(w http.ResponseWriter, r *http.Request) {
	// Load config
	config, _ := auth.LoadConfig()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract and clean input variables
	webnovelName := strings.TrimSpace(r.FormValue("wn_name"))
	alternateName := strings.TrimSpace(r.FormValue("wn_alternate_name"))

	// If empty, set to "Null"
	if webnovelName == "" {
		webnovelName = "Null"
	}
	if alternateName == "" {
		alternateName = "Null"
	}

	// Open database connection
	dbConnection, _ := postgresqldb.OpenDatabase(config.PgServer, config.PgPort, config.PgUser, config.PgPassword, config.PgDbName)

	// Prepare the response
	var result string
	var searchResult []map[string]any

	// Declare error variable
	var err error

	// Search by lightNovelName or alternateName
	if webnovelName != "Null" {
		searchResult, err = postgresqldb.WebtoonSearchSubstring(dbConnection, "webnovel", "name", webnovelName)
		result = fmt.Sprintf("Search Result for Webnovel Name: %s", webnovelName)
	} else if alternateName != "Null" {
		searchResult, err = postgresqldb.WebtoonSearchSubstring(dbConnection, "webnovel", "alt_name", alternateName)
		result = fmt.Sprintf("Search Result for Webnovel Alternate Name: %s", alternateName)
	}
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}

	// Prepare data for the template
	data := struct {
		Result       string
		SearchResult []map[string]any
	}{
		Result:       result,
		SearchResult: searchResult,
	}

	// Load the searchresult page template with the requested data
	tmpl, err := template.ParseFiles("./webfrontend/webnovel/webnovelSearchResult.html")
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

func addWebNovelEntryHandler(w http.ResponseWriter, r *http.Request) {
	// Load config
	config, _ := auth.LoadConfig()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract and clean input variables
	webnovelName := strings.TrimSpace(r.FormValue("wn_name"))
	alternateName := strings.TrimSpace(r.FormValue("wn_alternate_name"))
	url := strings.TrimSpace(r.FormValue("url"))

	// Extract boolean fields (use pointers so NULL can be stored)
	var completed *bool

	if r.FormValue("completed") == "on" {
		val := true
		completed = &val
	}

	// Validate input (ensure webnovel is provided)
	if webnovelName == "" {
		http.Error(w, "Webnovel name is required", http.StatusBadRequest)
		return
	}

	// Open database connection
	dbConnection, err := postgresqldb.OpenDatabase(config.PgServer, config.PgPort, config.PgUser, config.PgPassword, config.PgDbName)
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		log.Println("Database connection error:", err)
		return
	}
	defer dbConnection.Close()

	// Add entry to the database and get the new ID
	newID, err := postgresqldb.AddWebnovelRow(dbConnection, webnovelName, alternateName, url, completed)
	if err != nil {
		http.Error(w, "Error adding webnovel entry to the database", http.StatusInternalServerError)
		log.Println("Error adding entry:", err)
		return
	}
	fmt.Printf("New entry added with ID: %d\n", newID)

	// Query the database using the new ID
	newEntry, err := postgresqldb.LookupByID(dbConnection, "webnovel", fmt.Sprintf("%d", newID))
	if err != nil {
		http.Error(w, "Error retrieving the added webnovel entry from the database", http.StatusInternalServerError)
		log.Println("Error querying for added entry:", err)
		return
	}

	// Load the add anime db entry result.html template
	tmpl, err := template.ParseFiles("./webfrontend/webnovel/webnovelAddDbEntryResult.html")
	if err != nil {
		log.Println("Error loading template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	// Prepare data for the template
	data := struct {
		Message string
		Entry   map[string]any
	}{
		Message: fmt.Sprintf("Webnovel entry '%s' was added successfully!", webnovelName),
		Entry:   newEntry,
	}

	// Send the response to the user
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

package webfrontend

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// StartServer initializes and starts the web server on the given port.
func StartServer(port string) {
	http.HandleFunc("/", homePageHandler)
	http.HandleFunc("/query", queryHandler)
	http.HandleFunc("/update", updateHandler)

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

	// Extract input variables
	mangaName := r.FormValue("manga_name")
	alternateName := r.FormValue("alternate_name")
	id := r.FormValue("id")

	// Use a switch statement to handle empty or null values
	switch {
	case mangaName == "":
		mangaName = "Null"
	}
	switch {
	case alternateName == "":
		alternateName = "Null"
	}
	switch {
	case id == "":
		id = "Null"
	}

	// Format the result string
	result := fmt.Sprintf("Variables received - Manga Name: %s, Alternate Name: %s, ID: %s", mangaName, alternateName, id)

	// Respond back with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
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

package models

// NOTE this model is for the MangaDex API seach functionality NOT a mangadex table in the database

// MangadexManga represents a manga from the MangaDex API (no database fields)
type MangadexManga struct {
	MangadexID  string `json:"mangadex_id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	CoverURL    string `json:"cover_url"`
	Status      string `json:"status"` // ongoing, completed, hiatus, cancelled
}

// MangadexAPIResponse represents the raw response from MangaDex API
type MangadexAPIResponse struct {
	Data []MangadexAPIData `json:"data"`
}

// MangadexAPIData represents a single manga entry in the API response
type MangadexAPIData struct {
	ID            string                    `json:"id"`
	Attributes    MangadexAPIAttributes     `json:"attributes"`
	Relationships []MangadexAPIRelationship `json:"relationships"`
}

// MangadexAPIAttributes contains manga metadata
type MangadexAPIAttributes struct {
	Title       map[string]string `json:"title"`
	Description map[string]string `json:"description"`
	Status      string            `json:"status"`
}

// MangadexAPIRelationship contains related entities (author, artist, cover)
type MangadexAPIRelationship struct {
	Type       string                             `json:"type"`
	ID         string                             `json:"id"`
	Attributes *MangadexAPIRelationshipAttributes `json:"attributes,omitempty"`
}

// MangadexAPIRelationshipAttributes contains relationship metadata
type MangadexAPIRelationshipAttributes struct {
	Name     string `json:"name,omitempty"`     // for author/artist
	FileName string `json:"fileName,omitempty"` // for cover art
}

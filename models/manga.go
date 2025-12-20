package models

import "time"

// Manga represents a manga entry in the database (unified model)
type Manga struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	AltTitle    string    `json:"alt_title"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	CoverURL    string    `json:"cover_url"`
	URL         string    `json:"url"`
	MangadexID  string    `json:"mangadex_id"`
	Status      string    `json:"status"` // ongoing, completed, hiatus, cancelled
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

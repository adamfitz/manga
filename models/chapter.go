package models

import "time"

type Chapter struct {
	ID            int       `json:"id"`
	MangaID       int       `json:"manga_id"`
	ChapterNumber string    `json:"chapter_number"`
	Title         string    `json:"title"`
	URL           string    `json:"url"`
	Read          bool      `json:"read"`
	CreatedAt     time.Time `json:"created_at"`
}
